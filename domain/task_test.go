package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewTask_RootTask(t *testing.T) {
	description := "Root task description"
	
	task, err := NewTask(description, nil, 0)
	
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
	
	if task.IsRoot() != true {
		t.Error("expected IsRoot() to return true for root task")
	}
	
	// Check that ID is generated
	if task.ID().String() == "" {
		t.Error("expected non-empty task ID")
	}
	
	// Check timestamps are set
	if task.CreatedAt().IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	
	if task.UpdatedAt().IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewTask_ChildTask(t *testing.T) {
	description := "Child task description"
	parentID := NewTaskID()
	position := 2
	
	task, err := NewTask(description, &parentID, position)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if task == nil {
		t.Fatal("expected task to be created, got nil")
	}
	
	if task.Description() != description {
		t.Errorf("expected description %q, got %q", description, task.Description())
	}
	
	if task.Status() != StatusTODO {
		t.Errorf("expected status %v, got %v", StatusTODO, task.Status())
	}
	
	if task.ParentID() == nil {
		t.Error("expected non-nil parent ID for child task")
	} else if !task.ParentID().Equals(parentID) {
		t.Errorf("expected parent ID %v, got %v", parentID, *task.ParentID())
	}
	
	if task.Position() != position {
		t.Errorf("expected position %d, got %d", position, task.Position())
	}
	
	if task.IsRoot() != false {
		t.Error("expected IsRoot() to return false for child task")
	}
}

func TestNewTask_EmptyDescription(t *testing.T) {
	testCases := []struct {
		name        string
		description string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"tabs only", "\t\t"},
		{"newlines only", "\n\n"},
		{"mixed whitespace", " \t\n "},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.description, nil, 0)
			
			if err == nil {
				t.Error("expected error for empty/whitespace description, got nil")
			}
			
			if task != nil {
				t.Errorf("expected nil task for invalid description, got %v", task)
			}
			
			// Check that it's a ValidationError
			if _, ok := err.(ValidationError); !ok {
				t.Errorf("expected ValidationError, got %T", err)
			}
			
			// Check error message contains "description"
			if !strings.Contains(err.Error(), "description") {
				t.Errorf("expected error message to mention 'description', got %q", err.Error())
			}
		})
	}
}

func TestNewTask_NegativePosition(t *testing.T) {
	description := "Valid description"
	
	task, err := NewTask(description, nil, -1)
	
	if err == nil {
		t.Error("expected error for negative position, got nil")
	}
	
	if task != nil {
		t.Errorf("expected nil task for invalid position, got %v", task)
	}
	
	// Check that it's a ValidationError
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	
	// Check error message contains "position"
	if !strings.Contains(err.Error(), "position") {
		t.Errorf("expected error message to mention 'position', got %q", err.Error())
	}
}

func TestTask_UniqueIDs(t *testing.T) {
	task1, err := NewTask("Task 1", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task1: %v", err)
	}
	
	task2, err := NewTask("Task 2", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task2: %v", err)
	}
	
	if task1.ID().Equals(task2.ID()) {
		t.Error("expected different tasks to have unique IDs")
	}
}

