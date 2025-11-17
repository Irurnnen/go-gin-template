package middlewares

import (
	"github.com/Irurnnen/go-gin-template/internal/delivery/http/dto"
	"github.com/Irurnnen/go-gin-template/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func ErrorHandler(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			logger.Trace().Int("err_count", len(ctx.Errors)).Msg("One or more errors occured")
			if appErr, ok := err.Err.(errors.AppErrorInterface); ok {
				ctx.JSON(appErr.StatusCode(), dto.HTTPError{
					Error:   appErr.Error(),
					Message: appErr.Message(),
				})
				return
			} else {
				logger.Error().Err(err.Err).Msg("Error handler got an error of unknown type")
			}
		}
	}
}
