package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHelloRepository_GetHelloMessage(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock query
	rows := sqlmock.NewRows([]string{"message"}).AddRow("Hello World")
	mock.ExpectQuery("SELECT 'Hello World' AS message").WillReturnRows(rows)

	// Initialize repository
	logger := zerolog.Nop()
	repo := NewHelloRepository(sqlx.NewDb(db, "sqlmock"), &logger)

	// Call method
	message, err := repo.GetHelloMessage()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
