package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name string
		pc   PostgresConfig
		want string
	}{
		{
			name: "insecure",
			pc: PostgresConfig{
				Address:  "localhost:8080",
				User:     "irc",
				Password: "pass",
				DBName:   "checkout",
				Secure:   false,
			},
			want: "postgresql://irc:pass@localhost:8080/checkout?sslmode=disable",
		},
		{
			name: "secure",
			pc: PostgresConfig{
				User:     "postgres",
				Password: "passw0rd",
				DBName:   "data",
				Secure:   true,
				Address:  "postgres:5432",
			},
			want: "postgresql://postgres:passw0rd@postgres:5432/data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pc.GetDSN()
			assert.Equal(t, got, tt.want)
		})
	}
}
