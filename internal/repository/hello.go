package repository

import (
	"context"

	"github.com/Irurnnen/go-gin-template/pkg/postgres"
	"github.com/rs/zerolog"
)

type (
	HelloRepository struct {
		db     postgres.PgxPoolInterface
		logger *zerolog.Logger
	}

	HelloRepositoryInterface interface {
		GetHelloMessage() (string, error)
	}
)

func NewHelloRepository(db postgres.PgxPoolInterface, logger *zerolog.Logger) *HelloRepository {
	return &HelloRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HelloRepository) GetHelloMessage() (string, error) {
	var message string
	query := "SELECT 'Hello World' AS message"
	row := r.db.QueryRow(context.Background(), query)
	err := row.Scan(&message)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute query")
		return "", err
	}
	return message, nil
}
