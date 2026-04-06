package repository

import (
	"errors"
	"fmt"
	"httpServ/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"database/sql"
)

type RepoPostgres struct {
	db *sqlx.DB
}

func NewRepoPostgres(db *sqlx.DB) *RepoPostgres {
	return &RepoPostgres{db: db}
}

func (r *RepoPostgres) Create(p model.Payment) (string, error) {
	id := uuid.New().String()

	_, err := r.db.Exec("INSERT INTO payments (id, amount) VALUES ($1, $2)", id, p.Amount)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *RepoPostgres) Get(id string) (model.Payment, error) {
	var p model.Payment

	err := r.db.Get(&p, "SELECT id, amount FROM payments WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("payment with id %s: %w", id, ErrNotFound)
		}
		return p, err
	}

	return p, nil
}

func (r *RepoPostgres) Update(p model.Payment) error {
	res, err := r.db.Exec("UPDATE payments SET amount = $1 where id = $2", p.Amount, p.ID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment with id %s: %w", p.ID, ErrNotFound)
	}

	return nil
}

func (r *RepoPostgres) Delete(id string) error {
	res, err := r.db.Exec("DELETE FROM payments WHERE id = $1", id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment with id %s: %w", id, ErrNotFound)
	}

	return nil
}
