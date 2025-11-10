package hello

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHelloServices_GetHelloMessage(t *testing.T) {
	origin_message := &Message{Message: "Hello World"}
	ctx := context.Background()

	// Mock Repository
	mockRepo := NewMockHelloRepositoryInterface(t)
	mockRepo.EXPECT().
		GetHelloMessage(ctx).
		Return(origin_message, nil).
		Once()

	logger := zerolog.Nop()
	service := NewHelloService(mockRepo, &logger)

	// Call method
	message, err := service.GetHelloMessage(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, origin_message, message)
	mockRepo.AssertExpectations(t)
}
