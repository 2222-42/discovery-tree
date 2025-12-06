# Implementation Plan

- [x] 1. Set up project structure and domain foundation
  - Create Go module and directory structure (domain, application, infrastructure)
  - Set up testing framework with gopter for property-based testing
  - Create base error types (ValidationError, NotFoundError, ConstraintViolationError)
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [x] 2. Implement TaskID value object
  - Create TaskID struct with string value
  - Implement NewTaskID() for generation
  - Implement TaskIDFromString() with validation
  - Implement String() and Equals() methods
  - _Requirements: 1.3_

- [ ]* 2.1 Write property test for TaskID uniqueness
  - **Property 1: Root task creation with valid description**
  - **Validates: Requirements 1.1, 1.3**

- [x] 3. Implement Status value object
  - Create Status type with constants (TODO, InProgress, DONE, Blocked, RootWorkItem)
  - Implement NewStatus() with validation
  - Implement String() and IsValid() methods
  - _Requirements: 3.1, 3.2_

- [ ]* 3.1 Write property test for status validation
  - **Property 9: Invalid status rejection**
  - **Validates: Requirements 3.2**

- [x] 4. Implement Task entity (aggregate root)
  - Create Task struct with all fields (id, description, status, parentID, position, timestamps)
  - Implement NewTask() constructor with validation
  - Implement getter methods (ID, Status, ParentID, Position, Description)
  - Enforce description validation (non-empty)
  - _Requirements: 1.1, 1.2, 2.1, 10.1_

- [ ]* 4.1 Write property test for task creation with valid descriptions
  - **Property 1: Root task creation with valid description**
  - **Validates: Requirements 1.1, 1.3**

- [ ]* 4.2 Write property test for invalid description rejection
  - **Property 2: Invalid description rejection**
  - **Validates: Requirements 1.2**

- [x] 5. Implement basic status change functionality
  - Add ChangeStatus() method to Task entity
  - Implement basic status validation
  - Add UpdateDescription() method with validation
  - _Requirements: 3.1, 3.2, 3.3_

- [ ]* 5.1 Write property test for valid status updates
  - **Property 8: Valid status update persistence**
  - **Validates: Requirements 3.1**

- [ ]* 5.2 Write property test for DONE status recording
  - **Property 10: DONE status recording**
  - **Validates: Requirements 3.3**

- [x] 6. Implement TaskRepository interface
  - Define repository interface with all methods (Save, FindByID, FindByParentID, FindRoot, FindAll, Delete, DeleteSubtree)
  - Create in-memory implementation for testing
  - _Requirements: 10.4_

- [x] 7. Implement child task creation logic
  - Add logic to calculate position for new children
  - Implement validation for parent existence
  - Add support for creating child tasks with proper parent references
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [ ]* 7.1 Write property test for child parent reference
  - **Property 4: Child task parent reference**
  - **Validates: Requirements 2.1**

- [ ]* 7.2 Write property test for position assignment
  - **Property 5: Child task position assignment**
  - **Validates: Requirements 2.2**

- [ ]* 7.3 Write property test for multiple children ordering
  - **Property 6: Multiple children ordering**
  - **Validates: Requirements 2.3**

- [ ]* 7.4 Write property test for non-existent parent rejection
  - **Property 7: Non-existent parent rejection**
  - **Validates: Requirements 2.4**

- [x] 8. Implement single root constraint
  - Add validation to prevent multiple root tasks
  - Update repository to check for existing root
  - _Requirements: 1.4_

- [ ]* 8.1 Write property test for single root constraint
  - **Property 3: Single root constraint**
  - **Validates: Requirements 1.4**

- [x] 9. Implement TreeNavigator domain service
  - Create TreeNavigator interface
  - Implement GetParent() method
  - Implement GetChildren() with ordering by position
  - Implement GetSiblings() method
  - Implement GetLeftSibling() and GetRightSibling() methods
  - Implement GetRoot() method
  - Implement GetSubtree() method
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ]* 9.1 Write property test for parent navigation
  - **Property 19: Parent navigation**
  - **Validates: Requirements 7.1**

- [ ]* 9.2 Write property test for children navigation
  - **Property 20: Children navigation with ordering**
  - **Validates: Requirements 7.2**

- [ ]* 9.3 Write property test for left sibling navigation
  - **Property 21: Left sibling navigation**
  - **Validates: Requirements 7.4**

- [ ]* 9.4 Write property test for right sibling navigation
  - **Property 22: Right sibling navigation**
  - **Validates: Requirements 7.5**

