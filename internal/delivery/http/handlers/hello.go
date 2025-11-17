package handlers

import (
	"context"
	"net/http"

	"github.com/Irurnnen/go-gin-template/internal/delivery/http/dto"
	"github.com/Irurnnen/go-gin-template/internal/services/hello"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	HelloHandler struct {
		service HelloServiceInterface
		logger  *zerolog.Logger
	}

	HelloHandlerInterface interface {
		GetHelloMessage(c *gin.Context)
	}

	HelloServiceInterface interface {
		GetHelloMessage(context.Context) (*hello.Message, error)
	}
)

func NewHelloHandler(service HelloServiceInterface, logger *zerolog.Logger) *HelloHandler {
	return &HelloHandler{
		service: service,
		logger:  logger,
	}
}

func (hh *HelloHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/hello", hh.GetHelloMessage)
}

// GetHelloMessage
//
//	@Summary		Get Hello World message using database
//	@Description	get hello world
//	@Tags			Hello
//	@Produce		json
//	@Success		200	{object}	dto.Message
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/hello [GET]
func (hh *HelloHandler) GetHelloMessage(c *gin.Context) {
	hh.logger.Debug().Msg("Get hello message in handler")

	message, err := hh.service.GetHelloMessage(c.Request.Context())
	switch err {
	case nil:
		break
	default:
		hh.logger.Error().Err(err).Msg("Failed to get hello message")
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.HTTPError{Error: "unknown error", Message: "Unknown internal error"})
		return
	}

	c.JSON(http.StatusOK, dto.Message{Message: message.Message})
}
