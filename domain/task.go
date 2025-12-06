package domain

import (
	"strings"
	"time"
)

// Task represents a work item in the discovery tree (Aggregate Root)
type Task struct {
	id          TaskID
	description string
	status      Status
	parentID    *TaskID // nil for root tasks
	position    int     // position among siblings (0-indexed)
	createdAt   time.Time
	updatedAt   time.Time
}

// NewTask creates a new Task with validation
// For root tasks, parentID should be nil
// For child tasks, parentID should point to the parent task
func NewTask(description string, parentID *TaskID, position int) (*Task, error) {
	// Validate description is not empty or whitespace-only
	if strings.TrimSpace(description) == "" {
		return nil, NewValidationError("description", "description cannot be empty")
	}

	// Validate position is non-negative
	if position < 0 {
		return nil, NewValidationError("position", "position must be non-negative")
	}

	now := time.Now()
	
	// Determine initial status based on whether this is a root task
	initialStatus := StatusTODO
	if parentID == nil {
		initialStatus = StatusRootWorkItem
	}

	task := &Task{
		id:          NewTaskID(),
		description: description,
		status:      initialStatus,
		parentID:    parentID,
		position:    position,
		createdAt:   now,
		updatedAt:   now,
	}

	return task, nil
}

// ID returns the task's unique identifier
func (t *Task) ID() TaskID {
	return t.id
}

// Description returns the task's description
func (t *Task) Description() string {
	return t.description
}

// Status returns the task's current status
func (t *Task) Status() Status {
	return t.status
}

// ParentID returns the task's parent ID (nil for root tasks)
func (t *Task) ParentID() *TaskID {
	return t.parentID
}

// Position returns the task's position among siblings
func (t *Task) Position() int {
	return t.position
}

// CreatedAt returns the task's creation timestamp
func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the task's last update timestamp
func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}

// IsRoot returns true if this is a root task (no parent)
func (t *Task) IsRoot() bool {
	return t.parentID == nil
}

// ChangeStatus updates the task's status with validation
// This method performs basic status validation (checking if the status is valid)
// More complex validation (e.g., bottom-to-top enforcement) will be handled by TaskValidator
func (t *Task) ChangeStatus(newStatus Status) error {
	// Validate that the new status is a valid status value
	if !newStatus.IsValid() {
		return NewValidationError("status", "invalid status value")
	}

	// Update the status and timestamp
	t.status = newStatus
	t.updatedAt = time.Now()

	return nil
}

// UpdateDescription updates the task's description with validation
func (t *Task) UpdateDescription(description string) error {
	// Validate description is not empty or whitespace-only
	if strings.TrimSpace(description) == "" {
		return NewValidationError("description", "description cannot be empty")
	}

	// Update the description and timestamp
	t.description = description
	t.updatedAt = time.Now()

	return nil
}

// Move updates the task's parent and position
// This method only updates the task's internal state
// Position adjustments for siblings should be handled by the caller (e.g., TaskService)
// Validation (cycle detection) should be performed before calling this method
func (t *Task) Move(newParentID *TaskID, newPosition int) error {
	// Validate position is non-negative
	if newPosition < 0 {
		return NewValidationError("position", "position must be non-negative")
	}

	// Update parent and position
	t.parentID = newParentID
	t.position = newPosition
	t.updatedAt = time.Now()

	// Update status if moving to/from root
	if newParentID == nil {
		t.status = StatusRootWorkItem
	} else if t.status == StatusRootWorkItem {
		// If moving from root to non-root, change status to TODO
		t.status = StatusTODO
	}

	return nil
}
