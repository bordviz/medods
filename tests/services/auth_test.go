package services

import (
	"context"
	"medods/internal/domain/dto"
	"medods/internal/domain/models"
	authservice "medods/internal/services/auth"
	"medods/internal/storage/userstorage/usererror"
	"medods/tests/services/suite"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var st *suite.Suite

func TestInitMigrations(t *testing.T) {
	var err error
	st, err = suite.NewSuite(t, "config/test.yaml")
	require.NoError(t, err)

	err = st.Migrator.Down()
	require.NoError(t, err)

	err = st.Migrator.Up()
	require.NoError(t, err)
}

type UserWithTokens struct {
	user      *models.User
	tokenPair *models.TokenPair
}

var fakeusers []*UserWithTokens = []*UserWithTokens{
	{
		user:      &models.User{Email: "test1@example.com"},
		tokenPair: &models.TokenPair{},
	},
	{
		user:      &models.User{Email: "test2@example.com"},
		tokenPair: &models.TokenPair{},
	},
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name string
		data *UserWithTokens
		err  error
	}{
		{name: "valid email", data: fakeusers[0]},
		{name: "second valid email", data: fakeusers[1]},
		{name: "exists email", data: fakeusers[0], err: usererror.ErrEmailExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			id, err := st.AuthService.CreateUser(ctx, &dto.User{Email: tt.data.user.Email}, tt.name)
			require.Equal(t, tt.err, err)

			if err == nil {
				require.NotEqual(t, uuid.Nil, id)
				tt.data.user.ID = id
			}
		})
	}
}

func TestCreateTokenPair(t *testing.T) {
	tests := []struct {
		name string
		data *UserWithTokens
		err  error
	}{
		{name: "valid user", data: fakeusers[0]},
		{name: "second valid user", data: fakeusers[1]},
		{name: "user not found", data: &UserWithTokens{user: &models.User{ID: uuid.New()}}, err: usererror.ErrUserNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tokenPair, err := st.AuthService.CreateTokenPair(ctx, tt.data.user.ID, "127.0.0.1", tt.name)
			require.Equal(t, tt.err, err)

			if err == nil {
				require.NotNil(t, tokenPair)
				tt.data.tokenPair = tokenPair
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		data      *UserWithTokens
		ipAddress string
		err       error
	}{
		{name: "valid user", data: fakeusers[0], ipAddress: "127.0.0.1"},
		{name: "second valid user", data: fakeusers[1], ipAddress: "127.0.0.1"},
		{name: "invalid token", data: &UserWithTokens{user: fakeusers[0].user, tokenPair: &models.TokenPair{RefreshToken: "dadada"}}, ipAddress: "127.0.0.1", err: authservice.ErrUnauthorized},
		{name: "invalid ip address", data: fakeusers[1], ipAddress: "192.163.1.1", err: authservice.ErrUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tokenPair, err := st.AuthService.RefreshTokens(ctx, tt.data.tokenPair.RefreshToken, tt.ipAddress, tt.name)
			require.Equal(t, tt.err, err)

			if err == nil {
				require.NotNil(t, tokenPair)
			}
		})
	}
}
