package authservice

import (
	"context"
	"fmt"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
)

func (s *AuthService) CreateUser(ctx context.Context, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError) {
	const op = "services.auth.CreateUser"
	log := with.WithOpAndRequestID(s.log, op, requestID)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return nil, ErrBeginTx
	}
	defer tx.Rollback(ctx)

	userID, cerr := s.storage.CreateUser(ctx, tx, user, requestID)
	if cerr != nil {
		log.Error("failed to create user", sl.Err(cerr))
		return nil, cerr
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return nil, ErrCommitTx
	}

	log.Debug("user created successfully", slog.String("id", userID.String()))
	return userID, nil
}

func (s *AuthService) CreateTokenPair(ctx context.Context, userID uuid.UUID, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError) {
	const op = "services.auth.CreateTokenPair"
	log := with.WithOpAndRequestID(s.log, op, requestID)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return nil, ErrBeginTx
	}
	defer tx.Rollback(ctx)

	jwtPair, refreshHash, err := s.jwtAuth.CreateTokenPair(userID, ipAddress)
	if err != nil {
		log.Error("failed to create token pair", sl.Err(err))
		return nil, ErrInternal
	}

	id, cerr := s.storage.UpdateUser(ctx, tx, userID, refreshHash, requestID)
	if cerr != nil {
		log.Error("failed to update user", sl.Err(cerr))
		return nil, cerr
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return nil, ErrCommitTx
	}

	log.Debug("token pair created successfully", slog.String("user_id", id.String()))
	return jwtPair, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string, ipAddress string, requestID string) (*models.TokenPair, customerror.CustomError) {
	const op = "sercices.auth.RefreshToken"
	log := with.WithOpAndRequestID(s.log, op, requestID)

	tokenData, refreshHash, err := s.jwtAuth.DecodeRefreshToken(refreshToken)
	if err != nil {
		log.Error("failed to decode refresh token", sl.Err(err))
		return nil, ErrUnauthorized
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return nil, ErrBeginTx
	}
	defer tx.Rollback(ctx)

	user, cerr := s.storage.GetUserByID(ctx, tx, tokenData.UserID, requestID)
	if cerr != nil {
		log.Error("failed to get user", sl.Err(cerr))
		return nil, cerr
	}

	if tokenData.IpAddress != ipAddress {
		log.Error("invalid IP address", slog.String("ip_address", ipAddress))
		// smtp mock request
		s.smptServer.SendEmail(user.Email, fmt.Sprintf("Attempted to refresh token from unknown ip address: %s", ipAddress))
		return nil, ErrUnauthorized
	}

	if refreshHash != user.RefreshToken {
		log.Error("invalid refresh token", slog.String("refresh_token", refreshToken))
		// smtp mock request
		s.smptServer.SendEmail(user.Email, fmt.Sprintf("Attempted to refresh token from unknown ip address: %s", ipAddress))
		return nil, ErrUnauthorized
	}

	jwtPair, newRefreshHash, err := s.jwtAuth.CreateTokenPair(user.ID, ipAddress)
	if err != nil {
		log.Error("failed to create token pair", sl.Err(err))
		return nil, ErrInternal
	}

	id, cerr := s.storage.UpdateUser(ctx, tx, user.ID, newRefreshHash, requestID)
	if cerr != nil {
		log.Error("failed to update user", sl.Err(cerr))
		return nil, ErrInternal
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", sl.Err(err))
		return nil, ErrCommitTx
	}

	log.Debug("token pair refreshed successfully", slog.String("user_id", id.String()))
	return jwtPair, nil
}
