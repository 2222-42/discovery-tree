package domain

import (
	"strings"
	"testing"
)

func TestTaskService_CreateRootTask_Success(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	description := "Root task"
	task, err := service.CreateRootTask(description)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task == nil {
		t.Fatal("expected task to be created, got nil")
	}

	if task.Description() != description {
		t.Errorf("expected description %q, got %q", description, task.Description())
	}

	if task.Status() != StatusRootWorkItem {
		t.Errorf("expected status %v, got %v", StatusRootWorkItem, task.Status())
	}

	if task.ParentID() != nil {
		t.Errorf("expected nil parent ID for root task, got %v", task.ParentID())
	}

	if task.Position() != 0 {
		t.Errorf("expected position 0, got %d", task.Position())
	}

	// Verify task was saved to repository
	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Errorf("task not found in repository: %v", err)
	}
	if !retrieved.ID().Equals(task.ID()) {
		t.Error("retrieved task ID does not match")
	}
}

func TestTaskService_CreateRootTask_SingleRootConstraint(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create first root task
	_, err := service.CreateRootTask("First root")
	if err != nil {
		t.Fatalf("failed to create first root task: %v", err)
	}

	// Attempt to create second root task
	_, err = service.CreateRootTask("Second root")
	if err == nil {
		t.Fatal("expected error when creating second root task, got nil")
	}

	// Check that it's a ConstraintViolationError
	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("expected ConstraintViolationError, got %T", err)
	}

	// Check error message mentions root
	if !strings.Contains(err.Error(), "root") {
		t.Errorf("expected error message to mention 'root', got %q", err.Error())
	}
}

func TestTaskService_CreateRootTask_EmptyDescription(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	task, err := service.CreateRootTask("")
	if err == nil {
		t.Error("expected error for empty description, got nil")
	}

	if task != nil {
		t.Errorf("expected nil task for invalid description, got %v", task)
	}

	// Check that it's a ValidationError
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTaskService_CreateChildTask_Success(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent task
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent task: %v", err)
	}

	// Create child task
	description := "Child task"
	child, err := service.CreateChildTask(description, parent.ID())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if child == nil {
		t.Fatal("expected child task to be created, got nil")
	}

	if child.Description() != description {
		t.Errorf("expected description %q, got %q", description, child.Description())
	}

	if child.Status() != StatusTODO {
		t.Errorf("expected status %v, got %v", StatusTODO, child.Status())
	}

	if child.ParentID() == nil {
		t.Fatal("expected non-nil parent ID for child task")
	}

	if !child.ParentID().Equals(parent.ID()) {
		t.Errorf("expected parent ID %v, got %v", parent.ID(), *child.ParentID())
	}

	if child.Position() != 0 {
		t.Errorf("expected position 0 for first child, got %d", child.Position())
	}

	// Verify task was saved to repository
	retrieved, err := repo.FindByID(child.ID())
	if err != nil {
		t.Errorf("child task not found in repository: %v", err)
	}
	if !retrieved.ID().Equals(child.ID()) {
		t.Error("retrieved child task ID does not match")
	}
}

func TestTaskService_CreateChildTask_AutomaticPositionCalculation(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent task
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent task: %v", err)
	}

	// Create multiple child tasks
	child1, err := service.CreateChildTask("Child 1", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 1: %v", err)
	}

	child2, err := service.CreateChildTask("Child 2", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 2: %v", err)
	}

	child3, err := service.CreateChildTask("Child 3", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 3: %v", err)
	}

	// Verify positions are assigned sequentially
	if child1.Position() != 0 {
		t.Errorf("expected child1 position 0, got %d", child1.Position())
	}

	if child2.Position() != 1 {
		t.Errorf("expected child2 position 1, got %d", child2.Position())
	}

	if child3.Position() != 2 {
		t.Errorf("expected child3 position 2, got %d", child3.Position())
	}

	// Verify all children have the same parent
	if !child1.ParentID().Equals(parent.ID()) {
		t.Error("child1 has wrong parent ID")
	}
	if !child2.ParentID().Equals(parent.ID()) {
		t.Error("child2 has wrong parent ID")
	}
	if !child3.ParentID().Equals(parent.ID()) {
		t.Error("child3 has wrong parent ID")
	}

	// Verify children can be retrieved in order
	parentID := parent.ID()
	children, err := repo.FindByParentID(&parentID)
	if err != nil {
		t.Fatalf("failed to find children: %v", err)
	}

	if len(children) != 3 {
		t.Errorf("expected 3 children, got %d", len(children))
	}

	// Verify ordering
	if children[0].Position() != 0 || children[1].Position() != 1 || children[2].Position() != 2 {
		t.Error("children not in correct order")
	}
}

