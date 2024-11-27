package userstorage

import (
	"context"
	"log/slog"
	"medods/internal/lib/logger/query"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *UserStorage) UpdateUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, refreshToken string, requestID string) error {
	const op = "storage.userstorage.UpdateUser"
	db.log = with.WithOpAndRequestID(db.log, op, requestID)

	q := `
		UPDATE public.user
		SET refresh_token = $1
		WHERE id = $2
		RETURNING id;
	`
	db.log.Debug("update user query", query.QueryToString(q))

	var id uuid.UUID
	if err := tx.QueryRow(ctx, q, refreshToken, userID).Scan(&id); err != nil {
		if db.errors.IsNotFound(err) {
			db.log.Error("user not found", slog.String("id", userID.String()))
			return db.errors.ErrUserNotFound
		}

		db.log.Error("failed to update user", sl.Err(err))
		return err
	}

	db.log.Debug("user successfully updated", slog.String("id", userID.String()))
	return nil
}
