package domain

import "fmt"

// ValidationError represents an error when input validation fails
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NotFoundError represents an error when a requested entity is not found
type NotFoundError struct {
	EntityType string
	ID         string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.EntityType, e.ID)
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(entityType, id string) NotFoundError {
	return NotFoundError{
		EntityType: entityType,
		ID:         id,
	}
}

// ConstraintViolationError represents an error when a business rule is violated
type ConstraintViolationError struct {
	Constraint string
	Message    string
}

func (e ConstraintViolationError) Error() string {
	return fmt.Sprintf("constraint violation '%s': %s", e.Constraint, e.Message)
}

// NewConstraintViolationError creates a new ConstraintViolationError
func NewConstraintViolationError(constraint, message string) ConstraintViolationError {
	return ConstraintViolationError{
		Constraint: constraint,
		Message:    message,
	}
}
