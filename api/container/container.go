package container

import (
	"discovery-tree/api/handlers"
	"discovery-tree/domain"
	"discovery-tree/infrastructure"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler interfaces for loose coupling between container and handlers
// These interfaces ensure that the container doesn't depend on concrete handler implementations

// TaskHandlerInterface defines the contract for task-related HTTP handlers
type TaskHandlerInterface interface {
	CreateRootTask(c *gin.Context)
	CreateChildTask(c *gin.Context)
	GetTask(c *gin.Context)
	GetAllTasks(c *gin.Context)
	GetRootTask(c *gin.Context)
	GetTaskChildren(c *gin.Context)
	UpdateTask(c *gin.Context)
	UpdateTaskStatus(c *gin.Context)
	MoveTask(c *gin.Context)
	DeleteTask(c *gin.Context)
}

// HealthHandlerInterface defines the contract for health check handlers
type HealthHandlerInterface interface {
	HealthCheck(c *gin.Context)
}

// Handler implementations are now in the handlers package

// Config holds configuration settings for the API server
type Config struct {
	Port         string `json:"port"`
	DataPath     string `json:"dataPath"`
	LogLevel     string `json:"logLevel"`
	EnableCORS   bool   `json:"enableCORS"`
	EnableSwagger bool  `json:"enableSwagger"`
}

// LoadConfigFromEnv loads configuration from environment variables with defaults
func LoadConfigFromEnv() *Config {
	config := &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		DataPath:     getEnvOrDefault("DATA_PATH", "./data/tasks.json"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		EnableCORS:   getEnvBoolOrDefault("ENABLE_CORS", true),
		EnableSwagger: getEnvBoolOrDefault("ENABLE_SWAGGER", true),
	}
	return config
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolOrDefault returns environment variable as bool or default if not set/invalid
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// Container holds all application dependencies and provides dependency injection
// It implements singleton pattern for services to ensure single instances
type Container struct {
	config         *Config
	taskRepository domain.TaskRepository
	taskService    *domain.TaskService
	
	// Singleton instances for handlers (created on first access)
	taskHandler   TaskHandlerInterface
	healthHandler HealthHandlerInterface
	
	// Service lifecycle management
	initialized bool
	shutdown    bool
}

// configureSlog sets up structured logging based on the configuration
func configureSlog(config *Config) error {
	// Parse log level from configuration
	var level slog.Level
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default to info level
	}

	// Create a new slog handler with the configured level
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Use JSON handler for structured logging in production, text handler for development
	var handler slog.Handler
	if gin.Mode() == gin.ReleaseMode {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// Set the default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Logging configured",
		slog.String("level", level.String()),
		slog.String("mode", gin.Mode()),
	)

	return nil
}

// NewContainer creates and initializes a new dependency injection container
// It sets up the repository, services, and handlers with proper dependency injection
func NewContainer(config *Config) (*Container, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Configure structured logging
	if err := configureSlog(config); err != nil {
		return nil, fmt.Errorf("failed to configure logging: %w", err)
	}

	slog.Info("Initializing container",
		slog.String("data_path", config.DataPath),
		slog.String("port", config.Port),
	)

	// Initialize the task repository with the configured data path
	taskRepository, err := infrastructure.NewFileTaskRepository(config.DataPath)
	if err != nil {
		slog.Error("Failed to initialize task repository", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to initialize task repository: %w", err)
	}

	// Initialize the task service with the repository dependency
	taskService := domain.NewTaskService(taskRepository)

	// Create the container with all dependencies
	container := &Container{
		config:         config,
		taskRepository: taskRepository,
		taskService:    taskService,
		initialized:    true,
		shutdown:       false,
	}

	slog.Info("Container initialized successfully")
	return container, nil
}

// NewContainerWithDefaults creates a container with default configuration loaded from environment
func NewContainerWithDefaults() (*Container, error) {
	config := LoadConfigFromEnv()
	return NewContainer(config)
}

// Config returns the configuration
func (c *Container) Config() *Config {
	return c.config
}

// TaskRepository returns the task repository instance
func (c *Container) TaskRepository() domain.TaskRepository {
	return c.taskRepository
}

// TaskService returns the task service instance
func (c *Container) TaskService() *domain.TaskService {
	return c.taskService
}

// GetTaskHandler returns the singleton task handler instance with injected dependencies
// This method implements proper singleton service lifetime management
func (c *Container) GetTaskHandler() TaskHandlerInterface {
	if err := c.ensureNotShutdown(); err != nil {
		panic(err) // Service access after shutdown is a programming error
	}
	
	if c.taskHandler == nil {
		c.taskHandler = handlers.NewTaskHandler(c.taskService, c.taskRepository)
	}
	return c.taskHandler
}

// GetHealthHandler returns the singleton health handler instance
// This method demonstrates singleton service creation without dependencies
func (c *Container) GetHealthHandler() HealthHandlerInterface {
	if err := c.ensureNotShutdown(); err != nil {
		panic(err) // Service access after shutdown is a programming error
	}
	
	if c.healthHandler == nil {
		c.healthHandler = handlers.NewHealthHandler()
	}
	return c.healthHandler
}

// CreateTaskHandler creates a new task handler instance (non-singleton)
// This method is provided for cases where a new instance is explicitly needed
func (c *Container) CreateTaskHandler() TaskHandlerInterface {
	if err := c.ensureNotShutdown(); err != nil {
		panic(err)
	}
	
	return handlers.NewTaskHandler(c.taskService, c.taskRepository)
}

// CreateHealthHandler creates a new health handler instance (non-singleton)
// This method is provided for cases where a new instance is explicitly needed
func (c *Container) CreateHealthHandler() HealthHandlerInterface {
	if err := c.ensureNotShutdown(); err != nil {
		panic(err)
	}
	
	return handlers.NewHealthHandler()
}

// Validate ensures all required dependencies are properly initialized
// This method can be called to verify container state before use
func (c *Container) Validate() error {
	if !c.initialized {
		return fmt.Errorf("container is not initialized")
	}
	if c.shutdown {
		return fmt.Errorf("container has been shut down")
	}
	if c.config == nil {
		return fmt.Errorf("container config is nil")
	}
	if c.taskRepository == nil {
		return fmt.Errorf("task repository is nil")
	}
	if c.taskService == nil {
		return fmt.Errorf("task service is nil")
	}
	return nil
}

// IsInitialized returns whether the container has been properly initialized
func (c *Container) IsInitialized() bool {
	return c.initialized
}

// IsShutdown returns whether the container has been shut down
func (c *Container) IsShutdown() bool {
	return c.shutdown
}

// ensureNotShutdown checks if the container is still active and returns error if shut down
func (c *Container) ensureNotShutdown() error {
	if c.shutdown {
		return fmt.Errorf("container has been shut down, cannot access services")
	}
	if !c.initialized {
		return fmt.Errorf("container is not initialized")
	}
	return nil
}

// GetServiceStatus returns the current status of all managed services
func (c *Container) GetServiceStatus() map[string]interface{} {
	status := map[string]interface{}{
		"initialized":    c.initialized,
		"shutdown":       c.shutdown,
		"taskHandler":    c.taskHandler != nil,
		"healthHandler":  c.healthHandler != nil,
		"taskRepository": c.taskRepository != nil,
		"taskService":    c.taskService != nil,
	}
	return status
}

// ResetSingletons clears all singleton instances, forcing recreation on next access
// This method is useful for testing or when service configuration changes
func (c *Container) ResetSingletons() error {
	if err := c.ensureNotShutdown(); err != nil {
		return err
	}
	
	c.taskHandler = nil
	c.healthHandler = nil
	return nil
}

// Shutdown performs cleanup operations for graceful shutdown
// This method handles resource cleanup and service lifecycle management
func (c *Container) Shutdown() error {
	if c.shutdown {
		return nil // Already shut down
	}
	
	// Mark as shut down to prevent further service access
	c.shutdown = true
	
	// Clear singleton references to help with garbage collection
	c.taskHandler = nil
	c.healthHandler = nil
	
	// Currently no cleanup needed for file repository
	// This method provides extension point for future cleanup needs
	// (e.g., database connections, background workers, etc.)
	
	return nil
}