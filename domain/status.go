package domain

import "fmt"

// Status represents the current state of a task
type Status int

const (
	StatusTODO Status = iota
	StatusInProgress
	StatusDONE
	StatusBlocked
	StatusRootWorkItem
)

// String returns the string representation of the Status
func (s Status) String() string {
	switch s {
	case StatusTODO:
		return "TODO"
	case StatusInProgress:
		return "In Progress"
	case StatusDONE:
		return "DONE"
	case StatusBlocked:
		return "Blocked"
	case StatusRootWorkItem:
		return "Root Work Item"
	default:
		return "Unknown"
	}
}

// IsValid checks if the status value is valid
func (s Status) IsValid() bool {
	return s >= StatusTODO && s <= StatusRootWorkItem
}

// NewStatus creates a Status from a string value with validation
func NewStatus(s string) (Status, error) {
	switch s {
	case "TODO":
		return StatusTODO, nil
	case "In Progress":
		return StatusInProgress, nil
	case "DONE":
		return StatusDONE, nil
	case "Blocked":
		return StatusBlocked, nil
	case "Root Work Item":
		return StatusRootWorkItem, nil
	default:
		return Status(-1), NewValidationError("status", fmt.Sprintf("invalid status value: %s", s))
	}
}
