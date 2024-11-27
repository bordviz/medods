package authservice

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
)

func (s *AuthService) CreateUser(ctx context.Context, user *dto.User, requestID string) (uuid.UUID, error) {
	const op = "services.auth.CreateUser"
	s.log = with.WithOpAndRequestID(s.log, op, requestID)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", sl.Err(err))
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	userID, err := s.userDB.CreateUser(ctx, tx, user, requestID)
	if err != nil {
		s.log.Error("failed to create user", sl.Err(err))
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", sl.Err(err))
		return uuid.Nil, err
	}

	s.log.Debug("user created successfully", slog.String("id", userID.String()))
	return userID, nil
}