func TestTaskService_CreateChildTask_NonExistentParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Try to create child with non-existent parent
	nonExistentParentID := NewTaskID()
	child, err := service.CreateChildTask("Child task", nonExistentParentID)

	if err == nil {
		t.Fatal("expected error when creating child with non-existent parent, got nil")
	}

	if child != nil {
		t.Errorf("expected nil task for non-existent parent, got %v", child)
	}

	// Check that it's a NotFoundError
	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestTaskService_CreateChildTask_EmptyDescription(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent task
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent task: %v", err)
	}

	// Try to create child with empty description
	child, err := service.CreateChildTask("", parent.ID())

	if err == nil {
		t.Error("expected error for empty description, got nil")
	}

	if child != nil {
		t.Errorf("expected nil task for invalid description, got %v", child)
	}

	// Check that it's a ValidationError
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTaskService_CreateChildTask_NestedChildren(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create root task
	root, err := service.CreateRootTask("Root task")
	if err != nil {
		t.Fatalf("failed to create root task: %v", err)
	}

	// Create child of root
	child, err := service.CreateChildTask("Child task", root.ID())
	if err != nil {
		t.Fatalf("failed to create child task: %v", err)
	}

	// Create grandchild (child of child)
	grandchild, err := service.CreateChildTask("Grandchild task", child.ID())
	if err != nil {
		t.Fatalf("failed to create grandchild task: %v", err)
	}

	// Verify grandchild has correct parent and position
	if !grandchild.ParentID().Equals(child.ID()) {
		t.Errorf("expected grandchild parent ID %v, got %v", child.ID(), *grandchild.ParentID())
	}

	if grandchild.Position() != 0 {
		t.Errorf("expected grandchild position 0, got %d", grandchild.Position())
	}

	// Verify child has correct parent
	if !child.ParentID().Equals(root.ID()) {
		t.Errorf("expected child parent ID %v, got %v", root.ID(), *child.ParentID())
	}
}

func TestTaskService_CreateChildTask_MultipleParents(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create root task
	root, err := service.CreateRootTask("Root task")
	if err != nil {
		t.Fatalf("failed to create root task: %v", err)
	}

	// Create two children of root
	child1, err := service.CreateChildTask("Child 1", root.ID())
	if err != nil {
		t.Fatalf("failed to create child 1: %v", err)
	}

	child2, err := service.CreateChildTask("Child 2", root.ID())
	if err != nil {
		t.Fatalf("failed to create child 2: %v", err)
	}

	// Create children under each child
	grandchild1a, err := service.CreateChildTask("Grandchild 1a", child1.ID())
	if err != nil {
		t.Fatalf("failed to create grandchild 1a: %v", err)
	}

	grandchild1b, err := service.CreateChildTask("Grandchild 1b", child1.ID())
	if err != nil {
		t.Fatalf("failed to create grandchild 1b: %v", err)
	}

	grandchild2a, err := service.CreateChildTask("Grandchild 2a", child2.ID())
	if err != nil {
		t.Fatalf("failed to create grandchild 2a: %v", err)
	}

	// Verify positions are independent per parent
	if grandchild1a.Position() != 0 {
		t.Errorf("expected grandchild1a position 0, got %d", grandchild1a.Position())
	}
	if grandchild1b.Position() != 1 {
		t.Errorf("expected grandchild1b position 1, got %d", grandchild1b.Position())
	}
	if grandchild2a.Position() != 0 {
		t.Errorf("expected grandchild2a position 0, got %d", grandchild2a.Position())
	}

	// Verify parent references
	if !grandchild1a.ParentID().Equals(child1.ID()) {
		t.Error("grandchild1a has wrong parent")
	}
	if !grandchild1b.ParentID().Equals(child1.ID()) {
		t.Error("grandchild1b has wrong parent")
	}
	if !grandchild2a.ParentID().Equals(child2.ID()) {
		t.Error("grandchild2a has wrong parent")
	}
}