- [x] 10. Implement ReadinessState value object
  - Create ReadinessState struct with fields (isReady, leftSiblingComplete, allChildrenComplete, reasons)
  - Implement IsReady() and Reasons() methods
  - _Requirements: 4.1, 4.2_

- [x] 11. Implement ReadinessEvaluator domain service
  - Create ReadinessEvaluator interface
  - Implement EvaluateReadiness() method
  - Check left sibling completion status
  - Check all children completion status
  - Compute overall readiness and reasons
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [ ]* 11.1 Write property test for left sibling readiness check
  - **Property 11: Left sibling readiness check**
  - **Validates: Requirements 4.1**

- [ ]* 11.2 Write property test for children completion check
  - **Property 12: Children completion readiness check**
  - **Validates: Requirements 4.2**

- [ ]* 11.3 Write property test for readiness not blocking status changes
  - **Property 13: Readiness does not block status changes**
  - **Validates: Requirements 4.5**

- [x] 12. Implement TaskValidator domain service for bottom-to-top enforcement
  - Create TaskValidator interface
  - Implement ValidateStatusChange() method
  - Add validation for DONE status requiring all children to be DONE
  - Allow non-DONE statuses regardless of children
  - _Requirements: 5.1, 5.2, 5.3_

- [ ]* 12.1 Write property test for parent-children DONE constraint
  - **Property 14: Parent cannot be DONE before children**
  - **Validates: Requirements 5.1**

- [ ]* 12.2 Write property test for non-DONE statuses ignoring children
  - **Property 15: Non-DONE statuses ignore children**
  - **Validates: Requirements 5.3**

- [x] 13. Integrate validation into Task.ChangeStatus()
  - Update ChangeStatus() to use TaskValidator
  - Ensure bottom-to-top enforcement is applied
  - Maintain flexibility for non-DONE statuses
  - _Requirements: 5.1, 5.3_

- [x] 14. Implement tree retrieval operations
  - Add GetTree() method to retrieve complete tree structure
  - Ensure parent-child relationships are preserved
  - Ensure left-to-right ordering is maintained
  - Add GetSubtree() support for partial tree retrieval
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [ ]* 14.1 Write property test for complete tree retrieval
  - **Property 16: Complete tree retrieval**
  - **Validates: Requirements 6.1**

- [ ]* 14.2 Write property test for tree structure preservation
  - **Property 17: Tree structure preservation**
  - **Validates: Requirements 6.2, 6.3**

- [ ]* 14.3 Write property test for subtree retrieval
  - **Property 18: Subtree retrieval**
  - **Validates: Requirements 6.4**

- [x] 15. Implement move operations
  - Add Move() method to Task entity
  - Implement position adjustment logic for old and new parents
  - Update TaskValidator to include ValidateMove() for cycle detection
  - Ensure subtree moves with parent task
  - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ]* 15.1 Write property test for move updates
  - **Property 23: Task move updates parent and positions**
  - **Validates: Requirements 8.1**

- [ ]* 15.2 Write property test for reordering
  - **Property 24: Sibling reordering maintains ordering**
  - **Validates: Requirements 8.2**

- [ ]* 15.3 Write property test for subtree preservation
  - **Property 25: Move preserves subtree**
  - **Validates: Requirements 8.3**

- [ ]* 15.4 Write property test for cycle prevention
  - **Property 26: Cycle prevention in moves**
  - **Validates: Requirements 8.4**

- [x] 16. Implement delete operations
  - Implement Delete() in repository for leaf tasks
  - Implement DeleteSubtree() for cascading deletion
  - Add position adjustment logic for remaining siblings
  - Handle root deletion (entire tree removal)
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [ ]* 16.1 Write property test for leaf deletion
  - **Property 27: Leaf deletion adjusts positions**
  - **Validates: Requirements 9.1**

- [ ]* 16.2 Write property test for cascading deletion
  - **Property 28: Cascading deletion**
  - **Validates: Requirements 9.2**

- [x] 17. Create property test generators
  - Implement generator for valid descriptions
  - Implement generator for valid tree structures
  - Implement generator for status values
  - Implement generator for positions
  - Ensure generators create diverse test cases
  - _Requirements: All property tests_

- [ ] 18. Final checkpoint - Ensure all tests pass
  - Run all unit tests and property tests
  - Verify all 28 correctness properties are validated
  - Ensure test coverage is comprehensive
  - Ask the user if questions arise
