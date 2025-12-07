package infrastructure

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	"discovery-tree/domain"
)

// TestNewFileTaskRepository_DefaultPath tests that default path is used when empty string is provided
func TestNewFileTaskRepository_DefaultPath(t *testing.T) {
	// Clean up before and after
	defaultPath := "./data/tasks.json"
	os.RemoveAll("./data")
	defer os.RemoveAll("./data")

	repo, err := NewFileTaskRepository("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.filePath != defaultPath {
		t.Errorf("expected filePath to be %s, got %s", defaultPath, repo.filePath)
	}

	// Verify directory was created
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

// TestNewFileTaskRepository_CustomPath tests that custom path is used when provided
func TestNewFileTaskRepository_CustomPath(t *testing.T) {
	customPath := "./test_data/custom_tasks.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(customPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.filePath != customPath {
		t.Errorf("expected filePath to be %s, got %s", customPath, repo.filePath)
	}

	// Verify directory was created
	if _, err := os.Stat("./test_data"); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

// TestNewFileTaskRepository_CreatesNestedDirectories tests that nested directories are created
func TestNewFileTaskRepository_CreatesNestedDirectories(t *testing.T) {
	nestedPath := "./test_data/nested/deep/tasks.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(nestedPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.filePath != nestedPath {
		t.Errorf("expected filePath to be %s, got %s", nestedPath, repo.filePath)
	}

	// Verify nested directories were created
	if _, err := os.Stat("./test_data/nested/deep"); os.IsNotExist(err) {
		t.Error("expected nested directories to be created")
	}
}

// TestNewFileTaskRepository_NonExistentFile tests initialization with non-existent file
func TestNewFileTaskRepository_NonExistentFile(t *testing.T) {
	testPath := "./test_data/nonexistent.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got %v", err)
	}

	// Should initialize with empty task map
	if len(repo.tasks) != 0 {
		t.Errorf("expected empty task map, got %d tasks", len(repo.tasks))
	}
}

// TestNewFileTaskRepository_EmptyFile tests initialization with empty file
func TestNewFileTaskRepository_EmptyFile(t *testing.T) {
	testPath := "./test_data/empty.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Create empty file
	os.MkdirAll("./test_data", 0755)
	os.WriteFile(testPath, []byte(""), 0644)

	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error for empty file, got %v", err)
	}

	// Should initialize with empty task map
	if len(repo.tasks) != 0 {
		t.Errorf("expected empty task map, got %d tasks", len(repo.tasks))
	}
}

// TestNewFileTaskRepository_LoadsExistingTasks tests loading existing tasks from file
func TestNewFileTaskRepository_LoadsExistingTasks(t *testing.T) {
	testPath := "./test_data/existing.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Create test tasks
	task1, _ := domain.NewTask("Task 1", nil, 0)
	task2ID := task1.ID()
	task2, _ := domain.NewTask("Task 2", &task2ID, 0)

	// Write tasks to file
	os.MkdirAll("./test_data", 0755)
	dtos := []TaskDTO{
		ToDTO(task1),
		ToDTO(task2),
	}
	data, _ := json.MarshalIndent(dtos, "", "  ")
	os.WriteFile(testPath, data, 0644)

	// Load repository
	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify tasks were loaded
	if len(repo.tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(repo.tasks))
	}

	// Verify task data
	loadedTask1, exists := repo.tasks[task1.ID().String()]
	if !exists {
		t.Error("expected task1 to be loaded")
	}
	if loadedTask1.Description() != "Task 1" {
		t.Errorf("expected description 'Task 1', got '%s'", loadedTask1.Description())
	}

	loadedTask2, exists := repo.tasks[task2.ID().String()]
	if !exists {
		t.Error("expected task2 to be loaded")
	}
	if loadedTask2.Description() != "Task 2" {
		t.Errorf("expected description 'Task 2', got '%s'", loadedTask2.Description())
	}
}

