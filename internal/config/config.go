package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	ServerConfig   *ServerConfig   `mapstructure:"server"`
	DatabaseConfig *DatabaseConfig `mapstructure:"database"`
	LogLevel       string          `mapstructure:"log_level"`
	Debug          bool            `mapstructure:"debug"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
	Secure   bool   `mapstructure:"secure"`
}

func (d *DatabaseConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func NewConfigExample() *Config {
	return &Config{
		ServerConfig: &ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		DatabaseConfig: &DatabaseConfig{
			Host:     "hostname",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBName:   "dbname",
		},
		LogLevel: "production",
	}
}

func NewConfig() *Config {
	viper.AddConfigPath("/run/secrets")
	viper.SetConfigName("go-gin-template")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Error("Failed to read config, using example config", zap.Error(err))
		return NewConfigExample()
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		zap.L().Error("Failed to parse config, using example config", zap.Error(err))
		return NewConfigExample()
	}

	return &config
}

func NewConfigDebug() *Config {
	Config := NewConfig()
	Config.Debug = true // Set debug mode to true
	return Config       // Return the modified config
}
