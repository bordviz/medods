package userstorage

import (
	"log/slog"
	"medods/internal/storage/userstorage/usererror"
)

type UserStorage struct {
	log    *slog.Logger
	errors *usererror.UserError
}

func NewUserStorage(log *slog.Logger) *UserStorage {
	return &UserStorage{
		log:    log,
		errors: usererror.NewUserErrorsHandler(),
	}
}
