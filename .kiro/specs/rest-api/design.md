# REST API Design Document

## Overview

This document describes the design for a REST API layer that provides HTTP endpoints for the Discovery Tree task management system. The API will expose the existing domain functionality through HTTP handlers, use dependency injection for service composition, and provide comprehensive OpenAPI documentation. The design follows clean architecture principles by keeping the API layer separate from domain logic while leveraging the existing TaskService and repository patterns.

## Architecture

The REST API will be implemented as a new layer in the existing clean architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    API Layer (New)                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   HTTP Server   │  │   HTTP Handlers │  │  OpenAPI    │ │
│  │   (Gin/Echo)    │  │   (Controllers) │  │   Schema    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Application Layer (Existing)                 │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  TaskService    │  │ Dependency      │                  │
│  │  (Domain Logic) │  │ Injection       │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                Domain Layer (Existing)                      │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │     Task        │  │ TaskRepository  │                  │
│  │   (Aggregate)   │  │  (Interface)    │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Infrastructure Layer (Existing)                │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │FileTaskRepository│  │    TaskDTO      │                  │
│  │ (JSON Storage)   │  │ (Serialization) │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### HTTP Server
- **Framework**: Gin (lightweight, fast HTTP framework for Go)
- **Port**: Configurable (default 8080)
- **Middleware**: CORS, logging, error handling, request validation
- **Content-Type**: application/json for all requests and responses

### HTTP Handlers
The API will provide the following handlers:

#### Task Handlers
```go
type TaskHandler struct {
    taskService *domain.TaskService
}

// POST /api/v1/tasks/root - Create root task
func (h *TaskHandler) CreateRootTask(c *gin.Context)

// POST /api/v1/tasks - Create child task  
func (h *TaskHandler) CreateChildTask(c *gin.Context)

// GET /api/v1/tasks/{id} - Get task by ID
func (h *TaskHandler) GetTask(c *gin.Context)

// GET /api/v1/tasks - Get all tasks
func (h *TaskHandler) GetAllTasks(c *gin.Context)

// GET /api/v1/tasks/root - Get root task
func (h *TaskHandler) GetRootTask(c *gin.Context)

// GET /api/v1/tasks/{id}/children - Get task children
func (h *TaskHandler) GetTaskChildren(c *gin.Context)

// PUT /api/v1/tasks/{id} - Update task
func (h *TaskHandler) UpdateTask(c *gin.Context)

// PUT /api/v1/tasks/{id}/status - Update task status
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context)

// PUT /api/v1/tasks/{id}/move - Move task
func (h *TaskHandler) MoveTask(c *gin.Context)

// DELETE /api/v1/tasks/{id} - Delete task
func (h *TaskHandler) DeleteTask(c *gin.Context)
```

#### Health Check Handler
```go
type HealthHandler struct{}

// GET /health - Health check endpoint
func (h *HealthHandler) HealthCheck(c *gin.Context)
```

### Dependency Injection Container
```go
type Container struct {
    TaskRepository domain.TaskRepository
    TaskService    *domain.TaskService
    TaskHandler    *TaskHandler
    HealthHandler  *HealthHandler
}

func NewContainer(repoPath string) (*Container, error)
func (c *Container) SetupRoutes() *gin.Engine
```

## Data Models

### API Request/Response Models
The API will use dedicated request/response models separate from domain entities:

#### Task Response Model
```go
type TaskResponse struct {
    ID          string     `json:"id"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    ParentID    *string    `json:"parentId"`
    Position    int        `json:"position"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
```

#### Create Root Task Request
```go
type CreateRootTaskRequest struct {
    Description string `json:"description" binding:"required,min=1"`
}
```

#### Create Child Task Request
```go
type CreateChildTaskRequest struct {
    Description string `json:"description" binding:"required,min=1"`
    ParentID    string `json:"parentId" binding:"required,uuid"`
}
```

#### Update Task Request
```go
type UpdateTaskRequest struct {
    Description string `json:"description" binding:"required,min=1"`
}
```

