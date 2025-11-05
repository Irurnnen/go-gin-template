package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type (
	Config struct {
		ServerConfig   *ServerConfig   `mapstructure:"server" validate:"required"`
		DatabaseConfig *DatabaseConfig `mapstructure:"database" validate:"required"`
		LogLevel       string          `mapstructure:"log_level" validate:"omitempty"`
	}

	DatabaseConfig struct {
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

func (d *DatabaseConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}

func NewConfig() *Config {
	viper.AddConfigPath("/run/secrets")
	viper.SetConfigName("go-gin-template")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
		return nil
	}

	config := new(Config)
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
		return nil
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to validate config")
		return nil
	}

	return config
}
