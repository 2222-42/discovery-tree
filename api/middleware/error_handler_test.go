package middleware

import (
	"discovery-tree/api/models"
	"discovery-tree/domain"
	"discovery-tree/infrastructure"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMapDomainError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedError  string
		expectedCode   string
	}{
		{
			name:           "ValidationError",
			err:            domain.NewValidationError("description", "cannot be empty"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "ValidationError",
			expectedCode:   "description",
		},
		{
			name:           "NotFoundError",
			err:            domain.NewNotFoundError("Task", "123"),
			expectedStatus: http.StatusNotFound,
			expectedError:  "NotFoundError",
			expectedCode:   "RESOURCE_NOT_FOUND",
		},
		{
			name:           "ConstraintViolationError",
			err:            domain.NewConstraintViolationError("UNIQUE_ROOT", "only one root task allowed"),
			expectedStatus: http.StatusConflict,
			expectedError:  "ConstraintViolationError",
			expectedCode:   "UNIQUE_ROOT",
		},
		{
			name:           "FileSystemError",
			err:            infrastructure.NewFileSystemError("read", "/path/to/file", assert.AnError),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "InternalServerError",
			expectedCode:   "FILESYSTEM_ERROR",
		},
		{
			name:           "UnknownError",
			err:            assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "InternalServerError",
			expectedCode:   "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, errorResp := MapDomainError(tt.err)
			
			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedError, errorResp.Error)
			assert.Equal(t, tt.expectedCode, errorResp.Code)
			assert.NotEmpty(t, errorResp.Message)
		})
	}
}

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	err := domain.NewValidationError("description", "cannot be empty")
	HandleError(c, err)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response models.ErrorResponse
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)
	assert.Equal(t, "ValidationError", response.Error)
	assert.Equal(t, "description", response.Code)
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(ErrorHandler())
	
	// Test panic recovery
	router.GET("/panic", func(c *gin.Context) {
		panic(domain.NewValidationError("test", "test error"))
	})
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response models.ErrorResponse
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, jsonErr)
	assert.Equal(t, "ValidationError", response.Error)
}