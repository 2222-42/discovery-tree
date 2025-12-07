package domain

// TaskValidator validates operations that span multiple tasks or require tree-wide knowledge
type TaskValidator interface {
	// ValidateStatusChange validates whether a status change is allowed
	// Returns an error if the status change violates business rules
	ValidateStatusChange(task *Task, newStatus Status) error

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
func (v *taskValidator) ValidateStatusChange(task *Task, newStatus Status) error {
	// Only enforce constraints when changing to DONE status
	if newStatus != StatusDONE {
		// Non-DONE statuses (TODO, In Progress, Blocked, Root Work Item) are always allowed
		return nil
	}

	// For DONE status, check that all children are DONE
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
	// Validate position is non-negative
	if newPosition < 0 {
		return NewValidationError("position", "position must be non-negative")
	}

	// Verify the task being moved exists
	task, err := v.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	// If moving to root (newParentID is nil), no cycle check needed
	if newParentID == nil {
		// Check if there's already a root task (unless we're moving the current root)
		if !task.IsRoot() {
			existingRoot, err := v.repo.FindRoot()
			if err == nil && existingRoot != nil {
				return NewConstraintViolationError(
					"single-root",
					"cannot move task to root: a root task already exists",
				)
			}
		}
		return nil
	}

	// Verify the new parent exists
	newParent, err := v.repo.FindByID(*newParentID)
	if err != nil {
		return err
	}

	// Prevent moving a task to itself
	if taskID.Equals(*newParentID) {
		return NewConstraintViolationError(
			"cycle-prevention",
			"cannot move task to itself",
		)
	}

	// Prevent creating a cycle: check if newParent is a descendant of task
	if v.isDescendant(taskID, *newParentID) {
		return NewConstraintViolationError(
			"cycle-prevention",
			"cannot move task to its own descendant",
		)
	}

	// Validate position is within valid range for the new parent
	siblings, err := v.repo.FindByParentID(newParentID)
	if err != nil {
		return err
	}

	// If moving within the same parent, max position is len(siblings) - 1
	// If moving to a different parent, max position is len(siblings)
	maxPosition := len(siblings)
	if task.ParentID() != nil && task.ParentID().Equals(newParent.ID()) {
		maxPosition = len(siblings) - 1
	}

	if newPosition > maxPosition {
		return NewValidationError("position", "position exceeds valid range")
	}

	return nil
}

// isDescendant checks if potentialDescendant is a descendant of ancestor
func (v *taskValidator) isDescendant(ancestor TaskID, potentialDescendant TaskID) bool {
	// Start from potentialDescendant and walk up the tree
	current := potentialDescendant
	
	for {
		task, err := v.repo.FindByID(current)
		if err != nil {
			// If we can't find the task, it's not a descendant
			return false
		}

		// If we reached the root, potentialDescendant is not a descendant of ancestor
		parentId := task.ParentID()
		if parentId == nil {
			return false
		}

		// If the parent is the ancestor, potentialDescendant is a descendant
		if parentId.Equals(ancestor) {
			return true
		}

		// Move up to the parent
		current = *parentId
	}
}

// ValidateDelete validates whether a delete operation is allowed
func (v *taskValidator) ValidateDelete(taskID TaskID) error {
	// TODO: Implement delete validation
	// This will be implemented in a future task
	return nil
}
