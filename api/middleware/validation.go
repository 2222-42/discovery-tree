package middleware

import (
	"discovery-tree/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ValidationErrorHandler middleware handles validation errors consistently
func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Handle validation errors specifically
			if validationErr, ok := err.Err.(validator.ValidationErrors); ok {
				errorResp := models.ErrorResponse{
					Error:   "ValidationError",
					Code:    "INVALID_REQUEST",
					Message: formatValidationError(validationErr),
				}
				c.JSON(http.StatusBadRequest, errorResp)
				c.Abort()
				return
			}

			// Handle binding errors
			errorResp := models.ErrorResponse{
				Error:   "ValidationError",
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			}
			c.JSON(http.StatusBadRequest, errorResp)
			c.Abort()
		}
	}
}

// BindJSON is a helper function that binds JSON and handles validation errors consistently
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		errorResp := models.ErrorResponse{
			Error:   "ValidationError",
			Code:    "INVALID_REQUEST",
			Message: formatBindingError(err),
		}
		c.JSON(http.StatusBadRequest, errorResp)
		return err
	}
	return nil
}

// formatValidationError formats validator.ValidationErrors into a user-friendly message
func formatValidationError(validationErrors validator.ValidationErrors) string {
	if len(validationErrors) == 0 {
		return "Invalid request data"
	}

	// Return the first validation error in a readable format
	err := validationErrors[0]
	switch err.Tag() {
	case "required":
		return "Field '" + err.Field() + "' is required"
	case "min":
		return "Field '" + err.Field() + "' must be at least " + err.Param() + " characters long"
	case "max":
		return "Field '" + err.Field() + "' must be at most " + err.Param() + " characters long"
	case "uuid":
		return "Field '" + err.Field() + "' must be a valid UUID"
	case "oneof":
		return "Field '" + err.Field() + "' must be one of: " + err.Param()
	default:
		return "Field '" + err.Field() + "' is invalid"
	}
}

// formatBindingError formats binding errors into user-friendly messages
func formatBindingError(err error) string {
	if err == nil {
		return "Invalid request data"
	}

	// Handle specific binding errors
	switch err.Error() {
	case "EOF":
		return "Request body is empty"
	default:
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(validationErrors)
		}
		return "Invalid JSON format: " + err.Error()
	}
}

// ValidateUUID is a helper function to validate UUID strings
func ValidateUUID(c *gin.Context, uuidStr string, fieldName string) error {
	if uuidStr == "" {
		errorResp := models.ErrorResponse{
			Error:   "ValidationError",
			Code:    "INVALID_UUID",
			Message: "Field '" + fieldName + "' cannot be empty",
		}
		c.JSON(http.StatusBadRequest, errorResp)
		return validator.ValidationErrors{}
	}

	// Use google/uuid package for proper UUID validation
	if _, err := uuid.Parse(uuidStr); err != nil {
		errorResp := models.ErrorResponse{
			Error:   "ValidationError",
			Code:    "INVALID_UUID",
			Message: "Field '" + fieldName + "' must be a valid UUID",
		}
		c.JSON(http.StatusBadRequest, errorResp)
		return validator.ValidationErrors{}
	}

	return nil
}