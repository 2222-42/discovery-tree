package domain

import (
	"github.com/google/uuid"
)

// TaskID represents a unique identifier for a task
type TaskID struct {
	value string
}

// NewTaskID generates a new unique TaskID
func NewTaskID() TaskID {
	return TaskID{
		value: uuid.New().String(),
	}
}

// TaskIDFromString creates a TaskID from a string value with validation
func TaskIDFromString(s string) (TaskID, error) {
	if s == "" {
		return TaskID{}, NewValidationError("taskID", "task ID cannot be empty")
	}
	
	// Validate that the string is a valid UUID format
	if _, err := uuid.Parse(s); err != nil {
		return TaskID{}, NewValidationError("taskID", "task ID must be a valid UUID")
	}
	
	return TaskID{value: s}, nil
}

// String returns the string representation of the TaskID
func (t TaskID) String() string {
	return t.value
}

// Equals checks if two TaskIDs are equal
func (t TaskID) Equals(other TaskID) bool {
	return t.value == other.value
}
