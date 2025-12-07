# Implementation Plan

- [x] 1. Create infrastructure layer structure and error types
  - Create `infrastructure` directory if it doesn't exist
  - Create `FileSystemError` type for file I/O errors
  - Add error wrapping utilities for OS errors
  - _Requirements: 3.4, 4.3, 4.4_

- [ ] 2. Implement TaskDTO for JSON serialization
  - Create `TaskDTO` struct with JSON tags
  - Implement `ToDTO(task *Task)` function to convert Task to DTO
  - Implement `FromDTO(dto TaskDTO)` function to convert DTO to Task
  - Handle nil parent ID conversion (nil <-> null)
  - Ensure timestamp formatting uses ISO 8601
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2_

- [ ]* 2.1 Write property test for serialization round trip
  - **Property 1: Serialization completeness**
  - **Validates: Requirements 1.1, 2.1**

- [ ]* 2.2 Write property test for null parent handling
  - **Property 2: Null parent handling**
  - **Validates: Requirements 1.2, 2.2**

- [ ]* 2.3 Write property test for timestamp format
  - **Property 3: Timestamp format consistency**
  - **Validates: Requirements 1.3**

- [ ]* 2.4 Write property test for invalid data rejection
  - **Property 4: Invalid data rejection**
  - **Validates: Requirements 2.3, 2.4**

- [ ] 3. Implement FileTaskRepository structure and initialization
  - Create `FileTaskRepository` struct with filePath, tasks map, and mutex
  - Implement `NewFileTaskRepository(filePath string)` constructor
  - Add default file path handling (use `./data/tasks.json` if empty)
  - Implement directory creation for file path
  - Implement `load()` method to read JSON from file
  - Handle non-existent file case (initialize empty)
  - Handle invalid JSON and invalid task data with errors
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 5.4, 6.1, 6.2, 6.3, 6.4_

- [ ]* 3.1 Write property test for directory creation
  - **Property 16: Directory creation**
  - **Validates: Requirements 6.3**

- [ ] 4. Implement atomic file persistence
  - Implement `persist()` method with atomic write pattern
  - Write to temporary file with `.tmp` suffix
  - Use `os.Rename()` for atomic replacement
  - Clean up temporary file on failure
  - Format JSON with indentation (2 spaces)
  - _Requirements: 3.1, 3.2, 3.3, 8.1_

- [ ]* 4.1 Write property test for save-load round trip
  - **Property 5: Save-load round trip**
  - **Validates: Requirements 3.1, 4.1**

- [ ]* 4.2 Write property test for atomic write safety
  - **Property 6: Atomic write safety**
  - **Validates: Requirements 3.2**

- [ ]* 4.3 Write property test for human-readable formatting
  - **Property 17: Human-readable formatting**
  - **Validates: Requirements 8.1**

- [ ] 5. Implement Save method
  - Implement `Save(task *Task)` method
  - Add task to in-memory map (or update if exists)
  - Call `persist()` to write to file
  - Use write lock (mutex.Lock) for thread safety
  - _Requirements: 3.1, 5.2, 5.3_

- [ ]* 5.1 Write property test for save-find consistency
  - **Property 9: Save then find consistency**
  - **Validates: Requirements 5.2**

- [ ] 6. Implement query methods
  - Implement `FindByID(id TaskID)` with read lock
  - Implement `FindAll()` to return all tasks with read lock
  - Implement `FindRoot()` to find task with nil parent with read lock
  - Implement `FindByParentID(parentID *TaskID)` with ordering by position with read lock
  - Return `NotFoundError` when task doesn't exist
  - _Requirements: 5.1, 5.2_

- [ ]* 6.1 Write property test for FindByParentID ordering
  - **Property 10: FindByParentID ordering**
  - **Validates: Requirements 5.2**

- [ ] 7. Implement Delete method
  - Implement `Delete(id TaskID)` for single task deletion
  - Remove task from in-memory map
  - Call `persist()` to write changes
  - Use write lock for thread safety
  - Return `NotFoundError` if task doesn't exist
  - _Requirements: 5.3_

- [ ]* 7.1 Write property test for delete removes task
  - **Property 11: Delete removes task**
  - **Validates: Requirements 5.3**

- [ ] 8. Implement DeleteSubtree method
  - Implement `DeleteSubtree(id TaskID)` for cascading deletion
  - Find all descendants recursively using parent relationships
  - Remove task and all descendants from in-memory map
  - Call `persist()` to write changes
  - Use write lock for thread safety
  - _Requirements: 5.3_

- [ ]* 8.1 Write property test for cascading deletion
  - **Property 12: DeleteSubtree cascades**
  - **Validates: Requirements 5.3**

- [ ] 9. Add concurrency safety tests
  - Verify concurrent read operations are safe
  - Verify concurrent write operations are safe
  - Use Go race detector during testing
  - _Requirements: 7.1, 7.2, 7.3_

- [ ]* 9.1 Write property test for concurrent read safety
  - **Property 13: Concurrent read safety**
  - **Validates: Requirements 7.1, 7.3**

- [ ]* 9.2 Write property test for concurrent write safety
  - **Property 14: Concurrent write safety**
  - **Validates: Requirements 7.1, 7.2**

- [ ] 10. Create property test generators
  - Implement generator for random tasks
  - Implement generator for random task trees
  - Implement generator for invalid JSON data
  - Implement generator for concurrent operation sequences
  - _Requirements: All property tests_

- [ ] 11. Update TaskService to use dependency injection
  - Modify `TaskService` to accept `TaskRepository` interface in constructor
  - Remove direct dependency on in-memory repository
  - Update all service methods to use injected repository
  - _Requirements: 5.1_

- [ ] 12. Create initialization example in main or application layer
  - Show how to create `FileTaskRepository` instance
  - Show how to inject repository into `TaskService`
  - Demonstrate configuration with custom file path
  - _Requirements: 5.1, 6.1_

- [ ] 13. Final checkpoint - Ensure all tests pass
  - Run all unit tests and property tests
  - Verify all 17 correctness properties are validated
  - Run with race detector enabled
  - Test with actual file system operations
  - Ask the user if questions arise
