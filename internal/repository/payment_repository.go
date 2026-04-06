package repository

import (
	"errors"
	"httpServ/internal/model"
)

type PaymentRepo interface {
	Create(p model.Payment) (string, error)
	Get(id string) (model.Payment, error)
	Update(p model.Payment) error
	Delete(id string) error
}

var ErrNotFound = errors.New("not found")
