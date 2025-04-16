package authservice

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	log        *slog.Logger
	pool       *pgxpool.Pool
	smptServer SMTPServer
	jwtAuth    JWTAuth
	storage    UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError)
	GetUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID, requestID string) (*models.User, customerror.CustomError)
	UpdateUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, refreshToken string, requestID string) (*uuid.UUID, customerror.CustomError)
}

type SMTPServer interface {
	SendEmail(email string, message string)
}

type JWTAuth interface {
	CreateTokenPair(userID uuid.UUID, ipAddress string) (*models.TokenPair, string, error)
	DecodeAccessToken(token string) (*models.Token, error)
	DecodeRefreshToken(token string) (*models.Token, string, error)
}

var (
	ErrBeginTx      = customerror.NewCustomError("failed to begin transaction", 500)
	ErrCommitTx     = customerror.NewCustomError("failed to commit transaction", 500)
	ErrUnauthorized = customerror.NewCustomError("unauthorized", 401)
	ErrInternal     = customerror.NewCustomError("internal error", 500)
)

func NewAuthService(
	log *slog.Logger,
	pool *pgxpool.Pool,
	jwtAuth JWTAuth,
	smptServer SMTPServer,
	storage UserStorage,
) *AuthService {
	return &AuthService{
		log:        log,
		pool:       pool,
		smptServer: smptServer,
		jwtAuth:    jwtAuth,
		storage:    storage,
	}
}
