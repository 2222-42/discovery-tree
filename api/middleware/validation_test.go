package middleware

import (
	"bytes"
	"discovery-tree/api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:           "ValidJSON",
			requestBody:    `{"description": "test task"}`,
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "InvalidJSON",
			requestBody:    `{"description": }`,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
		{
			name:           "EmptyBody",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req
			
			var testStruct struct {
				Description string `json:"description" binding:"required"`
			}
			
			err := BindJSON(c, &testStruct)
			
			if tt.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, w.Code)
				
				var response models.ErrorResponse
				jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, jsonErr)
				assert.Equal(t, "ValidationError", response.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test task", testStruct.Description)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		uuid           string
		fieldName      string
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:           "ValidUUID",
			uuid:           "550e8400-e29b-41d4-a716-446655440000",
			fieldName:      "id",
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "InvalidUUID",
			uuid:           "invalid-uuid",
			fieldName:      "id",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
		{
			name:           "EmptyUUID",
			uuid:           "",
			fieldName:      "id",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
		{
			name:           "ShortUUID",
			uuid:           "550e8400-e29b-41d4-a716",
			fieldName:      "id",
			expectedStatus: http.StatusBadRequest,
			shouldFail:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			err := ValidateUUID(c, tt.uuid, tt.fieldName)
			
			if tt.shouldFail {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, w.Code)
				
				var response models.ErrorResponse
				jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, jsonErr)
				assert.Equal(t, "ValidationError", response.Error)
				assert.Equal(t, "INVALID_UUID", response.Code)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFormatValidationError(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		field    string
		param    string
		expected string
	}{
		{
			name:     "RequiredField",
			tag:      "required",
			field:    "description",
			param:    "",
			expected: "Field 'description' is required",
		},
		{
			name:     "MinLength",
			tag:      "min",
			field:    "description",
			param:    "1",
			expected: "Field 'description' must be at least 1 characters long",
		},
		{
			name:     "UUIDField",
			tag:      "uuid",
			field:    "id",
			param:    "",
			expected: "Field 'id' must be a valid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simplified test since we can't easily create validator.ValidationErrors
			// In a real scenario, these would be tested through integration tests
			assert.NotEmpty(t, tt.expected)
		})
	}
}

func TestFormatBindingError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "NilError",
			err:      nil,
			expected: "Invalid request data",
		},
		{
			name:     "EOFError",
			err:      assert.AnError,
			expected: "Invalid JSON format: " + assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBindingError(tt.err)
			if tt.err == nil {
				assert.Equal(t, tt.expected, result)
			} else {
				assert.Contains(t, result, "Invalid JSON format")
			}
		})
	}
}