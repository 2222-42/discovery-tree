import React, { useState, useCallback } from 'react';

import { useTaskContext } from '../../context/TaskContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import './TaskForm.css';

interface TaskFormProps extends BaseComponentProps {
  /** Parent task ID for creating child tasks, null for root tasks */
  parentId?: string | null;
  /** Initial description value for editing */
  initialDescription?: string;
  /** Task ID when editing an existing task */
  taskId?: string;
  /** Form mode - create or edit */
  mode?: 'create' | 'edit';
  /** Callback when form is submitted successfully */
  onSubmit?: (taskId: string) => void;
  /** Callback when form is cancelled */
  onCancel?: () => void;
}

/**
 * TaskForm component for creating and editing tasks
 * Provides form validation and integrates with TaskContext for operations
 */
export function TaskForm({
  parentId = null,
  initialDescription = '',
  taskId,
  mode = 'create',
  onSubmit,
  onCancel,
  className = '',
  'data-testid': testId = 'task-form'
}: TaskFormProps): React.JSX.Element {
  const [description, setDescription] = useState(initialDescription);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { createRootTask, createChildTask, updateTask, loading } = useTaskContext();

  const validateDescription = (desc: string): string | null => {
    const trimmed = desc.trim();
    if (!trimmed) {
      return 'Task description is required';
    }
    if (trimmed.length < 3) {
      return 'Task description must be at least 3 characters';
    }
    if (trimmed.length > 500) {
      return 'Task description must be less than 500 characters';
    }
    return null;
  };

  const handleSubmit = useCallback(async (event: React.FormEvent): Promise<void> => {
    event.preventDefault();
    
    const validationError = validateDescription(description);
    if (validationError !== null) {
      setError(validationError);
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      if (mode === 'edit' && taskId !== undefined) {
        await updateTask(taskId, description.trim());
        if (onSubmit) {
          onSubmit(taskId);
        }
      } else {
        // Create mode
        if (parentId !== null) {
          await createChildTask(description.trim(), parentId);
        } else {
          await createRootTask(description.trim());
        }
        // For create mode, we don't have the task ID immediately
        // The onSubmit callback will be called without an ID
        if (onSubmit) {
          onSubmit('');
        }
        setDescription(''); // Clear form after successful creation
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save task');
    } finally {
      setIsSubmitting(false);
    }
  }, [description, mode, taskId, parentId, updateTask, createChildTask, createRootTask, onSubmit]);

  const handleCancel = useCallback((): void => {
    setDescription(initialDescription);
    setError(null);
    onCancel?.();
  }, [initialDescription, onCancel]);

  const handleDescriptionChange = useCallback((event: React.ChangeEvent<HTMLTextAreaElement>): void => {
    setDescription(event.target.value);
    if (error !== null) {
      setError(null); // Clear error when user starts typing
    }
  }, [error]);

  const isLoading = loading || isSubmitting;
  const isRootTask = parentId === null;
  const formTitle = mode === 'edit' 
    ? 'Edit Task' 
    : isRootTask 
      ? 'Create Root Task' 
      : 'Create Child Task';

  return (
    <form 
      className={`task-form ${className}`} 
      onSubmit={(e) => { void handleSubmit(e); }}
      data-testid={testId}
    >
      <div className="task-form__header">
        <h3 className="task-form__title">{formTitle}</h3>
      </div>

      <div className="task-form__body">
        <div className="task-form__field">
          <label htmlFor="task-description" className="task-form__label">
            Description *
          </label>
          <textarea
            id="task-description"
            className={`task-form__textarea ${error !== null ? 'task-form__textarea--error' : ''}`}
            value={description}
            onChange={handleDescriptionChange}
            placeholder="Enter task description..."
            rows={3}
            disabled={isLoading}
            data-testid={`${testId}-description`}
          />
          {error !== null && (
            <div className="task-form__error" data-testid={`${testId}-error`}>
              {error}
            </div>
          )}
        </div>

        {parentId !== null && (
          <div className="task-form__info">
            <small className="task-form__parent-info">
              Creating child task under parent: {parentId.slice(-6)}
            </small>
          </div>
        )}
      </div>

      <div className="task-form__actions">
        <button
          type="button"
          className="task-form__button task-form__button--secondary"
          onClick={handleCancel}
          disabled={isLoading}
          data-testid={`${testId}-cancel`}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="task-form__button task-form__button--primary"
          disabled={isLoading || description.trim().length === 0}
          data-testid={`${testId}-submit`}
        >
          {isLoading ? 'Saving...' : mode === 'edit' ? 'Update Task' : 'Create Task'}
        </button>
      </div>
    </form>
  );
}

export default TaskForm;