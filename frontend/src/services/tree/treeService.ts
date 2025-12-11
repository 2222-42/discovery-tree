/**
 * Tree Service for building and managing hierarchical task structures
 * Implements tree navigation, validation, and cycle detection utilities
 * Requirements: 3.2, 4.1
 */

import type { Task } from '../../types/task.js';
import type { TreeNode, TreePath, TreeOperationResult, TreeSearchResult, TreeDragOperation } from '../../types/tree.js';

/**
 * TreeService class for managing hierarchical task structures
 * Provides utilities for building trees from flat task lists and navigation
 */
export class TreeService {
  /**
   * Builds a tree structure from a flat list of tasks
   * @param tasks - Flat array of tasks to convert to tree structure
   * @returns Array of root TreeNode objects representing the tree
   */
  buildTreeFromTasks(tasks: Task[]): TreeNode[] {
    if (tasks.length === 0) {
      return [];
    }

    // Create a map for quick task lookup
    const taskMap = new Map<string, Task>();
    tasks.forEach(task => {
      taskMap.set(task.id, task);
    });

    // Create tree nodes map with temporary level tracking
    const nodeMap = new Map<string, { node: TreeNode; tempLevel: number }>();
    
    // First pass: create all nodes with temporary level 0
    tasks.forEach(task => {
      const treeNode: TreeNode = {
        task,
        children: [],
        isExpanded: false,
        level: 0 // Will be updated after hierarchy is built
      };
      nodeMap.set(task.id, { node: treeNode, tempLevel: 0 });
    });

    // Second pass: build parent-child relationships
    const rootNodes: TreeNode[] = [];
    
    for (const task of tasks) {
      const nodeData = nodeMap.get(task.id);
      if (!nodeData) {
        continue; // Skip if node not found (defensive programming)
      }
      
      if (task.parentId === null) {
        // Root node
        nodeData.tempLevel = 0;
        rootNodes.push(nodeData.node);
      } else {
        // Child node
        const parentNodeData = nodeMap.get(task.parentId);
        if (parentNodeData) {
          nodeData.tempLevel = parentNodeData.tempLevel + 1;
          parentNodeData.node.children.push(nodeData.node);
        } else {
          // Parent not found, treat as root (defensive programming)
          nodeData.tempLevel = 0;
          rootNodes.push(nodeData.node);
        }
      }
    }

    // Third pass: update the actual level property by recreating nodes
    const updateNodeLevels = (nodes: TreeNode[], level: number): TreeNode[] => {
      return nodes.map(node => {
        const updatedNode: TreeNode = {
          ...node,
          level,
          children: updateNodeLevels(node.children, level + 1)
        };
        return updatedNode;
      });
    };

    const finalRootNodes = updateNodeLevels(rootNodes, 0);

    // Sort children by position within each parent
    this.sortChildrenByPosition(finalRootNodes);
    
    return finalRootNodes;
  }

  /**
   * Recursively sorts children by their position property
   * @param nodes - Array of tree nodes to sort
   */
  private sortChildrenByPosition(nodes: TreeNode[]): void {
    nodes.forEach(node => {
      if (node.children.length > 0) {
        node.children.sort((a, b) => a.task.position - b.task.position);
        this.sortChildrenByPosition(node.children);
      }
    });
  }

  /**
   * Finds a task node in the tree by task ID
   * @param tree - Array of root tree nodes to search
   * @param taskId - ID of the task to find
   * @returns TreeSearchResult with node and path, or null if not found
   */
  findTaskInTree(tree: TreeNode[], taskId: string): TreeSearchResult | null {
    const path: string[] = [];
    
    const searchRecursive = (nodes: TreeNode[], currentPath: string[]): TreeNode | null => {
      for (const node of nodes) {
        const newPath = [...currentPath, node.task.id];
        
        if (node.task.id === taskId) {
          path.push(...newPath);
          return node;
        }
        
        const found = searchRecursive(node.children, newPath);
        if (found) {
          return found;
        }
      }
      return null;
    };

    const foundNode = searchRecursive(tree, []);
    
    if (foundNode) {
      return {
        node: foundNode,
        path: path as TreePath
      };
    }
    
    return null;
  }

  /**
   * Gets the path from root to a specific task
   * @param tree - Array of root tree nodes
   * @param taskId - ID of the target task
   * @returns Array of task IDs representing the path from root to target
   */
  getTaskPath(tree: TreeNode[], taskId: string): TreePath {
    const result = this.findTaskInTree(tree, taskId);
    return result ? result.path : [];
  }

  /**
   * Validates if a move operation would create a cycle in the tree
   * @param tree - Current tree structure
   * @param taskId - ID of task being moved
   * @param newParentId - ID of new parent (null for root level)
   * @returns TreeOperationResult indicating if move is valid
   */
  validateMove(tree: TreeNode[], taskId: string, newParentId: string | null | undefined): TreeOperationResult {
    // Cannot move a task to be its own parent
    if (taskId === newParentId) {
      return {
        isValid: false,
        error: 'Cannot move a task to be its own parent'
      };
    }

    // If moving to root level, it's always valid (no cycle possible)
    if (newParentId === null || newParentId === undefined) {
      return { isValid: true };
    }

    // Check if the new parent is a descendant of the task being moved
    const taskNode = this.findTaskInTree(tree, taskId);
    if (!taskNode) {
      return {
        isValid: false,
        error: 'Task not found in tree'
      };
    }

    // Check if newParentId is a descendant of taskId
    const isDescendant = this.isDescendant(taskNode.node, newParentId);
    
    if (isDescendant) {
      return {
        isValid: false,
        error: 'Cannot move a task under one of its descendants (would create a cycle)'
      };
    }

    return { isValid: true };
  }

