package worker

import (
	"context"
	"sync"
	"time"

	"httpServ/internal/model"
	"httpServ/internal/repository"
	"httpServ/pkg/db"

	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, key string, payload []byte) error
}

type Config struct {
	WorkerCount  int
	BatchSize    int
	PollInterval time.Duration
	MaxAttempts  int
	BaseBackoff  time.Duration
	MaxBackoff   time.Duration
}

type OutboxRelay struct {
	repo      repository.OutboxRepo
	publisher Publisher
	cfg       Config
	logger    *zap.Logger
}

func NewOutboxRelay(repo repository.OutboxRepo, publisher Publisher, cfg Config, logger *zap.Logger) *OutboxRelay {
	return &OutboxRelay{
		repo:      repo,
		publisher: publisher,
		cfg:       cfg,
		logger:    logger,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) {
	r.logger.Info("outbox relay starting", zap.Int("workers", r.cfg.WorkerCount))

	var wg sync.WaitGroup
	for i := 0; i < r.cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r.workerPool(ctx, id)
		}(i)
	}

	wg.Wait()
	r.logger.Info("outbox relay stopped")
}

func (r *OutboxRelay) workerPool(ctx context.Context, id int) {
	r.logger.Info("worker started", zap.Int("worker_id", id))
	defer r.logger.Info("worker stopped", zap.Int("worker_id", id))

	ticker := time.NewTicker(r.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.processBatch(ctx, id)
		}
	}
}

func (r *OutboxRelay) processBatch(ctx context.Context, workerID int) {
	tx, err := r.repo.BeginTx(ctx)
	if err != nil {
		r.logger.Error("begin tx", zap.Int("worker_id", workerID), zap.Error(err))
		return
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	events, err := r.repo.FetchBatch(ctx, tx, r.cfg.BatchSize)
	if err != nil {
		r.logger.Error("fetch batch", zap.Int("worker_id", workerID), zap.Error(err))
		return
	}

	for _, ev := range events {
		if err := r.processEvent(ctx, tx, ev); err != nil {
			r.logger.Error("process event",
				zap.String("event_id", ev.ID),
				zap.Error(err))
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("commit", zap.Int("worker_id", workerID), zap.Error(err))
		return
	}
	committed = true
}

func (r *OutboxRelay) processEvent(ctx context.Context, tx db.Tx, ev model.OutboxEvent) error {
	pubErr := r.publisher.Publish(ctx, ev.AggregateID, ev.Payload)
	if pubErr == nil {
		return r.repo.MarkSent(ctx, tx, ev.ID)
	}

	attempt := ev.Attempt + 1
	r.logger.Warn("publish failed",
		zap.String("event_id", ev.ID),
		zap.Int("attempt", attempt),
		zap.Error(pubErr))

	if attempt >= r.cfg.MaxAttempts {
		return r.repo.MarkFailed(ctx, tx, ev.ID, ev.AggregateID, pubErr.Error())
	}

	nextAt := time.Now().Add(NextBackoff(attempt, r.cfg.BaseBackoff, r.cfg.MaxBackoff))
	return r.repo.MarkRetry(ctx, tx, ev.ID, attempt, nextAt, pubErr.Error())
}
