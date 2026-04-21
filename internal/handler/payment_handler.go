package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"httpServ/internal/apperror"
	"httpServ/internal/client/exchangerate"
	"httpServ/internal/model"
	"httpServ/internal/service"
	"httpServ/pkg/validation"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Handler struct {
	service  *service.Service
	validate *validator.Validate
	logger   *zap.Logger
}

func NewHandler(s *service.Service, z *zap.Logger) *Handler {
	return &Handler{
		service:  s,
		validate: validation.Validate,
		logger:   z,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(p); err != nil {
		writeValidationError(w, err)
		return
	}

	h.logger.Info("payment parsed", zap.Any("payment", p))

	id, err := h.service.Create(r.Context(), p)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrAlreadyExists):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		h.logger.Error("internal error", zap.Error(err))
		return
	}

	p.ID = id

	writeJSON(w, h.logger, p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	p := model.Payment{}
	id := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(p); err != nil {
		writeValidationError(w, err)
		return
	}

	h.logger.Info("payment parsed", zap.Any("payment", p))

	p.ID = id

	if err := h.service.Update(r.Context(), p); err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, h.logger, p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	payment, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, h.logger, payment)
}

func (h *Handler) GetInCurrency(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	currency := strings.ToUpper(r.URL.Query().Get("currency"))

	if err := h.validate.Var(currency, "required,currency"); err != nil {
		http.Error(w, "currency must be 3 uppercase letters", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetInCurrency(r.Context(), id, currency)
	if err != nil {
		switch {
		case errors.Is(err, exchangerate.ErrUnsupportedCurrency):
			http.Error(w, "unsupported currency", http.StatusBadRequest)
		case errors.Is(err, exchangerate.ErrUpstream):
			http.Error(w, "rate provider unavailable", http.StatusBadGateway)
		default:
			writeAppError(w, err)
		}
		h.logger.Error("getInCurrency failed", zap.Error(err))
		return
	}

	writeJSON(w, h.logger, resp)
}

func writeJSON(w http.ResponseWriter, logger *zap.Logger, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.Error("encode response", zap.Error(err))
	}
}

func writeAppError(w http.ResponseWriter, err error) {
	var appErr apperror.ErrorResp
	if errors.As(err, &appErr) {
		http.Error(w, appErr.Message, appErr.Status)
		return
	}
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func writeValidationError(w http.ResponseWriter, err error) {
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
}
