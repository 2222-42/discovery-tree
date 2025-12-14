/**
 * Configuration module exports
 */

export * from './development.js';
export * from './production.js';

// Export the appropriate config based on environment
const isDevelopment = import.meta.env.DEV;

export const config = isDevelopment 
  ? await import('./development.js').then(m => m.getDevelopmentConfig())
  : await import('./production.js').then(m => m.productionConfig);

export const getConfig = (): typeof config => config;