func TestTaskService_ChangeTaskStatus_ValidStatusChange(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create a task
	task, err := service.CreateRootTask("Test task")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	// Change status to InProgress
	err = service.ChangeTaskStatus(task.ID(), StatusInProgress)
	if err != nil {
		t.Errorf("expected no error when changing to valid status, got %v", err)
	}

	// Verify status was changed
	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.Status() != StatusInProgress {
		t.Errorf("expected status %v, got %v", StatusInProgress, retrieved.Status())
	}
}

func TestTaskService_ChangeTaskStatus_NonDONEStatusesAllowed(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with incomplete children
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create child that is not DONE
	_, err = service.CreateChildTask("Child task", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child: %v", err)
	}

	// Try to change parent to non-DONE statuses - should all succeed
	nonDoneStatuses := []Status{StatusTODO, StatusInProgress, StatusBlocked}
	for _, status := range nonDoneStatuses {
		err = service.ChangeTaskStatus(parent.ID(), status)
		if err != nil {
			t.Errorf("expected no error when changing to %v with incomplete children, got %v", status, err)
		}

		// Verify status was changed
		retrieved, err := repo.FindByID(parent.ID())
		if err != nil {
			t.Fatalf("failed to retrieve task: %v", err)
		}

		if retrieved.Status() != status {
			t.Errorf("expected status %v, got %v", status, retrieved.Status())
		}
	}
}

func TestTaskService_ChangeTaskStatus_DONEWithIncompleteChildren(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create child that is not DONE
	child, err := service.CreateChildTask("Child task", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child: %v", err)
	}

	// Try to mark parent as DONE - should fail
	err = service.ChangeTaskStatus(parent.ID(), StatusDONE)
	if err == nil {
		t.Fatal("expected error when marking parent DONE with incomplete children, got nil")
	}

	// Check that it's a ConstraintViolationError
	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("expected ConstraintViolationError, got %T", err)
	}

	// Verify parent status didn't change
	retrieved, err := repo.FindByID(parent.ID())
	if err != nil {
		t.Fatalf("failed to retrieve parent: %v", err)
	}

	if retrieved.Status() == StatusDONE {
		t.Error("parent status should not have changed to DONE")
	}

	// Now mark child as DONE
	err = child.ChangeStatus(StatusDONE)
	if err != nil {
		t.Fatalf("failed to change child status: %v", err)
	}
	err = repo.Save(child)
	if err != nil {
		t.Fatalf("failed to save child: %v", err)
	}

	// Now marking parent as DONE should succeed
	err = service.ChangeTaskStatus(parent.ID(), StatusDONE)
	if err != nil {
		t.Errorf("expected no error when marking parent DONE with all children DONE, got %v", err)
	}

	// Verify parent status changed
	retrieved, err = repo.FindByID(parent.ID())
	if err != nil {
		t.Fatalf("failed to retrieve parent: %v", err)
	}

	if retrieved.Status() != StatusDONE {
		t.Errorf("expected parent status %v, got %v", StatusDONE, retrieved.Status())
	}
}

