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
