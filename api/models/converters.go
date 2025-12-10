package models

import (
	"discovery-tree/domain"
)

// TaskToResponse converts a domain Task to a TaskResponse
func TaskToResponse(task *domain.Task) TaskResponse {
	var parentID *string
	if task.ParentID() != nil {
		parentIDStr := task.ParentID().String()
		parentID = &parentIDStr
	}

	return TaskResponse{
		ID:          task.ID().String(),
		Description: task.Description(),
		Status:      task.Status().String(),
		ParentID:    parentID,
		Position:    task.Position(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}
}

// ErrorToResponse converts a domain error to an ErrorResponse
func ErrorToResponse(err error) ErrorResponse {
	switch e := err.(type) {
	case domain.ValidationError:
		return ErrorResponse{
			Error:   "ValidationError",
			Code:    e.Field,
			Message: e.Message,
		}
	case domain.NotFoundError:
		return ErrorResponse{
			Error:   "NotFoundError",
			Code:    "RESOURCE_NOT_FOUND",
			Message: e.Error(),
		}
	case domain.ConstraintViolationError:
		return ErrorResponse{
			Error:   "ConstraintViolationError",
			Code:    e.Constraint,
			Message: e.Message,
		}
	default:
		return ErrorResponse{
			Error:   "InternalServerError",
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred",
		}
	}
}