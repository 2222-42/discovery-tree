package domain

import (
	"testing"
)

func TestInMemoryTaskRepository_Save(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task, err := NewTask("Test task", nil, 0)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = repo.Save(task)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Verify task was saved
	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Errorf("FindByID failed: %v", err)
	}
	if !retrieved.ID().Equals(task.ID()) {
		t.Errorf("Retrieved task ID mismatch: got %v, want %v", retrieved.ID(), task.ID())
	}
}

func TestInMemoryTaskRepository_SaveNil(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	err := repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil task")
	}
}

func TestInMemoryTaskRepository_FindByID(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task, _ := NewTask("Test task", nil, 0)
	repo.Save(task)

	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Errorf("FindByID failed: %v", err)
	}
	if !retrieved.ID().Equals(task.ID()) {
		t.Errorf("Retrieved task ID mismatch")
	}
}

func TestInMemoryTaskRepository_FindByIDNotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	nonExistentID := NewTaskID()
	_, err := repo.FindByID(nonExistentID)
	if err == nil {
		t.Error("Expected NotFoundError for non-existent task")
	}
	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("Expected NotFoundError, got %T", err)
	}
}

func TestInMemoryTaskRepository_FindByParentID(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	// Create parent task
	parent, _ := NewTask("Parent", nil, 0)
	repo.Save(parent)

	// Create child tasks
	child1, _ := NewTask("Child 1", &parent.id, 0)
	child2, _ := NewTask("Child 2", &parent.id, 1)
	child3, _ := NewTask("Child 3", &parent.id, 2)
	repo.Save(child1)
	repo.Save(child2)
	repo.Save(child3)

	// Find children
	children, err := repo.FindByParentID(&parent.id)
	if err != nil {
		t.Errorf("FindByParentID failed: %v", err)
	}

	if len(children) != 3 {
		t.Errorf("Expected 3 children, got %d", len(children))
	}

	// Verify ordering by position
	if children[0].Position() != 0 || children[1].Position() != 1 || children[2].Position() != 2 {
		t.Error("Children not ordered by position")
	}
}

func TestInMemoryTaskRepository_FindRoot(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	// Create root task
	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	// Create child task
	child, _ := NewTask("Child", &root.id, 0)
	repo.Save(child)

	// Find root
	foundRoot, err := repo.FindRoot()
	if err != nil {
		t.Errorf("FindRoot failed: %v", err)
	}
	if !foundRoot.ID().Equals(root.ID()) {
		t.Error("Found wrong root task")
	}
}

func TestInMemoryTaskRepository_FindRootNotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	_, err := repo.FindRoot()
	if err == nil {
		t.Error("Expected NotFoundError when no root exists")
	}
}

func TestInMemoryTaskRepository_FindAll(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task1, _ := NewTask("Task 1", nil, 0)
	task2, _ := NewTask("Task 2", &task1.id, 0)
	task3, _ := NewTask("Task 3", &task1.id, 1)

	repo.Save(task1)
	repo.Save(task2)
	repo.Save(task3)

	all, err := repo.FindAll()
	if err != nil {
		t.Errorf("FindAll failed: %v", err)
	}

	if len(all) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(all))
	}
}

func TestInMemoryTaskRepository_Delete(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task, _ := NewTask("Task", nil, 0)
	repo.Save(task)

	err := repo.Delete(task.ID())
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify task was deleted
	_, err = repo.FindByID(task.ID())
	if err == nil {
		t.Error("Expected NotFoundError after deletion")
	}
}

func TestInMemoryTaskRepository_DeleteNotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	nonExistentID := NewTaskID()
	err := repo.Delete(nonExistentID)
	if err == nil {
		t.Error("Expected NotFoundError when deleting non-existent task")
	}
}

func TestInMemoryTaskRepository_DeleteSubtree(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	// Create a tree structure:
	//   root
	//   ├── child1
	//   │   ├── grandchild1
	//   │   └── grandchild2
	//   └── child2

	root, _ := NewTask("Root", nil, 0)
	repo.Save(root)

	child1, _ := NewTask("Child 1", &root.id, 0)
	child2, _ := NewTask("Child 2", &root.id, 1)
	repo.Save(child1)
	repo.Save(child2)

	grandchild1, _ := NewTask("Grandchild 1", &child1.id, 0)
	grandchild2, _ := NewTask("Grandchild 2", &child1.id, 1)
	repo.Save(grandchild1)
	repo.Save(grandchild2)

	// Delete child1 subtree
	err := repo.DeleteSubtree(child1.ID())
	if err != nil {
		t.Errorf("DeleteSubtree failed: %v", err)
	}

	// Verify child1 and its descendants are deleted
	_, err = repo.FindByID(child1.ID())
	if err == nil {
		t.Error("child1 should be deleted")
	}
	_, err = repo.FindByID(grandchild1.ID())
	if err == nil {
		t.Error("grandchild1 should be deleted")
	}
	_, err = repo.FindByID(grandchild2.ID())
	if err == nil {
		t.Error("grandchild2 should be deleted")
	}

	// Verify root and child2 still exist
	_, err = repo.FindByID(root.ID())
	if err != nil {
		t.Error("root should still exist")
	}
	_, err = repo.FindByID(child2.ID())
	if err != nil {
		t.Error("child2 should still exist")
	}
}

func TestInMemoryTaskRepository_DeleteSubtreeNotFound(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	nonExistentID := NewTaskID()
	err := repo.DeleteSubtree(nonExistentID)
	if err == nil {
		t.Error("Expected NotFoundError when deleting non-existent subtree")
	}
}
