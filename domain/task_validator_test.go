package domain

import (
	"testing"
)

func TestTaskValidator_ValidateStatusChange_NonDONEStatuses(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a parent task with incomplete children
	parent, _ := NewTask("Parent Task", nil, 0)
	repo.Save(parent)

	child1, _ := NewTask("Child 1", &parent.id, 0)
	child1.ChangeStatus(StatusTODO)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &parent.id, 1)
	child2.ChangeStatus(StatusInProgress)
	repo.Save(child2)

	// Test that non-DONE statuses are allowed regardless of children status
	nonDoneStatuses := []Status{StatusTODO, StatusInProgress, StatusBlocked, StatusRootWorkItem}

	for _, status := range nonDoneStatuses {
		t.Run("Allow_"+status.String(), func(t *testing.T) {
			err := validator.ValidateStatusChange(parent, status)
			if err != nil {
				t.Errorf("Expected non-DONE status %s to be allowed, but got error: %v", status.String(), err)
			}
		})
	}
}

func TestTaskValidator_ValidateStatusChange_DONEWithIncompleteChildren(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a parent task with incomplete children
	parent, _ := NewTask("Parent Task", nil, 0)
	repo.Save(parent)

	child1, _ := NewTask("Child 1", &parent.id, 0)
	child1.ChangeStatus(StatusDONE)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &parent.id, 1)
	child2.ChangeStatus(StatusTODO) // Not DONE
	repo.Save(child2)

	// Test that DONE status is rejected when children are not all DONE
	err := validator.ValidateStatusChange(parent, StatusDONE)
	if err == nil {
		t.Error("Expected error when marking parent as DONE with incomplete children, but got nil")
	}

	// Verify it's a ConstraintViolationError
	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError, but got: %T", err)
	}
}

func TestTaskValidator_ValidateStatusChange_DONEWithAllChildrenDone(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a parent task with all children DONE
	parent, _ := NewTask("Parent Task", nil, 0)
	repo.Save(parent)

	child1, _ := NewTask("Child 1", &parent.id, 0)
	child1.ChangeStatus(StatusDONE)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &parent.id, 1)
	child2.ChangeStatus(StatusDONE)
	repo.Save(child2)

	// Test that DONE status is allowed when all children are DONE
	err := validator.ValidateStatusChange(parent, StatusDONE)
	if err != nil {
		t.Errorf("Expected DONE status to be allowed when all children are DONE, but got error: %v", err)
	}
}

func TestTaskValidator_ValidateStatusChange_DONEWithNoChildren(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a leaf task (no children)
	leaf, _ := NewTask("Leaf Task", nil, 0)
	repo.Save(leaf)

	// Test that DONE status is allowed for tasks with no children
	err := validator.ValidateStatusChange(leaf, StatusDONE)
	if err != nil {
		t.Errorf("Expected DONE status to be allowed for leaf task, but got error: %v", err)
	}
}

// TestTaskValidator_ValidateStatusChange_TaskNotFound is no longer needed
// since the validator now receives a Task object instead of TaskID.
// The task lookup is now done in the service layer before calling the validator.

func TestTaskValidator_ValidateStatusChange_MultiLevelTree(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a multi-level tree
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	parent1, _ := NewTask("Parent 1", &root.id, 0)
	repo.Save(parent1)

	child1, _ := NewTask("Child 1", &parent1.id, 0)
	child1.ChangeStatus(StatusDONE)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &parent1.id, 1)
	child2.ChangeStatus(StatusDONE)
	repo.Save(child2)

	parent2, _ := NewTask("Parent 2", &root.id, 1)
	parent2.ChangeStatus(StatusTODO) // Not DONE
	repo.Save(parent2)

	// Test that parent1 can be marked DONE (all its children are DONE)
	err := validator.ValidateStatusChange(parent1, StatusDONE)
	if err != nil {
		t.Errorf("Expected parent1 to be markable as DONE, but got error: %v", err)
	}

	// Test that root cannot be marked DONE (parent2 is not DONE)
	err = validator.ValidateStatusChange(root, StatusDONE)
	if err == nil {
		t.Error("Expected error when marking root as DONE with incomplete children, but got nil")
	}
}

func TestTaskValidator_ValidateMove_ValidMove(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a simple tree
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	child1, _ := NewTask("Child 1", &root.id, 0)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &root.id, 1)
	repo.Save(child2)

	// Move child1 to position 1 (swap with child2)
	err := validator.ValidateMove(child1.ID(), &root.id, 1)
	if err != nil {
		t.Errorf("Expected valid move to be allowed, got error: %v", err)
	}
}