func TestTaskService_ChangeTaskStatus_DONEWithNoChildren(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create a task with no children
	task, err := service.CreateRootTask("Leaf task")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	// Mark as DONE - should succeed
	err = service.ChangeTaskStatus(task.ID(), StatusDONE)
	if err != nil {
		t.Errorf("expected no error when marking leaf task DONE, got %v", err)
	}

	// Verify status changed
	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.Status() != StatusDONE {
		t.Errorf("expected status %v, got %v", StatusDONE, retrieved.Status())
	}
}

func TestTaskService_ChangeTaskStatus_DONEWithAllChildrenDONE(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent
	parent, err := service.CreateRootTask("Parent task")
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	// Create multiple children
	child1, err := service.CreateChildTask("Child 1", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 1: %v", err)
	}

	child2, err := service.CreateChildTask("Child 2", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 2: %v", err)
	}

	child3, err := service.CreateChildTask("Child 3", parent.ID())
	if err != nil {
		t.Fatalf("failed to create child 3: %v", err)
	}

	// Mark all children as DONE
	for _, child := range []*Task{child1, child2, child3} {
		err = child.ChangeStatus(StatusDONE)
		if err != nil {
			t.Fatalf("failed to change child status: %v", err)
		}
		err = repo.Save(child)
		if err != nil {
			t.Fatalf("failed to save child: %v", err)
		}
	}

	// Now marking parent as DONE should succeed
	err = service.ChangeTaskStatus(parent.ID(), StatusDONE)
	if err != nil {
		t.Errorf("expected no error when marking parent DONE with all children DONE, got %v", err)
	}

	// Verify parent status changed
	retrieved, err := repo.FindByID(parent.ID())
	if err != nil {
		t.Fatalf("failed to retrieve parent: %v", err)
	}

	if retrieved.Status() != StatusDONE {
		t.Errorf("expected parent status %v, got %v", StatusDONE, retrieved.Status())
	}
}

func TestTaskService_ChangeTaskStatus_InvalidStatus(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create a task
	task, err := service.CreateRootTask("Test task")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	initialStatus := task.Status()

	// Try to change to invalid status
	invalidStatus := Status(-1)
	err = service.ChangeTaskStatus(task.ID(), invalidStatus)
	if err == nil {
		t.Fatal("expected error when changing to invalid status, got nil")
	}

	// Check that it's a ValidationError
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Verify status didn't change
	retrieved, err := repo.FindByID(task.ID())
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrieved.Status() != initialStatus {
		t.Errorf("expected status to remain %v, got %v", initialStatus, retrieved.Status())
	}
}

func TestTaskService_ChangeTaskStatus_NonExistentTask(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Try to change status of non-existent task
	nonExistentID := NewTaskID()
	err := service.ChangeTaskStatus(nonExistentID, StatusInProgress)

	if err == nil {
		t.Fatal("expected error when changing status of non-existent task, got nil")
	}

	// Check that it's a NotFoundError
	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestTaskService_ChangeTaskStatus_BottomToTopEnforcement(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create a tree: root -> child1 -> grandchild
	//                     -> child2
	root, err := service.CreateRootTask("Root")
	if err != nil {
		t.Fatalf("failed to create root: %v", err)
	}

	child1, err := service.CreateChildTask("Child 1", root.ID())
	if err != nil {
		t.Fatalf("failed to create child1: %v", err)
	}

	child2, err := service.CreateChildTask("Child 2", root.ID())
	if err != nil {
		t.Fatalf("failed to create child2: %v", err)
	}

	grandchild, err := service.CreateChildTask("Grandchild", child1.ID())
	if err != nil {
		t.Fatalf("failed to create grandchild: %v", err)
	}

	// Try to mark child1 as DONE - should fail (grandchild not DONE)
	err = service.ChangeTaskStatus(child1.ID(), StatusDONE)
	if err == nil {
		t.Error("expected error when marking child1 DONE with incomplete grandchild")
	}

	// Mark grandchild as DONE
	err = grandchild.ChangeStatus(StatusDONE)
	if err != nil {
		t.Fatalf("failed to change grandchild status: %v", err)
	}
	err = repo.Save(grandchild)
	if err != nil {
		t.Fatalf("failed to save grandchild: %v", err)
	}

	// Now child1 can be marked DONE
	err = service.ChangeTaskStatus(child1.ID(), StatusDONE)
	if err != nil {
		t.Errorf("expected no error when marking child1 DONE with grandchild DONE, got %v", err)
	}

	// Try to mark root as DONE - should fail (child2 not DONE)
	err = service.ChangeTaskStatus(root.ID(), StatusDONE)
	if err == nil {
		t.Error("expected error when marking root DONE with incomplete child2")
	}

	// Mark child2 as DONE
	err = child2.ChangeStatus(StatusDONE)
	if err != nil {
		t.Fatalf("failed to change child2 status: %v", err)
	}
	err = repo.Save(child2)
	if err != nil {
		t.Fatalf("failed to save child2: %v", err)
	}

	// Now root can be marked DONE
	err = service.ChangeTaskStatus(root.ID(), StatusDONE)
	if err != nil {
		t.Errorf("expected no error when marking root DONE with all children DONE, got %v", err)
	}

	// Verify all tasks are DONE
	allTasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("failed to find all tasks: %v", err)
	}

	for _, task := range allTasks {
		if task.Status() != StatusDONE {
			t.Errorf("expected all tasks to be DONE, but task %v has status %v", task.ID(), task.Status())
		}
	}
}

