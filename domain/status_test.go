package domain

import "testing"

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{"TODO status", StatusTODO, "TODO"},
		{"InProgress status", StatusInProgress, "In Progress"},
		{"DONE status", StatusDONE, "DONE"},
		{"Blocked status", StatusBlocked, "Blocked"},
		{"RootWorkItem status", StatusRootWorkItem, "Root Work Item"},
		{"Invalid status", Status(-1), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Status.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"TODO is valid", StatusTODO, true},
		{"InProgress is valid", StatusInProgress, true},
		{"DONE is valid", StatusDONE, true},
		{"Blocked is valid", StatusBlocked, true},
		{"RootWorkItem is valid", StatusRootWorkItem, true},
		{"Negative value is invalid", Status(-1), false},
		{"Out of range value is invalid", Status(100), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			if result != tt.expected {
				t.Errorf("Status.IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Status
		expectError bool
	}{
		{"Valid TODO", "TODO", StatusTODO, false},
		{"Valid In Progress", "In Progress", StatusInProgress, false},
		{"Valid DONE", "DONE", StatusDONE, false},
		{"Valid Blocked", "Blocked", StatusBlocked, false},
		{"Valid Root Work Item", "Root Work Item", StatusRootWorkItem, false},
		{"Invalid status", "InvalidStatus", Status(-1), true},
		{"Empty string", "", Status(-1), true},
		{"Case sensitive", "todo", Status(-1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewStatus(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("NewStatus(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("NewStatus(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("NewStatus(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}
