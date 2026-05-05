package service

import (
	"context"
	"time"

	"httpServ/internal/model"
	"httpServ/internal/repository"

	"github.com/google/uuid"
)

type RateProvider interface {
	GetRate(ctx context.Context, from, to string) (model.Rate, error)
}

type Service struct {
	Repo repository.PaymentRepo
	Rate RateProvider
}

func NewService(repo repository.PaymentRepo, rate RateProvider) *Service {
	return &Service{Repo: repo, Rate: rate}
}

func (s *Service) GetInCurrency(ctx context.Context, id, target string) (model.PaymentInCurrency, error) {
	p, err := s.Repo.Get(ctx, id)
	if err != nil {
		return model.PaymentInCurrency{}, err
	}

	if p.Currency == target {
		return model.PaymentInCurrency{
			ID:               p.ID,
			OriginalAmount:   p.Amount,
			OriginalCurrency: p.Currency,
			TargetCurrency:   target,
			ConvertedAmount:  float64(p.Amount),
			Rate:             1.0,
		}, nil
	}

	rate, err := s.Rate.GetRate(ctx, p.Currency, target)
	if err != nil {
		return model.PaymentInCurrency{}, err
	}

	return model.PaymentInCurrency{
		ID:               p.ID,
		OriginalAmount:   p.Amount,
		OriginalCurrency: p.Currency,
		TargetCurrency:   target,
		ConvertedAmount:  float64(p.Amount) * rate.Value,
		Rate:             rate.Value,
		RateFetchedAt:    rate.FetchedAt,
	}, nil
}

func (s *Service) Create(ctx context.Context, p model.Payment) (string, error) {
	event := model.PaymentCreatedEvent{
		EventID:       uuid.NewString(),
		EventType:     model.PaymentCreatedType,
		OccurredAt:    time.Now().UTC(),
		SchemaVersion: 1,
		Amount:        p.Amount,
		Currency:      p.Currency,
	}

	return s.Repo.Create(ctx, p, event)
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
