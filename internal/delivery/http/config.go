package http

type (
	ServerConfig struct {
		Address string `mapstructure:"address" validate:"required,hostname_port"`
	}
)
