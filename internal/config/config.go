package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	DefaultConfigPath = "/run/secrets"
	DefaultConfigName = "shop-api"
	ConfigType        = "yml"
	EnvPrefix         = "shop-api"
	EnvConfigPath     = "CONFIG_PATH"
	EnvConfigName     = "CONFIG_NAME"
)

type (
	Config struct {
		ServerConfig   *ServerConfig   `mapstructure:"server" validate:"required"`
		PostgresConfig *PostgresConfig `mapstructure:"database" validate:"required"`
		Logger         *LoggerConfig   `mapstructure:"log_level" validate:"omitempty"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host" validate:"required,hostname"`
		Port     int    `mapstructure:"port" validate:"required,port"`
		User     string `mapstructure:"user" validate:"required"`     // TODO: add custom validator
		Password string `mapstructure:"password" validate:"required"` // TODO: add custom validator
		DBName   string `mapstructure:"db_name" validate:"required"`
		Secure   bool   `mapstructure:"secure" validate:"omitempty"`
	}

	ServerConfig struct {
		Host string `mapstructure:"host" validate:"required,hostname"`
		Port int    `mapstructure:"port" validate:"required,port"`
	}

	LoggerConfig struct {
		Default ComponentLoggerConfig            `mapstructure:"default" validate:"omitempty"`
		Modules map[string]ComponentLoggerConfig `mapstructure:"modules" validate:"omitempty"`
	}

	ComponentLoggerConfig struct {
		Level string `mapstructure:"level" validate:"required,oneof=trace debug info warn error fatal panic"`
	}
)

func (c *Config) GetLoggerConfig(module string) ComponentLoggerConfig {
	if loggerCfg, ok := c.Logger.Modules[module]; ok {
		return loggerCfg
	}
	log.Warn().Str("module", module).Msg("Not found logger config for module")
	return c.Logger.Default
}

func (d *PostgresConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}

func Load() *Config {
	// Initialize Viper
	viperConfig := viper.New()

	// Set env overriding
	viperConfig.AutomaticEnv()
	viperConfig.SetEnvPrefix(EnvPrefix)
	viperConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults for config initialization
	viperConfig.SetDefault("config.name", DefaultConfigName)
	viperConfig.SetDefault("config.path", DefaultConfigPath)
	viperConfig.SetDefault("config.type", ConfigType)

	// Set config path
	viperConfig.AddConfigPath(viperConfig.GetString("config.path"))
	viperConfig.SetConfigType(viperConfig.GetString("config.type"))
	viperConfig.SetConfigName(viperConfig.GetString("config.name"))

	// Read raw config
	if err := viperConfig.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
		return nil
	}

	// Unmarshal config
	config := new(Config)
	if err := viperConfig.Unmarshal(&viperConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
		return nil
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to validate config")
	}

	return config
}