// TestNewFileTaskRepository_InvalidJSON tests error handling for invalid JSON
func TestNewFileTaskRepository_InvalidJSON(t *testing.T) {
	testPath := "./test_data/invalid.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write invalid JSON
	os.MkdirAll("./test_data", 0755)
	os.WriteFile(testPath, []byte("{invalid json}"), 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	// Should be a FileSystemError
	if _, ok := err.(FileSystemError); !ok {
		t.Errorf("expected FileSystemError, got %T", err)
	}
}

// TestNewFileTaskRepository_InvalidTaskData tests error handling for invalid task data
func TestNewFileTaskRepository_InvalidTaskData(t *testing.T) {
	testPath := "./test_data/invalid_task.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write JSON with invalid task data (empty description)
	os.MkdirAll("./test_data", 0755)
	invalidDTO := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "", // Invalid: empty description
		Status:      "TODO",
		ParentID:    nil,
		Position:    0,
	}
	data, _ := json.Marshal([]TaskDTO{invalidDTO})
	os.WriteFile(testPath, data, 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for invalid task data, got nil")
	}

	// Should be a ValidationError
	if _, ok := err.(domain.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

// TestNewFileTaskRepository_InvalidTaskID tests error handling for invalid task ID
func TestNewFileTaskRepository_InvalidTaskID(t *testing.T) {
	testPath := "./test_data/invalid_id.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write JSON with invalid task ID
	os.MkdirAll("./test_data", 0755)
	invalidDTO := TaskDTO{
		ID:          "not-a-uuid",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    nil,
		Position:    0,
	}
	data, _ := json.Marshal([]TaskDTO{invalidDTO})
	os.WriteFile(testPath, data, 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for invalid task ID, got nil")
	}
}

// TestNewFileTaskRepository_InvalidStatus tests error handling for invalid status
func TestNewFileTaskRepository_InvalidStatus(t *testing.T) {
	testPath := "./test_data/invalid_status.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write JSON with invalid status
	os.MkdirAll("./test_data", 0755)
	invalidDTO := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "INVALID_STATUS",
		ParentID:    nil,
		Position:    0,
	}
	data, _ := json.Marshal([]TaskDTO{invalidDTO})
	os.WriteFile(testPath, data, 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

// TestNewFileTaskRepository_NegativePosition tests error handling for negative position
func TestNewFileTaskRepository_NegativePosition(t *testing.T) {
	testPath := "./test_data/negative_position.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write JSON with negative position
	os.MkdirAll("./test_data", 0755)
	invalidDTO := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    nil,
		Position:    -1, // Invalid: negative position
	}
	data, _ := json.Marshal([]TaskDTO{invalidDTO})
	os.WriteFile(testPath, data, 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for negative position, got nil")
	}

	// Should be a ValidationError
	if _, ok := err.(domain.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

// TestNewFileTaskRepository_InvalidParentID tests error handling for invalid parent ID
func TestNewFileTaskRepository_InvalidParentID(t *testing.T) {
	testPath := "./test_data/invalid_parent.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Write JSON with invalid parent ID
	os.MkdirAll("./test_data", 0755)
	invalidParentID := "not-a-uuid"
	invalidDTO := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    &invalidParentID,
		Position:    0,
	}
	data, _ := json.Marshal([]TaskDTO{invalidDTO})
	os.WriteFile(testPath, data, 0644)

	// Should return error
	_, err := NewFileTaskRepository(testPath)
	if err == nil {
		t.Fatal("expected error for invalid parent ID, got nil")
	}
}

// TestNewFileTaskRepository_InvalidPath tests error handling for invalid file path
func TestNewFileTaskRepository_InvalidPath(t *testing.T) {
	// Try to create a file in a location that requires root permissions
	// This test may behave differently on different systems
	invalidPath := "/root/cannot_create/tasks.json"

	_, err := NewFileTaskRepository(invalidPath)
	if err == nil {
		// If no error, the system allowed the creation (unlikely but possible)
		// Clean up
		os.RemoveAll("/root/cannot_create")
		return
	}

	// Should be a FileSystemError
	if _, ok := err.(FileSystemError); !ok {
		t.Errorf("expected FileSystemError, got %T", err)
	}
}

// TestLoad_PreservesTaskRelationships tests that parent-child relationships are preserved
func TestLoad_PreservesTaskRelationships(t *testing.T) {
	testPath := "./test_data/relationships.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Create task tree: root -> child1, child2
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child1, _ := domain.NewTask("Child 1", &rootID, 0)
	child2, _ := domain.NewTask("Child 2", &rootID, 1)

	// Write to file
	os.MkdirAll("./test_data", 0755)
	dtos := []TaskDTO{
		ToDTO(root),
		ToDTO(child1),
		ToDTO(child2),
	}
	data, _ := json.MarshalIndent(dtos, "", "  ")
	os.WriteFile(testPath, data, 0644)

	// Load repository
	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify relationships
	loadedRoot := repo.tasks[root.ID().String()]
	if loadedRoot.ParentID() != nil {
		t.Error("expected root to have nil parent")
	}

	loadedChild1 := repo.tasks[child1.ID().String()]
	if loadedChild1.ParentID() == nil {
		t.Error("expected child1 to have parent")
	} else if loadedChild1.ParentID().String() != rootID.String() {
		t.Error("expected child1 parent to match root ID")
	}

	loadedChild2 := repo.tasks[child2.ID().String()]
	if loadedChild2.ParentID() == nil {
		t.Error("expected child2 to have parent")
	} else if loadedChild2.ParentID().String() != rootID.String() {
		t.Error("expected child2 parent to match root ID")
	}

	// Verify positions
	if loadedChild1.Position() != 0 {
		t.Errorf("expected child1 position 0, got %d", loadedChild1.Position())
	}
	if loadedChild2.Position() != 1 {
		t.Errorf("expected child2 position 1, got %d", loadedChild2.Position())
	}
}

// TestSave_CreatesNewTask tests that Save creates a new task
func TestSave_CreatesNewTask(t *testing.T) {
	testPath := "./test_data/save_new.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error creating repository, got %v", err)
	}

	// Create and save a task
	task, _ := domain.NewTask("New Task", nil, 0)
	err = repo.Save(task)
	if err != nil {
		t.Fatalf("expected no error saving task, got %v", err)
	}

	// Verify task is in memory
	if len(repo.tasks) != 1 {
		t.Errorf("expected 1 task in memory, got %d", len(repo.tasks))
	}

	savedTask, exists := repo.tasks[task.ID().String()]
	if !exists {
		t.Fatal("expected task to be in memory")
	}
	if savedTask.Description() != "New Task" {
		t.Errorf("expected description 'New Task', got '%s'", savedTask.Description())
	}

	// Verify task was persisted to file
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}

	var dtos []TaskDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if len(dtos) != 1 {
		t.Errorf("expected 1 task in file, got %d", len(dtos))
	}
	if dtos[0].Description != "New Task" {
		t.Errorf("expected description 'New Task', got '%s'", dtos[0].Description)
	}
}

