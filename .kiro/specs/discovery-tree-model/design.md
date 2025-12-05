# Discovery Tree Model Design

## Overview

This design document describes a Discovery Tree data model implementation in Go using Domain-Driven Design (DDD) principles. The Discovery Tree is a hierarchical task management structure that enforces bottom-to-top completion (children must be done before parents) and provides readiness evaluation based on left-to-right sibling ordering.

The implementation follows DDD tactical patterns including Entities, Value Objects, Aggregates, Repositories, and Domain Services to create a rich domain model that encapsulates business logic and maintains invariants.

## Architecture

### Domain-Driven Design Layers

The system follows a layered architecture:

1. **Domain Layer**: Contains the core business logic, entities, value objects, and domain services
2. **Application Layer**: Orchestrates use cases and coordinates domain objects
3. **Infrastructure Layer**: Implements repositories and external concerns (persistence, etc.)

### Aggregate Design

The Discovery Tree follows the Aggregate pattern where:
- **Task** is the Aggregate Root
- Each Task aggregate maintains its own consistency boundary
- Tree-wide operations are coordinated through Domain Services
- The Repository pattern provides access to Task aggregates

## Components and Interfaces

### Domain Model Components

#### 1. Task Entity (Aggregate Root)

The Task entity represents a work item in the discovery tree.

```go
type Task struct {
    id          TaskID
    description string
    status      Status
    parentID    *TaskID  // nil for root
    position    int      // position among siblings (0-indexed)
    createdAt   time.Time
    updatedAt   time.Time
}
```

**Responsibilities:**
- Maintain task identity and basic attributes
- Enforce status transition rules
- Validate description constraints
- Provide access to task properties

**Key Methods:**
- `NewTask(description string, parentID *TaskID, position int) (*Task, error)`
- `ChangeStatus(newStatus Status) error`
- `UpdateDescription(description string) error`
- `Move(newParentID *TaskID, newPosition int) error`
- `ID() TaskID`
- `Status() Status`
- `ParentID() *TaskID`
- `Position() int`

#### 2. TaskID Value Object

Represents a unique identifier for a task.

```go
type TaskID struct {
    value string
}
```

**Responsibilities:**
- Ensure unique identification
- Provide value equality semantics

**Key Methods:**
- `NewTaskID() TaskID`
- `TaskIDFromString(s string) (TaskID, error)`
- `String() string`
- `Equals(other TaskID) bool`

#### 3. Status Value Object

Represents the current state of a task.

```go
type Status int

const (
    StatusTODO Status = iota
    StatusInProgress
    StatusDONE
    StatusBlocked
    StatusRootWorkItem
)
```

**Responsibilities:**
- Enumerate valid status values
- Provide status validation
- Support status comparison

**Key Methods:**
- `NewStatus(s string) (Status, error)`
- `String() string`
- `IsValid() bool`

#### 4. ReadinessState Value Object

Represents whether a task is ready to be worked on based on ordering constraints.

```go
type ReadinessState struct {
    isReady              bool
    leftSiblingComplete  bool
    allChildrenComplete  bool
    reasons              []string
}
```

**Responsibilities:**
- Evaluate task readiness
- Provide reasons for non-readiness
- Support decision-making for task selection

**Key Methods:**
- `IsReady() bool`
- `Reasons() []string`

#### 5. TreeNavigator Domain Service

Provides navigation operations across the tree structure.

```go
type TreeNavigator interface {
    GetParent(taskID TaskID) (*Task, error)
    GetChildren(taskID TaskID) ([]*Task, error)
    GetSiblings(taskID TaskID) ([]*Task, error)
    GetLeftSibling(taskID TaskID) (*Task, error)
    GetRightSibling(taskID TaskID) (*Task, error)
    GetRoot() (*Task, error)
    GetSubtree(taskID TaskID) ([]*Task, error)
}
```

**Responsibilities:**
- Navigate parent-child relationships
- Navigate sibling relationships
- Retrieve tree and subtree structures

#### 6. ReadinessEvaluator Domain Service

Evaluates whether a task is ready to be worked on.

```go
type ReadinessEvaluator interface {
    EvaluateReadiness(taskID TaskID) (ReadinessState, error)
}
```

**Responsibilities:**
- Check left sibling completion status
- Check children completion status
- Compute overall readiness state

#### 7. TaskValidator Domain Service

Validates operations that span multiple tasks or require tree-wide knowledge.

```go
type TaskValidator interface {
    ValidateStatusChange(taskID TaskID, newStatus Status) error
    ValidateMove(taskID TaskID, newParentID *TaskID, newPosition int) error
    ValidateDelete(taskID TaskID) error
}
```

**Responsibilities:**
- Validate bottom-to-top completion constraint
- Validate move operations (prevent cycles)
- Validate deletion operations

#### 8. TaskRepository Interface

