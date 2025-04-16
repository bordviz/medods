package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"medods/internal/config"
	"medods/internal/lib/logger/sl"
	"medods/internal/lib/logger/with"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(ctx context.Context, log *slog.Logger, cfg *config.Database) (*pgxpool.Pool, error) {
	const op = "storage.postgres.NewPostgresConnection"

	logger := with.WithOp(log, op)

	dsn := createDSN(cfg)

	pool, err := doWithTries(cfg, logger, dsn)
	if err != nil {
		logger.Error("failed connect to database", sl.Err(err))
		return nil, err
	}

	logger.Info("database connection established")
	return pool, nil
}

func createDSN(cfg *config.Database) string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	return dsn
}

func doWithTries(cfg *config.Database, log *slog.Logger, dsn string) (*pgxpool.Pool, error) {
	var err error
	var pool *pgxpool.Pool

	if cfg.MaxAttempts == 0 {
		return nil, fmt.Errorf("max attempts must be greater than 0")
	}

	for range cfg.MaxAttempts {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.AttempTimeout)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Warn("postgresql connection attemt failed", sl.Err(err))
			time.Sleep(cfg.AttempDelay)
			continue
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Warn("failed to ping connection", sl.Err(err))
			time.Sleep(cfg.AttempDelay)
			continue
		}
		return pool, nil
	}

	return nil, err
}
