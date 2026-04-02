package repository

import (
	"fmt"
	"httpServ/internal/model"

	"sync"

	"github.com/google/uuid"
)

type RepoInMemory struct {
	Payments map[string]model.Payment
	mu       sync.RWMutex
}

func NewRepoInMemory() *RepoInMemory {
	return &RepoInMemory{
		Payments: make(map[string]model.Payment),
	}
}

func (r *RepoInMemory) Create(p model.Payment) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	p.ID = id
	r.Payments[id] = p

	return id, nil
}

func (r *RepoInMemory) Get(id string) (model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.Payments[id]

	if !ok {
		return p, fmt.Errorf("payment with id %s not found", id)
	}

	return p, nil
}

func (r *RepoInMemory) Update(p model.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.Payments[p.ID]

	if !ok {
		return fmt.Errorf("payment with id %s not found", p.ID)
	}

	r.Payments[p.ID] = p

	return nil
}

func (r *RepoInMemory) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.Payments[id]

	if !ok {
		return fmt.Errorf("payment with id %s not found", id)
	}

	delete(r.Payments, id)

	return nil
}