func TestTaskService_MoveTask_WithinSameParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with 3 children
	parent, _ := service.CreateRootTask("Parent")
	child1, _ := service.CreateChildTask("Child 1", parent.ID())
	child2, _ := service.CreateChildTask("Child 2", parent.ID())
	child3, _ := service.CreateChildTask("Child 3", parent.ID())

	// Move child1 from position 0 to position 2
	err := service.MoveTask(child1.ID(), &parent.id, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify new positions: child2(0), child3(1), child1(2)
	children, _ := repo.FindByParentID(&parent.id)
	if len(children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(children))
	}

	// Check positions
	if children[0].ID().Equals(child2.ID()) && children[0].Position() != 0 {
		t.Errorf("expected child2 at position 0, got %d", children[0].Position())
	}
	if children[1].ID().Equals(child3.ID()) && children[1].Position() != 1 {
		t.Errorf("expected child3 at position 1, got %d", children[1].Position())
	}
	if children[2].ID().Equals(child1.ID()) && children[2].Position() != 2 {
		t.Errorf("expected child1 at position 2, got %d", children[2].Position())
	}
}

func TestTaskService_MoveTask_ToDifferentParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create two parents with children
	parent1, _ := service.CreateRootTask("Parent 1")
	parent2, _ := service.CreateChildTask("Parent 2", parent1.ID())
	
	child1, _ := service.CreateChildTask("Child 1", parent1.ID())
	child2, _ := service.CreateChildTask("Child 2", parent1.ID())
	_, _ = service.CreateChildTask("Child 3", parent2.ID())

	// Move child1 from parent1 to parent2
	err := service.MoveTask(child1.ID(), &parent2.id, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify child1 is now under parent2
	retrieved, _ := repo.FindByID(child1.ID())
	if !retrieved.ParentID().Equals(parent2.ID()) {
		t.Errorf("expected child1 parent to be parent2, got %v", retrieved.ParentID())
	}
	if retrieved.Position() != 1 {
		t.Errorf("expected child1 position 1, got %d", retrieved.Position())
	}

	// Verify parent1 children positions adjusted
	// parent1 should have parent2 and child2
	parent1Children, _ := repo.FindByParentID(&parent1.id)
	if len(parent1Children) != 2 {
		t.Fatalf("expected 2 children under parent1, got %d", len(parent1Children))
	}
	// parent2 should be at position 0, child2 at position 1
	if parent1Children[0].ID().Equals(parent2.ID()) && parent1Children[0].Position() != 0 {
		t.Errorf("expected parent2 at position 0, got %d", parent1Children[0].Position())
	}
	if parent1Children[1].ID().Equals(child2.ID()) && parent1Children[1].Position() != 1 {
		t.Errorf("expected child2 at position 1, got %d", parent1Children[1].Position())
	}

	// Verify parent2 children positions adjusted
	parent2Children, _ := repo.FindByParentID(&parent2.id)
	if len(parent2Children) != 2 {
		t.Fatalf("expected 2 children under parent2, got %d", len(parent2Children))
	}
}

