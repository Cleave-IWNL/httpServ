package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"httpServ/internal/model"
	"httpServ/internal/repository"
	"httpServ/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Printf("Распарсили в структу p %+v\n", p)

	id, err := h.service.Create(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	p.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}
	id := chi.URLParam(r, "id")

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	fmt.Printf("Распарсили в структу p %+v\n", p)

	p.ID = id

	err = h.service.Update(p)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}
	id := chi.URLParam(r, "id")

	payment, err := h.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	p = payment

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}
