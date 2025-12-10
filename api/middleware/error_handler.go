package middleware

import (
	"discovery-tree/api/models"
	"discovery-tree/domain"
	"discovery-tree/infrastructure"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware handles errors consistently across all endpoints
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			statusCode, errorResp := MapDomainError(err)
			
			// Log the error with structured logging
			slog.Error("Request error recovered",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("error", err.Error()),
				slog.Int("status", statusCode),
				slog.String("client_ip", c.ClientIP()),
			)
			
			c.JSON(statusCode, errorResp)
		} else {
			// Handle non-error panics
			slog.Error("Non-error panic recovered",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Any("recovered", recovered),
				slog.String("client_ip", c.ClientIP()),
			)
			
			errorResp := models.ErrorResponse{
				Error:   "InternalServerError",
				Code:    "PANIC_RECOVERED",
				Message: "An unexpected error occurred",
			}
			c.JSON(http.StatusInternalServerError, errorResp)
		}
		c.Abort()
	})
}

// MapDomainError converts domain errors to HTTP status codes and error responses
func MapDomainError(err error) (int, models.ErrorResponse) {
	switch e := err.(type) {
	case domain.ValidationError:
		return http.StatusBadRequest, models.ErrorResponse{
			Error:   "ValidationError",
			Code:    e.Field,
			Message: e.Message,
		}
	case domain.NotFoundError:
		return http.StatusNotFound, models.ErrorResponse{
			Error:   "NotFoundError",
			Code:    "RESOURCE_NOT_FOUND",
			Message: e.Error(),
		}
	case domain.ConstraintViolationError:
		return http.StatusConflict, models.ErrorResponse{
			Error:   "ConstraintViolationError",
			Code:    e.Constraint,
			Message: e.Message,
		}
	case infrastructure.FileSystemError:
		return http.StatusInternalServerError, models.ErrorResponse{
			Error:   "InternalServerError",
			Code:    "FILESYSTEM_ERROR",
			Message: "A storage error occurred",
		}
	default:
		return http.StatusInternalServerError, models.ErrorResponse{
			Error:   "InternalServerError",
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred",
		}
	}
}

// HandleError is a helper function that handlers can use to consistently handle errors
func HandleError(c *gin.Context, err error) {
	statusCode, errorResp := MapDomainError(err)
	
	// Log the error with structured logging
	slog.Error("Handler error",
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("error", err.Error()),
		slog.Int("status", statusCode),
		slog.String("client_ip", c.ClientIP()),
	)
	
	c.JSON(statusCode, errorResp)
}