func TestTaskService_MoveTask_SubtreeMovesWithParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create tree: root -> parent1 -> child -> grandchild
	//                   -> parent2
	root, _ := service.CreateRootTask("Root")
	parent1, _ := service.CreateChildTask("Parent 1", root.ID())
	parent2, _ := service.CreateChildTask("Parent 2", root.ID())
	child, _ := service.CreateChildTask("Child", parent1.ID())
	grandchild, _ := service.CreateChildTask("Grandchild", child.ID())

	// Move parent1 (with its subtree) under parent2
	err := service.MoveTask(parent1.ID(), &parent2.id, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify parent1 is now under parent2
	retrievedParent1, _ := repo.FindByID(parent1.ID())
	if !retrievedParent1.ParentID().Equals(parent2.ID()) {
		t.Errorf("expected parent1 under parent2")
	}

	// Verify child is still under parent1
	retrievedChild, _ := repo.FindByID(child.ID())
	if !retrievedChild.ParentID().Equals(parent1.ID()) {
		t.Errorf("expected child still under parent1")
	}

	// Verify grandchild is still under child
	retrievedGrandchild, _ := repo.FindByID(grandchild.ID())
	if !retrievedGrandchild.ParentID().Equals(child.ID()) {
		t.Errorf("expected grandchild still under child")
	}
}

func TestTaskService_MoveTask_PreventCycle(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create tree: root -> parent -> child
	root, _ := service.CreateRootTask("Root")
	parent, _ := service.CreateChildTask("Parent", root.ID())
	child, _ := service.CreateChildTask("Child", parent.ID())

	// Try to move parent under its own child (creates cycle)
	err := service.MoveTask(parent.ID(), &child.id, 0)
	if err == nil {
		t.Fatal("expected error when creating cycle, got nil")
	}

	if _, ok := err.(ConstraintViolationError); !ok {
		t.Errorf("expected ConstraintViolationError, got %T", err)
	}

	// Verify parent is still under root
	retrievedParent, _ := repo.FindByID(parent.ID())
	if !retrievedParent.ParentID().Equals(root.ID()) {
		t.Errorf("expected parent still under root after failed move")
	}
}

func TestTaskService_MoveTask_NoOpMove(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with child
	parent, _ := service.CreateRootTask("Parent")
	child, _ := service.CreateChildTask("Child", parent.ID())

	// Move child to same parent and position (no-op)
	err := service.MoveTask(child.ID(), &parent.id, 0)
	if err != nil {
		t.Errorf("expected no error for no-op move, got %v", err)
	}

	// Verify child is still at position 0
	retrieved, _ := repo.FindByID(child.ID())
	if retrieved.Position() != 0 {
		t.Errorf("expected position 0, got %d", retrieved.Position())
	}
}

