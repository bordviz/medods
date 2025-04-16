package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"medods/docs"
	"medods/internal/config"
	"medods/internal/handlers"
	"medods/internal/lib/jwt"
	"medods/internal/lib/logger/sl"
	mwLogger "medods/internal/lib/middleware"
	"medods/internal/lib/smtp"
	"medods/internal/logger"
	"medods/internal/services"
	"medods/internal/storage"
	"medods/internal/storage/migrations"
	"medods/internal/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Medods API Service
//	@version		2.0.0
//
// @BasePath	/
func main() {
	configPath := flag.String("config", "config/docker.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.MustLoad(*configPath)
	if err != nil {
		fmt.Printf("failed to load config: %s\n", err.Error())
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfg.Env)
	if err != nil {
		fmt.Printf("failed to create logger: %s\n", err.Error())
		os.Exit(1)
	}

	log.Debug("Debug messages are available")
	log.Info("Info messages are available")
	log.Warn("Warn messages are available")
	log.Error("Error messages are available")

	pool, err := postgres.NewPostgresConnection(context.TODO(), log, &cfg.Database)
	if err != nil {
		log.Error("failed to connect to postgres", sl.Err(err))
		os.Exit(1)
	}
	defer pool.Close()

	migrationsHandler, err := migrations.NewMigrationHandler(&cfg.Database, log, cfg.MigrationsPath)
	if err != nil {
		log.Error("failed to create migrations handler", sl.Err(err))
		os.Exit(1)
	}
	if err := migrationsHandler.Up(); err != nil {
		log.Error("failed to up migrations", sl.Err(err))
		os.Exit(1)
	}

	storage := storage.NewStorage(log)
	jwtAuth := jwt.NewJWTAuth(&cfg.Auth)
	smtpServer := smtp.NewSMTPServer()
	services := services.NewServices(log, pool, jwtAuth, smtpServer, storage)
	handlers := handlers.NewHandler(log, services)

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
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

	handlers.Register(router)
	if cfg.Env != "prod" {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
		router.Mount("/swagger", httpSwagger.WrapHandler)
	}
	log.Info("handlers successfully conected")

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Info("starting server", slog.String("addr", fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)))
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
	log.Info("server was stopped")
}
