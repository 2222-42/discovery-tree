package domain

import (
	"testing"
)

func TestNewTaskID(t *testing.T) {
	// Test that NewTaskID generates a valid ID
	id := NewTaskID()
	if id.String() == "" {
		t.Error("NewTaskID should generate a non-empty ID")
	}
	
	// Test that multiple calls generate unique IDs
	id2 := NewTaskID()
	if id.Equals(id2) {
		t.Error("NewTaskID should generate unique IDs")
	}
}

func TestTaskIDFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid UUID format",
			input:   "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "random string",
			input:   "abc123",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := TaskIDFromString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("TaskIDFromString(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("TaskIDFromString(%q) unexpected error: %v", tt.input, err)
				}
				if id.String() != tt.input {
					t.Errorf("TaskIDFromString(%q).String() = %q, want %q", tt.input, id.String(), tt.input)
				}
			}
		})
	}
}

func TestTaskIDString(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id, err := TaskIDFromString(validUUID)
	if err != nil {
		t.Fatalf("TaskIDFromString failed: %v", err)
	}
	
	if id.String() != validUUID {
		t.Errorf("String() = %q, want %q", id.String(), validUUID)
	}
}

func TestTaskIDEquals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	id1, _ := TaskIDFromString(validUUID)
	id2, _ := TaskIDFromString(validUUID)
	id3, _ := TaskIDFromString("650e8400-e29b-41d4-a716-446655440000")
	
	if !id1.Equals(id2) {
		t.Error("TaskIDs with same value should be equal")
	}
	
	if id1.Equals(id3) {
		t.Error("TaskIDs with different values should not be equal")
	}
}