// TestSave_UpdatesExistingTask tests that Save updates an existing task
func TestSave_UpdatesExistingTask(t *testing.T) {
	testPath := "./test_data/save_update.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error creating repository, got %v", err)
	}

	// Create and save a task
	task, _ := domain.NewTask("Original Task", nil, 0)
	repo.Save(task)

	// Update the task
	task.UpdateDescription("Updated Task")
	err = repo.Save(task)
	if err != nil {
		t.Fatalf("expected no error updating task, got %v", err)
	}

	// Verify task is updated in memory
	savedTask := repo.tasks[task.ID().String()]
	if savedTask.Description() != "Updated Task" {
		t.Errorf("expected description 'Updated Task', got '%s'", savedTask.Description())
	}

	// Verify only one task exists
	if len(repo.tasks) != 1 {
		t.Errorf("expected 1 task in memory, got %d", len(repo.tasks))
	}

	// Verify task was persisted to file
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}

	var dtos []TaskDTO
	json.Unmarshal(data, &dtos)

	if len(dtos) != 1 {
		t.Errorf("expected 1 task in file, got %d", len(dtos))
	}
	if dtos[0].Description != "Updated Task" {
		t.Errorf("expected description 'Updated Task', got '%s'", dtos[0].Description)
	}
}

// TestSave_MultipleTasks tests saving multiple tasks
func TestSave_MultipleTasks(t *testing.T) {
	testPath := "./test_data/save_multiple.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error creating repository, got %v", err)
	}

	// Create and save multiple tasks
	task1, _ := domain.NewTask("Task 1", nil, 0)
	task2, _ := domain.NewTask("Task 2", nil, 1)
	task3, _ := domain.NewTask("Task 3", nil, 2)

	repo.Save(task1)
	repo.Save(task2)
	repo.Save(task3)

	// Verify all tasks are in memory
	if len(repo.tasks) != 3 {
		t.Errorf("expected 3 tasks in memory, got %d", len(repo.tasks))
	}

	// Verify all tasks were persisted to file
	data, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}

	var dtos []TaskDTO
	json.Unmarshal(data, &dtos)

	if len(dtos) != 3 {
		t.Errorf("expected 3 tasks in file, got %d", len(dtos))
	}
}

// TestSave_PersistsToFile tests that Save writes to file
func TestSave_PersistsToFile(t *testing.T) {
	testPath := "./test_data/save_persist.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	// Create repository and save a task
	repo1, _ := NewFileTaskRepository(testPath)
	task, _ := domain.NewTask("Persisted Task", nil, 0)
	repo1.Save(task)

	// Create a new repository instance (loads from file)
	repo2, err := NewFileTaskRepository(testPath)
	if err != nil {
		t.Fatalf("expected no error loading repository, got %v", err)
	}

	// Verify task was loaded from file
	if len(repo2.tasks) != 1 {
		t.Errorf("expected 1 task loaded from file, got %d", len(repo2.tasks))
	}

	loadedTask, exists := repo2.tasks[task.ID().String()]
	if !exists {
		t.Fatal("expected task to be loaded from file")
	}
	if loadedTask.Description() != "Persisted Task" {
		t.Errorf("expected description 'Persisted Task', got '%s'", loadedTask.Description())
	}
}

