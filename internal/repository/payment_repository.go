package repository

import (
	"context"
	"httpServ/internal/model"
)

type PaymentRepo interface {
	Create(ctx context.Context, p model.Payment, event model.PaymentCreatedEvent) (string, error)
	Get(ctx context.Context, id string) (model.Payment, error)
	Update(ctx context.Context, p model.Payment) error
	Delete(ctx context.Context, id string) error
}
