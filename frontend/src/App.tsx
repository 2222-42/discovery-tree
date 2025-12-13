import React, { useState, useEffect } from 'react';

import ErrorBoundary from './components/ErrorBoundary/ErrorBoundary.js';
import TaskDetails from './components/TaskDetails/TaskDetails.js';
import TaskForm from './components/TaskForm/TaskForm.js';
import TreeView from './components/TreeView/TreeView.js';
import { TaskProvider, TreeProvider } from './context/index.js';
import { useTaskContext } from './context/TaskContext.js';

import './App.css';

/**
 * Main application layout component
 * Provides the overall structure and layout for the discovery tree interface
 */
function AppLayout(): React.JSX.Element {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const { selectedTask, selectTask, loading, error, refreshTasks } = useTaskContext();

  // Load initial data when the app starts
  useEffect(() => {
    void refreshTasks();
  }, [refreshTasks]);

  const handleNodeSelect = (nodeId: string | null): void => {
    selectTask(nodeId);
  };

  const handleCreateTask = (): void => {
    setShowCreateForm(true);
  };

  const handleFormSubmit = (): void => {
    setShowCreateForm(false);
    // Refresh tasks after creation to ensure tree is up to date
    void refreshTasks();
  };

  const handleFormCancel = (): void => {
    setShowCreateForm(false);
  };

  const handleTaskUpdate = (): void => {
    // Task updates are handled through the context
    // Refresh tasks to ensure consistency
    void refreshTasks();
  };

  const handleTaskDelete = (): void => {
    // Task deletion is handled through the context
    // Refresh tasks to ensure tree is updated
    void refreshTasks();
  };

  const handleCloseDetails = (): void => {
    selectTask(null);
  };

  // Show global error state if there's an error and no tasks loaded
  if (error !== null && error !== '' && !loading) {
    return (
      <div className="app app--error">
        <header className="app__header">
          <h1 className="app__title">Discovery Tree</h1>
        </header>
        <main className="app__main">
          <div className="app__error-state">
            <h2>Unable to load tasks</h2>
            <p>{error}</p>
            <button 
              className="app__retry-button"
              onClick={() => { void refreshTasks(); }}
              data-testid="retry-button"
            >
              Retry
            </button>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="app__header">
        <h1 className="app__title">Discovery Tree</h1>
        <button 
          className="app__create-button"
          onClick={handleCreateTask}
          disabled={loading}
          data-testid="create-task-button"
        >
          {loading ? 'Loading...' : 'Create Root Task'}
        </button>
        {(error !== null && error !== '') && (
          <div className="app__error-banner" data-testid="error-banner">
            {error}
          </div>
        )}
      </header>

      <main className="app__main">
        <div className="app__layout">
          <div className="app__tree-panel">
            <TreeView 
              onNodeSelect={handleNodeSelect}
              loading={loading}
              error={error}
              data-testid="main-tree-view"
            />
          </div>

          <div className="app__details-panel">
            {showCreateForm ? (
              <TaskForm
                mode="create"
                onSubmit={handleFormSubmit}
                onCancel={handleFormCancel}
                data-testid="create-task-form"
              />
            ) : (
              <TaskDetails
                task={selectedTask}
                onTaskUpdate={handleTaskUpdate}
                onTaskDelete={handleTaskDelete}
                onClose={handleCloseDetails}
                data-testid="task-details-panel"
              />
            )}
          </div>
        </div>
      </main>
    </div>
  );
}

/**
 * Root App component with error boundaries and context providers
 * Sets up the application structure with proper error handling and state management
 */
function App(): React.JSX.Element {
  return (
    <ErrorBoundary>
      <TaskProvider>
        <TreeProvider>
          <AppLayout />
        </TreeProvider>
      </TaskProvider>
    </ErrorBoundary>
  );
}

export default App;
