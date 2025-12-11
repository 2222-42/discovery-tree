import React, { useState } from 'react';

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
  const { selectedTask, selectTask } = useTaskContext();

  const handleNodeSelect = (nodeId: string | null): void => {
    selectTask(nodeId);
  };

  const handleCreateTask = (): void => {
    setShowCreateForm(true);
  };

  const handleFormSubmit = (): void => {
    setShowCreateForm(false);
  };

  const handleFormCancel = (): void => {
    setShowCreateForm(false);
  };

  const handleTaskUpdate = (): void => {
    // Task updates are handled through the context
    // This callback can be used for additional UI updates if needed
  };

  const handleTaskDelete = (): void => {
    // Task deletion is handled through the context
    // This callback can be used for additional UI updates if needed
  };

  const handleCloseDetails = (): void => {
    selectTask(null);
  };

  return (
    <div className="app">
      <header className="app__header">
        <h1 className="app__title">Discovery Tree</h1>
        <button 
          className="app__create-button"
          onClick={handleCreateTask}
          data-testid="create-task-button"
        >
          Create Root Task
        </button>
      </header>

      <main className="app__main">
        <div className="app__layout">
          <div className="app__tree-panel">
            <TreeView 
              onNodeSelect={handleNodeSelect}
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
