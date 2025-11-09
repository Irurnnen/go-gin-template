package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Success(t *testing.T) {
	// Prepare temp directory and config file
	dir := t.TempDir()

	content := `server:
  address: localhost:8080
postgres:
  address: db.local:5432
  user: "user"
  password: "pass"
  dbname: "dbname"
  secure: false
logger:
  default:
    level: "info"
  modules:
    handler:
      level: "debug"
`

	cfgPath := filepath.Join(dir, "shop-api.yml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	keyPath := EnvPrefix + "_CONFIG_PATH"
	keyName := EnvPrefix + "_CONFIG_NAME"
	keyType := EnvPrefix + "_CONFIG_TYPE"

	if err := os.Setenv(keyPath, dir); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyPath)
	if err := os.Setenv(keyName, "shop-api"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyName)
	if err := os.Setenv(keyType, "yml"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyType)

	// Call Load and assert values
	cfg, err := Load()

	// Config
	assert.Nil(t, err)

	// ServerConfig
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.ServerConfig)
	assert.Equal(t, cfg.ServerConfig.Address, "localhost:8080")

	// PostgresConfig
	assert.NotNil(t, cfg.PostgresConfig)
	assert.Equal(t, cfg.PostgresConfig.Address, "db.local:5432")
	assert.Equal(t, cfg.PostgresConfig.DBName, "dbname")

	// logger default
	assert.NotNil(t, cfg.Logger)
	assert.Equal(t, cfg.Logger.Default.Level, "info")

}

func TestLoad_ReadConfigFail(t *testing.T) {
	// point to an empty temp dir (no config file)
	dir := t.TempDir()

	keyPath := EnvPrefix + "_CONFIG_PATH"
	keyName := EnvPrefix + "_CONFIG_NAME"
	keyType := EnvPrefix + "_CONFIG_TYPE"

	if err := os.Setenv(keyPath, dir); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyPath)
	if err := os.Setenv(keyName, "shop-api"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyName)
	if err := os.Setenv(keyType, "yml"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyType)

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoad_UnmarshalFail(t *testing.T) {
	// Write syntactically valid config where 'server' is a string (should be mapping)
	dir := t.TempDir()

	content := `server: "this-should-be-a-map"
database:
  host: "db.local"
  port: 5432
  user: "user"
  password: "pass"
  db_name: "dbname"
log_level:
  default:
    level: "info"
`

	cfgPath := filepath.Join(dir, "shop-api.yml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	keyPath := EnvPrefix + "_CONFIG_PATH"
	keyName := EnvPrefix + "_CONFIG_NAME"
	keyType := EnvPrefix + "_CONFIG_TYPE"

	if err := os.Setenv(keyPath, dir); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyPath)
	if err := os.Setenv(keyName, "shop-api"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyName)
	if err := os.Setenv(keyType, "yml"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyType)

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoad_ValidateFail(t *testing.T) {
	// Missing/invalid server.host to trigger validation error
	dir := t.TempDir()

	content := `server:
  host: ""
  port: 8080
database:
  host: "db.local"
  port: 5432
  user: "user"
  password: "pass"
  db_name: "dbname"
log_level:
  default:
    level: "info"
`

	cfgPath := filepath.Join(dir, "shop-api.yml")
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	keyPath := EnvPrefix + "_CONFIG_PATH"
	keyName := EnvPrefix + "SHOP-API_CONFIG_NAME"
	keyType := EnvPrefix + "SHOP-API_CONFIG_TYPE"

	if err := os.Setenv(keyPath, dir); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyPath)
	if err := os.Setenv(keyName, "shop-api"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyName)
	if err := os.Setenv(keyType, "yml"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.Unsetenv(keyType)

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
}

func TestLoggerConfig_GetLoggerConfig(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		lc     LoggerConfig
		module string
		want   ComponentLoggerConfig
	}{
		{
			name: "Get existing module",
			lc: LoggerConfig{
				Default: ComponentLoggerConfig{
					Level: "panic",
				},
				Modules: map[string]ComponentLoggerConfig{
					"example_trace": {
						Level: "trace",
					},
				},
			},
			module: "example_trace",
			want: ComponentLoggerConfig{
				Level: "trace",
			},
		},
		{
			name: "Get non-existing module",
			lc: LoggerConfig{
				Default: ComponentLoggerConfig{
					Level: "panic",
				},
				Modules: map[string]ComponentLoggerConfig{
					"example_trace": {
						Level: "trace",
					},
				},
			},
			module: "non_existed_module_name",
			want: ComponentLoggerConfig{
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