Provides persistence operations for Task aggregates.

```go
type TaskRepository interface {
    Save(task *Task) error
    FindByID(id TaskID) (*Task, error)
    FindByParentID(parentID *TaskID) ([]*Task, error)
    FindRoot() (*Task, error)
    FindAll() ([]*Task, error)
    Delete(id TaskID) error
    DeleteSubtree(id TaskID) error
}
```

**Responsibilities:**
- Persist and retrieve Task aggregates
- Query tasks by relationships
- Handle deletion operations

## Data Models

### Task Aggregate

The Task is the primary aggregate in the system. Each task maintains:

- **Identity**: Unique TaskID
- **Attributes**: Description, Status, timestamps
- **Relationships**: ParentID (optional), Position among siblings
- **Invariants**:
  - Description must not be empty
  - Position must be non-negative
  - Root tasks have nil ParentID
  - Status must be valid

### Tree Structure

The tree is represented through parent-child relationships:

- Each task (except root) has exactly one parent
- Each task has zero or more children
- Siblings are ordered by position (0-indexed)
- The tree has exactly one root task

### Persistence Model

For persistence, tasks can be stored as a flat collection with parent references:

```
tasks table:
- id (primary key)
- description
- status
- parent_id (nullable, foreign key to tasks.id)
- position
- created_at
- updated_at
```

Indexes:
- Primary key on id
- Index on parent_id for efficient child queries
- Unique constraint on (parent_id, position) to prevent duplicate positions


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Task Creation Properties

**Property 1: Root task creation with valid description**
*For any* valid non-empty description, creating a root task should result in a task with status "Root Work Item", no parent reference, and a unique identifier.
**Validates: Requirements 1.1, 1.3**

**Property 2: Invalid description rejection**
*For any* empty or whitespace-only description, attempting to create a root task should be rejected with an error.
**Validates: Requirements 1.2**

**Property 3: Single root constraint**
*For any* tree with an existing root task, attempting to create another root task should be rejected with an error.
**Validates: Requirements 1.4**

### Child Task Properties

**Property 4: Child task parent reference**
*For any* existing parent task and valid child description, creating a child task should result in a task with the correct parent reference.
**Validates: Requirements 2.1**

**Property 5: Child task position assignment**
*For any* parent task with N children, adding a new child should assign it position N (0-indexed).
**Validates: Requirements 2.2**

**Property 6: Multiple children ordering**
*For any* parent task, adding multiple children in sequence should result in children with sequential positions (0, 1, 2, ...).
**Validates: Requirements 2.3**

**Property 7: Non-existent parent rejection**
*For any* non-existent parent ID, attempting to create a child task should be rejected with an error.
**Validates: Requirements 2.4**

### Status Update Properties

**Property 8: Valid status update persistence**
*For any* task and valid status value, updating the task's status should persist the new status value.
**Validates: Requirements 3.1**

**Property 9: Invalid status rejection**
*For any* task and invalid status value, attempting to update the status should be rejected with an error.
**Validates: Requirements 3.2**

**Property 10: DONE status recording**
*For any* task that can be marked DONE (no incomplete children), marking it as DONE should result in the status being recorded as DONE.
**Validates: Requirements 3.3**

### Readiness Evaluation Properties

**Property 11: Left sibling readiness check**
*For any* task, evaluating its readiness should correctly identify whether its left sibling is DONE or does not exist.
**Validates: Requirements 4.1**

**Property 12: Children completion readiness check**
*For any* task, evaluating its readiness should correctly identify whether all its children are DONE or it has no children.
**Validates: Requirements 4.2**

**Property 13: Readiness does not block status changes**
*For any* task regardless of its readiness state, changing its status to any valid value should be allowed.
**Validates: Requirements 4.5**

### Bottom-to-Top Enforcement Properties

**Property 14: Parent cannot be DONE before children**
*For any* task with at least one child that is not DONE, attempting to mark the task as DONE should be rejected with an error.
**Validates: Requirements 5.1**

**Property 15: Non-DONE statuses ignore children**
*For any* task with incomplete children, marking it as "In Progress", "TODO", or "Blocked" should be allowed.
**Validates: Requirements 5.3**

### Tree Retrieval Properties

**Property 16: Complete tree retrieval**
*For any* tree, requesting the tree structure should return the root task and all its descendants.
**Validates: Requirements 6.1**

**Property 17: Tree structure preservation**
*For any* tree, the returned tree structure should preserve all parent-child relationships and maintain left-to-right sibling ordering by position.
**Validates: Requirements 6.2, 6.3**

**Property 18: Subtree retrieval**
*For any* task in a tree, requesting its subtree should return that task and all its descendants.
**Validates: Requirements 6.4**

### Navigation Properties

**Property 19: Parent navigation**
*For any* task, requesting its parent should return the parent task if it exists, or null if the task is root.
**Validates: Requirements 7.1**

