import React, { useState, useCallback } from 'react';

import { useTaskContext } from '../../context/TaskContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import type { Task, TaskStatus } from '../../types/task.js';
import TaskForm from '../TaskForm/TaskForm.js';

interface TaskDetailsProps extends BaseComponentProps {
  /** The task to display details for */
  task: Task | null;
  /** Callback when task is updated */
  onTaskUpdate?: (task: Task) => void;
  /** Callback when task is deleted */
  onTaskDelete?: (taskId: string) => void;
  /** Callback when close is requested */
  onClose?: () => void;
}

/**
 * TaskDetails component that shows comprehensive task information
 * Provides editing capabilities and task status management
 */
export function TaskDetails({
  task,
  onTaskUpdate,
  onTaskDelete,
  onClose,
  className = '',
  'data-testid': testId = 'task-details'
}: TaskDetailsProps): React.JSX.Element {
  const [isEditing, setIsEditing] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  
  const { updateTaskStatus, deleteTask, loading } = useTaskContext();

  const handleStatusChange = useCallback(async (newStatus: TaskStatus): Promise<void> => {
    if (!task) return;
    
    try {
      await updateTaskStatus(task.id, newStatus);
      // The task will be updated through the context
      onTaskUpdate?.(task);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to update task status:', error);
    }
  }, [task, updateTaskStatus, onTaskUpdate]);

  const handleDelete = useCallback(async (): Promise<void> => {
    if (!task) return;
    
    const confirmed = window.confirm(
      `Are you sure you want to delete "${task.description}"? This action cannot be undone.`
    );
    
    if (!confirmed) return;

    setIsDeleting(true);
    try {
      await deleteTask(task.id);
      onTaskDelete?.(task.id);
      onClose?.();
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to delete task:', error);
    } finally {
      setIsDeleting(false);
    }
  }, [task, deleteTask, onTaskDelete, onClose]);

  const handleEditSubmit = useCallback((_taskId: string): void => {
    setIsEditing(false);
    if (task) {
      onTaskUpdate?.(task);
    }
  }, [task, onTaskUpdate]);

  const handleEditCancel = useCallback((): void => {
    setIsEditing(false);
  }, []);

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString();
  };

  const getStatusColor = (status: TaskStatus): string => {
    switch (status) {
      case 'TODO':
        return '#6b7280';
      case 'IN_PROGRESS':
        return '#f59e0b';
      case 'DONE':
        return '#10b981';
      case 'ROOT_WORK_ITEM':
        return '#8b5cf6';
      default:
        return '#6b7280';
    }
  };

  if (!task) {
    return (
      <div className={`task-details task-details--empty ${className}`} data-testid={testId}>
        <div className="task-details__empty-state">
          <p>Select a task to view its details</p>
        </div>
      </div>
    );
  }

  if (isEditing) {
    return (
      <div className={`task-details task-details--editing ${className}`} data-testid={testId}>
        <TaskForm
          mode="edit"
          taskId={task.id}
          initialDescription={task.description}
          onSubmit={handleEditSubmit}
          onCancel={handleEditCancel}
          data-testid={`${testId}-edit-form`}
        />
      </div>
    );
  }

  return (
    <div className={`task-details ${className}`} data-testid={testId}>
      <div className="task-details__header">
        <div className="task-details__title-row">
          <h2 className="task-details__title">Task Details</h2>
          {onClose && (
            <button
              className="task-details__close"
              onClick={onClose}
              aria-label="Close task details"
              data-testid={`${testId}-close`}
            >
              ×
            </button>
          )}
        </div>
        <div className="task-details__id">
          ID: {task.id}
        </div>
      </div>

      <div className="task-details__content">
        <div className="task-details__section">
          <h3 className="task-details__section-title">Description</h3>
          <p className="task-details__description" data-testid={`${testId}-description`}>
            {task.description}
          </p>
        </div>

        <div className="task-details__section">
          <h3 className="task-details__section-title">Status</h3>
          <div className="task-details__status-container">
            <select
              className="task-details__status-select"
              value={task.status}
              onChange={(e) => { 
              void handleStatusChange(e.target.value as TaskStatus); 
            }}
              disabled={loading}
              style={{ color: getStatusColor(task.status) }}
              data-testid={`${testId}-status-select`}
            >
              <option value="TODO">TODO</option>
              <option value="IN_PROGRESS">IN PROGRESS</option>
              <option value="DONE">DONE</option>
              <option value="ROOT_WORK_ITEM">ROOT WORK ITEM</option>
            </select>
          </div>
        </div>

        <div className="task-details__section">
          <h3 className="task-details__section-title">Hierarchy</h3>
          <div className="task-details__hierarchy">
            <div className="task-details__field">
              <span className="task-details__field-label">Parent ID:</span>
              <span className="task-details__field-value" data-testid={`${testId}-parent-id`}>
                {task.parentId ?? 'None (Root Task)'}
              </span>
            </div>
            <div className="task-details__field">
              <span className="task-details__field-label">Position:</span>
              <span className="task-details__field-value" data-testid={`${testId}-position`}>
                {task.position}
              </span>
            </div>
          </div>
        </div>

        <div className="task-details__section">
          <h3 className="task-details__section-title">Timestamps</h3>
          <div className="task-details__timestamps">
            <div className="task-details__field">
              <span className="task-details__field-label">Created:</span>
              <span className="task-details__field-value" data-testid={`${testId}-created-at`}>
                {formatDate(task.createdAt)}
              </span>
            </div>
            <div className="task-details__field">
              <span className="task-details__field-label">Updated:</span>
              <span className="task-details__field-value" data-testid={`${testId}-updated-at`}>
                {formatDate(task.updatedAt)}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="task-details__actions">
        <button
          className="task-details__button task-details__button--primary"
          onClick={() => {
            setIsEditing(true);
          }}
          disabled={loading}
          data-testid={`${testId}-edit`}
        >
          Edit Task
        </button>
        <button
          className="task-details__button task-details__button--danger"
          onClick={() => {
            handleDelete().catch((error: unknown) => {
              // eslint-disable-next-line no-console
              console.error('Delete failed:', error);
            });
          }}
          disabled={loading || isDeleting}
          data-testid={`${testId}-delete`}
        >
          {isDeleting ? 'Deleting...' : 'Delete Task'}
        </button>
      </div>
    </div>
  );
}

export default TaskDetails;