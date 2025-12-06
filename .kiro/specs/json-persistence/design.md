# JSON Persistence Design

## Overview

This design document describes the implementation of JSON-based persistence for the Discovery Tree task management system. The persistence layer provides a lightweight storage solution that serializes the task tree to JSON files, enabling data persistence between application sessions without requiring a database system.

The implementation follows the Repository pattern established in the domain model, providing a file-based implementation of the TaskRepository interface. This allows for dependency injection and makes it easy to swap persistence implementations in the future.

## Architecture

### Layered Architecture

The JSON persistence implementation fits into the existing layered architecture:

1. **Domain Layer**: Defines the TaskRepository interface (already exists)
2. **Infrastructure Layer**: Implements the FileTaskRepository that persists to JSON files
3. **Application Layer**: Uses the repository through dependency injection

### Key Design Decisions

1. **File-based storage**: Store all tasks in a single JSON file for simplicity
2. **In-memory cache**: Maintain an in-memory map of tasks for fast access, synchronized with file storage
3. **Atomic writes**: Use atomic file operations (write to temp file, then rename) to prevent corruption
4. **Mutex protection**: Use sync.RWMutex to handle concurrent access safely
5. **Lazy loading**: Load data from file on first access or explicit initialization

## Components and Interfaces

### 1. TaskDTO (Data Transfer Object)

A struct for JSON serialization/deserialization that exposes Task fields.

```go
type TaskDTO struct {
    ID          string     `json:"id"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    ParentID    *string    `json:"parentId"`  // pointer to handle null
    Position    int        `json:"position"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}
```

**Responsibilities:**
- Provide JSON-serializable representation of Task
- Handle conversion between domain Task and DTO
- Support null values for optional fields (parentID)

**Key Methods:**
- `ToDTO(task *Task) TaskDTO`
- `FromDTO(dto TaskDTO) (*Task, error)`

### 2. FileTaskRepository

Implements the TaskRepository interface with JSON file persistence.

```go
type FileTaskRepository struct {
    filePath string
    tasks    map[string]*Task  // in-memory cache, keyed by task ID
    mu       sync.RWMutex      // protects concurrent access
}
```

**Responsibilities:**
- Implement all TaskRepository interface methods
- Maintain in-memory cache synchronized with file storage
- Handle JSON serialization/deserialization
- Ensure atomic file writes
- Provide thread-safe access

**Key Methods:**
- `NewFileTaskRepository(filePath string) (*FileTaskRepository, error)`
- `Save(task *Task) error`
- `FindByID(id TaskID) (*Task, error)`
- `FindByParentID(parentID *TaskID) ([]*Task, error)`
- `FindRoot() (*Task, error)`
- `FindAll() ([]*Task, error)`
- `Delete(id TaskID) error`
- `DeleteSubtree(id TaskID) error`
- `load() error` (private - loads from file)
- `persist() error` (private - saves to file)

### 3. Repository Factory (Optional)

A factory function to create repository instances with configuration.

```go
func NewRepository(config RepositoryConfig) (TaskRepository, error)
```

**Responsibilities:**
- Create appropriate repository implementation based on configuration
- Handle default configuration values
- Validate configuration parameters

## Data Models

### JSON File Format

The JSON file stores an array of task DTOs:

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "description": "Root task",
    "status": "Root Work Item",
    "parentId": null,
    "position": 0,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "description": "Child task 1",
    "status": "TODO",
    "parentId": "550e8400-e29b-41d4-a716-446655440000",
    "position": 0,
    "createdAt": "2024-01-15T10:31:00Z",
    "updatedAt": "2024-01-15T10:31:00Z"
  }
]
```

### Field Mappings

| Domain Field | JSON Field | Type | Notes |
|--------------|------------|------|-------|
| id | id | string | UUID format |
| description | description | string | Non-empty |
| status | status | string | Human-readable status name |
| parentID | parentId | string or null | null for root tasks |
| position | position | int | 0-indexed |
| createdAt | createdAt | string | ISO 8601 format |
| updatedAt | updatedAt | string | ISO 8601 format |

### In-Memory Structure

The repository maintains a map for O(1) lookups:

```go
tasks map[string]*Task  // key: task ID string
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Serialization Properties

