/**
 * API-related type definitions
 * Defines request/response types for communication with the backend REST API
 */

import type { Task, TaskStatus } from './task.js';

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T> {
  data: T;
  success: boolean;
  message?: string;
}

/**
 * API error response structure
 */
export interface ApiError {
  error: string;
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

/**
 * Task response from the backend API
 * Matches the backend TaskResponse structure
 */
export interface TaskResponse {
  id: string;
  description: string;
  status: string;
  parentId: string | null;
  position: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * Create root task request
 */
export interface CreateRootTaskRequest {
  description: string;
}

/**
 * Create child task request
 */
export interface CreateChildTaskRequest {
  description: string;
  parentId: string;
}

/**
 * Update task request
 */
export interface UpdateTaskRequest {
  description?: string;
  status?: TaskStatus;
  parentId?: string;
  position?: number;
}

/**
 * Update task status request
 */
export interface UpdateTaskStatusRequest {
  status: TaskStatus;
}

/**
 * Move task request
 */
export interface MoveTaskRequest {
  parentId?: string;
  position: number;
}

/**
 * HTTP method types
 */
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

/**
 * API request configuration
 */
export interface ApiRequestConfig {
  method: HttpMethod;
  url: string;
  data?: unknown;
  params?: Record<string, string | number>;
  headers?: Record<string, string>;
}

/**
 * API client interface defining all available operations
 */
export interface ApiClient {
  // Task CRUD operations
  createRootTask(description: string): Promise<Task>;
  createChildTask(description: string, parentId: string): Promise<Task>;
  getTask(id: string): Promise<Task>;
  getAllTasks(): Promise<Task[]>;
  getRootTask(): Promise<Task>;
  getTaskChildren(id: string): Promise<Task[]>;
  updateTask(id: string, description: string): Promise<Task>;
  updateTaskStatus(id: string, status: TaskStatus): Promise<Task>;
  moveTask(id: string, parentId: string | undefined, position: number): Promise<Task>;
  deleteTask(id: string): Promise<void>;
}

/**
 * HTTP client interface for making requests
 */
export interface HttpClient {
  get<T>(url: string, config?: Partial<ApiRequestConfig>): Promise<T>;
  post<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T>;
  put<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T>;
  delete<T>(url: string, config?: Partial<ApiRequestConfig>): Promise<T>;
  patch<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T>;
}