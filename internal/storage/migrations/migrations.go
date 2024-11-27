package migrations

import (
	"errors"
	"fmt"
	"log/slog"
	"medods/internal/config"
	"medods/internal/lib/logger/sl"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationsHandler struct {
	log            *slog.Logger
	migrationsPath string
	migrator       *migrate.Migrate
}

func NewMigrationHandler(cfg *config.Database, log *slog.Logger, migrationsPath string) (*MigrationsHandler, error) {
	dsn := createDSN(cfg)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return nil, err
	}
	return &MigrationsHandler{
		log:            log,
		migrationsPath: migrationsPath,
		migrator:       m,
	}, nil
}

func createDSN(cfg *config.Database) string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=migrations",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	return dsn
}

func (m *MigrationsHandler) Up() error {
	if err := m.migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.log.Info("no migrations to apply")
			return nil
		}
		m.log.Error("failed to apply migrations", sl.Err(err))
		return err
	}
	m.log.Info("migrations applied successfully")
	return nil
}

func (m *MigrationsHandler) Down() error {
	if err := m.migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.log.Info("no migrations to revert")
			return nil
		}
		m.log.Error("failed to revert migrations", sl.Err(err))
		return err
	}
	m.log.Info("migrations reverted successfully")
	return nil
}
