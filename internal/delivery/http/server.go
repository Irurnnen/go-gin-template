package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	ServerInterface interface {
		Start() error
		Shutdown() error
	}

	Server struct {
		cfg        *ServerConfig
		logger     *zerolog.Logger
		engine     *gin.Engine
		httpServer *http.Server
	}
)

func New(cfg *ServerConfig, logger *zerolog.Logger, globalMW ...gin.HandlerFunc) *Server {
	// Set gin mode
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()

	//
	e.Use(globalMW...)

	AddDocsForDebugVersion(e)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: e,
	}

	return &Server{
		logger:     logger,
		cfg:        cfg,
		engine:     e,
		httpServer: srv,
	}
}

func (s *Server) RegisterRoutes(regiserFns ...func(*gin.RouterGroup)) {
	api := s.engine.Group("/")
	for _, fn := range regiserFns {
		fn(api)
	}
}

func (s *Server) Start() error {
	// TODO: add tls
	s.logger.Info().Str("address", s.cfg.Address).Msg("HTTP server starting")
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