// TestFindByID_ExistingTask tests finding a task by ID
func TestFindByID_ExistingTask(t *testing.T) {
	testPath := "./test_data/find_by_id.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create and save a task
	task, _ := domain.NewTask("Test Task", nil, 0)
	repo.Save(task)

	// Find the task
	found, err := repo.FindByID(task.ID())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID().String() != task.ID().String() {
		t.Errorf("expected task ID %s, got %s", task.ID().String(), found.ID().String())
	}
	if found.Description() != "Test Task" {
		t.Errorf("expected description 'Test Task', got '%s'", found.Description())
	}
}

// TestFindByID_NonExistentTask tests finding a non-existent task
func TestFindByID_NonExistentTask(t *testing.T) {
	testPath := "./test_data/find_by_id_not_found.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Try to find a non-existent task
	nonExistentID := domain.NewTaskID()
	_, err := repo.FindByID(nonExistentID)
	
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}

	if _, ok := err.(domain.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

// TestFindAll_EmptyRepository tests finding all tasks in empty repository
func TestFindAll_EmptyRepository(t *testing.T) {
	testPath := "./test_data/find_all_empty.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

// TestFindAll_MultipleTasks tests finding all tasks
func TestFindAll_MultipleTasks(t *testing.T) {
	testPath := "./test_data/find_all_multiple.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create and save multiple tasks
	task1, _ := domain.NewTask("Task 1", nil, 0)
	task2, _ := domain.NewTask("Task 2", nil, 1)
	task3, _ := domain.NewTask("Task 3", nil, 2)
	
	repo.Save(task1)
	repo.Save(task2)
	repo.Save(task3)

	// Find all tasks
	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

// TestFindRoot_ExistingRoot tests finding the root task
func TestFindRoot_ExistingRoot(t *testing.T) {
	testPath := "./test_data/find_root.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create root and child tasks
	root, _ := domain.NewTask("Root Task", nil, 0)
	rootID := root.ID()
	child, _ := domain.NewTask("Child Task", &rootID, 0)
	
	repo.Save(root)
	repo.Save(child)

	// Find root
	found, err := repo.FindRoot()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID().String() != root.ID().String() {
		t.Errorf("expected root ID %s, got %s", root.ID().String(), found.ID().String())
	}
	if found.ParentID() != nil {
		t.Error("expected root to have nil parent")
	}
}

// TestFindRoot_NoRoot tests finding root when none exists
func TestFindRoot_NoRoot(t *testing.T) {
	testPath := "./test_data/find_root_none.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Try to find root in empty repository
	_, err := repo.FindRoot()
	
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}

	if _, ok := err.(domain.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

// TestFindByParentID_RootChildren tests finding children of root
func TestFindByParentID_RootChildren(t *testing.T) {
	testPath := "./test_data/find_by_parent.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create root and children
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child1, _ := domain.NewTask("Child 1", &rootID, 0)
	child2, _ := domain.NewTask("Child 2", &rootID, 1)
	child3, _ := domain.NewTask("Child 3", &rootID, 2)
	
	repo.Save(root)
	repo.Save(child1)
	repo.Save(child2)
	repo.Save(child3)

	// Find children of root
	children, err := repo.FindByParentID(&rootID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(children))
	}

	// Verify ordering by position
	if children[0].Position() != 0 {
		t.Errorf("expected first child position 0, got %d", children[0].Position())
	}
	if children[1].Position() != 1 {
		t.Errorf("expected second child position 1, got %d", children[1].Position())
	}
	if children[2].Position() != 2 {
		t.Errorf("expected third child position 2, got %d", children[2].Position())
	}

	// Verify descriptions match expected order
	if children[0].Description() != "Child 1" {
		t.Errorf("expected 'Child 1', got '%s'", children[0].Description())
	}
	if children[1].Description() != "Child 2" {
		t.Errorf("expected 'Child 2', got '%s'", children[1].Description())
	}
	if children[2].Description() != "Child 3" {
		t.Errorf("expected 'Child 3', got '%s'", children[2].Description())
	}
}

// TestFindByParentID_NoChildren tests finding children when none exist
func TestFindByParentID_NoChildren(t *testing.T) {
	testPath := "./test_data/find_by_parent_none.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create a task with no children
	task, _ := domain.NewTask("Leaf Task", nil, 0)
	repo.Save(task)

	taskID := task.ID()
	children, err := repo.FindByParentID(&taskID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 0 {
		t.Errorf("expected 0 children, got %d", len(children))
	}
}

// TestFindByParentID_OrderingWithGaps tests that ordering works even with position gaps
func TestFindByParentID_OrderingWithGaps(t *testing.T) {
	testPath := "./test_data/find_by_parent_gaps.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create root and children with gaps in positions
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child1, _ := domain.NewTask("Child at 0", &rootID, 0)
	child2, _ := domain.NewTask("Child at 5", &rootID, 5)
	child3, _ := domain.NewTask("Child at 2", &rootID, 2)
	
	repo.Save(root)
	repo.Save(child1)
	repo.Save(child2)
	repo.Save(child3)

	// Find children of root
	children, err := repo.FindByParentID(&rootID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(children))
	}

	// Verify ordering by position (should be 0, 2, 5)
	if children[0].Description() != "Child at 0" {
		t.Errorf("expected 'Child at 0' first, got '%s'", children[0].Description())
	}
	if children[1].Description() != "Child at 2" {
		t.Errorf("expected 'Child at 2' second, got '%s'", children[1].Description())
	}
	if children[2].Description() != "Child at 5" {
		t.Errorf("expected 'Child at 5' third, got '%s'", children[2].Description())
	}
}

// TestDelete_ExistingTask tests deleting an existing task
func TestDelete_ExistingTask(t *testing.T) {
	testPath := "./test_data/delete_existing.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create and save a task
	task, _ := domain.NewTask("Task to Delete", nil, 0)
	repo.Save(task)

	// Verify task exists
	if len(repo.tasks) != 1 {
		t.Fatalf("expected 1 task before delete, got %d", len(repo.tasks))
	}

	// Delete the task
	err := repo.Delete(task.ID())
	if err != nil {
		t.Fatalf("expected no error deleting task, got %v", err)
	}

	// Verify task is removed from memory
	if len(repo.tasks) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(repo.tasks))
	}

	// Verify task cannot be found
	_, err = repo.FindByID(task.ID())
	if err == nil {
		t.Error("expected NotFoundError after delete, got nil")
	}
	if _, ok := err.(domain.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}

	// Verify deletion was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	if len(repo2.tasks) != 0 {
		t.Errorf("expected 0 tasks in file after delete, got %d", len(repo2.tasks))
	}
}

