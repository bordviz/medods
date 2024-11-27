package usererror

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailExists   = errors.New("user with this email already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrInternalError = errors.New("internal error")
)

type UserError struct {
	ErrEmailExists   error
	ErrUserNotFound  error
	ErrInternalError error
}

func NewUserErrorsHandler() *UserError {
	return &UserError{
		ErrEmailExists:   ErrEmailExists,
		ErrUserNotFound:  ErrUserNotFound,
		ErrInternalError: ErrInternalError,
	}
}

func (u *UserError) IsNotFound(err error) bool {
	return err == pgx.ErrNoRows
}

func (u *UserError) IsExists(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
		return true
	}
	return false
}
