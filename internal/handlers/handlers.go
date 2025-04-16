package handlers

import (
	"log/slog"
	"medods/internal/handlers/authhandler"
	"medods/internal/services"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	log      *slog.Logger
	services *services.Services
}

func NewHandler(log *slog.Logger, services *services.Services) *Handlers {
	return &Handlers{
		log:      log,
		services: services,
	}
}

func (h *Handlers) Register(r *chi.Mux) {
	r.Route("/auth", authhandler.AddAuthHandlers(h.log, h.services.AuthService))
}
