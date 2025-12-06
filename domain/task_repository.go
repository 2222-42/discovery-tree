package domain

// TaskRepository provides persistence operations for Task aggregates
type TaskRepository interface {
	// Save persists a task (create or update)
	Save(task *Task) error

	// FindByID retrieves a task by its ID
	FindByID(id TaskID) (*Task, error)

	// FindByParentID retrieves all tasks with the given parent ID, ordered by position
	FindByParentID(parentID *TaskID) ([]*Task, error)

	// FindRoot retrieves the root task (task with no parent)
	FindRoot() (*Task, error)

	// FindAll retrieves all tasks
	FindAll() ([]*Task, error)

	// Delete removes a task by its ID (should only be used for leaf tasks)
	Delete(id TaskID) error

	// DeleteSubtree removes a task and all its descendants
	DeleteSubtree(id TaskID) error
}
