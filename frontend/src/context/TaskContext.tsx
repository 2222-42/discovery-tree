import React, { createContext, useContext, useReducer, useCallback } from 'react';

import type { TaskContextValue, TaskContextState } from '../types/app.js';
import type { Task, TaskStatus } from '../types/task.js';

// Task context state reducer
type TaskAction =
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_TASKS'; payload: Task[] }
  | { type: 'ADD_TASK'; payload: Task }
  | { type: 'UPDATE_TASK'; payload: Task }
  | { type: 'DELETE_TASK'; payload: string }
  | { type: 'SELECT_TASK'; payload: string | null };

const initialState: TaskContextState = {
  tasks: [],
  selectedTask: null,
  loading: false,
  error: null,
};

function taskReducer(state: TaskContextState, action: TaskAction): TaskContextState {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'SET_TASKS':
      return { ...state, tasks: action.payload, loading: false, error: null };
    case 'ADD_TASK':
      return { 
        ...state, 
        tasks: [...state.tasks, action.payload],
        loading: false,
        error: null 
      };
    case 'UPDATE_TASK':
      return {
        ...state,
        tasks: state.tasks.map(task => 
          task.id === action.payload.id ? action.payload : task
        ),
        selectedTask: (state.selectedTask?.id === action.payload.id) ? action.payload : state.selectedTask,
        loading: false,
        error: null
      };
    case 'DELETE_TASK':
      return {
        ...state,
        tasks: state.tasks.filter(task => task.id !== action.payload),
        selectedTask: (state.selectedTask?.id === action.payload) ? null : state.selectedTask,
        loading: false,
        error: null
      };
    case 'SELECT_TASK':
      return {
        ...state,
        selectedTask: (action.payload !== null) ? (state.tasks.find(task => task.id === action.payload) ?? null) : null
      };
    default:
      return state;
  }
}

const TaskContext = createContext<TaskContextValue | undefined>(undefined);

interface TaskProviderProps {
  children: React.ReactNode;
}

/**
 * TaskProvider component that provides task state and operations to child components
 * Manages the global task state including CRUD operations and selection
 */
export function TaskProvider({ children }: TaskProviderProps): React.JSX.Element {
  const [state, dispatch] = useReducer(taskReducer, initialState);

  // Task operations - these will be implemented when the API client is available
  const createRootTask = useCallback(async (description: string): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // const task = await apiClient.createRootTask(description);
      // dispatch({ type: 'ADD_TASK', payload: task });
      
      // Placeholder implementation for now - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      const mockTask: Task = {
        id: `task-${String(Date.now())}`,
        description,
        status: 'ROOT_WORK_ITEM' as TaskStatus,
        parentId: null,
        position: 0,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      dispatch({ type: 'ADD_TASK', payload: mockTask });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to create task' });
    }
  }, []);

  const createChildTask = useCallback(async (description: string, parentId: string): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // const task = await apiClient.createChildTask(description, parentId);
      // dispatch({ type: 'ADD_TASK', payload: task });
      
      // Placeholder implementation for now - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      const mockTask: Task = {
        id: `task-${String(Date.now())}`,
        description,
        status: 'TODO' as TaskStatus,
        parentId,
        position: 0,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      dispatch({ type: 'ADD_TASK', payload: mockTask });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to create task' });
    }
  }, []);

  const updateTask = useCallback(async (id: string, description: string): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // const task = await apiClient.updateTask(id, description);
      // dispatch({ type: 'UPDATE_TASK', payload: task });
      
      // Placeholder implementation for now - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      const existingTask = state.tasks.find(task => task.id === id);
      if (existingTask) {
        const updatedTask: Task = {
          ...existingTask,
          description,
          updatedAt: new Date().toISOString(),
        };
        dispatch({ type: 'UPDATE_TASK', payload: updatedTask });
      }
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to update task' });
    }
  }, [state.tasks]);

  const updateTaskStatus = useCallback(async (id: string, status: TaskStatus): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // const task = await apiClient.updateTaskStatus(id, status);
      // dispatch({ type: 'UPDATE_TASK', payload: task });
      
      // Placeholder implementation for now - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      const existingTask = state.tasks.find(task => task.id === id);
      if (existingTask) {
        const updatedTask: Task = {
          ...existingTask,
          status,
          updatedAt: new Date().toISOString(),
        };
        dispatch({ type: 'UPDATE_TASK', payload: updatedTask });
      }
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to update task status' });
    }
  }, [state.tasks]);

  const deleteTask = useCallback(async (id: string): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // await apiClient.deleteTask(id);
      
      // Placeholder implementation for now - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      dispatch({ type: 'DELETE_TASK', payload: id });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to delete task' });
    }
  }, []);

  const selectTask = useCallback((id: string | null): void => {
    dispatch({ type: 'SELECT_TASK', payload: id });
  }, []);

  const refreshTasks = useCallback(async (): Promise<void> => {
    dispatch({ type: 'SET_LOADING', payload: true });
    try {
      // TODO: Implement API call when ApiClient is available
      // const tasks = await apiClient.getAllTasks();
      // dispatch({ type: 'SET_TASKS', payload: tasks });
      
      // Placeholder - simulate async operation
      await new Promise(resolve => setTimeout(resolve, 100));
      dispatch({ type: 'SET_LOADING', payload: false });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: error instanceof Error ? error.message : 'Failed to refresh tasks' });
    }
  }, []);

  const contextValue: TaskContextValue = {
    ...state,
    createRootTask,
    createChildTask,
    updateTask,
    updateTaskStatus,
    deleteTask,
    selectTask,
    refreshTasks,
  };

  return (
    <TaskContext.Provider value={contextValue}>
      {children}
    </TaskContext.Provider>
  );
}

/**
 * Hook to access the TaskContext
 * Throws an error if used outside of TaskProvider
 */
// eslint-disable-next-line react-refresh/only-export-components
export function useTaskContext(): TaskContextValue {
  const context = useContext(TaskContext);
  if (context === undefined) {
    throw new Error('useTaskContext must be used within a TaskProvider');
  }
  return context;
}