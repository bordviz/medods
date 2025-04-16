package userstorage

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"
	"medods/internal/lib/logger/query"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"medods/internal/storage/storageerror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *UserStorage) CreateUser(ctx context.Context, tx pgx.Tx, user *dto.User, requestID string) (*uuid.UUID, customerror.CustomError) {
	const op = "storage.userstorage.CreateUser"
	log := with.WithOpAndRequestID(db.log, op, requestID)

	q := `
		INSERT INTO public.user 
		(email)
		VALUES ($1)
		RETURNING id;
	`

	log.Debug("create user query", query.QueryToString(q))

	var id uuid.UUID
	if err := tx.QueryRow(ctx, q, user.Email).Scan(&id); err != nil {
		log.Error("failed to create user", sl.Err(err))
		return nil, storageerror.ErrorHandler(err, db.errorMap)
	}

	db.log.Debug("new user successfully created", slog.String("id", id.String()))
	return &id, nil
}

func (db *UserStorage) GetUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID, requestID string) (*models.User, customerror.CustomError) {
	const op = "storage.userstorage.GetUserByID"
	log := with.WithOpAndRequestID(db.log, op, requestID)

	q := `
        SELECT id, email, COALESCE(refresh_token, '') as refresh_token
        FROM public.user
        WHERE id = $1;
    `

	log.Debug("get user by ID query", query.QueryToString(q))

	var user models.User
	if err := tx.QueryRow(ctx, q, userID).Scan(&user.ID, &user.Email, &user.RefreshToken); err != nil {
		log.Error("failed to get user by ID", sl.Err(err))
		return nil, storageerror.ErrorHandler(err, db.errorMap)
	}

	log.Debug("user successfully fetched", slog.String("id", userID.String()))
	return &user, nil
}

func (db *UserStorage) UpdateUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, refreshToken string, requestID string) (*uuid.UUID, customerror.CustomError) {
	const op = "storage.userstorage.UpdateUser"
	log := with.WithOpAndRequestID(db.log, op, requestID)

	q := `
		UPDATE public.user
		SET refresh_token = $1
		WHERE id = $2
		RETURNING id;
	`
	log.Debug("update user query", query.QueryToString(q))

	var id uuid.UUID
	if err := tx.QueryRow(ctx, q, refreshToken, userID).Scan(&id); err != nil {
		log.Error("failed to update user", sl.Err(err))
		return nil, storageerror.ErrorHandler(err, db.errorMap)
	}

	log.Debug("user successfully updated", slog.String("id", userID.String()))
	return &id, nil
}
