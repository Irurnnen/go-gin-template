package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type HelloRepository struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

type HelloRepositoryInterface interface {
	GetHelloMessage() (string, error)
}

func NewHelloRepository(db *sqlx.DB, logger *zerolog.Logger) *HelloRepository {
	return &HelloRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HelloRepository) GetHelloMessage() (string, error) {
	var message string
	query := "SELECT 'Hello World' AS message"
	err := r.db.Get(&message, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute query")
		return "", err
	}
	return message, nil
}
