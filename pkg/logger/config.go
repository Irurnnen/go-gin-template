package logger

import "github.com/rs/zerolog/log"

type (
	LoggerConfig struct {
		Default *ComponentLoggerConfig            `mapstructure:"default" validate:"required"`
		Modules map[string]*ComponentLoggerConfig `mapstructure:"modules" validate:"omitempty"`
	}

	ComponentLoggerConfig struct {
		Level string `mapstructure:"level" validate:"required,oneof=trace debug info warn error fatal panic"`
	}
)

func (lc *LoggerConfig) GetLoggerConfig(module string) *ComponentLoggerConfig {
	if loggerCfg, ok := lc.Modules[module]; ok {
		return loggerCfg
	}
	log.Warn().Str("module", module).Msg("Not found logger config for module")
	return lc.Default
}
