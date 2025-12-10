package handlers

import (
	"bytes"
	"discovery-tree/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler_CreateRootTask_Success(t *testing.T) {
	// Setup
	repo := domain.NewInMemoryTaskRepository()
	service := domain.NewTaskService(repo)
	handler := NewTaskHandler(service, repo)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request
	requestBody := map[string]interface{}{
		"description": "Test root task",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest("POST", "/api/v1/tasks/root", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateRootTask(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "Test root task", response["description"])
	assert.Equal(t, "Root Work Item", response["status"])
	assert.Nil(t, response["parentId"])
	assert.Equal(t, float64(0), response["position"])
}

func TestTaskHandler_CreateRootTask_ValidationError(t *testing.T) {
	// Setup
	repo := domain.NewInMemoryTaskRepository()
	service := domain.NewTaskService(repo)
	handler := NewTaskHandler(service, repo)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request with empty description
	requestBody := map[string]interface{}{
		"description": "",
	}
	jsonBody, _ := json.Marshal(requestBody)
	c.Request = httptest.NewRequest("POST", "/api/v1/tasks/root", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.CreateRootTask(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "ValidationError", response["error"])
}

func TestTaskHandler_GetTask_Success(t *testing.T) {
	// Setup
	repo := domain.NewInMemoryTaskRepository()
	service := domain.NewTaskService(repo)
	handler := NewTaskHandler(service, repo)

	// Create a task first
	task, err := service.CreateRootTask("Test task")
	require.NoError(t, err)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: task.ID().String()}}
	c.Request = httptest.NewRequest("GET", "/api/v1/tasks/"+task.ID().String(), nil)

	// Execute
	handler.GetTask(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, task.ID().String(), response["id"])
	assert.Equal(t, "Test task", response["description"])
}

func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	// Setup
	repo := domain.NewInMemoryTaskRepository()
	service := domain.NewTaskService(repo)
	handler := NewTaskHandler(service, repo)

	// Create Gin context with non-existent ID
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	nonExistentID := domain.NewTaskID()
	c.Params = gin.Params{{Key: "id", Value: nonExistentID.String()}}
	c.Request = httptest.NewRequest("GET", "/api/v1/tasks/"+nonExistentID.String(), nil)

	// Execute
	handler.GetTask(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "NotFoundError", response["error"])
}