package repository

import (
	"context"
	"fmt"
	"httpServ/internal/model"
	"httpServ/pkg/db"
	"time"
)

type OutboxPostgres struct {
	db db.DB
}

func NewOutboxPostgres(d db.DB) *OutboxPostgres {
	return &OutboxPostgres{
		db: d,
	}
}

func (r *OutboxPostgres) BeginTx(ctx context.Context) (db.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

const fetchQuery = `
SELECT id, aggregate_id, event_type, payload, status, attempt, next_attempt_at, last_error
FROM outbox_events
WHERE status = 'pending' AND next_attempt_at <= NOW()
ORDER BY next_attempt_at
LIMIT $1
FOR UPDATE SKIP LOCKED
`

func (r *OutboxPostgres) FetchBatch(ctx context.Context, tx db.Tx, limit int) ([]model.OutboxEvent, error) {
	var events []model.OutboxEvent
	if err := tx.SelectContext(ctx, &events, fetchQuery, limit); err != nil {
		return nil, fmt.Errorf("fetch outbox batch: %w", err)
	}
	return events, nil
}

func (r *OutboxPostgres) MarkSent(ctx context.Context, tx db.Tx, id string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE outbox_events SET status='sent', sent_at=NOW() WHERE id=$1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("mark sent: %w", err)
	}
	return nil
}

func (r *OutboxPostgres) MarkRetry(ctx context.Context, tx db.Tx, id string, attempt int, nextAttemptAt time.Time, errMsg string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE outbox_events SET attempt=$2, next_attempt_at=$3, last_error=$4 WHERE id=$1`,
		id, attempt, nextAttemptAt, errMsg,
	)
	if err != nil {
		return fmt.Errorf("mark retry: %w", err)
	}
	return nil
}

func (r *OutboxPostgres) MarkFailed(ctx context.Context, tx db.Tx, id, aggregateID, errMsg string) error {
	if _, err := tx.ExecContext(ctx,
		`UPDATE outbox_events SET status='failed', last_error=$2 WHERE id=$1`,
		id, errMsg,
	); err != nil {
		return fmt.Errorf("mark outbox failed: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE payments SET status='failed' WHERE id=$1`,
		aggregateID,
	); err != nil {
		return fmt.Errorf("mark payment failed: %w", err)
	}
	return nil
}
