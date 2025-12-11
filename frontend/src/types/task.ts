/**
 * Task-related type definitions
 * Based on the backend API Task model and requirements 1.2, 1.4
 */

/**
 * Task status enumeration matching the backend API
 */
export const TaskStatus = {
  TODO: 'TODO',
  IN_PROGRESS: 'IN_PROGRESS',
  DONE: 'DONE',
  ROOT_WORK_ITEM: 'ROOT_WORK_ITEM'
} as const;

export type TaskStatus = typeof TaskStatus[keyof typeof TaskStatus];

/**
 * Core Task interface representing a task in the discovery tree
 * This matches the backend API Task model structure
 */
export interface Task {
  /** Unique identifier for the task */
  readonly id: string;
  /** Human-readable description of the task */
  description: string;
  /** Current status of the task */
  status: TaskStatus;
  /** ID of the parent task, null for root tasks */
  parentId: string | null;
  /** Position within the parent's children list */
  position: number;
  /** ISO timestamp when the task was created */
  readonly createdAt: string;
  /** ISO timestamp when the task was last updated */
  readonly updatedAt: string;
}

/**
 * Task creation request payload
 */
export interface CreateTaskRequest {
  description: string;
  parentId?: string;
}