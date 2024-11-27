package userstorage

import (
	"context"
	"log/slog"
	"medods/internal/domain/dto"
	"medods/internal/lib/logger/query"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *UserStorage) CreateUser(ctx context.Context, tx pgx.Tx, user *dto.User, requestID string) (uuid.UUID, error) {
	const op = "srotage.userstorage.CreateUser"
	db.log = with.WithOpAndRequestID(db.log, op, requestID)

	q := `
		INSERT INTO public.user 
		(email)
		VALUES ($1)
		RETURNING id;
	`

	db.log.Debug("create user query", query.QueryToString(q))

	var id uuid.UUID
	if err := tx.QueryRow(ctx, q, user.Email).Scan(&id); err != nil {
		if db.errors.IsExists(err) {
			db.log.Error("user with this email already exists", slog.String("email", user.Email))
			return uuid.Nil, db.errors.ErrEmailExists
		}
		db.log.Error("failed to create user", sl.Err(err))
		return uuid.Nil, db.errors.ErrInternalError
	}

	db.log.Debug("new user successfully created", slog.String("id", id.String()))
	return id, nil
}
