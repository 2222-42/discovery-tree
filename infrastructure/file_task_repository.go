package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"discovery-tree/domain"
)

// FileTaskRepository implements TaskRepository with JSON file persistence
type FileTaskRepository struct {
	filePath string
	tasks    map[string]*domain.Task // in-memory cache, keyed by task ID string
	mu       sync.RWMutex            // protects concurrent access
}

// NewFileTaskRepository creates a new FileTaskRepository
// If filePath is empty, uses default path "./data/tasks.json"
// Creates necessary directories if they don't exist
// Loads existing data from file if it exists
func NewFileTaskRepository(filePath string) (*FileTaskRepository, error) {
	// Use default path if empty
	if filePath == "" {
		filePath = "./data/tasks.json"
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, WrapFileSystemError("create directory", dir, err)
	}

	// Initialize repository
	repo := &FileTaskRepository{
		filePath: filePath,
		tasks:    make(map[string]*domain.Task),
	}

	// Load existing data from file
	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

// load reads tasks from the JSON file and populates the in-memory cache
// If the file doesn't exist, initializes with an empty collection
// Returns error if file contains invalid JSON or invalid task data
func (r *FileTaskRepository) load() error {
	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		// File doesn't exist, initialize with empty collection
		return nil
	}

	// Read file contents
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return WrapFileSystemError("read", r.filePath, err)
	}

	// Handle empty file
	if len(data) == 0 {
		return nil
	}

	// Parse JSON
	var dtos []TaskDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return WrapFileSystemError("parse JSON", r.filePath, err)
	}

	// Convert DTOs to tasks and populate cache
	for _, dto := range dtos {
		task, err := FromDTO(dto)
		if err != nil {
			// Invalid task data
			return err
		}
		r.tasks[task.ID().String()] = task
	}

	return nil
}

// Save persists a task (create or update)
func (r *FileTaskRepository) Save(task *domain.Task) error {
	// TODO: Implement in task 5
	return nil
}

// FindByID retrieves a task by its ID
func (r *FileTaskRepository) FindByID(id domain.TaskID) (*domain.Task, error) {
	// TODO: Implement in task 6
	return nil, nil
}

// FindByParentID retrieves all tasks with the given parent ID, ordered by position
func (r *FileTaskRepository) FindByParentID(parentID *domain.TaskID) ([]*domain.Task, error) {
	// TODO: Implement in task 6
	return nil, nil
}

// FindRoot retrieves the root task (task with no parent)
func (r *FileTaskRepository) FindRoot() (*domain.Task, error) {
	// TODO: Implement in task 6
	return nil, nil
}

// FindAll retrieves all tasks
func (r *FileTaskRepository) FindAll() ([]*domain.Task, error) {
	// TODO: Implement in task 6
	return nil, nil
}

// Delete removes a task by its ID
func (r *FileTaskRepository) Delete(id domain.TaskID) error {
	// TODO: Implement in task 7
	return nil
}

// DeleteSubtree removes a task and all its descendants
func (r *FileTaskRepository) DeleteSubtree(id domain.TaskID) error {
	// TODO: Implement in task 8
	return nil
}
