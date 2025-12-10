// Package main provides the entry point for the Discovery Tree REST API
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
	"discovery-tree/api/container"
	"discovery-tree/api/server"
	"log"
	"os"
	"os/signal"
	"context"
	"syscall"
	"time"
)

func main() {
	// Initialize the dependency injection container with default configuration
	c, err := container.NewContainerWithDefaults()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Create and start the HTTP server
	srv := server.NewServer(c)
	
	log.Printf("Starting Discovery Tree API server on port %s", c.Config().Port)
	log.Printf("Using data path: %s", c.Config().DataPath)
	log.Printf("Swagger documentation enabled: %v", c.Config().EnableSwagger)
	
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Shutdown container
	if err := c.Shutdown(); err != nil {
		log.Printf("Error during container shutdown: %v", err)
	}

	log.Println("Server exited")
}