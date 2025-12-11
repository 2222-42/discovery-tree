/**
 * HTTP client implementation using Axios
 * Provides a wrapper around Axios with interceptors for error handling
 * Requirements: 4.5 - Error handling for API operations
 */

import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

import type { HttpClient, ApiRequestConfig, ApiError } from '../../types/api.js';

/**
 * HTTP client configuration
 */
interface HttpClientConfig {
  baseURL: string;
  timeout?: number;
  headers?: Record<string, string>;
}

/**
 * Default configuration for the HTTP client
 */
const DEFAULT_CONFIG: HttpClientConfig = {
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
};

/**
 * HTTP client implementation using Axios
 */
export class AxiosHttpClient implements HttpClient {
  private readonly axiosInstance: AxiosInstance;

  constructor(config: Partial<HttpClientConfig> = {}) {
    const finalConfig = { ...DEFAULT_CONFIG, ...config };
    
    this.axiosInstance = axios.create({
      baseURL: finalConfig.baseURL,
      ...(finalConfig.timeout !== undefined && { timeout: finalConfig.timeout }),
      ...(finalConfig.headers && { headers: finalConfig.headers }),
    });

    this.setupInterceptors();
  }

  /**
   * Set up request and response interceptors for error handling
   */
  private setupInterceptors(): void {
    // Request interceptor for logging and adding auth headers if needed
    this.axiosInstance.interceptors.request.use(
      (config) => {
        // Only log in development mode
        if (import.meta.env.DEV) {
          // eslint-disable-next-line no-console
          console.debug('HTTP Request:', {
            method: config.method?.toUpperCase(),
            url: config.url ?? '',
          });
        }
        return config;
      },
      (error: unknown) => {
        if (import.meta.env.DEV) {
          // eslint-disable-next-line no-console
          console.error('HTTP Request Error:', error);
        }
        const apiError = this.transformError(error);
        return Promise.reject(new Error(apiError.message));
      }
    );

    // Response interceptor for error handling and logging
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => {
        if (import.meta.env.DEV) {
          // eslint-disable-next-line no-console
          console.debug('HTTP Response:', {
            status: response.status,
            url: response.config.url ?? '',
          });
        }
        return response;
      },
      (error: unknown) => {
        const apiError = this.transformError(error);
        if (import.meta.env.DEV) {
          // eslint-disable-next-line no-console
          console.error('HTTP Response Error:', apiError);
        }
        return Promise.reject(new Error(apiError.message));
      }
    );
  }

  /**
   * Transform Axios errors into standardized API errors
   */
  private transformError(error: unknown): ApiError {
    // Type guard for axios errors
    if (axios.isAxiosError(error)) {
      if (error.response) {
        // Server responded with error status
        const status = error.response.status;
        const data: unknown = error.response.data;
        
        // Try to extract structured error from response
        if (this.isErrorResponse(data)) {
          return {
            error: data.error,
            code: data.code,
            message: data.message,
            ...(data.details && { details: data.details }),
          };
        }
        
        // Otherwise, create a generic error based on status code
        return {
          error: 'HTTP_ERROR',
          code: `HTTP_${String(status)}`,
          message: this.getStatusMessage(status),
        };
      } else if (error.request !== undefined) {
        // Network error - no response received
        return {
          error: 'NETWORK_ERROR',
          code: 'NETWORK_ERROR',
          message: 'Network error - unable to reach server',
        };
      }
    }
    
    // Request setup error or other error
    const message = error instanceof Error ? error.message : 'Request configuration error';
    return {
      error: 'REQUEST_ERROR',
      code: 'REQUEST_ERROR',
      message,
    };
  }

  /**
   * Type guard to check if response data is a structured error
   */
  private isErrorResponse(data: unknown): data is ApiError {
    return (
      typeof data === 'object' &&
      data !== null &&
      'error' in data &&
      'code' in data &&
      'message' in data &&
      typeof (data as Record<string, unknown>)['error'] === 'string' &&
      typeof (data as Record<string, unknown>)['code'] === 'string' &&
      typeof (data as Record<string, unknown>)['message'] === 'string'
    );
  }

  /**
   * Get user-friendly message for HTTP status codes
   */
  private getStatusMessage(status: number): string {
    switch (status) {
      case 400:
        return 'Bad request - invalid data provided';
      case 401:
        return 'Unauthorized - authentication required';
      case 403:
        return 'Forbidden - insufficient permissions';
      case 404:
        return 'Not found - resource does not exist';
      case 409:
        return 'Conflict - resource already exists or operation not allowed';
      case 422:
        return 'Validation error - invalid data format';
      case 500:
        return 'Internal server error - please try again later';
      case 502:
        return 'Bad gateway - server temporarily unavailable';
      case 503:
        return 'Service unavailable - server temporarily down';
      case 504:
        return 'Gateway timeout - server took too long to respond';
      default:
        return `HTTP ${String(status)} error`;
    }
  }

  /**
   * Convert internal config to Axios config
   */
  private toAxiosConfig(config?: Partial<ApiRequestConfig>): AxiosRequestConfig {
    if (!config) return {};
    
    return {
      ...(config.params && { params: config.params }),
      ...(config.headers && { headers: config.headers }),
    };
  }

  /**
   * GET request
   */
  async get<T>(url: string, config?: Partial<ApiRequestConfig>): Promise<T> {
    const response = await this.axiosInstance.get<T>(url, this.toAxiosConfig(config));
    return response.data;
  }

  /**
   * POST request
   */
  async post<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T> {
    const response = await this.axiosInstance.post<T>(url, data, this.toAxiosConfig(config));
    return response.data;
  }

  /**
   * PUT request
   */
  async put<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T> {
    const response = await this.axiosInstance.put<T>(url, data, this.toAxiosConfig(config));
    return response.data;
  }

  /**
   * DELETE request
   */
  async delete<T>(url: string, config?: Partial<ApiRequestConfig>): Promise<T> {
    const response = await this.axiosInstance.delete<T>(url, this.toAxiosConfig(config));
    return response.data;
  }

  /**
   * PATCH request
   */
  async patch<T>(url: string, data?: unknown, config?: Partial<ApiRequestConfig>): Promise<T> {
    const response = await this.axiosInstance.patch<T>(url, data, this.toAxiosConfig(config));
    return response.data;
  }
}

/**
 * Create a default HTTP client instance
 */
export const createHttpClient = (config?: Partial<HttpClientConfig>): HttpClient => {
  return new AxiosHttpClient(config);
};