// TestDelete_NonExistentTask tests deleting a non-existent task
func TestDelete_NonExistentTask(t *testing.T) {
	testPath := "./test_data/delete_nonexistent.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Try to delete a non-existent task
	nonExistentID := domain.NewTaskID()
	err := repo.Delete(nonExistentID)
	
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}

	if _, ok := err.(domain.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

// TestDelete_OneOfMultipleTasks tests deleting one task when multiple exist
func TestDelete_OneOfMultipleTasks(t *testing.T) {
	testPath := "./test_data/delete_one_of_many.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create and save multiple tasks
	task1, _ := domain.NewTask("Task 1", nil, 0)
	task2, _ := domain.NewTask("Task 2", nil, 1)
	task3, _ := domain.NewTask("Task 3", nil, 2)
	
	repo.Save(task1)
	repo.Save(task2)
	repo.Save(task3)

	// Delete task2
	err := repo.Delete(task2.ID())
	if err != nil {
		t.Fatalf("expected no error deleting task, got %v", err)
	}

	// Verify only task2 is removed
	if len(repo.tasks) != 2 {
		t.Errorf("expected 2 tasks after delete, got %d", len(repo.tasks))
	}

	// Verify task1 and task3 still exist
	_, err = repo.FindByID(task1.ID())
	if err != nil {
		t.Error("expected task1 to still exist")
	}

	_, err = repo.FindByID(task3.ID())
	if err != nil {
		t.Error("expected task3 to still exist")
	}

	// Verify task2 is gone
	_, err = repo.FindByID(task2.ID())
	if err == nil {
		t.Error("expected task2 to be deleted")
	}

	// Verify deletion was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	if len(repo2.tasks) != 2 {
		t.Errorf("expected 2 tasks in file after delete, got %d", len(repo2.tasks))
	}
}

// TestDelete_DoesNotDeleteChildren tests that Delete only removes the specified task, not its children
func TestDelete_DoesNotDeleteChildren(t *testing.T) {
	testPath := "./test_data/delete_parent_keeps_children.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create parent and children
	parent, _ := domain.NewTask("Parent", nil, 0)
	parentID := parent.ID()
	child1, _ := domain.NewTask("Child 1", &parentID, 0)
	child2, _ := domain.NewTask("Child 2", &parentID, 1)
	
	repo.Save(parent)
	repo.Save(child1)
	repo.Save(child2)

	// Delete parent
	err := repo.Delete(parent.ID())
	if err != nil {
		t.Fatalf("expected no error deleting parent, got %v", err)
	}

	// Verify parent is removed but children remain
	if len(repo.tasks) != 2 {
		t.Errorf("expected 2 tasks (children) after delete, got %d", len(repo.tasks))
	}

	// Verify parent is gone
	_, err = repo.FindByID(parent.ID())
	if err == nil {
		t.Error("expected parent to be deleted")
	}

	// Verify children still exist
	_, err = repo.FindByID(child1.ID())
	if err != nil {
		t.Error("expected child1 to still exist")
	}

	_, err = repo.FindByID(child2.ID())
	if err != nil {
		t.Error("expected child2 to still exist")
	}
}

