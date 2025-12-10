package handlers

import (
	"discovery-tree/api/middleware"
	"discovery-tree/api/models"
	"discovery-tree/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles HTTP requests for task operations
type TaskHandler struct {
	taskService    *domain.TaskService
	taskRepository domain.TaskRepository
}

// NewTaskHandler creates a new TaskHandler with injected dependencies
func NewTaskHandler(taskService *domain.TaskService, taskRepository domain.TaskRepository) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		taskRepository: taskRepository,
	}
}

// CreateRootTask creates a new root task
// POST /api/v1/tasks/root
func (h *TaskHandler) CreateRootTask(c *gin.Context) {
	var req models.CreateRootTaskRequest
	
	// Bind and validate the request
	if err := middleware.BindJSON(c, &req); err != nil {
		return
	}

	// Create the root task using the service
	task, err := h.taskService.CreateRootTask(req.Description)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusCreated, response)
}

// CreateChildTask creates a new child task
// POST /api/v1/tasks
func (h *TaskHandler) CreateChildTask(c *gin.Context) {
	var req models.CreateChildTaskRequest
	
	// Bind and validate the request
	if err := middleware.BindJSON(c, &req); err != nil {
		return
	}

	// Convert parent ID string to TaskID
	parentID, err := domain.TaskIDFromString(req.ParentID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Create the child task using the service
	task, err := h.taskService.CreateChildTask(req.Description, parentID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusCreated, response)
}

// GetTask retrieves a specific task by ID
// GET /api/v1/tasks/{id}
func (h *TaskHandler) GetTask(c *gin.Context) {
	idParam := c.Param("id")
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Find the task using the repository
	task, err := h.taskRepository.FindByID(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// GetAllTasks retrieves all tasks
// GET /api/v1/tasks
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	// Find all tasks using the repository
	tasks, err := h.taskRepository.FindAll()
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert all tasks to response models
	responses := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = models.TaskToResponse(task)
	}

	c.JSON(http.StatusOK, responses)
}

// GetRootTask retrieves the root task
// GET /api/v1/tasks/root
func (h *TaskHandler) GetRootTask(c *gin.Context) {
	// Find the root task using the repository
	task, err := h.taskRepository.FindRoot()
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// GetTaskChildren retrieves children of a specific task
// GET /api/v1/tasks/{id}/children
func (h *TaskHandler) GetTaskChildren(c *gin.Context) {
	idParam := c.Param("id")
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// First verify the parent task exists
	_, err = h.taskRepository.FindByID(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Find children using the repository
	children, err := h.taskRepository.FindByParentID(&taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert all children to response models
	responses := make([]models.TaskResponse, len(children))
	for i, child := range children {
		responses[i] = models.TaskToResponse(child)
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateTask updates a task's description
// PUT /api/v1/tasks/{id}
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	var req models.UpdateTaskRequest
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Bind and validate the request
	if err := middleware.BindJSON(c, &req); err != nil {
		return
	}

	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Find the task first
	task, err := h.taskRepository.FindByID(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Update the description
	err = task.UpdateDescription(req.Description)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Save the updated task
	err = h.taskRepository.Save(task)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// UpdateTaskStatus updates a task's status
// PUT /api/v1/tasks/{id}/status
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	idParam := c.Param("id")
	var req models.UpdateStatusRequest
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Bind and validate the request
	if err := middleware.BindJSON(c, &req); err != nil {
		return
	}

	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert status string to Status
	status, err := domain.NewStatus(req.Status)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Update the status using the service (includes validation)
	err = h.taskService.ChangeTaskStatus(taskID, status)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Retrieve the updated task to return
	task, err := h.taskRepository.FindByID(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// MoveTask moves a task to a new position or parent
// PUT /api/v1/tasks/{id}/move
func (h *TaskHandler) MoveTask(c *gin.Context) {
	idParam := c.Param("id")
	var req models.MoveTaskRequest
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Bind and validate the request
	if err := middleware.BindJSON(c, &req); err != nil {
		return
	}

	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert parent ID string to TaskID if provided
	var newParentID *domain.TaskID
	if req.ParentID != nil {
		// Validate parent UUID format if provided
		if err := middleware.ValidateUUID(c, *req.ParentID, "parentId"); err != nil {
			return
		}
		
		parentID, err := domain.TaskIDFromString(*req.ParentID)
		if err != nil {
			middleware.HandleError(c, err)
			return
		}
		newParentID = &parentID
	}

	// Move the task using the service (includes validation and position adjustments)
	err = h.taskService.MoveTask(taskID, newParentID, req.Position)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Retrieve the updated task to return
	task, err := h.taskRepository.FindByID(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Convert to response model and return
	response := models.TaskToResponse(task)
	c.JSON(http.StatusOK, response)
}

// DeleteTask deletes a task
// DELETE /api/v1/tasks/{id}
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	
	// Validate UUID format
	if err := middleware.ValidateUUID(c, idParam, "id"); err != nil {
		return
	}
	
	// Convert ID string to TaskID
	taskID, err := domain.TaskIDFromString(idParam)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Delete the task using the service (includes cascading deletion and position adjustments)
	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		middleware.HandleError(c, err)
		return
	}

	// Return 204 No Content for successful deletion
	c.Status(http.StatusNoContent)
}

