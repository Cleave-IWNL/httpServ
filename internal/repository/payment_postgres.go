package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"httpServ/internal/apperror"
	"httpServ/internal/model"
	"httpServ/pkg/db"

	"github.com/google/uuid"
)

type RepoPostgres struct {
	db db.DB
}

func NewRepoPostgres(db db.DB) *RepoPostgres {
	return &RepoPostgres{db: db}
}

func (r *RepoPostgres) Create(ctx context.Context, p model.Payment, event model.PaymentCreatedEvent) (id string, err error) {
	id = uuid.New().String()
	event.PaymentID = id

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO payments (id, amount, currency) VALUES ($1, $2, $3)",
		id, p.Amount, p.Currency,
	)

	if err != nil {
		return "", fmt.Errorf("insert payment: %w", err)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("marshal outbox payload: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO outbox_events (aggregate_id, event_type, payload) VALUES ($1, $2, $3)",
		id, event.EventType, payload,
	)

	if err != nil {
		return "", fmt.Errorf("insert outbox event: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("commit tx: %w", err)
	}

	return id, nil
}

func (r *RepoPostgres) Get(ctx context.Context, id string) (model.Payment, error) {
	var p model.Payment

	err := r.db.GetContext(ctx, &p, "SELECT id, amount, currency FROM payments WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("payment with id %s: %w", id, apperror.ErrNotFound)
		}
		return p, err
	}

	return p, nil
}

func (r *RepoPostgres) Update(ctx context.Context, p model.Payment) error {
	res, err := r.db.ExecContext(ctx, "UPDATE payments SET amount = $1, currency = $2 where id = $3", p.Amount, p.Currency, p.ID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment with id %s: %w", p.ID, apperror.ErrNotFound)
	}

	return nil
}

func (r *RepoPostgres) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM payments WHERE id = $1", id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment with id %s: %w", id, apperror.ErrNotFound)
	}

	return nil
}
