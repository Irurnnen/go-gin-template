package handler

import (
	"net/http"

	"github.com/Irurnnen/go-gin-template/internal/models"
	"github.com/Irurnnen/go-gin-template/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	HelloHandler struct {
		service services.HelloServiceInterface
		logger  *zerolog.Logger
	}

	HelloHandlerInterface interface {
		GetHelloMessage(c *gin.Context)
	}
)

func NewHelloHandler(service services.HelloServiceInterface, logger *zerolog.Logger) *HelloHandler {
	return &HelloHandler{
		service: service,
		logger:  logger,
	}
}

// GetHelloMessage
//
//	@Summary		Get Hello World message using database
//	@Description	get hello world
//	@Tags			Hello
//	@Produce		json
//	@Success		200	{object}	models.Message
//	@Failure		500	{object}	models.HTTPError
//	@Router			/hello [GET]
func (hh *HelloHandler) GetHelloMessage(c *gin.Context) {
	hh.logger.Debug().Msg("Get hello message in handler")

	message, err := hh.service.GetHelloMessage()
	switch err {
	case nil:
		break
	default:
		hh.logger.Error().Err(err).Msg("Failed to get hello message")
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.HTTPError{Error: "unknown error", Message: "Unknown internal error"})
		return
	}

	c.JSON(http.StatusOK, models.Message{Message: message})
}
