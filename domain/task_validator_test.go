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
			err := validator.ValidateStatusChange(parent.ID(), status)
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
	err := validator.ValidateStatusChange(parent.ID(), StatusDONE)
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
	err := validator.ValidateStatusChange(parent.ID(), StatusDONE)
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
	err := validator.ValidateStatusChange(leaf.ID(), StatusDONE)
	if err != nil {
		t.Errorf("Expected DONE status to be allowed for leaf task, but got error: %v", err)
	}
}

func TestTaskValidator_ValidateStatusChange_TaskNotFound(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	validator := NewTaskValidator(repo)

	// Create a non-existent task ID
	nonExistentID := NewTaskID()

	// Test that validation returns an error for non-existent task
	err := validator.ValidateStatusChange(nonExistentID, StatusDONE)
	if err == nil {
		t.Error("Expected error when validating status change for non-existent task, but got nil")
	}

	// Verify it's a NotFoundError
	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("Expected NotFoundError, but got: %T", err)
	}
}

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
	err := validator.ValidateStatusChange(parent1.ID(), StatusDONE)
	if err != nil {
		t.Errorf("Expected parent1 to be markable as DONE, but got error: %v", err)
	}

	// Test that root cannot be marked DONE (parent2 is not DONE)
	err = validator.ValidateStatusChange(root.ID(), StatusDONE)
	if err == nil {
		t.Error("Expected error when marking root as DONE with incomplete children, but got nil")
	}
}
