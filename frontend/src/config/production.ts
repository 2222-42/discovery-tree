/**
 * Production configuration
 * Optimized settings for production deployment
 */

export const productionConfig = {
  // API configuration
  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
    timeout: 10000, // 10 seconds timeout for production
    retries: 3, // Retry failed requests up to 3 times
  },
  
  // Performance settings
  performance: {
    // Enable service worker caching in production
    enableCaching: true,
    // Lazy load components after this delay (ms)
    lazyLoadDelay: 100,
    // Virtual scrolling threshold for large trees
    virtualScrollThreshold: 100,
  },
  
  // Error handling
  errorHandling: {
    // Show detailed errors in production (set to false for security)
    showDetailedErrors: false,
    // Enable error reporting to external service
    enableErrorReporting: true,
    // Log errors to console in production
    logErrors: false,
  },
  
  // Feature flags
  features: {
    // Enable debug tools in production
    enableDebugTools: false,
    // Enable performance monitoring
    enablePerformanceMonitoring: true,
    // Enable analytics
    enableAnalytics: true,
  },
  
  // Build information
  build: {
    version: import.meta.env.VITE_APP_VERSION ?? '1.0.0',
    buildTime: import.meta.env['VITE_BUILD_TIME'] ?? new Date().toISOString(),
    environment: 'production',
  },
} as const;

export default productionConfig;