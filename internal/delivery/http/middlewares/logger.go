package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewLoggerConfig(logger *zerolog.Logger) gin.LoggerConfig {
	return gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			logger.Info().
				Time("time", params.TimeStamp).
				Int("status_code", params.StatusCode).
				Dur("latency", params.Latency).
				Str("ip", params.ClientIP).
				Str("method", params.Method).
				Str("path", params.Path).
				Str("err", params.ErrorMessage).
				Msg("gin")
			return ""
		},
		Output: zerolog.Nop(),
	}
}
