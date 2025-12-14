# Development Server Setup

This document describes the development server configuration and hot reloading setup for the Discovery Tree frontend application.

## Features

### Development Server
- **Port**: 3000 (configurable via `VITE_DEV_PORT`)
- **Host**: Allows external connections
- **Auto-open**: Automatically opens browser on start
- **Hot Module Replacement (HMR)**: Enabled with error overlay

### API Proxy
- **Target**: `http://localhost:8080` (backend server)
- **Path**: `/api/*` requests are proxied to the backend
- **Debug Logging**: Proxy requests/responses logged when debug is enabled

### Environment Variables
The application supports environment-specific configuration through `.env` files:

#### `.env` (Base configuration)
```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=10000
VITE_DEV_PORT=3000
VITE_DEV_HOST=localhost
VITE_APP_NAME=Discovery Tree Frontend
VITE_APP_VERSION=1.0.0
VITE_ENABLE_DEBUG_LOGGING=false
VITE_ENABLE_MOCK_API=false
```

#### `.env.development` (Development overrides)
```bash
VITE_ENABLE_DEBUG_LOGGING=true
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_HMR_OVERLAY=true
VITE_HMR_PORT=3000
```

#### `.env.production` (Production overrides)
```bash
VITE_API_BASE_URL=/api/v1
VITE_ENABLE_DEBUG_LOGGING=false
VITE_ENABLE_MOCK_API=false
```

## Usage

### Starting the Development Server
```bash
npm run dev
```

This will:
1. Start the Vite development server on port 3000
2. Enable hot module replacement
3. Proxy API requests to the backend server
4. Automatically open the browser
5. Show compilation errors in an overlay

### Testing the Setup
```bash
npm run dev:test
```

This script validates:
- Environment files exist and contain required variables
- Vite configuration includes all required features
- Creates a test component for HMR verification

### Development Workflow

1. **Start Backend Server**: Ensure the Go backend is running on port 8080
2. **Start Frontend**: Run `npm run dev`
3. **Make Changes**: Edit any source file
4. **See Updates**: Changes are hot-reloaded without full page refresh
5. **API Calls**: All `/api/*` requests are automatically proxied to the backend

### Hot Module Replacement (HMR)

HMR is configured to:
- Preserve component state when possible
- Show error overlay for compilation errors
- Update styles without page refresh
- Reload the page only when necessary

### Debug Logging

When `VITE_ENABLE_DEBUG_LOGGING=true`:
- HTTP requests/responses are logged to console
- Proxy operations are logged
- Development configuration is logged on app start

### Configuration Files

#### `vite.config.ts`
- Server configuration (port, host, HMR)
- API proxy setup
- Build optimization settings
- Path aliases for clean imports

#### `src/config/development.ts`
- Development utilities and configuration helpers
- Environment variable parsing
- HMR utilities
- Development logging functions

#### `src/vite-env.d.ts`
- TypeScript definitions for environment variables
- Ensures type safety for `import.meta.env` usage

## Troubleshooting

### Common Issues

1. **Port 3000 already in use**
   - Change `VITE_DEV_PORT` in `.env.development`
   - Or kill the process using port 3000

2. **API requests failing**
   - Ensure backend server is running on port 8080
   - Check proxy configuration in `vite.config.ts`
   - Verify `VITE_API_BASE_URL` environment variable

3. **HMR not working**
   - Check browser console for errors
   - Verify `VITE_HMR_OVERLAY` is not set to `false`
   - Restart the development server

4. **Environment variables not loading**
   - Ensure variables start with `VITE_`
   - Check file names (`.env`, `.env.development`, etc.)
   - Restart the development server after changes

### Performance Tips

1. **Large Projects**: Use `noWarnOnMultipleProjects` in TypeScript config
2. **Slow HMR**: Check for circular dependencies
3. **Memory Issues**: Increase Node.js memory limit if needed

## Requirements Validation

This setup fulfills the following requirements:

- **5.4**: Development server with hot reloading capabilities
- **4.1-4.5**: API integration through proxy configuration
- **1.1-1.5**: TypeScript compilation and build tooling
- **2.1-2.4**: ESLint integration with development workflow

The development server provides a modern, efficient development experience with fast feedback loops and seamless integration with the backend API.