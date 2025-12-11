/**
 * API services exports
 * Provides clean exports for all API-related functionality
 */

export { AxiosHttpClient, createHttpClient } from './httpClient.js';
export { DiscoveryTreeApiClient, createApiClient, apiClient } from './apiClient.js';

// Re-export types for convenience
export type { 
  ApiClient, 
  HttpClient, 
  ApiError,
  TaskResponse,
  CreateRootTaskRequest,
  CreateChildTaskRequest,
  UpdateTaskRequest,
  UpdateTaskStatusRequest,
  MoveTaskRequest
} from '../../types/api.js';