// TestDeleteSubtree_ExistingTask tests deleting a task and all its descendants
func TestDeleteSubtree_ExistingTask(t *testing.T) {
	testPath := "./test_data/delete_subtree.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create a tree: root -> child1 -> grandchild1, child2
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child1, _ := domain.NewTask("Child 1", &rootID, 0)
	child1ID := child1.ID()
	grandchild1, _ := domain.NewTask("Grandchild 1", &child1ID, 0)
	child2, _ := domain.NewTask("Child 2", &rootID, 1)
	
	if err := repo.Save(root); err != nil {
		t.Errorf("Failed to save, err: %v", err)
	}
	if err := repo.Save(child1); err != nil {
		t.Errorf("Failed to save, err: %v", err)
	}
	if err := repo.Save(grandchild1); err != nil {
		t.Errorf("Failed to save, err: %v", err)
	}
	if err := repo.Save(child2); err != nil {
		t.Errorf("Failed to save, err: %v", err)
	}

	// Verify all tasks exist
	if len(repo.tasks) != 4 {
		t.Fatalf("expected 4 tasks before delete, got %d", len(repo.tasks))
	}

	// Delete child1 subtree (should delete child1 and grandchild1)
	err := repo.DeleteSubtree(child1.ID())
	if err != nil {
		t.Fatalf("expected no error deleting subtree, got %v", err)
	}

	// Verify child1 and grandchild1 are removed, but root and child2 remain
	if len(repo.tasks) != 2 {
		t.Errorf("expected 2 tasks after delete, got %d", len(repo.tasks))
	}

	// Verify child1 is gone
	_, err = repo.FindByID(child1.ID())
	if err == nil {
		t.Error("expected child1 to be deleted")
	}

	// Verify grandchild1 is gone
	_, err = repo.FindByID(grandchild1.ID())
	if err == nil {
		t.Error("expected grandchild1 to be deleted")
	}

	// Verify root still exists
	_, err = repo.FindByID(root.ID())
	if err != nil {
		t.Error("expected root to still exist")
	}

	// Verify child2 still exists
	_, err = repo.FindByID(child2.ID())
	if err != nil {
		t.Error("expected child2 to still exist")
	}

	// Verify deletion was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	if len(repo2.tasks) != 2 {
		t.Errorf("expected 2 tasks in file after delete, got %d", len(repo2.tasks))
	}
}

