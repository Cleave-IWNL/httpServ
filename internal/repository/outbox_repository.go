package repository

import (
	"context"
	"httpServ/internal/model"
	"httpServ/pkg/db"
	"time"
)

type OutboxRepo interface {
	BeginTx(ctx context.Context) (db.Tx, error)
	FetchBatch(ctx context.Context, tx db.Tx, limit int) ([]model.OutboxEvent, error)
	MarkSent(ctx context.Context, tx db.Tx, id string) error
	MarkRetry(ctx context.Context, tx db.Tx, id string, attempt int, nextAttemptAt time.Time, errMsg string) error
	MarkFailed(ctx context.Context, tx db.Tx, id, aggregateId, errMsg string) error
}
