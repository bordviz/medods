package storageerror

import (
	"medods/internal/lib/customerror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}

var (
	ErrUnexpected = customerror.NewCustomError("unexpected error", 500)
	ErrNotFound   = customerror.NewCustomError("not found", 404)
)

func ErrorHandler(err error, errorMap map[string]customerror.CustomError) customerror.CustomError {
	if isNotFound(err) {
		return ErrNotFound
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return ErrUnexpected
	}

	val, not := errorMap[pgErr.Code]
	if !not {
		return ErrUnexpected
	}
	return val
}