#### Update Status Request
```go
type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=TODO IN_PROGRESS DONE ROOT_WORK_ITEM"`
}
```

#### Move Task Request
```go
type MoveTaskRequest struct {
    ParentID *string `json:"parentId" binding:"omitempty,uuid"`
    Position int     `json:"position" binding:"min=0"`
}
```

#### Error Response Model
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Model Conversion Functions
```go
// Convert domain Task to API response
func TaskToResponse(task *domain.Task) TaskResponse

// Convert domain error to API error response
func ErrorToResponse(err error) ErrorResponse
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property Reflection

Analyzing the acceptance criteria to identify testable properties and eliminate redundancy:

**Redundancy Analysis:**
- Requirements 1.1, 1.2, 1.3, 1.4, 1.5 all test HTTP status code behavior - can be unified into one comprehensive property
- Requirements 8.1, 8.2, 8.3, 8.5 all test JSON serialization aspects - can be combined into a round-trip property  
- Requirements 3.3, 5.4, and other validation errors test similar patterns - can be unified into validation property
- Requirements 4.1, 4.2, 4.3, 4.4 all test successful retrieval - can be combined into retrieval consistency property
- Requirements 6.1, 6.2, 6.5 test deletion with position adjustment - can be unified into deletion consistency property

**Unique Properties Identified:**
Each property below provides distinct validation without logical overlap.

### Property 1: HTTP Status Code Consistency
*For any* API request, the HTTP response status code should correctly reflect the operation outcome: 2xx for success, 4xx for client errors, 5xx for server errors, and 404 for non-existent resources
**Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 4.5, 6.4**

### Property 2: JSON Serialization Round Trip  
*For any* task returned by the API, serializing and deserializing through JSON should preserve all task properties including consistent field names, ISO 8601 timestamps, string UUIDs, and parent-child relationships
**Validates: Requirements 8.1, 8.2, 8.3, 8.5**

### Property 3: Task Creation Consistency
*For any* valid task creation request, the created task should have the specified description, correct parent relationship, appropriate initial status, and return 201 Created status
**Validates: Requirements 3.1, 3.2**

### Property 4: Input Validation Consistency
*For any* invalid request data (empty descriptions, non-existent parents, invalid status values), the API should reject the request with 400 Bad Request and consistent error format
**Validates: Requirements 3.3, 3.5, 5.4, 8.4**

### Property 5: Task Retrieval Consistency
*For any* existing task or task collection, retrieval operations should return complete and accurate data with 200 OK status, including proper ordering for child tasks
**Validates: Requirements 4.1, 4.2, 4.3, 4.4**

### Property 6: Task Hierarchy Operations
*For any* task move operation, the system should maintain position consistency, prevent cycles, update hierarchy correctly, and return 200 OK for valid moves
**Validates: Requirements 5.1, 5.2, 5.3, 5.5**

### Property 7: Task Deletion Consistency
*For any* task deletion, the system should remove the task (and descendants if applicable), adjust sibling positions to maintain consistency, and return 204 No Content
**Validates: Requirements 6.1, 6.2, 6.5**

## Error Handling

### HTTP Status Code Mapping
- **200 OK**: Successful GET, PUT operations
- **201 Created**: Successful POST operations (task creation)
- **204 No Content**: Successful DELETE operations
- **400 Bad Request**: Invalid request data, validation errors
- **404 Not Found**: Resource not found (task, parent)
- **409 Conflict**: Business rule violations (duplicate root task)
- **500 Internal Server Error**: Unexpected server errors

### Error Response Format
All error responses will follow a consistent JSON structure:
```json
{
  "error": "ValidationError",
  "code": "INVALID_DESCRIPTION", 
  "message": "Task description cannot be empty"
}
```

