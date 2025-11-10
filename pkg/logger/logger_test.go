package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew_LogLevel(t *testing.T) {
	tests := []struct {
		name      string
		logConfig *ComponentLoggerConfig
		want      zerolog.Level
	}{
		{
			name: "Trace",
			logConfig: &ComponentLoggerConfig{
				Level: "trace",
			},
			want: zerolog.TraceLevel,
		},
		{
			name: "Debug",
			logConfig: &ComponentLoggerConfig{
				Level: "debug",
			},
			want: zerolog.DebugLevel,
		},
		{
			name: "Info",
			logConfig: &ComponentLoggerConfig{
				Level: "info",
			},
			want: zerolog.InfoLevel,
		},
		{
			name: "Warning",
			logConfig: &ComponentLoggerConfig{
				Level: "warn",
			},
			want: zerolog.WarnLevel,
		},
		{
			name: "Error",
			logConfig: &ComponentLoggerConfig{
				Level: "error",
			},
			want: zerolog.ErrorLevel,
		},
		{
			name: "Fatal",
			logConfig: &ComponentLoggerConfig{
				Level: "fatal",
			},
			want: zerolog.FatalLevel,
		},
		{
			name: "Panic",
			logConfig: &ComponentLoggerConfig{
				Level: "panic",
			},
			want: zerolog.PanicLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.logConfig)

			assert.Equal(t, got.GetLevel(), tt.want)
		})
	}
}
