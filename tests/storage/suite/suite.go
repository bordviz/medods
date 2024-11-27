package suite

import (
	"context"
	"medods/internal/config"
	"medods/internal/lib/logger/slogdiscard"
	"medods/internal/storage/migrations"
	"medods/internal/storage/postgres"
	"medods/internal/storage/userstorage"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Suite struct {
	DB       *pgxpool.Pool
	Migrator *migrations.MigrationsHandler
	UserDB   *userstorage.UserStorage
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

	return &Suite{
		DB:       db,
		Migrator: migrator,
		UserDB:   userDB,
	}, nil
}
