package domain

// TaskValidator validates operations that span multiple tasks or require tree-wide knowledge
type TaskValidator interface {
	// ValidateStatusChange validates whether a status change is allowed
	// Returns an error if the status change violates business rules
	ValidateStatusChange(taskID TaskID, newStatus Status) error

	// ValidateMove validates whether a move operation is allowed
	// Returns an error if the move would create a cycle or violate constraints
	ValidateMove(taskID TaskID, newParentID *TaskID, newPosition int) error

	// ValidateDelete validates whether a delete operation is allowed
	// Returns an error if the delete violates constraints
	ValidateDelete(taskID TaskID) error
}

// taskValidator is the concrete implementation of TaskValidator
type taskValidator struct {
	repo TaskRepository
}

// NewTaskValidator creates a new TaskValidator instance
func NewTaskValidator(repo TaskRepository) TaskValidator {
	return &taskValidator{
		repo: repo,
	}
}

// ValidateStatusChange validates whether a status change is allowed
// Enforces bottom-to-top completion: a task can only be marked DONE if all children are DONE
// Non-DONE statuses are allowed regardless of children status
func (v *taskValidator) ValidateStatusChange(taskID TaskID, newStatus Status) error {
	// Only enforce constraints when changing to DONE status
	if newStatus != StatusDONE {
		// Non-DONE statuses (TODO, In Progress, Blocked, Root Work Item) are always allowed
		return nil
	}

	// For DONE status, check that all children are DONE
	task, err := v.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	// Get all children of this task
	children, err := v.repo.FindByParentID(&task.id)
	if err != nil {
		return err
	}

	// Check if any child is not DONE
	for _, child := range children {
		if child.Status() != StatusDONE {
			return NewConstraintViolationError(
				"bottom-to-top-completion",
				"cannot mark task as DONE when children are not all DONE",
			)
		}
	}

	// All children are DONE (or task has no children), so DONE status is allowed
	return nil
}

// ValidateMove validates whether a move operation is allowed
// Prevents cycles and validates that the new parent exists
func (v *taskValidator) ValidateMove(taskID TaskID, newParentID *TaskID, newPosition int) error {
	// TODO: Implement move validation (cycle detection)
	// This will be implemented in a future task
	return nil
}

// ValidateDelete validates whether a delete operation is allowed
func (v *taskValidator) ValidateDelete(taskID TaskID) error {
	// TODO: Implement delete validation
	// This will be implemented in a future task
	return nil
}
