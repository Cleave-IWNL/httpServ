package worker

import (
	"context"
	"httpServ/pkg/db"
	"time"

	"go.uber.org/zap"
)

type OutboxRelay struct {
	db     db.DB
	logger *zap.Logger
}

func NewOutboxRelay(db db.DB, zapLog *zap.Logger) *OutboxRelay {
	return &OutboxRelay{db: db, logger: zapLog}
}

func (r *OutboxRelay) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	r.logger.Info("outbox relay started")

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("outbox relay stopped")
			return
		case <-ticker.C:
			r.process(ctx)
		}
	}
}

func (r *OutboxRelay) process(ctx context.Context) {
	r.logger.Debug("checking outbox for pending messages")

	// TODO: SELECT * FROM outbox WHERE sent_at IS NULL LIMIT 10
	// TODO: for each row → отправить в Kafka
	// TODO: UPDATE outbox SET sent_at = NOW() WHERE id = ...

	r.logger.Info("message delivered to kafka",
		zap.String("topic", "payments"),
		zap.String("message_id", "stub-id"))
}
