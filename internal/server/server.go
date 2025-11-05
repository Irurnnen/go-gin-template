package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Irurnnen/go-gin-template/internal/config"
	"github.com/Irurnnen/go-gin-template/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ServerInterface interface {
	Start() error
	Shutdown() error
}

type Server struct {
	logger     *zerolog.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.ServerConfig, logger *zerolog.Logger, helloHandler handler.HelloHandlerInterface) *Server {
	// Create a new Gin router instance

	router := gin.New()

	// Add middlewares
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		// TODO: Add tracer
	)

	// Setup routes
	v1 := router.Group("/v1")
	{
		internal := v1.Group("/internal")
		{
			internal.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "pong"})
			})
		}

		v1.GET("/hello", helloHandler.GetHelloMessage)
	}

	logger.Debug().Msg("Router initialized")

	return &Server{
		logger: logger,
		httpServer: &http.Server{
			Addr:    cfg.Host + ":" + strconv.Itoa(cfg.Port),
			Handler: router.Handler(),
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info().Str("address", s.httpServer.Addr).Msg("Starting server")
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error().Err(err).Msg("Server failed to start")
		return err
	}
	s.logger.Info().Msg("Server started successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server failed to shutdown")
		return err
	}
	return nil
}