func TestTaskValidator_ValidateMove_ToNewParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a tree with two branches
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	parent1, _ := NewTask("Parent 1", &root.id, 0)
	repo.Save(parent1)

	parent2, _ := NewTask("Parent 2", &root.id, 1)
	repo.Save(parent2)

	child, _ := NewTask("Child", &parent1.id, 0)
	repo.Save(child)

	// Move child from parent1 to parent2
	err := validator.ValidateMove(child.ID(), &parent2.id, 0)
	if err != nil {
		t.Errorf("Expected move to different parent to be allowed, got error: %v", err)
	}
}

func TestTaskValidator_ValidateMove_PreventCycleToSelf(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a task
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	task, _ := NewTask("Task", &root.id, 0)
	repo.Save(task)

	// Try to move task to itself as parent
	err := validator.ValidateMove(task.ID(), &task.id, 0)
	if err == nil {
		t.Error("Expected error when moving task to itself, got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_PreventCycleToDescendant(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a multi-level tree
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	parent, _ := NewTask("Parent", &root.id, 0)
	repo.Save(parent)

	child, _ := NewTask("Child", &parent.id, 0)
	repo.Save(child)

	grandchild, _ := NewTask("Grandchild", &child.id, 0)
	repo.Save(grandchild)

	// Try to move parent to be a child of its own grandchild (creates cycle)
	err := validator.ValidateMove(parent.ID(), &grandchild.id, 0)
	if err == nil {
		t.Error("Expected error when moving task to its descendant, got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_PreventCycleToDirectChild(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a simple parent-child relationship
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	parent, _ := NewTask("Parent", &root.id, 0)
	repo.Save(parent)

	child, _ := NewTask("Child", &parent.id, 0)
	repo.Save(child)

	// Try to move parent to be a child of its own child
	err := validator.ValidateMove(parent.ID(), &child.id, 0)
	if err == nil {
		t.Error("Expected error when moving task to its direct child, got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_NonExistentTask(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a parent
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	// Try to move a non-existent task
	nonExistentID := NewTaskID()
	err := validator.ValidateMove(nonExistentID, &root.id, 0)
	if err == nil {
		t.Error("Expected error when moving non-existent task, got nil")
	}

	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("Expected NotFoundError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_NonExistentParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a task
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	task, _ := NewTask("Task", &root.id, 0)
	repo.Save(task)

	// Try to move to a non-existent parent
	nonExistentParentID := NewTaskID()
	err := validator.ValidateMove(task.ID(), &nonExistentParentID, 0)
	if err == nil {
		t.Error("Expected error when moving to non-existent parent, got nil")
	}

	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("Expected NotFoundError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_NegativePosition(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a simple tree
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	task, _ := NewTask("Task", &root.id, 0)
	repo.Save(task)

	// Try to move to negative position
	err := validator.ValidateMove(task.ID(), &root.id, -1)
	if err == nil {
		t.Error("Expected error for negative position, got nil")
	}

	if _, ok := err.(ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_ToRoot(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a root and a child
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	child, _ := NewTask("Child", &root.id, 0)
	repo.Save(child)

	// Try to move child to root (should fail because root already exists)
	err := validator.ValidateMove(child.ID(), nil, 0)
	if err == nil {
		t.Error("Expected error when moving to root when root exists, got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_RootToChild(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a root and another task
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	// Create a separate tree (this would be invalid in practice, but for testing)
	// Actually, we can't have two roots, so let's create a child first
	child, _ := NewTask("Child", &root.id, 0)
	repo.Save(child)

	// Move root to be a child (should be allowed if we're moving the current root)
	err := validator.ValidateMove(root.ID(), &child.id, 0)
	if err == nil {
		t.Error("Expected error when moving root to its own child (cycle), got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("Expected ConstraintViolationError for cycle, got %T", err)
	}
}

func TestTaskValidator_ValidateMove_PositionExceedsRange(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a parent with 2 children
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	parent, _ := NewTask("Parent", &root.id, 0)
	repo.Save(parent)

	child1, _ := NewTask("Child 1", &parent.id, 0)
	repo.Save(child1)

	child2, _ := NewTask("Child 2", &parent.id, 1)
	repo.Save(child2)

	// Create a task in a different parent
	otherParent, _ := NewTask("Other Parent", &root.id, 1)
	repo.Save(otherParent)

	task, _ := NewTask("Task", &otherParent.id, 0)
	repo.Save(task)

	// Try to move task to parent at position 3 (parent has 2 children, so max is 2)
	err := validator.ValidateMove(task.ID(), &parent.id, 3)
	if err == nil {
		t.Error("Expected error for position exceeding range, got nil")
	}

	if _, ok := err.(ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}
