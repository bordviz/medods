package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"medods/internal/domain/models"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string, ipAddress string, requestID string) (*models.TokenPair, error) {
	const op = "sercices.auth.RefreshToken"
	s.log = with.WithOpAndRequestID(s.log, op, requestID)

	tokenData, err := s.jwtAuth.DecodeRefreshToken(refreshToken)
	if err != nil {
		s.log.Error("failed to decode refresh token", sl.Err(err))
		return nil, ErrUnauthorized
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", sl.Err(err))
		return nil, ErrInternal
	}

	user, err := s.userDB.GetUserByID(ctx, tx, tokenData.UserID, requestID)
	if err != nil {
		s.log.Error("failed to get user", sl.Err(err))
		return nil, ErrUnauthorized
	}

	if tokenData.IpAddress != ipAddress {
		s.log.Error("invalid IP address", slog.String("ip_address", ipAddress))
		// smtp request
		s.smptServer.SendEmail(user.Email, fmt.Sprintf("Attempted to refresh token from unknown ip address: %s", ipAddress))
		return nil, ErrUnauthorized
	}

	bcryptToken := s.jwtAuth.HashRefreshToken(refreshToken)

	if bcryptToken != user.RefreshToken {
		s.log.Error("invalid refresh token", slog.String("refresh_token", refreshToken))
		// smtp request
		s.smptServer.SendEmail(user.Email, "Attempted to refresh token from unknown token")
		return nil, ErrUnauthorized
	}

	jwtPair, err := s.jwtAuth.CreateTokenPair(user.ID, ipAddress)
	if err != nil {
		s.log.Error("failed to create token pair", sl.Err(err))
		return nil, ErrInternal
	}

	newBcryptToken := s.jwtAuth.HashRefreshToken(jwtPair.RefreshToken)

	if err := s.userDB.UpdateUser(ctx, tx, user.ID, newBcryptToken, requestID); err != nil {
		s.log.Error("failed to update user", sl.Err(err))
		return nil, ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("failed to commit transaction", sl.Err(err))
		return nil, ErrInternal
	}

	s.log.Debug("token pair refreshed successfully")
	return jwtPair, nil
}
