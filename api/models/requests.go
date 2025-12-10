package models

// CreateRootTaskRequest represents the request to create a root task
type CreateRootTaskRequest struct {
	Description string `json:"description" binding:"required,min=1"`
}

// CreateChildTaskRequest represents the request to create a child task
type CreateChildTaskRequest struct {
	Description string `json:"description" binding:"required,min=1"`
	ParentID    string `json:"parentId" binding:"required,uuid"`
}

// UpdateTaskRequest represents the request to update a task's description
type UpdateTaskRequest struct {
	Description string `json:"description" binding:"required,min=1"`
}

// UpdateStatusRequest represents the request to update a task's status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof='TODO' 'In Progress' 'DONE' 'Blocked' 'Root Work Item'"`
}

// MoveTaskRequest represents the request to move a task to a new position or parent
type MoveTaskRequest struct {
	ParentID *string `json:"parentId" binding:"omitempty,uuid"`
	Position int     `json:"position" binding:"min=0"`
}