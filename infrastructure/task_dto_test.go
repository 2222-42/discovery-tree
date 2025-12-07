package infrastructure

import (
	"discovery-tree/domain"
	"testing"
	"time"
)

func TestToDTO_WithParent(t *testing.T) {
	// Create a parent task
	parentTask, err := domain.NewTask("Parent task", nil, 0)
	if err != nil {
		t.Fatalf("Failed to create parent task: %v", err)
	}
	parentID := parentTask.ID()

	// Create a child task
	childTask, err := domain.NewTask("Child task", &parentID, 0)
	if err != nil {
		t.Fatalf("Failed to create child task: %v", err)
	}

	// Convert to DTO
	dto := ToDTO(childTask)

	// Verify all fields
	if dto.ID != childTask.ID().String() {
		t.Errorf("Expected ID %s, got %s", childTask.ID().String(), dto.ID)
	}
	if dto.Description != childTask.Description() {
		t.Errorf("Expected description %s, got %s", childTask.Description(), dto.Description)
	}
	if dto.Status != childTask.Status().String() {
		t.Errorf("Expected status %s, got %s", childTask.Status().String(), dto.Status)
	}
	if dto.ParentID == nil {
		t.Error("Expected non-nil ParentID")
	} else if *dto.ParentID != parentID.String() {
		t.Errorf("Expected ParentID %s, got %s", parentID.String(), *dto.ParentID)
	}
	if dto.Position != childTask.Position() {
		t.Errorf("Expected position %d, got %d", childTask.Position(), dto.Position)
	}
	if !dto.CreatedAt.Equal(childTask.CreatedAt()) {
		t.Errorf("Expected createdAt %v, got %v", childTask.CreatedAt(), dto.CreatedAt)
	}
	if !dto.UpdatedAt.Equal(childTask.UpdatedAt()) {
		t.Errorf("Expected updatedAt %v, got %v", childTask.UpdatedAt(), dto.UpdatedAt)
	}
}

func TestToDTO_WithoutParent(t *testing.T) {
	// Create a root task
	rootTask, err := domain.NewTask("Root task", nil, 0)
	if err != nil {
		t.Fatalf("Failed to create root task: %v", err)
	}

	// Convert to DTO
	dto := ToDTO(rootTask)

	// Verify ParentID is nil
	if dto.ParentID != nil {
		t.Errorf("Expected nil ParentID for root task, got %v", *dto.ParentID)
	}

	// Verify other fields
	if dto.ID != rootTask.ID().String() {
		t.Errorf("Expected ID %s, got %s", rootTask.ID().String(), dto.ID)
	}
	if dto.Description != rootTask.Description() {
		t.Errorf("Expected description %s, got %s", rootTask.Description(), dto.Description)
	}
}

func TestFromDTO_WithParent(t *testing.T) {
	// Create a DTO with parent
	parentIDStr := "550e8400-e29b-41d4-a716-446655440000"
	dto := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Description: "Test task",
		Status:      "TODO",
		ParentID:    &parentIDStr,
		Position:    1,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	// Convert from DTO
	task, err := FromDTO(dto)
	if err != nil {
		t.Fatalf("Failed to convert from DTO: %v", err)
	}

	// Verify all fields
	if task.ID().String() != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, task.ID().String())
	}
	if task.Description() != dto.Description {
		t.Errorf("Expected description %s, got %s", dto.Description, task.Description())
	}
	if task.Status().String() != dto.Status {
		t.Errorf("Expected status %s, got %s", dto.Status, task.Status().String())
	}
	if task.ParentID() == nil {
		t.Error("Expected non-nil ParentID")
	} else if task.ParentID().String() != *dto.ParentID {
		t.Errorf("Expected ParentID %s, got %s", *dto.ParentID, task.ParentID().String())
	}
	if task.Position() != dto.Position {
		t.Errorf("Expected position %d, got %d", dto.Position, task.Position())
	}
	if !task.CreatedAt().Equal(dto.CreatedAt) {
		t.Errorf("Expected createdAt %v, got %v", dto.CreatedAt, task.CreatedAt())
	}
	if !task.UpdatedAt().Equal(dto.UpdatedAt) {
		t.Errorf("Expected updatedAt %v, got %v", dto.UpdatedAt, task.UpdatedAt())
	}
}

func TestFromDTO_WithoutParent(t *testing.T) {
	// Create a DTO without parent (root task)
	dto := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Root task",
		Status:      "Root Work Item",
		ParentID:    nil,
		Position:    0,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	// Convert from DTO
	task, err := FromDTO(dto)
	if err != nil {
		t.Fatalf("Failed to convert from DTO: %v", err)
	}

	// Verify ParentID is nil
	if task.ParentID() != nil {
		t.Errorf("Expected nil ParentID for root task, got %v", task.ParentID().String())
	}

	// Verify other fields
	if task.ID().String() != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, task.ID().String())
	}
	if task.Description() != dto.Description {
		t.Errorf("Expected description %s, got %s", dto.Description, task.Description())
	}
}

func TestFromDTO_RoundTrip(t *testing.T) {
	// Create a task
	parentTask, _ := domain.NewTask("Parent", nil, 0)
	parentID := parentTask.ID()
	originalTask, err := domain.NewTask("Original task", &parentID, 2)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Convert to DTO and back
	dto := ToDTO(originalTask)
	reconstructedTask, err := FromDTO(dto)
	if err != nil {
		t.Fatalf("Failed to reconstruct task: %v", err)
	}

	// Verify all fields match
	if reconstructedTask.ID().String() != originalTask.ID().String() {
		t.Errorf("ID mismatch: expected %s, got %s", originalTask.ID().String(), reconstructedTask.ID().String())
	}
	if reconstructedTask.Description() != originalTask.Description() {
		t.Errorf("Description mismatch: expected %s, got %s", originalTask.Description(), reconstructedTask.Description())
	}
	if reconstructedTask.Status().String() != originalTask.Status().String() {
		t.Errorf("Status mismatch: expected %s, got %s", originalTask.Status().String(), reconstructedTask.Status().String())
	}
	if reconstructedTask.Position() != originalTask.Position() {
		t.Errorf("Position mismatch: expected %d, got %d", originalTask.Position(), reconstructedTask.Position())
	}
	if !reconstructedTask.CreatedAt().Equal(originalTask.CreatedAt()) {
		t.Errorf("CreatedAt mismatch: expected %v, got %v", originalTask.CreatedAt(), reconstructedTask.CreatedAt())
	}
	if !reconstructedTask.UpdatedAt().Equal(originalTask.UpdatedAt()) {
		t.Errorf("UpdatedAt mismatch: expected %v, got %v", originalTask.UpdatedAt(), reconstructedTask.UpdatedAt())
	}
}

func TestFromDTO_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		dto  TaskDTO
	}{
		{
			name: "empty ID",
			dto: TaskDTO{
				ID:          "",
				Description: "Test",
				Status:      "TODO",
				Position:    0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "empty description",
			dto: TaskDTO{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				Description: "",
				Status:      "TODO",
				Position:    0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "negative position",
			dto: TaskDTO{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				Description: "Test",
				Status:      "TODO",
				Position:    -1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "invalid status",
			dto: TaskDTO{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				Description: "Test",
				Status:      "INVALID_STATUS",
				Position:    0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "invalid ID format",
			dto: TaskDTO{
				ID:          "not-a-uuid",
				Description: "Test",
				Status:      "TODO",
				Position:    0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromDTO(tt.dto)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}
