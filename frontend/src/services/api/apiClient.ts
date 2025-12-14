/**
 * API client implementation for the discovery tree backend
 * Provides all CRUD operations for tasks through REST endpoints
 * Requirements: 4.1, 4.2, 4.3, 4.4, 4.5
 */

import type { 
  ApiClient, 
  HttpClient, 
  TaskResponse,
  CreateRootTaskRequest,
  CreateChildTaskRequest,
  UpdateTaskRequest,
  UpdateTaskStatusRequest,
  MoveTaskRequest
} from '../../types/api.js';
import type { Task, TaskStatus } from '../../types/task.js';

import { createHttpClient } from './httpClient.js';

/**
 * API client configuration
 */
interface ApiClientConfig {
  baseURL?: string;
  timeout?: number;
}

/**
 * API client implementation for the discovery tree backend
 */
export class DiscoveryTreeApiClient implements ApiClient {
  private readonly httpClient: HttpClient;

  constructor(config?: ApiClientConfig) {
    this.httpClient = createHttpClient(config);
  }

  /**
   * Convert TaskResponse from backend to frontend Task model
   */
  private taskResponseToTask(response: TaskResponse): Task {
    return {
      id: response.id,
      description: response.description,
      status: response.status as TaskStatus,
      parentId: response.parentId,
      position: response.position,
      createdAt: response.createdAt,
      updatedAt: response.updatedAt,
    };
  }

  /**
   * Create a new root task
   * POST /api/v1/tasks/root
   * Requirements: 4.1 - Task creation through API
   */
  async createRootTask(description: string): Promise<Task> {
    const request: CreateRootTaskRequest = { description };
    const response = await this.httpClient.post<TaskResponse>('/tasks/root', request);
    return this.taskResponseToTask(response);
  }

  /**
   * Create a new child task under a parent
   * POST /api/v1/tasks
   * Requirements: 4.1 - Task creation through API
   */
  async createChildTask(description: string, parentId: string): Promise<Task> {
    const request: CreateChildTaskRequest = { description, parentId };
    const response = await this.httpClient.post<TaskResponse>('/tasks', request);
    return this.taskResponseToTask(response);
  }

  /**
   * Get a specific task by ID
   * GET /api/v1/tasks/{id}
   * Requirements: 4.2 - Task detail viewing
   */
  async getTask(id: string): Promise<Task> {
    const response = await this.httpClient.get<TaskResponse>(`/tasks/${id}`);
    return this.taskResponseToTask(response);
  }

  /**
   * Get all tasks in the system
   * GET /api/v1/tasks
   * Requirements: 4.2 - Task viewing and tree display
   */
  async getAllTasks(): Promise<Task[]> {
    const responses = await this.httpClient.get<TaskResponse[]>('/tasks');
    return responses.map(response => this.taskResponseToTask(response));
  }

  /**
   * Get the root task
   * GET /api/v1/tasks/root
   * Requirements: 4.2 - Tree structure access
   */
  async getRootTask(): Promise<Task> {
    const response = await this.httpClient.get<TaskResponse>('/tasks/root');
    return this.taskResponseToTask(response);
  }

  /**
   * Get children of a specific task
   * GET /api/v1/tasks/{id}/children
   * Requirements: 4.2 - Tree navigation and hierarchy display
   */
  async getTaskChildren(id: string): Promise<Task[]> {
    const responses = await this.httpClient.get<TaskResponse[]>(`/tasks/${id}/children`);
    return responses.map(response => this.taskResponseToTask(response));
  }

  /**
   * Update a task's description
   * PUT /api/v1/tasks/{id}
   * Requirements: 4.3 - Task editing and updates
   */
  async updateTask(id: string, description: string): Promise<Task> {
    const request: UpdateTaskRequest = { description };
    const response = await this.httpClient.put<TaskResponse>(`/tasks/${id}`, request);
    return this.taskResponseToTask(response);
  }

  /**
   * Update a task's status
   * PUT /api/v1/tasks/{id}/status
   * Requirements: 4.3 - Task status management
   */
  async updateTaskStatus(id: string, status: TaskStatus): Promise<Task> {
    const request: UpdateTaskStatusRequest = { status };
    const response = await this.httpClient.put<TaskResponse>(`/tasks/${id}/status`, request);
    return this.taskResponseToTask(response);
  }

  /**
   * Move a task to a new position or parent
   * PUT /api/v1/tasks/{id}/move
   * Requirements: 4.3 - Task hierarchy management
   */
  async moveTask(id: string, parentId: string | undefined, position: number): Promise<Task> {
    const request: MoveTaskRequest = { 
      position,
      ...(parentId !== undefined && { parentId }),
    };
    const response = await this.httpClient.put<TaskResponse>(`/tasks/${id}/move`, request);
    return this.taskResponseToTask(response);
  }

  /**
   * Delete a task and all its descendants
   * DELETE /api/v1/tasks/{id}
   * Requirements: 4.4 - Task deletion
   */
  async deleteTask(id: string): Promise<void> {
    await this.httpClient.delete(`/tasks/${id}`);
  }
}

/**
 * Create a default API client instance
 */
export const createApiClient = (config?: ApiClientConfig): ApiClient => {
  return new DiscoveryTreeApiClient(config);
};

/**
 * Default API client instance for use throughout the application
 */
export const apiClient = createApiClient();