package domain

// TreeNavigator provides navigation operations across the tree structure
type TreeNavigator interface {
	// GetParent returns the parent task of the given task, or nil if the task is root
	GetParent(taskID TaskID) (*Task, error)

	// GetChildren returns all child tasks of the given task, ordered by position (left-to-right)
	GetChildren(taskID TaskID) ([]*Task, error)

	// GetSiblings returns all tasks sharing the same parent as the given task (including the task itself), ordered by position
	GetSiblings(taskID TaskID) ([]*Task, error)

	// GetLeftSibling returns the immediate left sibling of the given task, or nil if none exists
	GetLeftSibling(taskID TaskID) (*Task, error)

	// GetRightSibling returns the immediate right sibling of the given task, or nil if none exists
	GetRightSibling(taskID TaskID) (*Task, error)

	// GetRoot returns the root task of the tree
	GetRoot() (*Task, error)

	// GetTree returns the complete tree structure (root task and all descendants)
	GetTree() ([]*Task, error)

	// GetSubtree returns the given task and all its descendants
	GetSubtree(taskID TaskID) ([]*Task, error)
}

// TreeNavigatorService implements TreeNavigator using a TaskRepository
type TreeNavigatorService struct {
	repo TaskRepository
}

// NewTreeNavigatorService creates a new TreeNavigatorService
func NewTreeNavigatorService(repo TaskRepository) *TreeNavigatorService {
	return &TreeNavigatorService{
		repo: repo,
	}
}

// GetParent returns the parent task of the given task, or nil if the task is root
func (s *TreeNavigatorService) GetParent(taskID TaskID) (*Task, error) {
	// First, find the task to get its parent ID
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// If the task has no parent (is root), return nil
	if task.ParentID() == nil {
		return nil, nil
	}

	// Find and return the parent task
	parent, err := s.repo.FindByID(*task.ParentID())
	if err != nil {
		return nil, err
	}

	return parent, nil
}

// GetChildren returns all child tasks of the given task, ordered by position (left-to-right)
func (s *TreeNavigatorService) GetChildren(taskID TaskID) ([]*Task, error) {
	// First, verify the task exists
	_, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// Find all children by parent ID (repository returns them ordered by position)
	children, err := s.repo.FindByParentID(&taskID)
	if err != nil {
		return nil, err
	}

	return children, nil
}

// GetSiblings returns all tasks sharing the same parent as the given task, ordered by position
func (s *TreeNavigatorService) GetSiblings(taskID TaskID) ([]*Task, error) {
	// First, find the task to get its parent ID
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// Find all tasks with the same parent (repository returns them ordered by position)
	siblings, err := s.repo.FindByParentID(task.ParentID())
	if err != nil {
		return nil, err
	}

	return siblings, nil
}

// GetLeftSibling returns the immediate left sibling of the given task, or nil if none exists
func (s *TreeNavigatorService) GetLeftSibling(taskID TaskID) (*Task, error) {
	// First, find the task to get its position and parent
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// If this is the leftmost sibling (position 0), there is no left sibling
	if task.Position() == 0 {
		return nil, nil
	}

	// Get all siblings
	siblings, err := s.repo.FindByParentID(task.ParentID())
	if err != nil {
		return nil, err
	}

	// Find the sibling with position-1
	targetPosition := task.Position() - 1
	for _, sibling := range siblings {
		if sibling.Position() == targetPosition {
			return sibling, nil
		}
	}

	// No left sibling found (shouldn't happen if positions are consistent)
	return nil, nil
}

// GetRightSibling returns the immediate right sibling of the given task, or nil if none exists
func (s *TreeNavigatorService) GetRightSibling(taskID TaskID) (*Task, error) {
	// First, find the task to get its position and parent
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// Get all siblings
	siblings, err := s.repo.FindByParentID(task.ParentID())
	if err != nil {
		return nil, err
	}

	// Find the sibling with position+1
	targetPosition := task.Position() + 1
	for _, sibling := range siblings {
		if sibling.Position() == targetPosition {
			return sibling, nil
		}
	}

	// No right sibling found (this is the rightmost sibling)
	return nil, nil
}

// GetRoot returns the root task of the tree
func (s *TreeNavigatorService) GetRoot() (*Task, error) {
	return s.repo.FindRoot()
}

// GetTree returns the complete tree structure (root task and all descendants)
func (s *TreeNavigatorService) GetTree() ([]*Task, error) {
	// Get the root task
	root, err := s.repo.FindRoot()
	if err != nil {
		return nil, err
	}

	// Get the subtree starting from the root (which is the entire tree)
	return s.GetSubtree(root.ID())
}

// GetSubtree returns the given task and all its descendants
func (s *TreeNavigatorService) GetSubtree(taskID TaskID) ([]*Task, error) {
	// First, verify the task exists and get it
	task, err := s.repo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	// Start with the task itself
	result := []*Task{task}

	// Recursively collect all descendants
	descendants, err := s.collectDescendants(taskID)
	if err != nil {
		return nil, err
	}

	result = append(result, descendants...)

	return result, nil
}

// collectDescendants recursively collects all descendant tasks
func (s *TreeNavigatorService) collectDescendants(parentID TaskID) ([]*Task, error) {
	var descendants []*Task

	// Get all direct children
	children, err := s.repo.FindByParentID(&parentID)
	if err != nil {
		return nil, err
	}

	// For each child, add it and its descendants
	for _, child := range children {
		descendants = append(descendants, child)

		// Recursively get descendants of this child
		childDescendants, err := s.collectDescendants(child.ID())
		if err != nil {
			return nil, err
		}

		descendants = append(descendants, childDescendants...)
	}

	return descendants, nil
}
