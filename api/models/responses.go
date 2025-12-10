package models

import (
	"time"
)

// TaskResponse represents the API response for a task
type TaskResponse struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ParentID    *string    `json:"parentId"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ErrorResponse represents the API response for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}