  /**
   * Checks if a given task ID is a descendant of a tree node
   * @param node - The potential ancestor node
   * @param taskId - ID of the potential descendant task
   * @returns true if taskId is a descendant of node
   */
  private isDescendant(node: TreeNode, taskId: string): boolean {
    for (const child of node.children) {
      if (child.task.id === taskId) {
        return true;
      }
      if (this.isDescendant(child, taskId)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Validates the entire tree structure for cycles and consistency
   * @param tasks - Flat array of tasks to validate
   * @returns TreeOperationResult indicating if tree structure is valid
   */
  validateTreeStructure(tasks: Task[]): TreeOperationResult {
    if (tasks.length === 0) {
      return { isValid: true };
    }

    const taskMap = new Map<string, Task>();
    const visited = new Set<string>();
    const recursionStack = new Set<string>();

    // Build task map
    tasks.forEach(task => {
      taskMap.set(task.id, task);
    });

    // Check for cycles using DFS
    const hasCycle = (taskId: string): boolean => {
      if (recursionStack.has(taskId)) {
        return true; // Cycle detected
      }
      
      if (visited.has(taskId)) {
        return false; // Already processed
      }

      visited.add(taskId);
      recursionStack.add(taskId);

      const task = taskMap.get(taskId);
      if (task !== undefined && task.parentId !== null) {
        if (hasCycle(task.parentId)) {
          return true;
        }
      }

      recursionStack.delete(taskId);
      return false;
    };

    // Check each task for cycles
    for (const task of tasks) {
      if (!visited.has(task.id)) {
        if (hasCycle(task.id)) {
          return {
            isValid: false,
            error: `Cycle detected in task hierarchy involving task ${task.id}`
          };
        }
      }
    }

    // Validate parent references exist
    for (const task of tasks) {
      if (task.parentId !== null && !taskMap.has(task.parentId)) {
        return {
          isValid: false,
          error: `Task ${task.id} references non-existent parent ${task.parentId}`
        };
      }
    }

    return { isValid: true };
  }

  /**
   * Validates a drag and drop operation
   * @param tree - Current tree structure
   * @param operation - The drag operation to validate
   * @returns TreeOperationResult indicating if operation is valid
   */
  validateDragOperation(tree: TreeNode[], operation: TreeDragOperation): TreeOperationResult {
    const { draggedTaskId, targetParentId, targetPosition } = operation;

    // Validate the move doesn't create a cycle
    const moveValidation = this.validateMove(tree, draggedTaskId, targetParentId);
    if (!moveValidation.isValid) {
      return moveValidation;
    }

    // Validate target position is reasonable
    if (targetPosition < 0) {
      return {
        isValid: false,
        error: 'Target position cannot be negative'
      };
    }

    // If moving to a parent, validate the parent exists
    if (targetParentId !== null) {
      const parentNode = this.findTaskInTree(tree, targetParentId);
      if (!parentNode) {
        return {
          isValid: false,
          error: 'Target parent not found in tree'
        };
      }

      // Validate position is within bounds (allowing one past the end for append)
      if (targetPosition > parentNode.node.children.length) {
        return {
          isValid: false,
          error: 'Target position is beyond the end of parent\'s children'
        };
      }
    } else {
      // Moving to root level - validate position against root nodes
      if (targetPosition > tree.length) {
        return {
          isValid: false,
          error: 'Target position is beyond the end of root nodes'
        };
      }
    }

    return { isValid: true };
  }

  /**
   * Gets all ancestor task IDs for a given task
   * @param tree - Tree structure to search
   * @param taskId - ID of the task to get ancestors for
   * @returns Array of ancestor task IDs from root to immediate parent
   */
  getAncestors(tree: TreeNode[], taskId: string): string[] {
    const path = this.getTaskPath(tree, taskId);
    // Remove the task itself from the path to get only ancestors
    return path.slice(0, -1);
  }

  /**
   * Gets all descendant task IDs for a given task
   * @param tree - Tree structure to search
   * @param taskId - ID of the task to get descendants for
   * @returns Array of all descendant task IDs
   */
  getDescendants(tree: TreeNode[], taskId: string): string[] {
    const taskNode = this.findTaskInTree(tree, taskId);
    if (!taskNode) {
      return [];
    }

    const descendants: string[] = [];
    
    const collectDescendants = (node: TreeNode): void => {
      for (const child of node.children) {
        descendants.push(child.task.id);
        collectDescendants(child);
      }
    };

    collectDescendants(taskNode.node);
    return descendants;
  }

  /**
   * Calculates the depth of the tree (maximum level)
   * @param tree - Tree structure to analyze
   * @returns Maximum depth of the tree
   */
  getTreeDepth(tree: TreeNode[]): number {
    if (tree.length === 0) {
      return 0;
    }

    let maxDepth = 0;
    
    const calculateDepth = (nodes: TreeNode[], currentDepth: number): void => {
      for (const node of nodes) {
        maxDepth = Math.max(maxDepth, currentDepth);
        if (node.children.length > 0) {
          calculateDepth(node.children, currentDepth + 1);
        }
      }
    };

    calculateDepth(tree, 1);
    return maxDepth;
  }

  /**
   * Counts total number of nodes in the tree
   * @param tree - Tree structure to count
   * @returns Total number of nodes
   */
  getNodeCount(tree: TreeNode[]): number {
    let count = 0;
    
    const countNodes = (nodes: TreeNode[]): void => {
      for (const node of nodes) {
        count++;
        countNodes(node.children);
      }
    };

    countNodes(tree);
    return count;
  }
}

/**
 * Default tree service instance
 */
export const treeService = new TreeService();