# Implementation Plan

- [x] 1. Set up API project structure and dependencies
  - Create api/ directory with subdirectories for handlers, models, middleware, container, and server
  - Add required dependencies: gin-gonic/gin, swaggo/swag, swaggo/gin-swagger for OpenAPI
  - Set up basic Go module configuration for API layer
  - _Requirements: 7.1_

- [x] 2. Implement API request/response models
  - [x] 2.1 Create request models for task operations
    - Define CreateRootTaskRequest, CreateChildTaskRequest, UpdateTaskRequest, UpdateStatusRequest, MoveTaskRequest structs
    - Add JSON binding tags and validation rules using gin binding
    - _Requirements: 2.2, 8.1_

  - [x] 2.2 Create response models and converters
    - Define TaskResponse, ErrorResponse structs with proper JSON tags
    - Implement TaskToResponse converter function from domain.Task to TaskResponse
    - Implement ErrorToResponse converter for domain errors to API error responses
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

  - [ ]* 2.3 Write property test for JSON serialization round trip
    - **Property 2: JSON Serialization Round Trip**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.5**

- [x] 3. Implement dependency injection container
  - [x] 3.1 Create container structure and initialization
    - Define Container struct with TaskRepository, TaskService, and handler dependencies
    - Implement NewContainer function that initializes repository and services
    - Add configuration struct for API settings (port, data path, CORS, etc.)
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [x] 3.2 Implement container service management
    - Add methods to container for creating and managing service instances
    - Ensure proper interface usage for loose coupling between components
    - Implement service lifetime management for singleton services
    - _Requirements: 7.2, 7.4, 7.5_

  - [ ]* 3.3 Write unit tests for dependency injection container
    - Test container initialization with different configurations
    - Test service creation and dependency injection
    - Test interface usage and loose coupling
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [x] 4. Implement HTTP handlers for task operations
  - [x] 4.1 Create TaskHandler struct and basic CRUD operations
    - Define TaskHandler with injected TaskService dependency
    - Implement CreateRootTask handler with request validation and 201 response
    - Implement CreateChildTask handler with parent validation
    - Implement GetTask, GetAllTasks, GetRootTask, GetTaskChildren handlers
    - _Requirements: 3.1, 3.2, 4.1, 4.2, 4.3, 4.4_

  - [ ]* 4.2 Write property test for task creation consistency
    - **Property 3: Task Creation Consistency**
    - **Validates: Requirements 3.1, 3.2**

  - [ ]* 4.3 Write property test for task retrieval consistency
    - **Property 5: Task Retrieval Consistency**
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.4**

  - [x] 4.4 Implement task update and move operations
    - Implement UpdateTask handler for description updates
    - Implement UpdateTaskStatus handler with status validation
    - Implement MoveTask handler with hierarchy validation and cycle prevention
    - _Requirements: 5.1, 5.2, 5.3_

  - [ ]* 4.5 Write property test for task hierarchy operations
    - **Property 6: Task Hierarchy Operations**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.5**

  - [x] 4.6 Implement task deletion operations
    - Implement DeleteTask handler with cascading deletion support
    - Handle leaf task deletion with sibling position adjustment
    - Handle root task deletion with entire tree removal
    - _Requirements: 6.1, 6.2, 6.3_

  - [ ]* 4.7 Write property test for task deletion consistency
    - **Property 7: Task Deletion Consistency**
    - **Validates: Requirements 6.1, 6.2, 6.5**

- [-] 5. Implement error handling and validation
  - [x] 5.1 Create error handling middleware and domain error mapping
    - Implement MapDomainError function for converting domain errors to HTTP responses
    - Create error handling middleware for consistent error response format
    - Add request validation middleware using gin binding
    - _Requirements: 1.2, 1.3, 3.3, 3.5, 5.4_

  - [ ]* 5.2 Write property test for input validation consistency
    - **Property 4: Input Validation Consistency**
    - **Validates: Requirements 3.3, 3.5, 5.4, 8.4**

  - [ ]* 5.3 Write property test for HTTP status code consistency
    - **Property 1: HTTP Status Code Consistency**
    - **Validates: Requirements 1.1, 1.2, 1.3, 1.4, 1.5, 4.5, 6.4**

- [ ] 6. Implement HTTP server and routing
  - [ ] 6.1 Create server setup and route configuration
    - Implement server initialization with Gin engine
    - Set up API routes with proper HTTP methods and paths (/api/v1/tasks/*)
    - Add middleware for CORS, logging, and error handling
    - Configure server with dependency injection container
    - _Requirements: 1.1, 7.5_

  - [ ] 6.2 Add health check endpoint and middleware
    - Implement HealthHandler with health check endpoint
    - Add CORS middleware for cross-origin requests
    - Add request logging middleware for debugging and monitoring
    - _Requirements: 1.1, 1.4_

  - [ ]* 6.3 Write integration tests for HTTP server
    - Test server startup and route registration
    - Test middleware functionality (CORS, logging, error handling)
    - Test health check endpoint
    - _Requirements: 1.1, 1.4_

- [ ] 7. Implement OpenAPI documentation
  - [ ] 7.1 Add Swagger annotations to handlers
    - Add swaggo annotations to all handler functions with proper documentation
    - Define request/response schemas in annotations
    - Include error response documentation with status codes
    - Add API metadata (title, version, description)
    - _Requirements: 2.1, 2.2, 2.3, 2.5_

  - [ ] 7.2 Set up OpenAPI schema generation and serving
    - Configure swaggo to generate OpenAPI JSON schema
    - Add endpoint to serve OpenAPI schema at /api/v1/swagger.json
    - Add Swagger UI endpoint at /api/v1/docs for interactive documentation
    - Ensure schema includes all endpoints, models, and validation rules
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [ ]* 7.3 Write unit tests for OpenAPI schema completeness
    - Test that generated schema includes all implemented endpoints
    - Test that request/response schemas are properly documented
    - Test that validation rules are included in schema
    - _Requirements: 2.1, 2.2, 2.3, 2.5_

- [ ] 8. Create main application entry point
  - [ ] 8.1 Implement main function and application startup
    - Create cmd/api/main.go with application entry point
    - Load configuration from environment variables with defaults
    - Initialize dependency injection container with configuration
    - Start HTTP server with graceful shutdown handling
    - _Requirements: 7.1, 1.1_

  - [ ] 8.2 Add configuration management and environment setup
    - Implement configuration loading from environment variables
    - Add default values for development environment
    - Add validation for required configuration parameters
    - Document configuration options in README or comments
    - _Requirements: 7.1_

- [ ] 9. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 10. Integration and final testing
  - [ ] 10.1 Create end-to-end API tests
    - Test complete API workflows: create root task, add children, update, move, delete
    - Test error scenarios across different endpoints
    - Test OpenAPI schema accessibility and correctness
    - Verify all HTTP status codes and response formats
    - _Requirements: 1.1, 1.4, 2.4_

  - [ ]* 10.2 Write performance and load tests
    - Test API performance with multiple concurrent requests
    - Test memory usage and resource management
    - Test file persistence performance with large task trees
    - _Requirements: 1.1_

- [ ] 11. Final Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.