// TestDeleteSubtree_NonExistentTask tests deleting a non-existent task
func TestDeleteSubtree_NonExistentTask(t *testing.T) {
	testPath := "./test_data/delete_subtree_nonexistent.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Try to delete a non-existent task
	nonExistentID := domain.NewTaskID()
	err := repo.DeleteSubtree(nonExistentID)
	
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}

	if _, ok := err.(domain.NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

// TestDeleteSubtree_LeafTask tests deleting a leaf task (no children)
func TestDeleteSubtree_LeafTask(t *testing.T) {
	testPath := "./test_data/delete_subtree_leaf.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create a simple tree: root -> child
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child, _ := domain.NewTask("Child", &rootID, 0)
	
	repo.Save(root)
	repo.Save(child)

	// Delete leaf task (child)
	err := repo.DeleteSubtree(child.ID())
	if err != nil {
		t.Fatalf("expected no error deleting leaf, got %v", err)
	}

	// Verify only child is removed
	if len(repo.tasks) != 1 {
		t.Errorf("expected 1 task after delete, got %d", len(repo.tasks))
	}

	// Verify child is gone
	_, err = repo.FindByID(child.ID())
	if err == nil {
		t.Error("expected child to be deleted")
	}

	// Verify root still exists
	_, err = repo.FindByID(root.ID())
	if err != nil {
		t.Error("expected root to still exist")
	}
}

// TestDeleteSubtree_DeepTree tests deleting a deep tree
func TestDeleteSubtree_DeepTree(t *testing.T) {
	testPath := "./test_data/delete_subtree_deep.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)
	
	// Create a deep tree: root -> child -> grandchild -> great-grandchild
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child, _ := domain.NewTask("Child", &rootID, 0)
	childID := child.ID()
	grandchild, _ := domain.NewTask("Grandchild", &childID, 0)
	grandchildID := grandchild.ID()
	greatGrandchild, _ := domain.NewTask("Great-Grandchild", &grandchildID, 0)
	
	repo.Save(root)
	repo.Save(child)
	repo.Save(grandchild)
	repo.Save(greatGrandchild)

	// Delete child subtree (should delete child, grandchild, and great-grandchild)
	err := repo.DeleteSubtree(child.ID())
	if err != nil {
		t.Fatalf("expected no error deleting subtree, got %v", err)
	}

	// Verify only root remains
	if len(repo.tasks) != 1 {
		t.Errorf("expected 1 task after delete, got %d", len(repo.tasks))
	}

	// Verify root still exists
	_, err = repo.FindByID(root.ID())
	if err != nil {
		t.Error("expected root to still exist")
	}

	// Verify all descendants are gone
	_, err = repo.FindByID(child.ID())
	if err == nil {
		t.Error("expected child to be deleted")
	}

	_, err = repo.FindByID(grandchild.ID())
	if err == nil {
		t.Error("expected grandchild to be deleted")
	}

	_, err = repo.FindByID(greatGrandchild.ID())
	if err == nil {
		t.Error("expected great-grandchild to be deleted")
	}
}

// TestConcurrentReads tests that concurrent read operations are safe
func TestConcurrentReads(t *testing.T) {
	testPath := "./test_data/concurrent_reads.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create and save some tasks
	root, _ := domain.NewTask("Root", nil, 0)
	rootID := root.ID()
	child1, _ := domain.NewTask("Child 1", &rootID, 0)
	child2, _ := domain.NewTask("Child 2", &rootID, 1)

	repo.Save(root)
	repo.Save(child1)
	repo.Save(child2)

	// Perform concurrent reads
	const numGoroutines = 10
	const numReadsPerGoroutine = 100
	done := make(chan bool, numGoroutines)
	var mu sync.Mutex
	errors := []error{}

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numReadsPerGoroutine; j++ {
				// Test FindByID
				_, err := repo.FindByID(root.ID())
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}

				// Test FindAll
				_, err = repo.FindAll()
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}

				// Test FindRoot
				_, err = repo.FindRoot()
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}

				// Test FindByParentID
				_, err = repo.FindByParentID(&rootID)
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	for _, err := range errors {
    	t.Errorf("FindByID failed: %v", err)
	}

	// Verify data integrity after concurrent reads
	tasks, _ := repo.FindAll()
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks after concurrent reads, got %d", len(tasks))
	}
}

