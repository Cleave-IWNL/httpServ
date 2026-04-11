package service

import (
	"context"
	"httpServ/internal/model"
	"httpServ/internal/repository"
)

type Service struct {
	Repo repository.PaymentRepo
}

func NewService(Repo repository.PaymentRepo) *Service {
	return &Service{Repo: Repo}
}

func (s *Service) Create(ctx context.Context, p model.Payment) (string, error) {
	return s.Repo.Create(ctx, p)
}

func (s *Service) Update(ctx context.Context, p model.Payment) error {
	return s.Repo.Update(ctx, p)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id string) (model.Payment, error) {
	return s.Repo.Get(ctx, id)
}
