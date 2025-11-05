package repository

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const DriverName = "pgx"

type (
	Repository struct {
		db              *sqlx.DB
		logger          *zerolog.Logger
		HelloRepository HelloRepositoryInterface
	}
)

func NewRepository(DSN string, logger *zerolog.Logger) (*Repository, error) {
	db, err := sqlx.Connect(DriverName, DSN)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:              db,
		logger:          logger,
		HelloRepository: NewHelloRepository(db, logger),
	}, nil
}

func NewRepositoryDB(db *sqlx.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:              db,
		logger:          logger,
		HelloRepository: NewHelloRepository(db, logger),
	}
}

func (r *Repository) Ping() error {
	return r.db.Ping()
}

func (r *Repository) Close() error {
	return r.db.Close()
}
