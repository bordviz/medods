package storage

import (
	"context"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	"medods/internal/lib/customerror"
	"medods/internal/storage/storageerror"
	"medods/internal/storage/userstorage"
	"medods/tests/storage/suite"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var st *suite.Suite

func TestInitMigrationsAndSuite(t *testing.T) {
	var err error
	st, err = suite.NewSuite(t, "config/test.yaml")
	require.NoError(t, err)

	err = st.Migrator.Down()
	require.NoError(t, err)

	err = st.Migrator.Up()
	require.NoError(t, err)
}

var fakeusers []*models.User = []*models.User{
	{Email: "test1@example.com"},
	{Email: "test2@example.com"},
}

func TestStorageCreateUser(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
		err  customerror.CustomError
	}{
		{name: "valid email", user: fakeusers[0], err: nil},
		{name: "second valid email", user: fakeusers[1], err: nil},
		{name: "exists email", user: fakeusers[0], err: userstorage.ErrEmailExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := st.DB.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			userID, cerr := st.UserDB.CreateUser(ctx, tx, &dto.User{Email: tt.user.Email}, tt.name)
			require.Equal(t, tt.err, cerr)

			if cerr == nil {
				err = tx.Commit(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, userID)
				tt.user.ID = *userID
			}
		})
	}
}

func TestStorageGetUserByID(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
		err  customerror.CustomError
	}{
		{name: "valid user", user: fakeusers[0], err: nil},
		{name: "not found user", user: &models.User{ID: uuid.New()}, err: storageerror.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := st.DB.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			user, cerr := st.UserDB.GetUserByID(ctx, tx, tt.user.ID, tt.name)
			require.Equal(t, tt.err, cerr)

			if cerr == nil {
				require.Equal(t, tt.user, user)
			}
		})
	}
}

func TestStorageUpdateUser(t *testing.T) {
	tests := []struct {
		name         string
		user         *models.User
		refreshToken string
		err          customerror.CustomError
	}{
		{name: "valid user", user: fakeusers[0], refreshToken: "new-token", err: nil},
		{name: "not found user", user: &models.User{ID: uuid.New()}, refreshToken: "new-token", err: storageerror.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := st.DB.Begin(ctx)
			defer tx.Rollback(ctx)
			require.NoError(t, err)

			userID, cerr := st.UserDB.UpdateUser(ctx, tx, tt.user.ID, tt.refreshToken, tt.name)
			require.Equal(t, tt.err, cerr)

			if cerr == nil {
				require.Equal(t, tt.user.ID, *userID)
				err = tx.Commit(ctx)
				require.NoError(t, err)
				tt.user.RefreshToken = tt.refreshToken
			}
		})
	}
}

func TestStorageGetUserByIDAfterUpdate(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
		err  customerror.CustomError
	}{
		{name: "valid user", user: fakeusers[0], err: nil},
		{name: "not found user", user: &models.User{ID: uuid.New()}, err: storageerror.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tx, err := st.DB.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			user, cerr := st.UserDB.GetUserByID(ctx, tx, tt.user.ID, tt.name)
			require.Equal(t, tt.err, cerr)

			if cerr == nil {
				require.Equal(t, tt.user, user)
			}
		})
	}
}
