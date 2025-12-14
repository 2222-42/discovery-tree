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
  build: {
    // Output directory for production build
    outDir: 'dist',
    // Generate source maps for production debugging
    sourcemap: true,
    // Minify using esbuild for better performance
    minify: 'esbuild',
    // Target modern browsers for better optimization
    target: 'esnext',
    // Enable CSS code splitting
    cssCodeSplit: true,
    // Rollup options for advanced bundling configuration
    rollupOptions: {
      output: {
        // Manual chunk splitting for better caching
        manualChunks: {
          // Vendor chunk for React and related libraries
          vendor: ['react', 'react-dom'],
          // API utilities chunk
          api: ['axios'],
          // Testing utilities (only if imported in production)
          ...(process.env.NODE_ENV !== 'production' && {
            testing: ['fast-check']
          })
        },
        // Naming pattern for chunks
        chunkFileNames: 'assets/js/[name]-[hash].js',
        entryFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          if (/\.(css)$/.test(assetInfo.name || '')) {
            return 'assets/css/[name]-[hash].[ext]';
          }
          if (/\.(png|jpe?g|svg|gif|tiff|bmp|ico)$/i.test(assetInfo.name || '')) {
            return 'assets/images/[name]-[hash].[ext]';
          }
          return 'assets/[name]-[hash].[ext]';
        }
      },
      // External dependencies (if any should not be bundled)
      external: []
    },
    // Chunk size warning limit (500kb)
    chunkSizeWarningLimit: 500,
    // Tree shaking is enabled by default in Vite
    // Report compressed size
    reportCompressedSize: true,
    // Emit manifest for deployment tools
    manifest: true
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
          proxy.on('proxyReq', (_proxyReq, req, _res) => {
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
