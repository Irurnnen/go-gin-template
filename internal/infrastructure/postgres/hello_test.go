package postgres

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHelloRepository_GetHelloMessage(t *testing.T) {
	// Mock database
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	// Mock query
	rows := pgxmock.NewRows([]string{"message"}).AddRow("Hello World")
	mock.ExpectQuery("SELECT 'Hello World' AS message").WillReturnRows(rows)

	// Initialize repository
	logger := zerolog.Nop()
	repo := NewHelloRepository(mock, &logger)

	// Call method
	message, err := repo.GetHelloMessage(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", message)
	assert.NoError(t, mock.ExpectationsWereMet())
}
