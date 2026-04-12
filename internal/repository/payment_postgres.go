package repository

import (
	"context"
	"errors"
	"fmt"
	"httpServ/internal/apperror"
	"httpServ/internal/model"
	"httpServ/pkg/db"

	"database/sql"

	"github.com/google/uuid"
)

type RepoPostgres struct {
	db db.DB
}

func NewRepoPostgres(db db.DB) *RepoPostgres {
	return &RepoPostgres{db: db}
}

func (r *RepoPostgres) Create(ctx context.Context, p model.Payment) (string, error) {
	id := uuid.New().String()

	var exists bool

	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM payments WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("payment with id %s: %w", id, apperror.ErrAlreadyExists)
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO payments (id, amount) VALUES ($1, $2)", id, p.Amount)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *RepoPostgres) Get(ctx context.Context, id string) (model.Payment, error) {
	var p model.Payment

	err := r.db.GetContext(ctx, &p, "SELECT id, amount FROM payments WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("payment with id %s: %w", id, apperror.ErrNotFound)
		}
		return p, err
	}

	return p, nil
}

func (r *RepoPostgres) Update(ctx context.Context, p model.Payment) error {
	res, err := r.db.ExecContext(ctx, "UPDATE payments SET amount = $1 where id = $2", p.Amount, p.ID)

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