**Property 1: Serialization completeness**
*For any* valid task, serializing it to JSON and then deserializing should produce a task with identical field values.
**Validates: Requirements 1.1, 2.1**

**Property 2: Null parent handling**
*For any* root task (with nil parent), serializing to JSON should produce a null parentId field, and deserializing should restore nil parentID.
**Validates: Requirements 1.2, 2.2**

**Property 3: Timestamp format consistency**
*For any* task with timestamps, serializing to JSON should use ISO 8601 format, and deserializing should restore the exact timestamp values.
**Validates: Requirements 1.3**

**Property 4: Invalid data rejection**
*For any* JSON with invalid task data (empty description, invalid status, negative position), deserialization should return an error.
**Validates: Requirements 2.3, 2.4**

### Persistence Properties

**Property 5: Save-load round trip**
*For any* collection of tasks, saving them to a file and then loading should restore the exact same collection.
**Validates: Requirements 3.1, 4.1**

**Property 6: Atomic write safety**
*For any* save operation, if the operation completes successfully, the file should contain valid JSON; if it fails, the previous file contents should remain intact.
**Validates: Requirements 3.2**

**Property 7: Empty file initialization**
*For any* non-existent file path, loading should initialize with an empty task collection without error.
**Validates: Requirements 4.2**

**Property 8: Invalid JSON rejection**
*For any* file containing invalid JSON, loading should return an error.
**Validates: Requirements 4.3**

### Repository Interface Properties

**Property 9: Save then find consistency**
*For any* task, after saving it, FindByID should return a task with identical field values.
**Validates: Requirements 5.2**

**Property 10: FindByParentID ordering**
*For any* parent ID, FindByParentID should return children ordered by position (ascending).
**Validates: Requirements 5.2**

**Property 11: Delete removes task**
*For any* task, after deleting it, FindByID should return a NotFoundError.
**Validates: Requirements 5.3**

**Property 12: DeleteSubtree cascades**
*For any* task with descendants, DeleteSubtree should remove the task and all its descendants.
**Validates: Requirements 5.3**

### Concurrency Properties

**Property 13: Concurrent read safety**
*For any* sequence of concurrent read operations, all reads should return consistent data without panics or race conditions.
**Validates: Requirements 7.1, 7.3**

**Property 14: Concurrent write safety**
*For any* sequence of concurrent write operations, all writes should complete without data corruption or race conditions.
**Validates: Requirements 7.1, 7.2**

### Configuration Properties

**Property 15: Custom path usage**
*For any* valid file path provided to the repository, the repository should use that path for storage.
**Validates: Requirements 6.1**

**Property 16: Directory creation**
*For any* file path where the directory does not exist, the repository should create the necessary directories.
**Validates: Requirements 6.3**

### Format Properties

**Property 17: Human-readable formatting**
*For any* saved JSON file, the content should be formatted with indentation for readability.
**Validates: Requirements 8.1**

## Error Handling

### Error Types

The persistence layer uses existing domain error types and adds infrastructure-specific errors:

1. **ValidationError** (domain): Invalid data during deserialization
2. **NotFoundError** (domain): Task not found in repository
3. **FileSystemError** (new): File I/O errors
   - File not readable
   - File not writable
   - Directory creation failure
   - Atomic write failure

### Error Handling Strategy

1. **Deserialization errors**: Return ValidationError with details about invalid fields
2. **File I/O errors**: Wrap OS errors in FileSystemError with context
3. **Not found errors**: Return NotFoundError when task doesn't exist
4. **Concurrent access**: Use mutex to prevent race conditions (no error, just blocking)

### Validation Points

1. **Load**: Validate JSON structure and task data
2. **Save**: Validate task before adding to collection
3. **File operations**: Check file permissions and disk space
4. **Configuration**: Validate file path format

## Testing Strategy

### Unit Testing

Unit tests will verify specific examples and edge cases:

1. **DTO Conversion Tests**:
   - Task to DTO conversion
   - DTO to Task conversion
   - Null parent handling
   - Timestamp formatting

2. **Repository Tests**:
   - Save and retrieve single task
   - Save and retrieve multiple tasks
   - FindByParentID with ordering
   - Delete operations
   - Empty repository initialization

3. **File Operations Tests**:
   - Atomic write behavior
   - Directory creation
   - Invalid file path handling
   - Corrupted JSON handling

4. **Edge Cases**:
   - Empty task collection
   - Large task trees
   - Special characters in descriptions
   - Concurrent access scenarios

### Property-Based Testing

Property-based testing will verify universal properties across all inputs using **gopter**.

**Configuration**:
- Each property test should run a minimum of 100 iterations
- Tests should use generators for valid task trees and random operations
- Each property test must be tagged with a comment referencing the design document property

**Tag Format**:
```go
// Feature: json-persistence, Property 1: Serialization completeness
```

**Property Test Coverage**:
- Each correctness property (1-17) will be implemented as a single property-based test
- Generators will create random task trees and operation sequences
- Properties will verify invariants hold across all generated inputs

**Test Organization**:
- Property tests will be in `infrastructure/file_task_repository_test.go`
- Generators will be in `infrastructure/generators_test.go`
- Each property test will focus on a single correctness property

### Integration Testing

Integration tests will verify:
- End-to-end persistence workflows (create, save, load, modify, save)
- Integration with existing domain services
- Dependency injection with application layer
- File system interactions with actual files

## Implementation Notes

### Atomic Write Implementation

To ensure atomic writes and prevent file corruption:

1. Write data to a temporary file (`.tmp` suffix)
2. Sync the temporary file to disk
3. Rename temporary file to target file (atomic operation on POSIX systems)
4. If any step fails, clean up temporary file

```go
func (r *FileTaskRepository) persist() error {
    tmpPath := r.filePath + ".tmp"
    
    // Write to temp file
    if err := writeJSONToFile(tmpPath, r.tasks); err != nil {
        return err
    }
    
    // Atomic rename
    if err := os.Rename(tmpPath, r.filePath); err != nil {
        os.Remove(tmpPath)  // cleanup
        return err
    }
    
    return nil
}
```

### Concurrency Strategy

Use `sync.RWMutex` for reader-writer locking:

- **Read operations** (FindByID, FindAll, etc.): Use `RLock()` to allow concurrent reads
- **Write operations** (Save, Delete, etc.): Use `Lock()` for exclusive access
- **Load/Persist**: Use `Lock()` since they modify the entire collection

### Default Configuration

Default values for repository configuration:

- **Default file path**: `./data/tasks.json`
- **Default directory permissions**: `0755`
- **Default file permissions**: `0644`
- **JSON indentation**: 2 spaces

### Performance Considerations

1. **In-memory cache**: All operations work with in-memory data for speed
2. **Persist on every write**: Trade-off between durability and performance (acceptable for MVP)
3. **Full file write**: Simple approach for MVP; can optimize later with incremental writes
4. **Index-free queries**: Use linear search for FindByParentID (acceptable for small-medium trees)

### Future Enhancements

Potential improvements beyond MVP:

1. **Incremental persistence**: Only write changed tasks
2. **Write-ahead logging**: For better durability and recovery
3. **Compression**: Compress JSON for large task trees
4. **Backup/versioning**: Automatic backups before writes
5. **Migration support**: Handle schema changes gracefully
6. **Performance optimization**: Add indexes for common queries
7. **Batch operations**: Optimize multiple saves/deletes

### Dependency Injection Pattern

The repository should be injected into services:

```go
// Application layer
type TaskService struct {
    repo TaskRepository
    // ... other dependencies
}

func NewTaskService(repo TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

// Main/initialization
func main() {
    repo, err := NewFileTaskRepository("./data/tasks.json")
    if err != nil {
        log.Fatal(err)
    }
    
    service := NewTaskService(repo)
    // ... use service
}
```

This allows easy testing with mock repositories and future swapping of implementations (e.g., to a database).
