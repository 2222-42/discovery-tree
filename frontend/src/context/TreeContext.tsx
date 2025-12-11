import React, { createContext, useContext, useReducer, useCallback, useMemo } from 'react';

import type { TreeContextValue } from '../types/app.js';
import type { TreeNode, TreeState } from '../types/tree.js';

import { useTaskContext } from './TaskContext.js';

// Tree context state reducer
type TreeAction =
  | { type: 'TOGGLE_NODE'; payload: string }
  | { type: 'EXPAND_NODE'; payload: string }
  | { type: 'COLLAPSE_NODE'; payload: string }
  | { type: 'SELECT_NODE'; payload: string | null }
  | { type: 'SET_ROOT_NODES'; payload: TreeNode[] };

const initialState: TreeState = {
  expandedNodes: new Set<string>(),
  selectedNodeId: null,
  rootNodes: [],
};

function treeReducer(state: TreeState, action: TreeAction): TreeState {
  switch (action.type) {
    case 'TOGGLE_NODE': {
      const newExpandedNodes = new Set(state.expandedNodes);
      if (newExpandedNodes.has(action.payload)) {
        newExpandedNodes.delete(action.payload);
      } else {
        newExpandedNodes.add(action.payload);
      }
      return { ...state, expandedNodes: newExpandedNodes };
    }
    case 'EXPAND_NODE': {
      const newExpandedNodes = new Set(state.expandedNodes);
      newExpandedNodes.add(action.payload);
      return { ...state, expandedNodes: newExpandedNodes };
    }
    case 'COLLAPSE_NODE': {
      const newExpandedNodes = new Set(state.expandedNodes);
      newExpandedNodes.delete(action.payload);
      return { ...state, expandedNodes: newExpandedNodes };
    }
    case 'SELECT_NODE':
      return { ...state, selectedNodeId: action.payload };
    case 'SET_ROOT_NODES':
      return { ...state, rootNodes: action.payload };
    default:
      return state;
  }
}

/**
 * Builds a tree structure from a flat array of tasks
 */
function buildTreeFromTasks(tasks: import('../types/task.js').Task[]): TreeNode[] {
  const taskMap = new Map<string, import('../types/task.js').Task>();
  const childrenMap = new Map<string, import('../types/task.js').Task[]>();
  
  // Build maps for efficient lookup
  tasks.forEach(task => {
    taskMap.set(task.id, task);
    if (task.parentId !== null) {
      if (!childrenMap.has(task.parentId)) {
        childrenMap.set(task.parentId, []);
      }
      const children = childrenMap.get(task.parentId);
      if (children) {
        children.push(task);
      }
    }
  });

  // Sort children by position
  childrenMap.forEach(children => {
    children.sort((a, b) => a.position - b.position);
  });

  // Build tree nodes recursively
  function buildNode(task: import('../types/task.js').Task, level: number): TreeNode {
    const children = childrenMap.get(task.id) ?? [];
    return {
      task,
      children: children.map(child => buildNode(child, level + 1)),
      isExpanded: false,
      level,
    };
  }

  // Find root tasks (no parent or parent doesn't exist)
  const rootTasks = tasks.filter(task => 
    (task.parentId === null) || !taskMap.has(task.parentId)
  );

  // Sort root tasks by position
  rootTasks.sort((a, b) => a.position - b.position);

  return rootTasks.map(task => buildNode(task, 0));
}

const TreeContext = createContext<TreeContextValue | undefined>(undefined);

interface TreeProviderProps {
  children: React.ReactNode;
}

/**
 * TreeProvider component that provides tree state and operations to child components
 * Manages the tree display state including expansion, selection, and navigation
 */
export function TreeProvider({ children }: TreeProviderProps): React.JSX.Element {
  const [state, dispatch] = useReducer(treeReducer, initialState);
  const { tasks } = useTaskContext();

  // Rebuild tree when tasks change
  const rootNodes = useMemo(() => {
    const newRootNodes = buildTreeFromTasks(tasks);
    // Update state if root nodes changed
    if (JSON.stringify(newRootNodes) !== JSON.stringify(state.rootNodes)) {
      dispatch({ type: 'SET_ROOT_NODES', payload: newRootNodes });
    }
    return newRootNodes;
  }, [tasks, state.rootNodes]);

  // Tree operations
  const toggleNode = useCallback((nodeId: string): void => {
    dispatch({ type: 'TOGGLE_NODE', payload: nodeId });
  }, []);

  const expandNode = useCallback((nodeId: string): void => {
    dispatch({ type: 'EXPAND_NODE', payload: nodeId });
  }, []);

  const collapseNode = useCallback((nodeId: string): void => {
    dispatch({ type: 'COLLAPSE_NODE', payload: nodeId });
  }, []);

  const selectNode = useCallback((nodeId: string | null): void => {
    dispatch({ type: 'SELECT_NODE', payload: nodeId });
  }, []);

  const moveTask = useCallback(async (
    taskId: string, 
    parentId: string | null, 
    position: number
  ): Promise<void> => {
    try {
      // TODO: Implement API call when ApiClient is available
      // await apiClient.moveTask(taskId, parentId, position);
      
      // For now, this is a placeholder - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      
      // eslint-disable-next-line no-console
      console.log(`Moving task ${taskId} to parent ${String(parentId)} at position ${String(position)}`);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to move task:', error);
      throw error;
    }
  }, []);

  const contextValue: TreeContextValue = {
    ...state,
    rootNodes,
    toggleNode,
    expandNode,
    collapseNode,
    selectNode,
    moveTask,
  };

  return (
    <TreeContext.Provider value={contextValue}>
      {children}
    </TreeContext.Provider>
  );
}

/**
 * Hook to access the TreeContext
 * Throws an error if used outside of TreeProvider
 */
// eslint-disable-next-line react-refresh/only-export-components
export function useTreeContext(): TreeContextValue {
  const context = useContext(TreeContext);
  if (context === undefined) {
    throw new Error('useTreeContext must be used within a TreeProvider');
  }
  return context;
}