func TestTaskService_MoveTask_PositionAdjustmentLeft(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with 4 children
	parent, _ := service.CreateRootTask("Parent")
	child1, _ := service.CreateChildTask("Child 1", parent.ID())
	child2, _ := service.CreateChildTask("Child 2", parent.ID())
	child3, _ := service.CreateChildTask("Child 3", parent.ID())
	child4, _ := service.CreateChildTask("Child 4", parent.ID())

	// Move child3 (position 2) to position 0
	err := service.MoveTask(child3.ID(), &parent.id, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify new order: child3(0), child1(1), child2(2), child4(3)
	children, _ := repo.FindByParentID(&parent.id)
	if len(children) != 4 {
		t.Fatalf("expected 4 children, got %d", len(children))
	}

	expectedOrder := []TaskID{child3.ID(), child1.ID(), child2.ID(), child4.ID()}
	for i, expected := range expectedOrder {
		if !children[i].ID().Equals(expected) {
			t.Errorf("position %d: expected %v, got %v", i, expected, children[i].ID())
		}
		if children[i].Position() != i {
			t.Errorf("child at index %d has wrong position: expected %d, got %d", i, i, children[i].Position())
		}
	}
}

func TestTaskService_MoveTask_PositionAdjustmentRight(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with 4 children
	parent, _ := service.CreateRootTask("Parent")
	child1, _ := service.CreateChildTask("Child 1", parent.ID())
	child2, _ := service.CreateChildTask("Child 2", parent.ID())
	child3, _ := service.CreateChildTask("Child 3", parent.ID())
	child4, _ := service.CreateChildTask("Child 4", parent.ID())

	// Move child1 (position 0) to position 3
	err := service.MoveTask(child1.ID(), &parent.id, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify new order: child2(0), child3(1), child4(2), child1(3)
	children, _ := repo.FindByParentID(&parent.id)
	if len(children) != 4 {
		t.Fatalf("expected 4 children, got %d", len(children))
	}

	expectedOrder := []TaskID{child2.ID(), child3.ID(), child4.ID(), child1.ID()}
	for i, expected := range expectedOrder {
		if !children[i].ID().Equals(expected) {
			t.Errorf("position %d: expected %v, got %v", i, expected, children[i].ID())
		}
		if children[i].Position() != i {
			t.Errorf("child at index %d has wrong position: expected %d, got %d", i, i, children[i].Position())
		}
	}
}

func TestTaskService_MoveTask_NonExistentTask(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	parent, _ := service.CreateRootTask("Parent")
	nonExistentID := NewTaskID()

	err := service.MoveTask(nonExistentID, &parent.id, 0)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}

	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestTaskService_MoveTask_NonExistentParent(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	root, _ := service.CreateRootTask("Root")
	child, _ := service.CreateChildTask("Child", root.ID())
	nonExistentParentID := NewTaskID()

	err := service.MoveTask(child.ID(), &nonExistentParentID, 0)
	if err == nil {
		t.Fatal("expected error for non-existent parent, got nil")
	}

	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestTaskService_MoveTask_NegativePosition(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	root, _ := service.CreateRootTask("Root")
	child, _ := service.CreateChildTask("Child", root.ID())

	err := service.MoveTask(child.ID(), &root.id, -1)
	if err == nil {
		t.Fatal("expected error for negative position, got nil")
	}

	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTaskService_MoveTask_ComplexReordering(t *testing.T) {
	repo := NewInMemoryTaskRepository()
	service := NewTaskService(repo)

	// Create parent with 5 children
	parent, _ := service.CreateRootTask("Parent")
	child1, _ := service.CreateChildTask("Child 1", parent.ID())
	child2, _ := service.CreateChildTask("Child 2", parent.ID())
	_, _ = service.CreateChildTask("Child 3", parent.ID())
	_, _ = service.CreateChildTask("Child 4", parent.ID())
	child5, _ := service.CreateChildTask("Child 5", parent.ID())

	// Perform multiple moves
	service.MoveTask(child2.ID(), &parent.id, 4) // Move child2 to end
	service.MoveTask(child5.ID(), &parent.id, 0) // Move child5 to start
	service.MoveTask(child1.ID(), &parent.id, 2) // Move child1 to middle

	// Verify final order
	children, _ := repo.FindByParentID(&parent.id)
	if len(children) != 5 {
		t.Fatalf("expected 5 children, got %d", len(children))
	}

	// Verify all positions are sequential
	for i, child := range children {
		if child.Position() != i {
			t.Errorf("child at index %d has wrong position: expected %d, got %d", i, i, child.Position())
		}
	}
}
