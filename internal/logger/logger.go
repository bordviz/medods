package logger

import (
	"log/slog"
	"medods/internal/lib/logger/slogpretty"
	"os"
)

func NewLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "docker":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		log = setupPrettyLogger()
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOptions: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
