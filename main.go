package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Irurnnen/go-gin-template/internal/config"
	"github.com/Irurnnen/go-gin-template/internal/handler"
	"github.com/Irurnnen/go-gin-template/internal/repository"
	"github.com/Irurnnen/go-gin-template/internal/server"
	"github.com/Irurnnen/go-gin-template/internal/services"
	"github.com/Irurnnen/go-gin-template/pkg/logger"
)

//	@title			go-gin-template
//	@version		0.0.1
//	@description	This is a sample server caller server.
//	@server			http://localhost:8080/v1

func main() {
	// Read config
	cfg := config.Load()

	// Setup logger
	log := logger.New(cfg.LogLevel)
	log.Info().Msg("Logger setup successfully")

	// Initialize database
	repo, err := repository.NewRepository(cfg.PostgresConfig.GetDSN(), log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	log.Info().Msg("Database setup successfully")
	defer repo.Close()

	// Ping database
	if err := repo.Ping(); err != nil {
		log.Fatal().Err(err).Str("host", cfg.PostgresConfig.Host).Msg("Failed to ping database")
	}
	log.Info().Msg("Database connection ping successfully")

	// Initialize Hello handler
	HelloService := services.NewHelloService(repo.HelloRepository, log)
	HelloHandler := handler.NewHelloHandler(HelloService, log)

	// Setup server
	srv := server.NewServer(cfg.ServerConfig, log, HelloHandler)
	log.Debug().Msg("Server created successfully")

	// Launch application
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("Application failed to run")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server shutdown")
	}

	<-ctx.Done()

	log.Warn().Msg("Timeout of 5 seconds")
	log.Info().Msg("Server exiting")
}
