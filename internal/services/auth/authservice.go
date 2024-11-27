package authservice

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/jwt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	log        *slog.Logger
	pool       *pgxpool.Pool
	smptServer SMTPServer
	jwtAuth    *jwt.JWTAuth
	userDB     UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *dto.User, requestID string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID, requestID string) (*models.User, error)
	UpdateUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, refreshToken string, requestID string) error
}

type SMTPServer interface {
	SendEmail(email string, message string)
}

func NewAuthService(
	log *slog.Logger,
	pool *pgxpool.Pool,
	jwtAuth *jwt.JWTAuth,
	smptServer SMTPServer,
	userDB UserStorage,
) *AuthService {
	return &AuthService{
		log:        log,
		pool:       pool,
		smptServer: smptServer,
		jwtAuth:    jwtAuth,
		userDB:     userDB,
	}
}
