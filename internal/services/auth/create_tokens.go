package authservice

import (
	"context"
	"medods/internal/domain/models"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
)

func (s *AuthService) CreateTokenPair(ctx context.Context, userID uuid.UUID, ipAddress string, requestID string) (*models.TokenPair, error) {
	const op = "services.auth.CreateTokenPair"
	s.log = with.WithOpAndRequestID(s.log, op, requestID)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", sl.Err(err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	jwtPair, err := s.jwtAuth.CreateTokenPair(userID, ipAddress)
	if err != nil {
		s.log.Error("failed to create token pair", sl.Err(err))
		return nil, err
	}

	cryptToken := s.jwtAuth.HashRefreshToken(jwtPair.RefreshToken)

	if err := s.userDB.UpdateUser(ctx, tx, userID, cryptToken, requestID); err != nil {
		s.log.Error("failed to update user", sl.Err(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", sl.Err(err))
		return nil, err
	}

	s.log.Debug("token pair created successfully")
	return jwtPair, nil
}
