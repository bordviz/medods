package logger

import (
	"fmt"
	"log/slog"
	"medods/internal/lib/logger/slogpretty"
	"os"
)

func NewLogger(env string) (*slog.Logger, error) {
	var log *slog.Logger

	switch env {
	case "test":
		log = slog.New(slog.DiscardHandler)
	case "local":
		log = setupPrettyLogger()
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return nil, fmt.Errorf("invalid environment: %s, expected: test, local, dev, prod", env)
	}

	return log, nil
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOptions: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
