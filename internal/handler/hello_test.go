package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHelloService struct {
	mock.Mock
}

func (m *MockHelloService) GetHelloMessage(ctx context.Context) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestHelloHandler_GetHelloMessage(t *testing.T) {
	// Mock service
	mockService := new(MockHelloService)
	mockService.On("GetHelloMessage").Return("Hello World", nil)

	// Init handler
	logger := zerolog.Nop()
	handler := NewHelloHandler(mockService, &logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/v1/hello", handler.GetHelloMessage)

	// Perform request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/hello", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"Hello World"}`, w.Body.String())
	mockService.AssertExpectations(t)
}
