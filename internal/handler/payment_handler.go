package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"httpServ/internal/apperror"
	"httpServ/internal/model"
	"httpServ/internal/service"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *service.Service
	validate *validator.Validate
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Struct(p); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			var sb strings.Builder

			for _, fe := range validationErrors {
				fmt.Fprintf(&sb, "поле %s: ошибка %s\n", fe.Field(), fe.Tag())
			}

			http.Error(w, sb.String(), http.StatusBadRequest)

			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	fmt.Printf("Распарсили в структу p %+v\n", p)

	id, err := h.service.Create(r.Context(), p)

	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	err := h.service.Delete(r.Context(), id)

	if err != nil {
		var appErr apperror.ErrorResp

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
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
		return
	}

	if err = h.validate.Struct(p); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			var sb strings.Builder

			for _, fe := range validationErrors {
				fmt.Fprintf(&sb, "поле %s: ошибка %s\n", fe.Field(), fe.Tag())
			}

			http.Error(w, sb.String(), http.StatusBadRequest)

			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	fmt.Printf("Распарсили в структу p %+v\n", p)

	p.ID = id

	err = h.service.Update(r.Context(), p)
	if err != nil {
		var appErr apperror.ErrorResp

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
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

	payment, err := h.service.Get(r.Context(), id)
	if err != nil {
		var appErr apperror.ErrorResp

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Message, appErr.Status)
		} else {
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
		service:  s,
		validate: validator.New(),
	}
}
