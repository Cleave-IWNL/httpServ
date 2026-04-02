package handler

import "github.com/go-chi/chi"

func NewRouter(h *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/payment", h.Create)
	r.Get("/payment/{id}", h.Get)
	r.Put("/payment/{id}", h.Update)
	r.Delete("/payment/{id}", h.Delete)
	return r
}
