package application

import (
	"discovery-tree/domain"
	"discovery-tree/infrastructure"
	"fmt"
	"log/slog"
	"os"
)

// Example 1: Basic initialization with default file path
// This example shows the simplest way to initialize the system
func ExampleBasicInitialization() {
	// Create a FileTaskRepository with default path (./data/tasks.json)
	repo, err := infrastructure.NewFileTaskRepository("")
	if err != nil {
		slog.Error("Failed to create repository", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Inject the repository into TaskService
	service := domain.NewTaskService(repo)

	// Now you can use the service for task operations
	fmt.Println("Task service initialized with default file path")
	_ = service // Use the service for your application logic
}

// Example 2: Custom file path configuration
// This example demonstrates how to specify a custom storage location
func ExampleCustomFilePath() {
	// Create a FileTaskRepository with a custom file path
	customPath := "./my-app-data/tasks.json"
	repo, err := infrastructure.NewFileTaskRepository(customPath)
	if err != nil {
		slog.Error("Failed to create repository with custom path", 
			slog.String("path", customPath),
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Inject the repository into TaskService
	service := domain.NewTaskService(repo)

	fmt.Printf("Task service initialized with custom path: %s\n", customPath)
	_ = service
}

// Example 3: Complete application initialization with error handling
// This example shows a production-ready initialization pattern
func InitializeApplication(configPath string) (*domain.TaskService, error) {
	// Determine the file path to use
	filePath := configPath
	if filePath == "" {
		// Use default path if none provided
		filePath = "./data/tasks.json"
	}

	// Create the repository
	repo, err := infrastructure.NewFileTaskRepository(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Create and return the service with injected repository
	service := domain.NewTaskService(repo)

	return service, nil
}

// Example 4: Using the initialized service
// This example demonstrates typical usage patterns after initialization
func ExampleUsage() {
	// Initialize the application
	service, err := InitializeApplication("./data/tasks.json")
	if err != nil {
		slog.Error("Failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create a root task
	rootTask, err := service.CreateRootTask("My Project")
	if err != nil {
		slog.Error("Failed to create root task", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Created root task", slog.String("description", rootTask.Description()))

	// Create child tasks
	child1, err := service.CreateChildTask("Phase 1", rootTask.ID())
	if err != nil {
		slog.Error("Failed to create child task", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Created child task", slog.String("description", child1.Description()))

	child2, err := service.CreateChildTask("Phase 2", rootTask.ID())
	if err != nil {
		slog.Error("Failed to create child task", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Created child task", slog.String("description", child2.Description()))

	// Change task status
	err = service.ChangeTaskStatus(child1.ID(), domain.StatusInProgress)
	if err != nil {
		slog.Error("Failed to change task status", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Changed task status", slog.String("status", "In Progress"))

	// All changes are automatically persisted to the JSON file
	fmt.Println("All changes persisted to disk")
}

// Example 5: Multiple repository instances (for testing or multi-tenant scenarios)
// This example shows how to work with multiple independent repositories
func ExampleMultipleRepositories() {
	// Create repository for user A
	repoA, err := infrastructure.NewFileTaskRepository("./data/user_a_tasks.json")
	if err != nil {
		slog.Error("Failed to create repository A", 
			slog.String("path", "./data/user_a_tasks.json"),
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	serviceA := domain.NewTaskService(repoA)

	// Create repository for user B
	repoB, err := infrastructure.NewFileTaskRepository("./data/user_b_tasks.json")
	if err != nil {
		slog.Error("Failed to create repository B", 
			slog.String("path", "./data/user_b_tasks.json"),
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	serviceB := domain.NewTaskService(repoB)

	// Each service operates independently with its own storage
	fmt.Println("Multiple independent task services initialized")
	_, _ = serviceA, serviceB
}

// Example 6: Dependency injection pattern for testing
// This example shows how the repository interface enables easy testing
func ExampleDependencyInjection() {
	// In production: use FileTaskRepository
	productionRepo, err := infrastructure.NewFileTaskRepository("./data/tasks.json")
	if err != nil {
		slog.Error("Failed to create production repository", 
			slog.String("path", "./data/tasks.json"),
			slog.String("error", err.Error()))
		os.Exit(1)
	}
	productionService := domain.NewTaskService(productionRepo)

	// In tests: use InMemoryTaskRepository (no file I/O)
	testRepo := domain.NewInMemoryTaskRepository()
	testService := domain.NewTaskService(testRepo)

	// Both services have the same interface and behavior
	// The only difference is where data is stored
	fmt.Println("Production and test services use the same interface")
	_, _ = productionService, testService
}
