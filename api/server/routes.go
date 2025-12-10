package server

import (
	"discovery-tree/api/container"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// RouteConfig holds configuration for route setup
type RouteConfig struct {
	EnableSwagger bool
	APIVersion    string
}

// SetupRoutes configures all API routes for the given engine and container
func SetupRoutes(engine *gin.Engine, container *container.Container) {
	config := &RouteConfig{
		EnableSwagger: container.Config().EnableSwagger,
		APIVersion:    "v1",
	}
	
	setupHealthRoutes(engine, container)
	setupAPIRoutes(engine, container, config)
	
	slog.Info("All routes configured successfully",
		slog.Int("total_routes", len(engine.Routes())),
		slog.String("api_version", config.APIVersion),
		slog.Bool("swagger_enabled", config.EnableSwagger),
	)
}

// setupHealthRoutes configures health check routes
func setupHealthRoutes(engine *gin.Engine, container *container.Container) {
	healthHandler := container.GetHealthHandler()
	
	// Health check endpoint (outside of API versioning for simplicity)
	engine.GET("/health", healthHandler.HealthCheck)
	
	slog.Debug("Health routes configured")
}

// setupAPIRoutes configures versioned API routes
func setupAPIRoutes(engine *gin.Engine, container *container.Container, config *RouteConfig) {
	// API version group
	apiGroup := engine.Group("/api/" + config.APIVersion)
	
	// Setup task routes
	setupTaskRoutes(apiGroup, container)
	
	// Future: Setup other resource routes here
	// setupUserRoutes(apiGroup, container)
	// setupProjectRoutes(apiGroup, container)
	
	slog.Debug("API routes configured", slog.String("version", config.APIVersion))
}

// setupTaskRoutes configures all task-related routes
func setupTaskRoutes(apiGroup *gin.RouterGroup, container *container.Container) {
	taskHandler := container.GetTaskHandler()
	
	// Task routes group
	tasks := apiGroup.Group("/tasks")
	
	// Root task operations (special endpoints)
	tasks.POST("/root", taskHandler.CreateRootTask)
	tasks.GET("/root", taskHandler.GetRootTask)
	
	// General task collection operations
	tasks.POST("", taskHandler.CreateChildTask)      // Create child task
	tasks.GET("", taskHandler.GetAllTasks)           // Get all tasks
	
	// Individual task operations (by ID)
	tasks.GET("/:id", taskHandler.GetTask)           // Get specific task
	tasks.PUT("/:id", taskHandler.UpdateTask)        // Update task description
	tasks.DELETE("/:id", taskHandler.DeleteTask)     // Delete task
	
	// Task status operations
	tasks.PUT("/:id/status", taskHandler.UpdateTaskStatus) // Update task status
	
	// Task hierarchy operations
	tasks.PUT("/:id/move", taskHandler.MoveTask)           // Move task
	tasks.GET("/:id/children", taskHandler.GetTaskChildren) // Get task children
	
	slog.Debug("Task routes configured",
		slog.Int("task_routes", 9), // Number of task-related routes
	)
}

// GetRoutesSummary returns a summary of all configured routes
func GetRoutesSummary(engine *gin.Engine) map[string]interface{} {
	routes := engine.Routes()
	
	summary := map[string]interface{}{
		"total_routes": len(routes),
		"routes_by_method": make(map[string]int),
		"routes_by_path_prefix": make(map[string]int),
	}
	
	// Count routes by method
	methodCounts := make(map[string]int)
	pathPrefixCounts := make(map[string]int)
	
	for _, route := range routes {
		methodCounts[route.Method]++
		
		// Count by path prefix
		if len(route.Path) > 0 {
			if route.Path == "/health" {
				pathPrefixCounts["health"]++
			} else if len(route.Path) >= 7 && route.Path[:7] == "/api/v1" {
				pathPrefixCounts["api_v1"]++
			} else {
				pathPrefixCounts["other"]++
			}
		}
	}
	
	summary["routes_by_method"] = methodCounts
	summary["routes_by_path_prefix"] = pathPrefixCounts
	
	return summary
}