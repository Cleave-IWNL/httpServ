package service

import (
	"context"
	"httpServ/internal/client/exchangerate"
	"httpServ/internal/model"
	"httpServ/internal/repository"
)

type RateProvider interface {
	GetRate(ctx context.Context, from, to string) (exchangerate.Rate, error)
}

type Service struct {
	Repo repository.PaymentRepo
	Rate RateProvider
}

func NewService(Repo repository.PaymentRepo, rate RateProvider) *Service {
	return &Service{Repo: Repo, Rate: rate}
}

func (s *Service) GetInCurrency(ctx context.Context, id, target string) (PaymentInCurrency, error) {
	p, err := s.Repo.Get(ctx, id)
	if err != nil {
		return PaymentInCurrency{}, err
	}

	if p.Currency == target {
		return PaymentInCurrency{
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
		return PaymentInCurrency{}, err
	}

	return PaymentInCurrency{
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
