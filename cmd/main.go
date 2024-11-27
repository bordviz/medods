package main

import (
	"context"
	"log/slog"
	"medods/internal/config"
	"medods/internal/handlers/authhandler"
	"medods/internal/lib/jwt"
	"medods/internal/lib/logger/sl"
	mwLogger "medods/internal/lib/middleware"
	"medods/internal/lib/smtp"
	"medods/internal/logger"
	authservice "medods/internal/services/auth"
	"medods/internal/storage/migrations"
	"medods/internal/storage/postgres"
	"medods/internal/storage/userstorage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.MustLoad("config/docker.yaml")
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Env)
	log.Debug("Debug messages are available")
	log.Info("Info messages are available")
	log.Warn("Warn messages are available")
	log.Error("Error messages are available")

	pool, err := postgres.NewPostgresConnection(context.TODO(), log, &cfg.Database)
	if err != nil {
		panic(err)
	}

	migrationsHandler, err := migrations.NewMigrationHandler(&cfg.Database, log, cfg.MigrationsPath)
	if err != nil {
		panic(err)
	}
	if err := migrationsHandler.Up(); err != nil {
		panic(err)
	}

	userDB := userstorage.NewUserStorage(log)
	jwtAuth := jwt.NewJWTAuth(&cfg.Auth)
	smtpServer := smtp.NewSMTPServer()
	authService := authservice.NewAuthService(log, pool, jwtAuth, smtpServer, userDB)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLogger.New(log))
	log.Info("middleware successfully conected")

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	log.Info("cors successfully conected")

	router.Route("/", authhandler.AddAuthHandlers(log, authService))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info("starting server", slog.String("addr", cfg.HTTPServer.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to listen and serve", sl.Err(err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	stopSignal := <-stop
	log.Info("stoppping server", slog.String("signal", stopSignal.String()))
	ctx, close := context.WithTimeout(context.Background(), time.Minute)
	defer close()
	srv.Shutdown(ctx)
	pool.Close()
	log.Info("server was stopped")
}
