package infrastructure

import (
	"discovery-tree/domain"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestTaskDTO_JSONSerialization(t *testing.T) {
	// Create a task
	task, err := domain.NewTask("Test task", nil, 0)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Convert to DTO
	dto := ToDTO(task)

	// Marshal to JSON
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal DTO to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify JSON contains expected fields
	if !strings.Contains(jsonStr, `"id"`) {
		t.Error("JSON should contain 'id' field")
	}
	if !strings.Contains(jsonStr, `"description"`) {
		t.Error("JSON should contain 'description' field")
	}
	if !strings.Contains(jsonStr, `"status"`) {
		t.Error("JSON should contain 'status' field")
	}
	if !strings.Contains(jsonStr, `"parentId"`) {
		t.Error("JSON should contain 'parentId' field")
	}
	if !strings.Contains(jsonStr, `"position"`) {
		t.Error("JSON should contain 'position' field")
	}
	if !strings.Contains(jsonStr, `"createdAt"`) {
		t.Error("JSON should contain 'createdAt' field")
	}
	if !strings.Contains(jsonStr, `"updatedAt"`) {
		t.Error("JSON should contain 'updatedAt' field")
	}

	// Verify null parent for root task
	if !strings.Contains(jsonStr, `"parentId":null`) {
		t.Error("JSON should contain 'parentId':null for root task")
	}

	// Unmarshal back
	var unmarshaledDTO TaskDTO
	err = json.Unmarshal(jsonBytes, &unmarshaledDTO)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify fields match
	if unmarshaledDTO.ID != dto.ID {
		t.Errorf("ID mismatch after JSON round trip: expected %s, got %s", dto.ID, unmarshaledDTO.ID)
	}
	if unmarshaledDTO.Description != dto.Description {
		t.Errorf("Description mismatch after JSON round trip")
	}
	if unmarshaledDTO.ParentID != nil {
		t.Error("ParentID should be nil after JSON round trip for root task")
	}
}

func TestTaskDTO_ISO8601Timestamps(t *testing.T) {
	// Create a DTO with specific timestamp
	createdAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 15, 11, 45, 30, 0, time.UTC)

	dto := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test task",
		Status:      "TODO",
		ParentID:    nil,
		Position:    0,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal DTO to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify ISO 8601 format (RFC3339 is Go's implementation of ISO 8601)
	// Should contain timestamps in format: "2024-01-15T10:30:00Z"
	if !strings.Contains(jsonStr, `"createdAt":"2024-01-15T10:30:00Z"`) {
		t.Errorf("CreatedAt should be in ISO 8601 format, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"updatedAt":"2024-01-15T11:45:30Z"`) {
		t.Errorf("UpdatedAt should be in ISO 8601 format, got: %s", jsonStr)
	}

	// Unmarshal and verify timestamps are preserved
	var unmarshaledDTO TaskDTO
	err = json.Unmarshal(jsonBytes, &unmarshaledDTO)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !unmarshaledDTO.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt mismatch: expected %v, got %v", createdAt, unmarshaledDTO.CreatedAt)
	}
	if !unmarshaledDTO.UpdatedAt.Equal(updatedAt) {
		t.Errorf("UpdatedAt mismatch: expected %v, got %v", updatedAt, unmarshaledDTO.UpdatedAt)
	}
}

func TestTaskDTO_WithParentJSON(t *testing.T) {
	// Create a DTO with parent
	parentIDStr := "550e8400-e29b-41d4-a716-446655440000"
	dto := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Description: "Child task",
		Status:      "TODO",
		ParentID:    &parentIDStr,
		Position:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal DTO to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify parent ID is a string, not null
	if strings.Contains(jsonStr, `"parentId":null`) {
		t.Error("ParentID should not be null for child task")
	}
	if !strings.Contains(jsonStr, `"parentId":"550e8400-e29b-41d4-a716-446655440000"`) {
		t.Errorf("ParentID should be a string in JSON, got: %s", jsonStr)
	}

	// Unmarshal and verify
	var unmarshaledDTO TaskDTO
	err = json.Unmarshal(jsonBytes, &unmarshaledDTO)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if unmarshaledDTO.ParentID == nil {
		t.Error("ParentID should not be nil after unmarshaling")
	} else if *unmarshaledDTO.ParentID != parentIDStr {
		t.Errorf("ParentID mismatch: expected %s, got %s", parentIDStr, *unmarshaledDTO.ParentID)
	}
}
