package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Irurnnen/go-gin-template/internal/services/hello"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHelloHandler_GetHelloMessage(t *testing.T) {
	// Mock service
	mockService := NewMockHelloServiceInterface(t)
	mockService.EXPECT().
		GetHelloMessage(mock.Anything).
		Return(&hello.Message{Message: "Hello World"}, nil).
		Once()

	// Init handler
	logger := zerolog.Nop()
	handler := NewHelloHandler(mockService, &logger)

	// Init gin server
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
