package infrastructure

import (
	"encoding/json"
	"os"
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
