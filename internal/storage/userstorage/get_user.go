package userstorage

import (
	"context"
	"log/slog"
	"medods/internal/domain/models"
	"medods/internal/lib/logger/query"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *UserStorage) GetUserByID(ctx context.Context, tx pgx.Tx, userID uuid.UUID, requestID string) (*models.User, error) {
	const op = "srotage.userstorage.GetUserByID"
	db.log = with.WithOpAndRequestID(db.log, op, requestID)

	q := `
        SELECT id, email, COALESCE(refresh_token, '') AS refresh_token
        FROM public.user
        WHERE id = $1;
    `

	db.log.Debug("get user by ID query", query.QueryToString(q))

	var user models.User
	if err := tx.QueryRow(ctx, q, userID).Scan(&user.ID, &user.Email, &user.RefreshToken); err != nil {
		if db.errors.IsNotFound(err) {
			db.log.Error("user not found", slog.String("id", userID.String()))
			return nil, db.errors.ErrUserNotFound
		}

		db.log.Error("failed to get user by ID", sl.Err(err))
		return nil, err
	}

	db.log.Debug("user successfully fetched", slog.String("id", userID.String()))
	return &user, nil
}
