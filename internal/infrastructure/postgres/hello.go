package postgres

import (
	"context"

	"github.com/rs/zerolog"
)

type (
	HelloRepository struct {
		db     PgxPoolInterface
		logger *zerolog.Logger
	}
)

func NewHelloRepository(db PgxPoolInterface, logger *zerolog.Logger) *HelloRepository {
	return &HelloRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HelloRepository) GetHelloMessage(ctx context.Context) (string, error) {
	var message string
	query := "SELECT 'Hello World' AS message"
	row := r.db.QueryRow(ctx, query)
	err := row.Scan(&message)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute query")
		return "", err
	}
	return message, nil
}
