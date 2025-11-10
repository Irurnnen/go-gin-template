package postgres

import "fmt"

type (
	PostgresConfig struct {
		Address  string `mapstructure:"address" validate:"required,hostname_port"`
		User     string `mapstructure:"user" validate:"required"`     // TODO: add custom validator
		Password string `mapstructure:"password" validate:"required"` // TODO: add custom validator
		DBName   string `mapstructure:"dbname" validate:"required"`
		Secure   bool   `mapstructure:"secure" validate:"omitempty"`
	}
)

func (d *PostgresConfig) GetDSN() string {
	DSN := fmt.Sprintf("postgresql://%s:%s@%s/%s", d.User, d.Password, d.Address, d.DBName)
	if d.Secure {
		return DSN
	}
	return DSN + "?sslmode=disable"
}
