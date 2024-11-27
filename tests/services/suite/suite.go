package suite

import (
	"context"
	"medods/internal/config"
	"medods/internal/lib/jwt"
	"medods/internal/lib/logger/slogdiscard"
	"medods/internal/lib/smtp"
	authservice "medods/internal/services/auth"
	"medods/internal/storage/migrations"
	"medods/internal/storage/postgres"
	"medods/internal/storage/userstorage"
	"testing"
)

type Suite struct {
	AuthService *authservice.AuthService
	Migrator    *migrations.MigrationsHandler
}

func NewSuite(t *testing.T, configPath string) (*Suite, error) {
	t.Helper()

	cfg, err := config.MustLoad("../../" + configPath)
	if err != nil {
		return nil, err
	}

	log := slogdiscard.NewDiscardLogger()
	db, err := postgres.NewPostgresConnection(context.Background(), log, &cfg.Database)
	if err != nil {
		return nil, err
	}

	migrator, err := migrations.NewMigrationHandler(&cfg.Database, log, "../../"+cfg.MigrationsPath)
	if err != nil {
		return nil, err
	}

	userDB := userstorage.NewUserStorage(log)
	authService := authservice.NewAuthService(log, db, jwt.NewJWTAuth(&cfg.Auth), smtp.NewSMTPServer(), userDB)

	return &Suite{
		AuthService: authService,
		Migrator:    migrator,
	}, nil
}