**Property 20: Children navigation with ordering**
*For any* task, requesting its children should return all child tasks ordered by position (left-to-right).
**Validates: Requirements 7.2**

**Property 21: Left sibling navigation**
*For any* task, requesting its left sibling should return the task with position-1 under the same parent, or null if it's the leftmost sibling.
**Validates: Requirements 7.4**

**Property 22: Right sibling navigation**
*For any* task, requesting its right sibling should return the task with position+1 under the same parent, or null if it's the rightmost sibling.
**Validates: Requirements 7.5**

### Move Operation Properties

**Property 23: Task move updates parent and positions**
*For any* task moved to a new parent, the task's parent reference should be updated and sibling positions should be adjusted in both old and new parent contexts.
**Validates: Requirements 8.1**

**Property 24: Sibling reordering maintains ordering**
*For any* set of sibling tasks that are reordered, their positions should be updated to maintain sequential left-to-right ordering.
**Validates: Requirements 8.2**

**Property 25: Move preserves subtree**
*For any* task with descendants, moving the task should move all its descendants with it, preserving the subtree structure.
**Validates: Requirements 8.3**

**Property 26: Cycle prevention in moves**
*For any* task, attempting to move it to become a descendant of itself (creating a cycle) should be rejected with an error.
**Validates: Requirements 8.4**

### Delete Operation Properties

**Property 27: Leaf deletion adjusts positions**
*For any* leaf task (no children), deleting it should remove the task and adjust the positions of its right siblings.
**Validates: Requirements 9.1**

**Property 28: Cascading deletion**
*For any* task with children, deleting it should remove the task and all its descendants.
**Validates: Requirements 9.2**

## Error Handling

### Error Types

The system defines specific error types for different failure scenarios:

1. **ValidationError**: Invalid input or constraint violation
   - Empty description
   - Invalid status value
   - Invalid position

2. **NotFoundError**: Referenced entity does not exist
   - Task not found
   - Parent not found

3. **ConstraintViolationError**: Business rule violation
   - Multiple root tasks
   - Parent marked DONE before children
   - Cycle creation in move operation

4. **ConcurrencyError**: Concurrent modification conflict
   - Optimistic locking failure

### Error Handling Strategy

- All domain operations return errors for invalid states
- Errors are propagated up through application layer
- Repository operations handle persistence errors
- Domain services validate cross-aggregate constraints before operations

### Validation Points

1. **Entity Creation**: Validate all required fields and constraints
2. **Status Changes**: Validate business rules (bottom-to-top enforcement)
3. **Move Operations**: Validate no cycles, valid parent exists
4. **Delete Operations**: Validate task exists

## Testing Strategy

### Unit Testing

Unit tests will verify specific examples and edge cases:

1. **Entity Tests**:
   - Task creation with various inputs
   - Status transitions
   - Edge cases: empty descriptions, nil parents

2. **Value Object Tests**:
   - TaskID generation and equality
   - Status validation
   - ReadinessState computation

3. **Domain Service Tests**:
   - Navigation with various tree structures
   - Readiness evaluation with different configurations
   - Validation logic for constraints

### Property-Based Testing

Property-based testing will verify universal properties across all inputs using **gopter** (a Go property-based testing library).

**Configuration**:
- Each property test should run a minimum of 100 iterations
- Tests should use custom generators for valid tree structures
- Each property test must be tagged with a comment referencing the design document property

**Tag Format**:
```go
// Feature: discovery-tree-model, Property 1: Root task creation with valid description
```

**Property Test Coverage**:
- Each correctness property (1-28) will be implemented as a single property-based test
- Generators will create random but valid tree structures
- Properties will verify invariants hold across all generated inputs

**Test Organization**:
- Property tests will be co-located with unit tests
- Generators will be defined in a separate `generators_test.go` file
- Each property test will focus on a single correctness property

### Integration Testing

Integration tests will verify:
- Repository implementations with actual persistence
- End-to-end workflows (create tree, navigate, modify, delete)
- Concurrent access scenarios

### Test Data Generators

For property-based testing, we need generators for:

1. **Valid Descriptions**: Non-empty strings
2. **Task Trees**: Valid tree structures with proper parent-child relationships
3. **Status Values**: Valid status enumerations
4. **Positions**: Valid position values for siblings

## Implementation Notes

### Concurrency Considerations

- Task aggregates are independently modifiable
- Repository implementations should use optimistic locking
- Tree-wide operations (like move) may require transaction support

### Performance Considerations

- Navigation operations should be efficient (indexed queries)
- Subtree operations may need optimization for deep trees
- Position adjustments on delete/move should be batched

### Future Extensions

Potential future enhancements:
- Task metadata (assignees, due dates, tags)
- Task history and audit trail
- Bulk operations (move multiple tasks)
- Tree visualization and export
- Search and filtering capabilities
