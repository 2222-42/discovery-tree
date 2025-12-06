package domain

import (
	"sort"
	"sync"
)

// InMemoryTaskRepository is an in-memory implementation of TaskRepository for testing
type InMemoryTaskRepository struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

// NewInMemoryTaskRepository creates a new in-memory task repository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*Task),
	}
}

// Save persists a task (create or update)
func (r *InMemoryTaskRepository) Save(task *Task) error {
	if task == nil {
		return NewValidationError("task", "task cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID().String()] = task
	return nil
}

// FindByID retrieves a task by its ID
func (r *InMemoryTaskRepository) FindByID(id TaskID) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id.String()]
	if !exists {
		return nil, NewNotFoundError("Task", id.String())
	}

	return task, nil
}

// FindByParentID retrieves all tasks with the given parent ID, ordered by position
func (r *InMemoryTaskRepository) FindByParentID(parentID *TaskID) ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Task

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

	// Sort by position
	sort.Slice(result, func(i, j int) bool {
		return result[i].Position() < result[j].Position()
	})

	return result, nil
}

// FindRoot retrieves the root task (task with no parent)
func (r *InMemoryTaskRepository) FindRoot() (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.ParentID() == nil {
			return task, nil
		}
	}

	return nil, NewNotFoundError("Root Task", "root")
}

// FindAll retrieves all tasks
func (r *InMemoryTaskRepository) FindAll() ([]*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		result = append(result, task)
	}

	return result, nil
}

// Delete removes a task by its ID
func (r *InMemoryTaskRepository) Delete(id TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id.String()]; !exists {
		return NewNotFoundError("Task", id.String())
	}

	delete(r.tasks, id.String())
	return nil
}

// DeleteSubtree removes a task and all its descendants
func (r *InMemoryTaskRepository) DeleteSubtree(id TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the task exists
	if _, exists := r.tasks[id.String()]; !exists {
		return NewNotFoundError("Task", id.String())
	}

	// Collect all tasks to delete (the task and all its descendants)
	toDelete := r.collectDescendants(id)
	toDelete = append(toDelete, id)

	// Delete all collected tasks
	for _, taskID := range toDelete {
		delete(r.tasks, taskID.String())
	}

	return nil
}

// collectDescendants recursively collects all descendant task IDs
// Note: This method assumes the lock is already held by the caller
func (r *InMemoryTaskRepository) collectDescendants(parentID TaskID) []TaskID {
	var descendants []TaskID

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