### Domain Error Mapping
```go
func MapDomainError(err error) (int, ErrorResponse) {
    switch e := err.(type) {
    case domain.ValidationError:
        return 400, ErrorResponse{
            Error: "ValidationError",
            Code: e.Field,
            Message: e.Message,
        }
    case domain.NotFoundError:
        return 404, ErrorResponse{
            Error: "NotFoundError", 
            Code: "RESOURCE_NOT_FOUND",
            Message: e.Message,
        }
    case domain.ConstraintViolationError:
        return 409, ErrorResponse{
            Error: "ConstraintViolationError",
            Code: e.Constraint,
            Message: e.Message,
        }
    default:
        return 500, ErrorResponse{
            Error: "InternalServerError",
            Code: "INTERNAL_ERROR", 
            Message: "An unexpected error occurred",
        }
    }
}
```

## Testing Strategy

### Dual Testing Approach
The API will use both unit testing and property-based testing to ensure comprehensive coverage:

**Unit Testing:**
- Test specific HTTP endpoints with known inputs and expected outputs
- Test error handling scenarios with invalid requests
- Test middleware functionality (CORS, logging, validation)
- Test dependency injection container setup
- Use Go's `net/http/httptest` package for HTTP testing

**Property-Based Testing:**
- Use `github.com/leanovate/gopter` (already in dependencies) for property-based tests
- Each property-based test will run a minimum of 100 iterations
- Generate random valid and invalid request data to test API behavior
- Test JSON serialization/deserialization properties
- Test HTTP status code consistency across different request types

**Property-Based Test Requirements:**
- Each correctness property must be implemented by a single property-based test
- Each test must be tagged with: `**Feature: rest-api, Property {number}: {property_text}**`
- Tests must reference the design document property they implement
- Configure tests to run 100+ iterations for thorough validation

**Testing Framework Integration:**
- Unit tests: Go's built-in `testing` package with `testify/assert` for assertions
- Property tests: `gopter` with custom generators for API request data
- HTTP testing: `httptest` for creating test servers and requests
- Test organization: Co-locate tests with handlers using `_test.go` suffix

## OpenAPI Schema

### Schema Generation
The OpenAPI 3.0 schema will be generated using `swaggo/swag` annotations in handler functions:

```go
// CreateRootTask creates a new root task
// @Summary Create root task
// @Description Creates a new root task for the discovery tree
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateRootTaskRequest true "Root task creation request"
// @Success 201 {object} TaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/tasks/root [post]
func (h *TaskHandler) CreateRootTask(c *gin.Context) {
    // Implementation
}
```

### Schema Structure
- **Info**: API title, version, description, contact information
- **Servers**: Base URL configurations for different environments
- **Paths**: All endpoint definitions with parameters, request/response schemas
- **Components**: Reusable schemas for request/response models
- **Security**: Authentication schemes (if needed in future)

### Schema Serving
- **Endpoint**: `GET /api/v1/swagger.json` - Returns OpenAPI JSON schema
- **Documentation UI**: `GET /api/v1/docs` - Swagger UI for interactive documentation
- **Validation**: Schema validates all request/response formats

## Implementation Architecture

### Project Structure
```
api/
├── handlers/
│   ├── task_handler.go
│   ├── task_handler_test.go
│   ├── health_handler.go
│   └── health_handler_test.go
├── models/
│   ├── requests.go
│   ├── responses.go
│   └── converters.go
├── middleware/
│   ├── cors.go
│   ├── logging.go
│   └── error_handler.go
├── container/
│   ├── container.go
│   └── container_test.go
└── server/
    ├── server.go
    └── routes.go
```

### Dependency Flow
1. **Container** initializes repository and services
2. **Handlers** receive injected services through container
3. **Server** sets up routes and middleware using container
4. **Middleware** handles cross-cutting concerns (CORS, logging, errors)
5. **Models** provide request/response serialization

### Configuration
```go
type Config struct {
    Port         string `env:"PORT" default:"8080"`
    DataPath     string `env:"DATA_PATH" default:"./data/tasks.json"`
    LogLevel     string `env:"LOG_LEVEL" default:"info"`
    EnableCORS   bool   `env:"ENABLE_CORS" default:"true"`
    EnableSwagger bool  `env:"ENABLE_SWAGGER" default:"true"`
}
```

This design provides a clean separation between the API layer and existing domain logic while maintaining the established patterns for dependency injection and error handling.