func TestTask_Timestamps(t *testing.T) {
	before := time.Now()
	time.Sleep(1 * time.Millisecond) // Ensure some time passes
	
	task, err := NewTask("Test task", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	time.Sleep(1 * time.Millisecond)
	after := time.Now()
	
	// Check CreatedAt is between before and after
	if task.CreatedAt().Before(before) || task.CreatedAt().After(after) {
		t.Errorf("expected CreatedAt to be between %v and %v, got %v", before, after, task.CreatedAt())
	}
	
	// Check UpdatedAt is between before and after
	if task.UpdatedAt().Before(before) || task.UpdatedAt().After(after) {
		t.Errorf("expected UpdatedAt to be between %v and %v, got %v", before, after, task.UpdatedAt())
	}
	
	// For a new task, CreatedAt and UpdatedAt should be very close (within milliseconds)
	timeDiff := task.UpdatedAt().Sub(task.CreatedAt())
	if timeDiff < 0 || timeDiff > 10*time.Millisecond {
		t.Errorf("expected CreatedAt and UpdatedAt to be very close for new task, difference: %v", timeDiff)
	}
}

func TestTask_GetterMethods(t *testing.T) {
	description := "Test description"
	parentID := NewTaskID()
	position := 5
	
	task, err := NewTask(description, &parentID, position)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	// Test all getter methods
	if task.Description() != description {
		t.Errorf("Description() = %q, want %q", task.Description(), description)
	}
	
	if task.Status() != StatusTODO {
		t.Errorf("Status() = %v, want %v", task.Status(), StatusTODO)
	}
	
	if task.ParentID() == nil || !task.ParentID().Equals(parentID) {
		t.Errorf("ParentID() = %v, want %v", task.ParentID(), parentID)
	}
	
	if task.Position() != position {
		t.Errorf("Position() = %d, want %d", task.Position(), position)
	}
	
	if task.ID().String() == "" {
		t.Error("ID() returned empty string")
	}
	
	if task.CreatedAt().IsZero() {
		t.Error("CreatedAt() returned zero time")
	}
	
	if task.UpdatedAt().IsZero() {
		t.Error("UpdatedAt() returned zero time")
	}
	
	if task.IsRoot() {
		t.Error("IsRoot() = true, want false for child task")
	}
}

func TestTask_ChangeStatus_ValidStatus(t *testing.T) {
	task, err := NewTask("Test task", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	initialStatus := task.Status()
	initialUpdatedAt := task.UpdatedAt()
	
	// Wait a bit to ensure timestamp changes
	time.Sleep(2 * time.Millisecond)
	
	// Change to a valid status
	err = task.ChangeStatus(StatusInProgress)
	if err != nil {
		t.Errorf("expected no error when changing to valid status, got %v", err)
	}
	
	if task.Status() != StatusInProgress {
		t.Errorf("expected status %v, got %v", StatusInProgress, task.Status())
	}
	
	// Verify the status actually changed from the initial value
	if task.Status() == initialStatus {
		t.Error("expected status to change from initial value")
	}
	
	// Check that UpdatedAt was updated
	if !task.UpdatedAt().After(initialUpdatedAt) {
		t.Error("expected UpdatedAt to be updated after status change")
	}
	
	// Verify we can change to other valid statuses
	validStatuses := []Status{StatusTODO, StatusDONE, StatusBlocked, StatusRootWorkItem}
	for _, status := range validStatuses {
		err = task.ChangeStatus(status)
		if err != nil {
			t.Errorf("expected no error when changing to %v, got %v", status, err)
		}
		if task.Status() != status {
			t.Errorf("expected status %v, got %v", status, task.Status())
		}
	}
}

func TestTask_ChangeStatus_InvalidStatus(t *testing.T) {
	task, err := NewTask("Test task", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	initialStatus := task.Status()
	
	// Try to change to an invalid status
	invalidStatus := Status(-1)
	err = task.ChangeStatus(invalidStatus)
	
	if err == nil {
		t.Error("expected error when changing to invalid status, got nil")
	}
	
	// Check that it's a ValidationError
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	
	// Check error message contains "status"
	if !strings.Contains(err.Error(), "status") {
		t.Errorf("expected error message to mention 'status', got %q", err.Error())
	}
	
	// Verify status didn't change
	if task.Status() != initialStatus {
		t.Errorf("expected status to remain %v after failed change, got %v", initialStatus, task.Status())
	}
	
	// Try another invalid status
	invalidStatus2 := Status(999)
	err = task.ChangeStatus(invalidStatus2)
	if err == nil {
		t.Error("expected error when changing to invalid status 999, got nil")
	}
}

func TestTask_UpdateDescription_ValidDescription(t *testing.T) {
	task, err := NewTask("Initial description", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	initialUpdatedAt := task.UpdatedAt()
	
	// Wait a bit to ensure timestamp changes
	time.Sleep(2 * time.Millisecond)
	
	newDescription := "Updated description"
	err = task.UpdateDescription(newDescription)
	
	if err != nil {
		t.Errorf("expected no error when updating to valid description, got %v", err)
	}
	
	if task.Description() != newDescription {
		t.Errorf("expected description %q, got %q", newDescription, task.Description())
	}
	
	// Check that UpdatedAt was updated
	if !task.UpdatedAt().After(initialUpdatedAt) {
		t.Error("expected UpdatedAt to be updated after description change")
	}
}

func TestTask_UpdateDescription_EmptyDescription(t *testing.T) {
	task, err := NewTask("Initial description", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	initialDescription := task.Description()
	
	testCases := []struct {
		name        string
		description string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
		{"tabs only", "\t\t"},
		{"newlines only", "\n\n"},
		{"mixed whitespace", " \t\n "},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := task.UpdateDescription(tc.description)
			
			if err == nil {
				t.Error("expected error for empty/whitespace description, got nil")
			}
			
			// Check that it's a ValidationError
			if _, ok := err.(ValidationError); !ok {
				t.Errorf("expected ValidationError, got %T", err)
			}
			
			// Check error message contains "description"
			if !strings.Contains(err.Error(), "description") {
				t.Errorf("expected error message to mention 'description', got %q", err.Error())
			}
			
			// Verify description didn't change
			if task.Description() != initialDescription {
				t.Errorf("expected description to remain %q after failed update, got %q", initialDescription, task.Description())
			}
		})
	}
}

func TestTask_UpdateDescription_PreservesWhitespace(t *testing.T) {
	task, err := NewTask("Initial description", nil, 0)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	
	// Description with leading/trailing whitespace but non-empty content
	newDescription := "  Description with spaces  "
	err = task.UpdateDescription(newDescription)
	
	if err != nil {
		t.Errorf("expected no error for description with whitespace, got %v", err)
	}
	
	// The description should be preserved as-is (not trimmed)
	if task.Description() != newDescription {
		t.Errorf("expected description %q, got %q", newDescription, task.Description())
	}
}
