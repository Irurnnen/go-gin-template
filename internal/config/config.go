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
	EnvPrefix         = "shop_api"
	EnvConfigPath     = "CONFIG_PATH"
	EnvConfigName     = "CONFIG_NAME"
)

type (
	Config struct {
		ServerConfig   *ServerConfig   `mapstructure:"server" validate:"required"`
		PostgresConfig *PostgresConfig `mapstructure:"database" validate:"required"`
		Logger         *LoggerConfig   `mapstructure:"logger" validate:"required"`
	}

	PostgresConfig struct {
		Address  string `mapstructure:"address" validate:"required,hostname_port"`
		User     string `mapstructure:"user" validate:"required"`     // TODO: add custom validator
		Password string `mapstructure:"password" validate:"required"` // TODO: add custom validator
		DBName   string `mapstructure:"dbname" validate:"required"`
		Secure   bool   `mapstructure:"secure" validate:"omitempty"`
	}

	ServerConfig struct {
		Address string `mapstructure:"address" validate:"required,hostname_port"`
	}

	LoggerConfig struct {
		Default ComponentLoggerConfig            `mapstructure:"default" validate:"required"`
		Modules map[string]ComponentLoggerConfig `mapstructure:"modules" validate:"omitempty"`
	}

	ComponentLoggerConfig struct {
		Level string `mapstructure:"level" validate:"required,oneof=trace debug info warn error fatal panic"`
	}
)

func (lc *LoggerConfig) GetLoggerConfig(module string) ComponentLoggerConfig {
	if loggerCfg, ok := lc.Modules[module]; ok {
		return loggerCfg
	}
	log.Warn().Str("module", module).Msg("Not found logger config for module")
	return lc.Default
}

func (d *PostgresConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s/%s", d.User, d.Password, d.Address, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}

func Load() (*Config, error) {
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
		return nil, err
	}

	// Unmarshal config
	config := new(Config)
	if err := viperConfig.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
