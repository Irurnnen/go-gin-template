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
		LogLevel       string          `mapstructure:"log_level" validate:"omitempty"`
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
)

func (d *PostgresConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}

func Load() *Config {
	// Initialize Viper
	viper_config := viper.New()

	// Set env overriding
	viper_config.AutomaticEnv()
	viper_config.SetEnvPrefix(EnvPrefix)
	viper_config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults for config initialization
	viper_config.SetDefault("config.name", DefaultConfigName)
	viper_config.SetDefault("config.path", DefaultConfigPath)
	viper_config.SetDefault("config.type", ConfigType)

	// Set config path
	viper_config.AddConfigPath(viper_config.GetString("config.path"))
	viper_config.SetConfigType(viper_config.GetString("config.type"))
	viper_config.SetConfigName(viper_config.GetString("config.name"))

	// Read raw config
	if err := viper_config.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
		return nil
	}

	// Unmarshal config
	config := new(Config)
	if err := viper_config.Unmarshal(&viper_config); err != nil {
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
