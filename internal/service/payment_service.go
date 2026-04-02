package service

import (
	"httpServ/internal/model"
	"httpServ/internal/repository"
)

type Service struct {
	Repo repository.PaymentRepo
}

func NewService(Repo repository.PaymentRepo) *Service {
	return &Service{Repo: Repo}
}

func (s *Service) Create(p model.Payment) (string, error) {
	return s.Repo.Create(p)
}

func (s *Service) Update(p model.Payment) error {
	return s.Repo.Update(p)
}

func (s *Service) Delete(id string) error {
	return s.Repo.Delete(id)
}

func (s *Service) Get(id string) (model.Payment, error) {
	return s.Repo.Get(id)
}
