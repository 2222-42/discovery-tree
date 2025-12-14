# Discovery Tree Frontend

A React-based web interface for managing hierarchical task trees. This frontend application provides an intuitive interface for creating, viewing, and managing tasks in a tree structure, connecting to the Discovery Tree API backend.

## Features

- **Interactive Tree View**: Visual representation of task hierarchies with expandable/collapsible nodes
- **Task Management**: Create, edit, delete, and update task status
- **Real-time Updates**: Hot module replacement for development and responsive UI updates
- **Error Handling**: Comprehensive error boundaries and user-friendly error messages
- **Responsive Design**: Clean, accessible interface that works across different screen sizes
- **TypeScript**: Full type safety with comprehensive type definitions

## Tech Stack

- **React 19** with TypeScript
- **Vite** for fast development and building
- **Axios** for API communication
- **Vitest** for testing
- **ESLint** for code quality

## Getting Started

### Prerequisites

- Node.js (version 18 or higher)
- npm or yarn
- Discovery Tree API server running (see main README)

### Installation

```bash
# Install dependencies
npm install

# Copy environment configuration
cp .env.example .env
```

### Environment Configuration

Create a `.env` file with the following variables:

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080
VITE_API_TIMEOUT=10000

# Development Configuration
VITE_ENABLE_DEBUG_LOGGING=true
```

### Development

```bash
# Start development server
npm run dev

# Test development setup
npm run dev:test

# Run tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage
```

The development server will start at `http://localhost:3000` and automatically proxy API requests to the backend server.

### Building

```bash
# Development build
npm run build

# Production build
npm run build:prod

# Analyze bundle size
npm run build:analyze

# Preview production build
npm run preview:prod
```

## Project Structure

```
src/
├── components/          # React components
│   ├── ErrorBoundary/   # Error handling
│   ├── TaskDetails/     # Task detail view
│   ├── TaskForm/        # Task creation/editing
│   ├── TaskNode/        # Individual tree nodes
│   └── TreeView/        # Main tree interface
├── context/             # React context providers
│   ├── TaskContext.tsx  # Task state management
│   └── TreeContext.tsx  # Tree state management
├── services/            # API and business logic
│   ├── api/             # HTTP client and API calls
│   └── tree/            # Tree manipulation utilities
├── types/               # TypeScript type definitions
├── utils/               # Utility functions and generators
└── config/              # Environment configuration
```

## API Integration

The frontend communicates with the Discovery Tree API backend:

- **Base URL**: Configured via `VITE_API_BASE_URL`
- **Endpoints**: RESTful API for task CRUD operations
- **Error Handling**: Comprehensive error handling with user feedback
- **Type Safety**: Full TypeScript integration with API models

## Development Features

### Hot Module Replacement

The application supports HMR for fast development:
- Component changes reload instantly
- State is preserved during updates
- Development logging shows HMR activity

### Testing

Comprehensive test suite using Vitest:
- Unit tests for components and utilities
- Integration tests for API services
- Property-based testing with fast-check
- Coverage reporting

### Code Quality

- **ESLint**: Strict linting rules with TypeScript support
- **Type Checking**: Full TypeScript coverage
- **Import Organization**: Automatic import sorting and validation
- **Accessibility**: JSX accessibility linting

## Scripts Reference

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server |
| `npm run dev:test` | Test development setup |
| `npm run build` | Build for development |
| `npm run build:prod` | Build for production |
| `npm run build:analyze` | Build and analyze bundle |
| `npm run test` | Run tests once |
| `npm run test:watch` | Run tests in watch mode |
| `npm run test:coverage` | Run tests with coverage |
| `npm run lint` | Lint and fix code |
| `npm run lint:check` | Check linting without fixing |
| `npm run preview` | Preview development build |
| `npm run preview:prod` | Preview production build |

## Contributing

1. Follow the existing code style and TypeScript patterns
2. Add tests for new features
3. Ensure all linting passes
4. Update documentation as needed

## Architecture Notes

- **Context-based State**: Uses React Context for global state management
- **Component Composition**: Modular component design with clear separation of concerns
- **Error Boundaries**: Graceful error handling at component boundaries
- **Type Safety**: Comprehensive TypeScript coverage with strict configuration
