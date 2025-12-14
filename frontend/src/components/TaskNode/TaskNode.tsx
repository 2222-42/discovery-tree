import React, { useCallback, useState, useRef, useEffect } from 'react';

import { useTaskContext } from '../../context/TaskContext.js';
import { useTreeContext } from '../../context/TreeContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import type { TaskStatus } from '../../types/task.js';
import type { TreeNode } from '../../types/tree.js';
import { InlineTaskForm } from '../InlineTaskForm/InlineTaskForm.js';
import './TaskNode.css';

interface TaskNodeProps extends BaseComponentProps {
  /** The tree node to render */
  node: TreeNode;
  /** Whether this node is currently selected */
  isSelected: boolean;
  /** Whether this node is currently expanded */
  isExpanded: boolean;
  /** Callback when the node is selected */
  onSelect?: (nodeId: string | null) => void;
}

interface ContextMenuPosition {
  x: number;
  y: number;
}

/**
 * TaskNode component that renders an individual task within the tree
 * Handles task display, selection, expand/collapse functionality, inline editing, and context menu operations
 */
export function TaskNode({
  node,
  isSelected,
  isExpanded,
  onSelect,
  className = '',
  'data-testid': testId = 'task-node'
}: TaskNodeProps): React.JSX.Element {
  const { 
    toggleNode, 
    selectNode, 
    expandedNodes, 
    selectedNodeId, 
    startInlineCreation,
    inlineCreationState 
  } = useTreeContext();
  const { updateTask, updateTaskStatus, deleteTask } = useTaskContext();
  const { task, children, level } = node;

  // Local state for inline editing and context menu
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(task.description);
  const [showContextMenu, setShowContextMenu] = useState(false);
  const [contextMenuPosition, setContextMenuPosition] = useState<ContextMenuPosition>({ x: 0, y: 0 });
  
  // Refs for managing focus and clicks
  const editInputRef = useRef<HTMLInputElement>(null);
  const contextMenuRef = useRef<HTMLDivElement>(null);
  const nodeRef = useRef<HTMLDivElement>(null);

  const hasChildren = children.length > 0;

  // Inline editing handlers - declare first to avoid hoisting issues
  const startEditing = useCallback((): void => {
    setIsEditing(true);
    setEditValue(task.description);
    setShowContextMenu(false);
  }, [task.description]);

  const cancelEditing = useCallback((): void => {
    setIsEditing(false);
    setEditValue(task.description);
  }, [task.description]);

  const saveEdit = useCallback(async (): Promise<void> => {
    if (editValue.trim() && editValue !== task.description) {
      try {
        await updateTask(task.id, editValue.trim());
      } catch {
        // Error handling is managed by TaskContext
        // Errors are displayed through the context's error state
      }
    }
    setIsEditing(false);
  }, [editValue, task.description, task.id, updateTask]);

  // Effect to focus input when editing starts
  useEffect(() => {
    if (isEditing && editInputRef.current) {
      editInputRef.current.focus();
      editInputRef.current.select();
    }
  }, [isEditing]);

  // Effect to handle clicks outside context menu
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent): void => {
      if (contextMenuRef.current && !contextMenuRef.current.contains(event.target as Node)) {
        setShowContextMenu(false);
      }
    };

    if (showContextMenu) {
      document.addEventListener('mousedown', handleClickOutside);
      return (): void => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
    
    return undefined;
  }, [showContextMenu]);

  const handleToggle = useCallback((event: React.MouseEvent): void => {
    event.stopPropagation();
    if (hasChildren) {
      toggleNode(task.id);
    }
  }, [hasChildren, toggleNode, task.id]);

  const handleSelect = useCallback((): void => {
    if (isEditing) return; // Don't select while editing
    const newSelectedId = isSelected ? null : task.id;
    selectNode(newSelectedId);
    onSelect?.(newSelectedId);
  }, [isSelected, task.id, selectNode, onSelect, isEditing]);

  const handleKeyDown = useCallback((event: React.KeyboardEvent): void => {
    if (isEditing) return; // Don't handle navigation while editing
    
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleSelect();
    } else if (event.key === 'ArrowRight' && hasChildren && !isExpanded) {
      event.preventDefault();
      toggleNode(task.id);
    } else if (event.key === 'ArrowLeft' && hasChildren && isExpanded) {
      event.preventDefault();
      toggleNode(task.id);
    } else if (event.key === 'F2' && isSelected) {
      event.preventDefault();
      startEditing();
    }
  }, [handleSelect, hasChildren, isExpanded, toggleNode, task.id, isEditing, isSelected, startEditing]);

  // Context menu handlers
  const handleContextMenu = useCallback((event: React.MouseEvent): void => {
    event.preventDefault();
    event.stopPropagation();
    
    setContextMenuPosition({ x: event.clientX, y: event.clientY });
    setShowContextMenu(true);
  }, []);

  const handleEditKeyDown = useCallback((event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.key === 'Enter') {
      event.preventDefault();
      void saveEdit();
    } else if (event.key === 'Escape') {
      event.preventDefault();
      cancelEditing();
    }
  }, [saveEdit, cancelEditing]);

  const handleEditBlur = useCallback((): void => {
    void saveEdit();
  }, [saveEdit]);

  // Context menu action handlers
  const handleEdit = useCallback((): void => {
    startEditing();
  }, [startEditing]);

  const handleDelete = useCallback(async (): Promise<void> => {
    if (window.confirm(`Are you sure you want to delete "${task.description}"?`)) {
      try {
        await deleteTask(task.id);
      } catch {
        // Error handling is managed by TaskContext
        // Errors are displayed through the context's error state
      }
    }
    setShowContextMenu(false);
  }, [task.description, task.id, deleteTask]);

  const handleAddChild = useCallback((): void => {
    startInlineCreation(task.id);
    setShowContextMenu(false);
  }, [startInlineCreation, task.id]);

  const handleStatusChange = useCallback(async (newStatus: TaskStatus): Promise<void> => {
    try {
      await updateTaskStatus(task.id, newStatus);
    } catch {
      // Error handling is managed by TaskContext
      // Errors are displayed through the context's error state
    }
    setShowContextMenu(false);
  }, [task.id, updateTaskStatus]);

  const handleDoubleClick = useCallback((event: React.MouseEvent): void => {
    event.preventDefault();
    event.stopPropagation();
    startEditing();
  }, [startEditing]);

  const getStatusIcon = (): string => {
    switch (task.status) {
      case 'TODO':
        return '○';
      case 'IN_PROGRESS':
        return '◐';
      case 'DONE':
        return '●';
      case 'ROOT_WORK_ITEM':
        return '◆';
      default:
        return '○';
    }
  };

  const getStatusColor = (): string => {
    switch (task.status) {
      case 'TODO':
        return '#6b7280';
      case 'IN_PROGRESS':
        return '#f59e0b';
      case 'DONE':
        return '#10b981';
      case 'ROOT_WORK_ITEM':
        return '#3b82f6';
      default:
        return '#6b7280';
    }
  };

  const nodeClasses = [
    'task-node',
    `task-node--level-${level.toString()}`,
    `task-node--status-${task.status.toLowerCase()}`,
    isSelected ? 'task-node--selected' : '',
    hasChildren ? 'task-node--has-children' : '',
    isEditing ? 'task-node--editing' : '',
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={nodeClasses} data-testid={testId} ref={nodeRef}>
      <div
        className="task-node__content"
        onClick={handleSelect}
        onDoubleClick={handleDoubleClick}
        onContextMenu={handleContextMenu}
        onKeyDown={handleKeyDown}
        tabIndex={0}
        role="treeitem"
        aria-selected={isSelected}
        aria-expanded={hasChildren ? isExpanded : undefined}
        style={{ paddingLeft: `${(level * 20).toString()}px` }}
      >
        {hasChildren && (
          <button
            className="task-node__toggle"
            onClick={handleToggle}
            aria-label={isExpanded ? 'Collapse' : 'Expand'}
            data-testid={`${testId}-toggle`}
          >
            {isExpanded ? '▼' : '▶'}
          </button>
        )}
        
        <span 
          className="task-node__status" 
          data-testid={`${testId}-status`}
          style={{ color: getStatusColor() }}
          title={task.status.replace('_', ' ')}
        >
          {getStatusIcon()}
        </span>
        
        {isEditing ? (
          <input
            ref={editInputRef}
            className="task-node__edit-input"
            value={editValue}
            onChange={(e): void => {
              setEditValue(e.target.value);
            }}
            onKeyDown={handleEditKeyDown}
            onBlur={handleEditBlur}
            data-testid={`${testId}-edit-input`}
          />
        ) : (
          <span className="task-node__description" data-testid={`${testId}-description`}>
            {task.description}
          </span>
        )}
        
        <span className="task-node__id" data-testid={`${testId}-id`}>
          #{task.id.slice(-6)}
        </span>
      </div>

      {/* Context Menu */}
      {showContextMenu && (
        <div
          ref={contextMenuRef}
          className="task-node__context-menu"
          style={{
            position: 'fixed',
            left: `${contextMenuPosition.x.toString()}px`,
            top: `${contextMenuPosition.y.toString()}px`,
            zIndex: 1000,
          }}
          data-testid={`${testId}-context-menu`}
        >
          <button
            className="task-node__context-menu-item"
            onClick={handleEdit}
            data-testid={`${testId}-context-edit`}
          >
            ✏️ Edit
          </button>
          
          <button
            className="task-node__context-menu-item"
            onClick={handleAddChild}
            data-testid={`${testId}-context-add-child`}
          >
            ➕ Add Child
          </button>
          
          <div className="task-node__context-menu-separator" />
          
          <div className="task-node__context-menu-submenu">
            <span className="task-node__context-menu-label">Status:</span>
            <button
              className="task-node__context-menu-item"
              onClick={() => {
                void handleStatusChange('TODO');
              }}
              data-testid={`${testId}-context-status-todo`}
            >
              ○ TODO
            </button>
            <button
              className="task-node__context-menu-item"
              onClick={() => {
                void handleStatusChange('IN_PROGRESS');
              }}
              data-testid={`${testId}-context-status-progress`}
            >
              ◐ In Progress
            </button>
            <button
              className="task-node__context-menu-item"
              onClick={() => {
                void handleStatusChange('DONE');
              }}
              data-testid={`${testId}-context-status-done`}
            >
              ● Done
            </button>
          </div>
          
          <div className="task-node__context-menu-separator" />
          
          <button
            className="task-node__context-menu-item task-node__context-menu-item--danger"
            onClick={() => {
              void handleDelete();
            }}
            data-testid={`${testId}-context-delete`}
          >
            🗑️ Delete
          </button>
        </div>
      )}

      {/* Inline task creation form */}
      {inlineCreationState.isCreating && inlineCreationState.activeParentId === task.id && (
        <InlineTaskForm
          parentId={task.id}
          level={level}
          data-testid={`${testId}-inline-form`}
        />
      )}

      {hasChildren && isExpanded && (
        <div className="task-node__children" data-testid={`${testId}-children`}>
          {children.map(childNode => (
            <TaskNode
              key={childNode.task.id}
              node={childNode}
              isSelected={selectedNodeId === childNode.task.id}
              isExpanded={expandedNodes.has(childNode.task.id)}
              {...(onSelect && { onSelect })}
              data-testid={`task-node-${childNode.task.id}`}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default TaskNode;