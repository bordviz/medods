package storage

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"
	"medods/internal/storage/userstorage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	UserStorage
}

type UserStorage interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError)
	GetUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID, requestID string) (*models.User, customerror.CustomError)
	UpdateUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, refreshToken string, requestID string) (*uuid.UUID, customerror.CustomError)
}

func NewStorage(log *slog.Logger) *Storage {
	return &Storage{
		UserStorage: userstorage.NewUserStorage(log),
	}
}
