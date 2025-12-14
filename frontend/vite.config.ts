import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@/components': resolve(__dirname, './src/components'),
      '@/services': resolve(__dirname, './src/services'),
      '@/types': resolve(__dirname, './src/types'),
      '@/hooks': resolve(__dirname, './src/hooks'),
      '@/context': resolve(__dirname, './src/context'),
      '@/utils': resolve(__dirname, './src/utils'),
    },
  },
  server: {
    port: Number(process.env.VITE_DEV_PORT) || 3000,
    host: process.env.VITE_DEV_HOST || true, // Allow external connections
    open: true, // Automatically open browser
    hmr: {
      overlay: process.env.VITE_HMR_OVERLAY !== 'false', // Show HMR overlay on errors
      port: Number(process.env.VITE_HMR_PORT) || 3000,
    },
    proxy: {
      // Proxy API requests to backend server
      '/api': {
        target: process.env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, _res) => {
            console.log('Proxy error:', err);
          });
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            if (process.env.VITE_ENABLE_DEBUG_LOGGING === 'true') {
              console.log('Sending Request to the Target:', req.method, req.url);
            }
          });
          proxy.on('proxyRes', (proxyRes, req, _res) => {
            if (process.env.VITE_ENABLE_DEBUG_LOGGING === 'true') {
              console.log('Received Response from the Target:', proxyRes.statusCode, req.url);
            }
          });
        },
      },
    },
  },
  preview: {
    port: Number(process.env.VITE_DEV_PORT) || 3000,
    host: process.env.VITE_DEV_HOST || true,
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    css: true,
  },
})
