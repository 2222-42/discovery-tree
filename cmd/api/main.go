// Package main provides the entry point for the Discovery Tree REST API
//
// Configuration:
// The application can be configured using environment variables:
//
//   - PORT: HTTP server port (default: 8080)
//   - DATA_PATH: Path to JSON file for task persistence (default: ./data/tasks.json)
//   - LOG_LEVEL: Logging level - debug, info, warn, error (default: info)
//   - ENABLE_CORS: Enable Cross-Origin Resource Sharing (default: true)
//   - ENABLE_SWAGGER: Enable Swagger/OpenAPI documentation (default: true)
//
// Example usage:
//   export PORT=3000
//   export DATA_PATH=/var/data/tasks.json
//   export LOG_LEVEL=debug
//   go run cmd/api/main.go
//
// @title Discovery Tree API
// @version 1.0
// @description REST API for the Discovery Tree task management system. Provides endpoints to create, read, update, and delete tasks in a hierarchical tree structure.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
package main

import (
	"context"
	"discovery-tree/api/container"
	"discovery-tree/api/server"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Load configuration from environment variables with defaults
	config := container.LoadConfigFromEnv()
	
	// Validate required configuration parameters
	if err := validateConfig(config); err != nil {
		slog.Error("Configuration validation failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Log configuration summary for debugging
	logConfigurationSummary(config)

	// Initialize the dependency injection container with configuration
	c, err := container.NewContainer(config)
	if err != nil {
		slog.Error("Failed to initialize container", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create and start the HTTP server
	srv := server.NewServer(c)
	
	slog.Info("Starting Discovery Tree API server",
		slog.String("port", config.Port),
		slog.String("data_path", config.DataPath),
		slog.String("log_level", config.LogLevel),
		slog.Bool("cors_enabled", config.EnableCORS),
		slog.Bool("swagger_enabled", config.EnableSwagger),
	)
	
	if err := srv.Start(); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	slog.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Stop(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	// Shutdown container
	if err := c.Shutdown(); err != nil {
		slog.Error("Error during container shutdown", slog.String("error", err.Error()))
	}

	slog.Info("Server exited")
}

// validateConfig validates required configuration parameters
func validateConfig(config *container.Config) error {
	if config == nil {
		return fmt.Errorf("configuration is nil")
	}
	
	// Validate port is not empty
	if config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	
	// Validate data path is not empty
	if config.DataPath == "" {
		return fmt.Errorf("data path cannot be empty")
	}
	
	// Validate log level is valid
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	
	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be one of: debug, info, warn, error)", config.LogLevel)
	}
	
	// Ensure data directory exists
	if err := ensureDataDirectory(config.DataPath); err != nil {
		return fmt.Errorf("failed to ensure data directory: %w", err)
	}
	
	return nil
}

// ensureDataDirectory creates the directory for the data path if it doesn't exist
func ensureDataDirectory(dataPath string) error {
	// Extract directory from file path
	dir := filepath.Dir(dataPath)
	
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory with appropriate permissions
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory %s: %w", dir, err)
		}
		slog.Info("Created data directory", slog.String("path", dir))
	}
	
	return nil
}

// logConfigurationSummary logs the current configuration for debugging purposes
func logConfigurationSummary(config *container.Config) {
	slog.Info("Configuration loaded",
		slog.String("port", config.Port),
		slog.String("data_path", config.DataPath),
		slog.String("log_level", config.LogLevel),
		slog.Bool("cors_enabled", config.EnableCORS),
		slog.Bool("swagger_enabled", config.EnableSwagger),
	)
}