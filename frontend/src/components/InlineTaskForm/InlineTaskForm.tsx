import React, { useCallback, useEffect, useRef } from 'react';

import { useTreeContext } from '../../context/TreeContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import './InlineTaskForm.css';

interface InlineTaskFormProps extends BaseComponentProps {
  /** The parent task ID for the new child task */
  parentId: string;
  /** The nesting level for proper indentation */
  level: number;
}

/**
 * InlineTaskForm component for creating child tasks directly within the tree interface
 * Provides a lightweight form that appears inline with proper visual hierarchy
 */
export function InlineTaskForm({
  parentId: _parentId,
  level,
  className = '',
  'data-testid': testId = 'inline-task-form'
}: InlineTaskFormProps): React.JSX.Element {
  const {
    inlineCreationState,
    updateInlineDescription,
    completeInlineCreation,
    cancelInlineCreation,
  } = useTreeContext();

  const inputRef = useRef<HTMLInputElement>(null);

  // Focus input when component mounts
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  const handleInputChange = useCallback((event: React.ChangeEvent<HTMLInputElement>): void => {
    updateInlineDescription(event.target.value);
  }, [updateInlineDescription]);

  const handleKeyDown = useCallback((event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.key === 'Enter') {
      event.preventDefault();
      void completeInlineCreation();
    } else if (event.key === 'Escape') {
      event.preventDefault();
      cancelInlineCreation();
    }
  }, [completeInlineCreation, cancelInlineCreation]);

  const handleSave = useCallback((): void => {
    void completeInlineCreation();
  }, [completeInlineCreation]);

  const handleCancel = useCallback((): void => {
    cancelInlineCreation();
  }, [cancelInlineCreation]);

  const handleBlur = useCallback((event: React.FocusEvent<HTMLDivElement>): void => {
    // Only cancel if focus is moving outside the form container
    if (event.relatedTarget !== null && !event.currentTarget.contains(event.relatedTarget as Node)) {
      cancelInlineCreation();
    }
  }, [cancelInlineCreation]);

  const formClasses = [
    'inline-task-form',
    `inline-task-form--level-${level.toString()}`,
    className
  ].filter(Boolean).join(' ');

  return (
    <div 
      className={formClasses} 
      data-testid={testId}
      style={{ paddingLeft: `${((level + 1) * 20).toString()}px` }}
      onBlur={handleBlur}
    >
      <div className="inline-task-form__content">
        <span className="inline-task-form__status">○</span>
        
        <input
          ref={inputRef}
          className="inline-task-form__input"
          type="text"
          value={inlineCreationState.description}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          placeholder="Enter task description..."
          data-testid={`${testId}-input`}
        />
        
        <div className="inline-task-form__actions">
          <button
            className="inline-task-form__button inline-task-form__button--save"
            onClick={handleSave}
            disabled={!inlineCreationState.description.trim()}
            data-testid={`${testId}-save`}
            title="Save (Enter)"
          >
            ✓
          </button>
          <button
            className="inline-task-form__button inline-task-form__button--cancel"
            onClick={handleCancel}
            data-testid={`${testId}-cancel`}
            title="Cancel (Escape)"
          >
            ✕
          </button>
        </div>
      </div>
      
      {inlineCreationState.error !== null && inlineCreationState.error !== '' && (
        <div className="inline-task-form__error" data-testid={`${testId}-error`}>
          {inlineCreationState.error}
        </div>
      )}
    </div>
  );
}

export default InlineTaskForm;