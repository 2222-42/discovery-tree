package server

import (
	"context"
	"discovery-tree/api/container"
	"discovery-tree/api/middleware"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server with all its dependencies
type Server struct {
	container  *container.Container
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer creates a new HTTP server with the given container
func NewServer(container *container.Container) *Server {
	// Validate container before proceeding
	if err := container.Validate(); err != nil {
		slog.Error("Invalid container provided to server", slog.String("error", err.Error()))
		panic(fmt.Sprintf("invalid container: %v", err))
	}

	server := &Server{
		container: container,
	}

	// Initialize the Gin engine with middleware and routes
	server.setupEngine()
	
	return server
}

// setupEngine initializes the Gin engine with middleware and routes
func (s *Server) setupEngine() {
	// Set Gin mode based on environment
	if s.container.Config().LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	s.engine = gin.New()

	// Add global middleware
	s.setupMiddleware()

	// Setup API routes
	s.setupRoutes()

	slog.Info("Gin engine configured successfully")
}

// setupMiddleware configures all middleware for the server
func (s *Server) setupMiddleware() {
	// Recovery middleware (should be first)
	s.engine.Use(middleware.ErrorHandler())

	// Logging middleware
	s.engine.Use(middleware.Logger())
	s.engine.Use(middleware.ErrorLogger())

	// CORS middleware (if enabled)
	if s.container.Config().EnableCORS {
		s.engine.Use(middleware.CORS())
		slog.Info("CORS middleware enabled")
	}

	slog.Info("Middleware configured successfully")
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Use the centralized route setup
	SetupRoutes(s.engine, s.container)
}

// Start starts the HTTP server on the configured port
func (s *Server) Start() error {
	config := s.container.Config()
	
	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         ":" + config.Port,
		Handler:      s.engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Starting HTTP server",
		slog.String("port", config.Port),
		slog.String("address", s.httpServer.Addr),
	)

	// Start server in a goroutine so it doesn't block
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start HTTP server", slog.String("error", err.Error()))
		}
	}()

	slog.Info("HTTP server started successfully")
	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	slog.Info("Shutting down HTTP server...")

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Failed to gracefully shutdown server", slog.String("error", err.Error()))
		return err
	}

	slog.Info("HTTP server shut down successfully")
	return nil
}

// Engine returns the Gin engine (useful for testing)
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Container returns the dependency injection container
func (s *Server) Container() *container.Container {
	return s.container
}

// GetRoutes returns information about all registered routes
func (s *Server) GetRoutes() []gin.RouteInfo {
	return s.engine.Routes()
}

// IsRunning returns true if the server is currently running
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}