/**
 * Development configuration utilities
 * Provides environment-specific configuration and development helpers
 */

/**
 * Development configuration interface
 */
export interface DevelopmentConfig {
  apiBaseUrl: string;
  apiTimeout: number;
  enableDebugLogging: boolean;
  enableMockApi: boolean;
  appName: string;
  appVersion: string;
  hmrOverlay: boolean;
}

/**
 * Get development configuration from environment variables
 */
export const getDevelopmentConfig = (): DevelopmentConfig => {
  return {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    apiTimeout: Number(import.meta.env.VITE_API_TIMEOUT) || 10000,
    enableDebugLogging: import.meta.env.VITE_ENABLE_DEBUG_LOGGING === 'true',
    enableMockApi: import.meta.env.VITE_ENABLE_MOCK_API === 'true',
    appName: import.meta.env.VITE_APP_NAME || 'Discovery Tree Frontend',
    appVersion: import.meta.env.VITE_APP_VERSION || '1.0.0',
    hmrOverlay: import.meta.env.VITE_HMR_OVERLAY !== 'false',
  };
};

/**
 * Check if running in development mode
 */
export const isDevelopment = (): boolean => {
  return import.meta.env.DEV;
};

/**
 * Check if running in production mode
 */
export const isProduction = (): boolean => {
  return import.meta.env.PROD;
};

/**
 * Development logging utility
 */
export const devLog = {
  debug: (message: string, ...args: unknown[]): void => {
    if (getDevelopmentConfig().enableDebugLogging) {
      // eslint-disable-next-line no-console
      console.debug(`[DEV] ${message}`, ...args);
    }
  },
  info: (message: string, ...args: unknown[]): void => {
    if (isDevelopment()) {
      // eslint-disable-next-line no-console
      console.info(`[DEV] ${message}`, ...args);
    }
  },
  warn: (message: string, ...args: unknown[]): void => {
    if (isDevelopment()) {
      // eslint-disable-next-line no-console
      console.warn(`[DEV] ${message}`, ...args);
    }
  },
  error: (message: string, ...args: unknown[]): void => {
    if (isDevelopment()) {
      // eslint-disable-next-line no-console
      console.error(`[DEV] ${message}`, ...args);
    }
  },
};

/**
 * Hot module replacement utilities
 */
export const hmrUtils = {
  /**
   * Accept HMR updates for a module
   */
  accept: (callback?: () => void): void => {
    if (import.meta.hot) {
      import.meta.hot.accept(callback);
    }
  },
  
  /**
   * Dispose HMR resources
   */
  dispose: (callback: () => void): void => {
    if (import.meta.hot) {
      import.meta.hot.dispose(callback);
    }
  },
  
  /**
   * Check if HMR is available
   */
  isAvailable: (): boolean => {
    return Boolean(import.meta.hot);
  },
};