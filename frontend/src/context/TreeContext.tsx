import React, { createContext, useContext, useReducer, useCallback, useMemo } from 'react';

import { treeService } from '../services/tree/index.js';
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
  const { tasks, refreshTasks } = useTaskContext();

  // Rebuild tree when tasks change using TreeService
  const rootNodes = useMemo(() => {
    const newRootNodes = treeService.buildTreeFromTasks(tasks);
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
      // Validate the move operation using TreeService
      const validation = treeService.validateMove(rootNodes, taskId, parentId);
      if (!validation.isValid) {
        throw new Error(validation.error ?? 'Invalid move operation');
      }

      // Validate drag operation
      const dragValidation = treeService.validateDragOperation(rootNodes, {
        draggedTaskId: taskId,
        targetParentId: parentId,
        targetPosition: position
      });
      if (!dragValidation.isValid) {
        throw new Error(dragValidation.error ?? 'Invalid drag operation');
      }
      
      // Import API client dynamically to avoid circular dependency
      const { apiClient } = await import('../services/api/apiClient.js');
      await apiClient.moveTask(taskId, parentId ?? undefined, position);
      
      // Refresh tasks to get updated tree structure
      await refreshTasks();
      
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to move task:', error);
      throw error;
    }
  }, [rootNodes, refreshTasks]);

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