// TestConcurrentWrites tests that concurrent write operations are safe
func TestConcurrentWrites(t *testing.T) {
	testPath := "./test_data/concurrent_writes.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create a root task
	root, _ := domain.NewTask("Root", nil, 0)
	repo.Save(root)
	rootID := root.ID()

	// Perform concurrent writes (Save operations)
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			// Each goroutine creates and saves a task
			task, _ := domain.NewTask("Task "+string(rune('A'+index)), &rootID, index)
			err := repo.Save(task)
			if err != nil {
				errCh <- err
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	close(errCh)
	if err := <-errCh; err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Verify all tasks were saved
	tasks, _ := repo.FindAll()
	expectedCount := numGoroutines + 1 // +1 for root
	if len(tasks) != expectedCount {
		t.Errorf("expected %d tasks after concurrent writes, got %d", expectedCount, len(tasks))
	}

	// Verify data was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	tasks2, _ := repo2.FindAll()
	if len(tasks2) != expectedCount {
		t.Errorf("expected %d tasks in file after concurrent writes, got %d", expectedCount, len(tasks2))
	}
}

// TestConcurrentSaves tests that concurrent save operations are safe
func TestConcurrentSaves(t *testing.T) {
	testPath := "./test_data/concurrent_saves.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create and save a root task
	root, _ := domain.NewTask("Root", nil, 0)
	repo.Save(root)
	rootID := root.ID()

	// Perform concurrent saves of different tasks
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			// Each goroutine creates and saves its own task
			task, _ := domain.NewTask("Task "+string(rune('A'+index)), &rootID, index)
			err := repo.Save(task)
			if err != nil {
				t.Errorf("Save failed: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all tasks exist
	allTasks, _ := repo.FindAll()
	expectedCount := numGoroutines + 1 // +1 for root
	if len(allTasks) != expectedCount {
		t.Errorf("expected %d tasks after concurrent saves, got %d", expectedCount, len(allTasks))
	}

	// Verify data was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	tasks2, _ := repo2.FindAll()
	if len(tasks2) != expectedCount {
		t.Errorf("expected %d tasks in file after concurrent saves, got %d", expectedCount, len(tasks2))
	}
}

// TestConcurrentDeletes tests that concurrent delete operations are safe
func TestConcurrentDeletes(t *testing.T) {
	testPath := "./test_data/concurrent_deletes.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create and save multiple tasks
	root, _ := domain.NewTask("Root", nil, 0)
	repo.Save(root)
	rootID := root.ID()

	const numTasks = 20
	taskIDs := make([]domain.TaskID, numTasks)
	for i := 0; i < numTasks; i++ {
		task, _ := domain.NewTask("Task "+string(rune('A'+i)), &rootID, i)
		repo.Save(task)
		taskIDs[i] = task.ID()
	}

	// Perform concurrent deletes
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			// Each goroutine deletes 2 tasks
			for j := 0; j < 2; j++ {
				taskIndex := index*2 + j
				if taskIndex < numTasks {
					err := repo.Delete(taskIDs[taskIndex])
					// Error is expected if another goroutine already deleted it
					if err != nil {
						if _, ok := err.(domain.NotFoundError); !ok {
							t.Errorf("unexpected error type: %T", err)
						}
					}
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all tasks were deleted (only root remains)
	tasks, _ := repo.FindAll()
	if len(tasks) != 1 {
		t.Errorf("expected 1 task (root) after concurrent deletes, got %d", len(tasks))
	}

	// Verify data was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	tasks2, _ := repo2.FindAll()
	if len(tasks2) != 1 {
		t.Errorf("expected 1 task in file after concurrent deletes, got %d", len(tasks2))
	}
}

// TestConcurrentMixedOperations tests concurrent reads and writes together
func TestConcurrentMixedOperations(t *testing.T) {
	testPath := "./test_data/concurrent_mixed.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create initial tasks
	root, _ := domain.NewTask("Root", nil, 0)
	repo.Save(root)
	rootID := root.ID()

	// Perform concurrent mixed operations
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	// 5 goroutines for reads
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				repo.FindAll()
				repo.FindRoot()
				repo.FindByParentID(&rootID)
			}
			done <- true
		}()
	}

	// 5 goroutines for writes (each creates different tasks)
	for i := 0; i < 5; i++ {
		go func(index int) {
			for j := 0; j < 10; j++ {
				task, _ := domain.NewTask("Task "+string(rune('A'+index))+"-"+string(rune('0'+j)), &rootID, index*10+j)
				repo.Save(task)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify data integrity
	tasks, _ := repo.FindAll()
	if len(tasks) < 1 {
		t.Error("expected at least root task after concurrent operations")
	}

	// Verify root still exists
	_, err := repo.FindByID(root.ID())
	if err != nil {
		t.Error("expected root to still exist")
	}

	// Verify data was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	tasks2, _ := repo2.FindAll()
	if len(tasks2) != len(tasks) {
		t.Errorf("expected %d tasks in file, got %d", len(tasks), len(tasks2))
	}
}

// TestConcurrentDeleteSubtree tests that concurrent DeleteSubtree operations are safe
func TestConcurrentDeleteSubtree(t *testing.T) {
	testPath := "./test_data/concurrent_delete_subtree.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	repo, _ := NewFileTaskRepository(testPath)

	// Create a tree structure
	root, _ := domain.NewTask("Root", nil, 0)
	repo.Save(root)
	rootID := root.ID()

	// Create multiple subtrees
	const numSubtrees = 5
	subtreeRoots := make([]domain.TaskID, numSubtrees)
	for i := 0; i < numSubtrees; i++ {
		subtreeRoot, _ := domain.NewTask("Subtree "+string(rune('A'+i)), &rootID, i)
		repo.Save(subtreeRoot)
		subtreeRoots[i] = subtreeRoot.ID()

		// Add children to each subtree
		subtreeRootID := subtreeRoot.ID()
		for j := 0; j < 3; j++ {
			child, _ := domain.NewTask("Child "+string(rune('A'+i))+"-"+string(rune('0'+j)), &subtreeRootID, j)
			repo.Save(child)
		}
	}

	// Perform concurrent DeleteSubtree operations
	done := make(chan bool, numSubtrees)
	errCh := make(chan error, numSubtrees)

	for i := 0; i < numSubtrees; i++ {
		go func(index int) {
			err := repo.DeleteSubtree(subtreeRoots[index])
			// Error is expected if another goroutine already deleted it
			if err != nil {
				if _, ok := err.(domain.NotFoundError); !ok {
					errCh <- err
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numSubtrees; i++ {
		<-done
	}
	close(errCh)
	if err := <-errCh; err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify only root remains
	tasks, _ := repo.FindAll()
	if len(tasks) != 1 {
		t.Errorf("expected 1 task (root) after concurrent DeleteSubtree, got %d", len(tasks))
	}

	// Verify root still exists
	_, err := repo.FindByID(root.ID())
	if err != nil {
		t.Error("expected root to still exist")
	}

	// Verify data was persisted to file
	repo2, _ := NewFileTaskRepository(testPath)
	tasks2, _ := repo2.FindAll()
	if len(tasks2) != 1 {
		t.Errorf("expected 1 task in file after concurrent DeleteSubtree, got %d", len(tasks2))
	}
}
