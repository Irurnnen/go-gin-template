package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerConfig_GetLoggerConfig(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		lc     LoggerConfig
		module string
		want   *ComponentLoggerConfig
	}{
		{
			name: "Get existing module",
			lc: LoggerConfig{
				Default: &ComponentLoggerConfig{
					Level: "panic",
				},
				Modules: map[string]*ComponentLoggerConfig{
					"example_trace": {
						Level: "trace",
					},
				},
			},
			module: "example_trace",
			want: &ComponentLoggerConfig{
				Level: "trace",
			},
		},
		{
			name: "Get non-existing module",
			lc: LoggerConfig{
				Default: &ComponentLoggerConfig{
					Level: "panic",
				},
				Modules: map[string]*ComponentLoggerConfig{
					"example_trace": {
						Level: "trace",
					},
				},
			},
			module: "non_existed_module_name",
			want: &ComponentLoggerConfig{
				Level: "panic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lc.GetLoggerConfig(tt.module)
			assert.True(t, assert.ObjectsAreEqual(got, tt.want), "got incorrect config")
		})
	}
}
