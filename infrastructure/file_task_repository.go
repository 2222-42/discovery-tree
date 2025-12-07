package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
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

// persist writes the in-memory task collection to the JSON file atomically
// Uses atomic write pattern: write to temp file, then rename
// Formats JSON with 2-space indentation for readability
func (r *FileTaskRepository) persist() error {
	// Convert tasks to DTOs
	dtos := make([]TaskDTO, 0, len(r.tasks))
	for _, task := range r.tasks {
		dtos = append(dtos, ToDTO(task))
	}

	// Marshal to JSON with indentation (2 spaces)
	data, err := json.MarshalIndent(dtos, "", "  ")
	if err != nil {
		return WrapFileSystemError("marshal JSON", r.filePath, err)
	}

	// Write to temporary file
	tmpPath := r.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return WrapFileSystemError("write temporary file", tmpPath, err)
	}

	// Atomic rename (replaces target file atomically on POSIX systems)
	if err := os.Rename(tmpPath, r.filePath); err != nil {
		// Clean up temporary file on failure
		os.Remove(tmpPath)
		return WrapFileSystemError("atomic rename", r.filePath, err)
	}

	return nil
}

// Save persists a task (create or update)
func (r *FileTaskRepository) Save(task *domain.Task) error {
	// Use write lock for thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// Add task to in-memory map (or update if exists)
	r.tasks[task.ID().String()] = task

	// Call persist() to write to file
	return r.persist()
}

// FindByID retrieves a task by its ID
func (r *FileTaskRepository) FindByID(id domain.TaskID) (*domain.Task, error) {
	// Use read lock for thread safety
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id.String()]
	if !exists {
		return nil, domain.NewNotFoundError("Task", id.String())
	}

	return task, nil
}

// FindByParentID retrieves all tasks with the given parent ID, ordered by position
func (r *FileTaskRepository) FindByParentID(parentID *domain.TaskID) ([]*domain.Task, error) {
	// Use read lock for thread safety
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Task

	for _, task := range r.tasks {
		// Check if this task has the matching parent
		if parentID == nil && task.ParentID() == nil {
			// Both are nil (root tasks)
			result = append(result, task)
		} else if parentID != nil && task.ParentID() != nil && parentID.Equals(*task.ParentID()) {
			// Both are non-nil and equal
			result = append(result, task)
		}
	}

	// Sort by position (ascending)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Position() < result[j].Position()
	})

	return result, nil
}

// FindRoot retrieves the root task (task with no parent)
func (r *FileTaskRepository) FindRoot() (*domain.Task, error) {
	// Use read lock for thread safety
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.ParentID() == nil {
			return task, nil
		}
	}

	return nil, domain.NewNotFoundError("Root Task", "root")
}

// FindAll retrieves all tasks
func (r *FileTaskRepository) FindAll() ([]*domain.Task, error) {
	// Use read lock for thread safety
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}

	return result, nil
}

// Delete removes a task by its ID
func (r *FileTaskRepository) Delete(id domain.TaskID) error {
	// Use write lock for thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if task exists
	idStr := id.String()
	if _, exists := r.tasks[idStr]; !exists {
		return domain.NewNotFoundError("Task", idStr)
	}

	// Remove task from in-memory map
	delete(r.tasks, idStr)

	// Call persist() to write changes
	return r.persist()
}

// DeleteSubtree removes a task and all its descendants
func (r *FileTaskRepository) DeleteSubtree(id domain.TaskID) error {
	// Use write lock for thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the task exists
	idStr := id.String()
	if _, exists := r.tasks[idStr]; !exists {
		return domain.NewNotFoundError("Task", idStr)
	}

	// Collect all tasks to delete (the task and all its descendants)
	toDelete := r.collectDescendants(id)
	toDelete = append(toDelete, id)

	// Delete all collected tasks from in-memory map
	for _, taskID := range toDelete {
		delete(r.tasks, taskID.String())
	}

	// Call persist() to write changes
	return r.persist()
}

// collectDescendants recursively collects all descendant task IDs
// Note: This method assumes the lock is already held by the caller
func (r *FileTaskRepository) collectDescendants(parentID domain.TaskID) []domain.TaskID {
	var descendants []domain.TaskID

	// Find all direct children
	for _, task := range r.tasks {
		if task.ParentID() != nil && task.ParentID().Equals(parentID) {
			descendants = append(descendants, task.ID())
			// Recursively collect descendants of this child
			childDescendants := r.collectDescendants(task.ID())
			descendants = append(descendants, childDescendants...)
		}
	}

	return descendants
}
