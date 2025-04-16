package services

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"
	authservice "medods/internal/services/auth"
	"medods/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	AuthService
}

type SMTPServer interface {
	SendEmail(email string, message string)
}

type JWTAuth interface {
	CreateTokenPair(userID uuid.UUID, ipAddress string) (*models.TokenPair, string, error)
	DecodeAccessToken(token string) (*models.Token, error)
	DecodeRefreshToken(token string) (*models.Token, string, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError)
	CreateTokenPair(ctx context.Context, userID uuid.UUID, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError)
	RefreshTokens(ctx context.Context, refreshToken string, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError)
}

func NewServices(
	log *slog.Logger,
	pool *pgxpool.Pool,
	jwtAuth JWTAuth,
	smptServer SMTPServer,
	storage *storage.Storage,
) *Services {
	return &Services{
		AuthService: authservice.NewAuthService(log, pool, jwtAuth, smptServer, storage),
	}
}
