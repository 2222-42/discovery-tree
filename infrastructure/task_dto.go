package infrastructure

import (
	"time"

	"discovery-tree/domain"
)

// TaskDTO is a data transfer object for JSON serialization of Task
type TaskDTO struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ParentID    *string    `json:"parentId"` // pointer to handle null
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ToDTO converts a domain Task to a TaskDTO for JSON serialization
func ToDTO(task *domain.Task) TaskDTO {
	dto := TaskDTO{
		ID:          task.ID().String(),
		Description: task.Description(),
		Status:      task.Status().String(),
		Position:    task.Position(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}

	// Handle nil parent ID conversion (nil -> null in JSON)
	if task.ParentID() != nil {
		parentIDStr := task.ParentID().String()
		dto.ParentID = &parentIDStr
	}

	return dto
}

// FromDTO converts a TaskDTO to a domain Task
// This function reconstructs a Task from persisted data
func FromDTO(dto TaskDTO) (*domain.Task, error) {
	// Validate required fields
	if dto.ID == "" {
		return nil, domain.NewValidationError("id", "task ID cannot be empty")
	}
	if dto.Description == "" {
		return nil, domain.NewValidationError("description", "description cannot be empty")
	}
	if dto.Position < 0 {
		return nil, domain.NewValidationError("position", "position must be non-negative")
	}

	// Parse TaskID
	taskID, err := domain.TaskIDFromString(dto.ID)
	if err != nil {
		return nil, err
	}

	// Parse Status
	status, err := domain.NewStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	// Parse ParentID (handle null -> nil conversion)
	var parentID *domain.TaskID
	if dto.ParentID != nil {
		pid, err := domain.TaskIDFromString(*dto.ParentID)
		if err != nil {
			return nil, err
		}
		parentID = &pid
	}

	// Reconstruct the task using reflection-like approach
	// Since Task fields are private, we need to create it and then set fields
	// For now, we'll use a helper function that creates a task with all fields
	task := reconstructTask(
		taskID,
		dto.Description,
		status,
		parentID,
		dto.Position,
		dto.CreatedAt,
		dto.UpdatedAt,
	)

	return task, nil
}

// reconstructTask creates a Task instance with all fields set
// This is used for deserializing tasks from persistent storage
// It delegates to the domain package's reconstruction function
func reconstructTask(
	id domain.TaskID,
	description string,
	status domain.Status,
	parentID *domain.TaskID,
	position int,
	createdAt time.Time,
	updatedAt time.Time,
) *domain.Task {
	return domain.ReconstructTask(id, description, status, parentID, position, createdAt, updatedAt)
}
