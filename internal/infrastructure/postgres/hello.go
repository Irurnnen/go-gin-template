package postgres

import (
	"context"

	"github.com/Irurnnen/go-gin-template/internal/services/hello"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *HelloRepository) GetHelloMessage(ctx context.Context) (*hello.Message, error) {
	message := new(Message)
	query := "SELECT 'Hello World' AS message"
	err := pgxscan.Get(ctx, r.db, message, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute query")
		return nil, err
	}
	serviceMessage := &hello.Message{
		Message: message.Message,
	}
	return serviceMessage, nil
}
