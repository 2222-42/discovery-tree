package container

import (
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewContainer_ConfiguresSlog(t *testing.T) {
	// Set gin to test mode to avoid release mode JSON logging
	gin.SetMode(gin.TestMode)
	
	config := &Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "debug",
		EnableCORS:   true,
		EnableSwagger: true,
	}

	container, err := NewContainer(config)
	assert.NoError(t, err)
	assert.NotNil(t, container)
	
	// Verify that slog is configured by checking if we can log at debug level
	assert.True(t, slog.Default().Enabled(nil, slog.LevelDebug))
	
	// Clean up test file
	os.Remove("./test_tasks.json")
}

func TestNewContainer_InvalidConfig(t *testing.T) {
	container, err := NewContainer(nil)
	assert.Error(t, err)
	assert.Nil(t, container)
	assert.Contains(t, err.Error(), "config cannot be nil")
}

func TestLoadConfigFromEnv_Defaults(t *testing.T) {
	// Clear environment variables to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATA_PATH")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ENABLE_CORS")
	os.Unsetenv("ENABLE_SWAGGER")
	
	config := LoadConfigFromEnv()
	
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "./data/tasks.json", config.DataPath)
	assert.Equal(t, "info", config.LogLevel)
	assert.True(t, config.EnableCORS)
	assert.True(t, config.EnableSwagger)
}

func TestLoadConfigFromEnv_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATA_PATH", "/custom/path/tasks.json")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENABLE_CORS", "false")
	os.Setenv("ENABLE_SWAGGER", "false")
	
	config := LoadConfigFromEnv()
	
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "/custom/path/tasks.json", config.DataPath)
	assert.Equal(t, "debug", config.LogLevel)
	assert.False(t, config.EnableCORS)
	assert.False(t, config.EnableSwagger)
	
	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("DATA_PATH")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("ENABLE_CORS")
	os.Unsetenv("ENABLE_SWAGGER")
}

func TestConfigureSlog_DifferentLevels(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name     string
		logLevel string
		expected slog.Level
	}{
		{"Debug level", "debug", slog.LevelDebug},
		{"Info level", "info", slog.LevelInfo},
		{"Warn level", "warn", slog.LevelWarn},
		{"Warning level", "warning", slog.LevelWarn},
		{"Error level", "error", slog.LevelError},
		{"Invalid level defaults to info", "invalid", slog.LevelInfo},
		{"Empty level defaults to info", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{LogLevel: tt.logLevel}
			err := configureSlog(config)
			assert.NoError(t, err)
			
			// Check if the expected level is enabled
			assert.True(t, slog.Default().Enabled(nil, tt.expected))
		})
	}
}