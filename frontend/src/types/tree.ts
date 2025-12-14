/**
 * Tree-related type definitions
 * Defines the hierarchical structure for displaying tasks in a tree format
 */

import type { Task } from './task.js';

/**
 * TreeNode represents a task within the hierarchical tree structure
 * Used for rendering and managing the tree view component
 */
export interface TreeNode {
  /** The task data associated with this tree node */
  readonly task: Task;
  /** Child nodes in the tree hierarchy */
  children: TreeNode[];
  /** Whether this node is currently expanded in the UI */
  isExpanded: boolean;
  /** Depth level in the tree (0 for root) */
  readonly level: number;
}

/**
 * Tree navigation path representing the route to a specific node
 */
export type TreePath = readonly string[];

/**
 * Tree operation result for validation
 */
export interface TreeOperationResult {
  /** Whether the operation is valid */
  isValid: boolean;
  /** Error message if operation is invalid */
  error?: string;
}

/**
 * Tree search result
 */
export interface TreeSearchResult {
  /** The found tree node */
  node: TreeNode;
  /** Path to the node from root */
  path: TreePath;
}

/**
 * Tree state for managing UI interactions
 */
export interface TreeState {
  /** Set of expanded node IDs */
  expandedNodes: Set<string>;
  /** Currently selected node ID */
  selectedNodeId: string | null;
  /** The root tree nodes */
  rootNodes: TreeNode[];
}

/**
 * Tree drag and drop operation
 */
export interface TreeDragOperation {
  /** ID of the task being dragged */
  draggedTaskId: string;
  /** ID of the target parent (null for root level) */
  targetParentId: string | null;
  /** Target position within the parent's children */
  targetPosition: number;
}