package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Irurnnen/go-gin-template/internal/config"
	"github.com/Irurnnen/go-gin-template/internal/delivery/http"
	"github.com/Irurnnen/go-gin-template/internal/delivery/http/handler"
	"github.com/Irurnnen/go-gin-template/internal/infrastructure/postgres"
	"github.com/Irurnnen/go-gin-template/internal/services/hello"
	"github.com/Irurnnen/go-gin-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

//	@title			go-gin-template
//	@version		0.0.1
//	@description	This is a sample server caller server.
//	@server			http://localhost:8080/v1

func main() {
	// Read config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}

	// Setup logger
	defaultLogger := logger.New(cfg.Logger.Default)
	defaultLogger.Info().Msg("Logger setup successfully")

	// Initialize database
	dbPool, err := pgxpool.New(context.Background(), cfg.PostgresConfig.GetDSN())
	if err != nil {
		defaultLogger.Fatal().Err(err).Msg("Failed to initialize database connection pool")
	}
	defaultLogger.Info().Msg("Database connection pool setup successfully")

	// Ping database
	if err := dbPool.Ping(context.Background()); err != nil {
		defaultLogger.Fatal().Err(err).Str("address", cfg.PostgresConfig.Address).Msg("Failed to ping database")
	}
	defaultLogger.Info().Msg("Database connection ping successfully")

	// Initialize Hello handler
	helloRepositoryLogger := logger.New(cfg.Logger.GetLoggerConfig("hello_repository"))
	helloRepository := postgres.NewHelloRepository(dbPool, helloRepositoryLogger)
	helloServiceLogger := logger.New(cfg.Logger.GetLoggerConfig("hello_service"))
	helloService := hello.NewHelloService(helloRepository, helloServiceLogger)

	helloHandlerLogger := logger.New(cfg.Logger.GetLoggerConfig("hello_handler"))
	helloHandler := handler.NewHelloHandler(helloService, helloHandlerLogger)

	srvLogger := logger.New(cfg.Logger.GetLoggerConfig("http_server"))
	srv := http.New(
		cfg.ServerConfig,
		srvLogger,
		gin.Logger(),   // TODO: write custom logger
		gin.Recovery(), // TODO: write custom recovery
	)
	srv.RegisterRoutes(
		helloHandler.RegisterRoutes,
	)

	// Launch application
	go func() {
		if err := srv.Start(); err != nil {
			defaultLogger.Fatal().Err(err).Msg("Application failed to run")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	defaultLogger.Info().Msg("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		defaultLogger.Error().Err(err).Msg("Failed shutdown server")
	}

	defaultLogger.Info().Msg("Server exciting")
}
