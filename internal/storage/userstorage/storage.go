package userstorage

import (
	"log/slog"
	"medods/internal/lib/customerror"

	"github.com/jackc/pgerrcode"
)

var (
	ErrEmailExists = customerror.NewCustomError("user with this email already exists", 400)
)

type UserStorage struct {
	log      *slog.Logger
	errorMap map[string]customerror.CustomError
}

func NewUserStorage(log *slog.Logger) *UserStorage {
	return &UserStorage{
		log: log,
		errorMap: map[string]customerror.CustomError{
			pgerrcode.UniqueViolation: ErrEmailExists,
		},
	}
}
