package authhandler

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	log         *slog.Logger
	authService AuthService
}

type AuthService interface {
	CreateUser(ctx context.Context, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError)
	CreateTokenPair(ctx context.Context, userID uuid.UUID, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError)
	RefreshTokens(ctx context.Context, refreshToken string, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError)
}

func AddAuthHandlers(log *slog.Logger, authService AuthService) func(r chi.Router) {
	handler := &AuthHandler{
		log:         log,
		authService: authService,
	}

	return func(r chi.Router) {
		r.Post("/create_user", handler.CreateUser())
		r.Get("/get_tokens/{userID}", handler.CreateTokenPair())
		r.Get("/refresh_tokens", handler.RefreshTokens())
	}
}
