package logger

import (
	"os"

	"github.com/Irurnnen/go-gin-template/internal/config"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

// New get new entity of zerolog logger by logLevel
func New(logConfig config.ComponentLoggerConfig) *zerolog.Logger {
	level := zerolog.InfoLevel

	switch logConfig.Level {
	case "trace":
		level = zerolog.TraceLevel
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel

	}

	logger := zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger()

	return &logger
}
