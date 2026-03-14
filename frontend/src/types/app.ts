/**
 * Application-wide type definitions
 * Defines global state, context types, and application-level interfaces
 */

import type { Task } from './task.js';
import type { TreeNode, TreeState } from './tree.js';

/**
 * Global application state
 */
export interface AppState {
  /** All tasks loaded from the API */
  tasks: Task[];
  /** Currently selected task ID */
  selectedTaskId: string | null;
  /** Set of expanded tree node IDs */
  expandedNodes: Set<string>;
  /** Loading state for async operations */
  loading: boolean;
  /** Current error message, if any */
  error: string | null;
  /** Computed tree data from tasks */
  treeData: TreeNode[];
}

/**
 * Application actions for state management
 */
export type AppAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_TASKS'; payload: Task[] }
  | { type: 'ADD_TASK'; payload: Task }
  | { type: 'UPDATE_TASK'; payload: Task }
  | { type: 'DELETE_TASK'; payload: string }
  | { type: 'SELECT_TASK'; payload: string | null }
  | { type: 'TOGGLE_NODE'; payload: string }
  | { type: 'EXPAND_NODE'; payload: string }
  | { type: 'COLLAPSE_NODE'; payload: string }
  | { type: 'SET_TREE_DATA'; payload: TreeNode[] };

/**
 * Task context state and actions
 */
export interface TaskContextState {
  tasks: Task[];
  selectedTask: Task | null;
  loading: boolean;
  error: string | null;
}

export interface TaskContextActions {
  createRootTask: (description: string) => Promise<void>;
  createChildTask: (description: string, parentId: string) => Promise<void>;
  updateTask: (id: string, description: string) => Promise<void>;
  updateTaskStatus: (id: string, status: import('./task.js').TaskStatus) => Promise<void>;
  deleteTask: (id: string) => Promise<void>;
  selectTask: (id: string | null) => void;
  refreshTasks: () => Promise<void>;
}

export interface TaskContextValue extends TaskContextState, TaskContextActions {}

/**
 * Tree context state and actions
 */
export type TreeContextState = TreeState;

export interface TreeContextActions {
  toggleNode: (nodeId: string) => void;
  expandNode: (nodeId: string) => void;
  collapseNode: (nodeId: string) => void;
  selectNode: (nodeId: string | null) => void;
  moveTask: (taskId: string, parentId: string | null, position: number) => Promise<void>;
  startInlineCreation: (parentId: string) => void;
  cancelInlineCreation: () => void;
  updateInlineDescription: (description: string) => void;
  completeInlineCreation: () => Promise<void>;
  startDrag: (taskId: string) => void;
  endDrag: () => void;
  setDragOver: (taskId: string | null, position: 'before' | 'after' | 'child' | null) => void;
  handleDrop: (targetTaskId: string, position: 'before' | 'after' | 'child') => Promise<void>;
}

export interface TreeContextValue extends TreeContextState, TreeContextActions {
  // Combines TreeContextState and TreeContextActions
}

/**
 * Loading states for different operations
 */
export interface LoadingState {
  tasks: boolean;
  creating: boolean;
  updating: boolean;
  deleting: boolean;
  moving: boolean;
}

/**
 * Form validation result
 */
export interface ValidationResult {
  isValid: boolean;
  errors: Record<string, string>;
}

/**
 * Component props for common patterns
 */
export interface BaseComponentProps {
  className?: string;
  'data-testid'?: string;
}

/**
 * Error boundary state
 */
export interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
  errorInfo?: React.ErrorInfo;
}