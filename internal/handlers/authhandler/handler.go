package authhandler

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
}

type AuthService interface {
	CreateUser(ctx context.Context, user *dto.User, requestID string) (uuid.UUID, error)
	CreateTokenPair(ctx context.Context, userID uuid.UUID, ipAddress string, requestID string) (*models.TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string, ipAddress string, requestID string) (*models.TokenPair, error)
}

func AddAuthHandlers(log *slog.Logger, authService AuthService) func(r chi.Router) {
	handler := &AuthHandler{
		log:         log,
		authService: authService,
	}
	ctx := context.TODO()

	return func(r chi.Router) {
		r.Post("/create", handler.CreateUser(ctx))
		r.Get("/tokens", handler.CreateTokenPair(ctx))
		r.Get("/refresh-tokens", handler.RefreshTokens(ctx))
	}
}
