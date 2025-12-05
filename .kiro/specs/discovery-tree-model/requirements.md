# Requirements Document

## Introduction

This document specifies the requirements for a Discovery Tree data model implementation in Go using Domain-Driven Design (DDD) principles. A Discovery Tree is a hierarchical task management structure where tasks are organized in a tree with specific ordering constraints. The tree structure resembles a B+ tree where tasks must be completed left-to-right at each level, ensuring dependencies are respected.

## Glossary

- **Discovery Tree**: A hierarchical tree structure for organizing and tracking work items with left-to-right completion ordering
- **Task**: A work item node in the Discovery Tree that can be a Root, High-Level Task, or any other granularity of work
- **Root Task**: The top-level goal or objective that serves as the tree's root node
- **Status**: The current state of a Task (TODO, In Progress, DONE, Blocked, Root Work Item)
- **Parent Task**: A Task that has one or more child Tasks beneath it in the hierarchy
- **Child Task**: A Task that belongs to a Parent Task
- **Left-to-Right Ordering**: A guideline where sibling Tasks should ideally be completed in order from left to right
- **Readiness State**: An evaluation of whether a Task is ready to be worked on based on its position and dependencies
- **System**: The Discovery Tree management system

## Requirements

### Requirement 1

**User Story:** As a user, I want to create a root task for my discovery tree, so that I can establish the top-level goal or objective.

#### Acceptance Criteria

1. WHEN a user creates a root task with a simple story(it may consists of a few words), THEN the System SHALL create a new Task with status "Root Work Item" and no parent
2. WHEN a user attempts to create a root task with an empty story, THEN the System SHALL reject the creation and return an error
3. WHEN a user creates a root task, THEN the System SHALL assign it a unique identifier
4. WHEN a root task exists, THEN the System SHALL prevent creation of additional root tasks in the same tree

### Requirement 2

**User Story:** As a user, I want to add child tasks under a parent task, so that I can break down work into smaller, manageable pieces.

#### Acceptance Criteria

1. WHEN a user adds a child task to an existing parent task, THEN the System SHALL create the child task with the specified parent reference
2. WHEN a user adds a child task, THEN the System SHALL assign it a position relative to existing siblings
3. WHEN a user adds multiple child tasks to the same parent, THEN the System SHALL maintain their left-to-right ordering
4. WHEN a user attempts to add a child task to a non-existent parent, THEN the System SHALL reject the operation and return an error

### Requirement 3

**User Story:** As a user, I want to update the status of a task, so that I can track progress through the workflow.

#### Acceptance Criteria

1. WHEN a user updates a task status to a valid status value, THEN the System SHALL persist the new status
2. WHEN a user attempts to update a task status to an invalid value, THEN the System SHALL reject the update and return an error
3. WHEN a user marks a task as "DONE", THEN the System SHALL record the completion
4. WHEN a user updates a task status, THEN the System SHALL validate the status transition is allowed

### Requirement 4

**User Story:** As a user, I want the system to evaluate task readiness based on ordering constraints, so that I can understand which tasks should ideally be worked on next while maintaining flexibility.

#### Acceptance Criteria

1. WHEN a user requests the readiness state of a task, THEN the System SHALL evaluate whether its left sibling is "DONE" or does not exist
2. WHEN a user requests the readiness state of a task, THEN the System SHALL evaluate whether all its children are "DONE" or it has no children
3. WHEN a task has no left sibling or its left sibling is "DONE", and all children are "DONE" or it has no children, THEN the System SHALL indicate the task is ready
4. WHEN a task has an incomplete left sibling or incomplete children, THEN the System SHALL indicate the task is not ready but still allow status changes
5. WHEN a user changes a task status to any value, THEN the System SHALL allow the change regardless of readiness state

### Requirement 5

**User Story:** As a user, I want the system to enforce bottom-to-top completion, so that parent tasks cannot be marked done before their children.

#### Acceptance Criteria

1. WHEN a user attempts to mark a task as "DONE" and any of its children are not "DONE", THEN the System SHALL reject the status change and return an error
2. WHEN a task has no children or all children are "DONE", THEN the System SHALL allow the task to be marked as "DONE"
3. WHEN a user marks a task as "In Progress", "TODO", or "Blocked", THEN the System SHALL allow the status change regardless of children status

### Requirement 6

**User Story:** As a user, I want to retrieve the tree structure, so that I can visualize the hierarchy and current state of all tasks.

#### Acceptance Criteria

1. WHEN a user requests the tree structure, THEN the System SHALL return the root task with all descendant tasks
2. WHEN the System returns the tree structure, THEN it SHALL preserve the parent-child relationships
3. WHEN the System returns the tree structure, THEN it SHALL maintain the left-to-right ordering of siblings
4. WHEN a user requests a subtree starting from a specific task, THEN the System SHALL return that task and all its descendants

### Requirement 7

**User Story:** As a user, I want to navigate the tree structure, so that I can find specific tasks and understand relationships.

#### Acceptance Criteria

1. WHEN a user requests a task's parent, THEN the System SHALL return the parent task or null if the task is root
2. WHEN a user requests a task's children, THEN the System SHALL return all child tasks in left-to-right order
3. WHEN a user requests a task's siblings, THEN the System SHALL return all tasks sharing the same parent in order
4. WHEN a user requests a task's left sibling, THEN the System SHALL return the immediate left sibling or null if none exists
5. WHEN a user requests a task's right sibling, THEN the System SHALL return the immediate right sibling or null if none exists

### Requirement 8

**User Story:** As a user, I want to move tasks within the tree, so that I can reorganize work as understanding evolves.

#### Acceptance Criteria

1. WHEN a user moves a task to a new parent, THEN the System SHALL update the parent reference and adjust sibling positions
2. WHEN a user reorders siblings, THEN the System SHALL update positions while maintaining left-to-right ordering
3. WHEN a user moves a task, THEN the System SHALL move all descendant tasks with it
4. WHEN a user attempts to move a task to become its own descendant, THEN the System SHALL reject the operation and return an error

### Requirement 8

**User Story:** As a user, I want to delete tasks from the tree, so that I can remove work that is no longer relevant.

#### Acceptance Criteria

1. WHEN a user deletes a task with no children, THEN the System SHALL remove the task and adjust sibling positions
2. WHEN a user deletes a task with children, THEN the System SHALL remove the task and all its descendants
3. WHEN a user deletes a task, THEN the System SHALL update the positions of remaining siblings to maintain ordering
4. WHEN a user attempts to delete the root task, THEN the System SHALL remove the entire tree

### Requirement 9

**User Story:** As a developer, I want the domain model to follow DDD principles, so that the codebase is maintainable and expresses business logic clearly.

#### Acceptance Criteria

1. WHEN implementing the Task entity, THEN the System SHALL encapsulate business rules within the entity
2. WHEN implementing status transitions, THEN the System SHALL use value objects for Status
3. WHEN implementing tree operations, THEN the System SHALL use domain services for complex operations spanning multiple entities
4. WHEN persisting the model, THEN the System SHALL use repository interfaces to abstract data access
5. WHEN validating operations, THEN the System SHALL enforce invariants within the domain model
