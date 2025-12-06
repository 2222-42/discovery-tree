package domain

// TaskService provides domain logic for task operations that require repository access
type TaskService struct {
	repo      TaskRepository
	validator TaskValidator
}

// NewTaskService creates a new TaskService
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{
		repo:      repo,
		validator: NewTaskValidator(repo),
	}
}

// CreateRootTask creates a new root task with validation
// Ensures only one root task exists in the tree
func (s *TaskService) CreateRootTask(description string) (*Task, error) {
	// Check if a root task already exists
	existingRoot, err := s.repo.FindRoot()
	if err == nil && existingRoot != nil {
		return nil, NewConstraintViolationError("single_root", "root task already exists")
	}
	// If error is NotFoundError, that's expected and we can proceed
	if err != nil {
		if _, ok := err.(NotFoundError); !ok {
			// Some other error occurred
			return nil, err
		}
	}

	// Create the root task with position 0
	task, err := NewTask(description, nil, 0)
	if err != nil {
		return nil, err
	}

	// Save the task
	err = s.repo.Save(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// CreateChildTask creates a new child task under the specified parent
// Automatically calculates the position based on existing children
// Validates that the parent exists
func (s *TaskService) CreateChildTask(description string, parentID TaskID) (*Task, error) {
	// Validate that the parent exists
	_, err := s.repo.FindByID(parentID)
	if err != nil {
		return nil, err
	}

	// Find existing children to calculate the next position
	children, err := s.repo.FindByParentID(&parentID)
	if err != nil {
		return nil, err
	}

	// Calculate the next position (number of existing children)
	nextPosition := len(children)

	// Create the child task
	task, err := NewTask(description, &parentID, nextPosition)
	if err != nil {
		return nil, err
	}

	// Save the task
	err = s.repo.Save(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// ChangeTaskStatus changes the status of a task with validation
// Enforces bottom-to-top completion: a task can only be marked DONE if all children are DONE
// Non-DONE statuses are allowed regardless of children status
func (s *TaskService) ChangeTaskStatus(taskID TaskID, newStatus Status) error {
	// Retrieve the task first
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	// Validate the status change using the validator
	// This enforces bottom-to-top completion for DONE status
	err = s.validator.ValidateStatusChange(task, newStatus)
	if err != nil {
		return err
	}

	// Change the status on the task entity
	// This performs basic validation (checking if status is valid)
	err = task.ChangeStatus(newStatus)
	if err != nil {
		return err
	}

	// Save the updated task
	err = s.repo.Save(task)
	if err != nil {
		return err
	}

	return nil
}

// MoveTask moves a task to a new parent and position
// Handles position adjustments for both old and new siblings
// Validates the move operation (prevents cycles)
// The entire subtree moves with the task
func (s *TaskService) MoveTask(taskID TaskID, newParentID *TaskID, newPosition int) error {
	// Retrieve the task being moved
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	// Validate the move operation (cycle detection, parent exists, etc.)
	err = s.validator.ValidateMove(taskID, newParentID, newPosition)
	if err != nil {
		return err
	}

	oldParentID := task.ParentID()
	oldPosition := task.Position()

	// Check if this is actually a move (parent or position changed)
	isSameParent := (oldParentID == nil && newParentID == nil) ||
		(oldParentID != nil && newParentID != nil && oldParentID.Equals(*newParentID))

	if isSameParent && oldPosition == newPosition {
		// No actual move needed
		return nil
	}

	if !isSameParent {
		// Moving to a different parent
		
		// Step 1: Adjust positions of old siblings (close the gap)
		oldSiblings, err := s.repo.FindByParentID(oldParentID)
		if err != nil {
			return err
		}

		for _, sibling := range oldSiblings {
			// Skip the task being moved
			if sibling.ID().Equals(taskID) {
				continue
			}

			// Shift left siblings that are to the right of the old position
			if sibling.Position() > oldPosition {
				err = sibling.Move(sibling.ParentID(), sibling.Position()-1)
				if err != nil {
					return err
				}
				err = s.repo.Save(sibling)
				if err != nil {
					return err
				}
			}
		}

		// Step 2: Adjust positions of new siblings (make room)
		newSiblings, err := s.repo.FindByParentID(newParentID)
		if err != nil {
			return err
		}

		for _, sibling := range newSiblings {
			// Shift right siblings that are at or after the new position
			if sibling.Position() >= newPosition {
				err = sibling.Move(sibling.ParentID(), sibling.Position()+1)
				if err != nil {
					return err
				}
				err = s.repo.Save(sibling)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Moving within the same parent (reordering)
		siblings, err := s.repo.FindByParentID(newParentID)
		if err != nil {
			return err
		}

		if newPosition > oldPosition {
			// Moving right: shift siblings between old and new position left
			for _, sibling := range siblings {
				if sibling.ID().Equals(taskID) {
					continue
				}
				if sibling.Position() > oldPosition && sibling.Position() <= newPosition {
					err = sibling.Move(sibling.ParentID(), sibling.Position()-1)
					if err != nil {
						return err
					}
					err = s.repo.Save(sibling)
					if err != nil {
						return err
					}
				}
			}
		} else {
			// Moving left: shift siblings between new and old position right
			for _, sibling := range siblings {
				if sibling.ID().Equals(taskID) {
					continue
				}
				if sibling.Position() >= newPosition && sibling.Position() < oldPosition {
					err = sibling.Move(sibling.ParentID(), sibling.Position()+1)
					if err != nil {
						return err
					}
					err = s.repo.Save(sibling)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// Step 3: Move the task itself
	err = task.Move(newParentID, newPosition)
	if err != nil {
		return err
	}

	err = s.repo.Save(task)
	if err != nil {
		return err
	}

	// Note: The subtree automatically moves with the task because
	// child tasks reference their parent by ID, and we're not changing
	// the task's ID, only its parent and position

	return nil
}

// DeleteTask deletes a task and adjusts sibling positions
// If the task has children, it performs cascading deletion
// If the task is the root, it removes the entire tree
func (s *TaskService) DeleteTask(taskID TaskID) error {
	// Retrieve the task to be deleted
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	// Check if this is the root task
	if task.IsRoot() {
		// Root deletion removes the entire tree
		return s.deleteEntireTree()
	}

	// Check if the task has children
	children, err := s.repo.FindByParentID(&taskID)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		// Task has children, perform cascading deletion
		return s.deleteSubtree(taskID)
	}

	// Task is a leaf, delete it and adjust sibling positions
	return s.deleteLeafTask(task)
}

// deleteLeafTask deletes a leaf task (no children) and adjusts sibling positions
func (s *TaskService) deleteLeafTask(task *Task) error {
	parentID := task.ParentID()
	position := task.Position()
	taskID := task.ID()

	// Delete the task
	err := s.repo.Delete(taskID)
	if err != nil {
		return err
	}

	// Adjust positions of right siblings (shift them left)
	siblings, err := s.repo.FindByParentID(parentID)
	if err != nil {
		return err
	}

	for _, sibling := range siblings {
		// Shift left siblings that were to the right of the deleted task
		if sibling.Position() > position {
			err = sibling.Move(sibling.ParentID(), sibling.Position()-1)
			if err != nil {
				return err
			}
			err = s.repo.Save(sibling)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// deleteSubtree deletes a task and all its descendants, then adjusts sibling positions
func (s *TaskService) deleteSubtree(taskID TaskID) error {
	// Retrieve the task to get its parent and position
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return err
	}

	parentID := task.ParentID()
	position := task.Position()

	// Delete the task and all its descendants using repository's DeleteSubtree
	err = s.repo.DeleteSubtree(taskID)
	if err != nil {
		return err
	}

	// Adjust positions of right siblings (shift them left)
	siblings, err := s.repo.FindByParentID(parentID)
	if err != nil {
		return err
	}

	for _, sibling := range siblings {
		// Shift left siblings that were to the right of the deleted task
		if sibling.Position() > position {
			err = sibling.Move(sibling.ParentID(), sibling.Position()-1)
			if err != nil {
				return err
			}
			err = s.repo.Save(sibling)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// deleteEntireTree deletes all tasks in the tree
func (s *TaskService) deleteEntireTree() error {
	// Find the root task
	root, err := s.repo.FindRoot()
	if err != nil {
		return err
	}

	// Delete the root and all its descendants
	return s.repo.DeleteSubtree(root.ID())
}
