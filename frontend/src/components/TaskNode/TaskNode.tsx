import React, { useCallback } from 'react';

import { useTreeContext } from '../../context/TreeContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import type { TreeNode } from '../../types/tree.js';

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

/**
 * TaskNode component that renders an individual task within the tree
 * Handles task display, selection, and expand/collapse functionality
 */
export function TaskNode({
  node,
  isSelected,
  isExpanded,
  onSelect,
  className = '',
  'data-testid': testId = 'task-node'
}: TaskNodeProps): React.JSX.Element {
  const { toggleNode, selectNode, expandedNodes, selectedNodeId } = useTreeContext();
  const { task, children, level } = node;

  const hasChildren = children.length > 0;

  const handleToggle = useCallback((event: React.MouseEvent): void => {
    event.stopPropagation();
    if (hasChildren) {
      toggleNode(task.id);
    }
  }, [hasChildren, toggleNode, task.id]);

  const handleSelect = useCallback((): void => {
    const newSelectedId = isSelected ? null : task.id;
    selectNode(newSelectedId);
    onSelect?.(newSelectedId);
  }, [isSelected, task.id, selectNode, onSelect]);

  const handleKeyDown = useCallback((event: React.KeyboardEvent): void => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      handleSelect();
    } else if (event.key === 'ArrowRight' && hasChildren && !isExpanded) {
      event.preventDefault();
      toggleNode(task.id);
    } else if (event.key === 'ArrowLeft' && hasChildren && isExpanded) {
      event.preventDefault();
      toggleNode(task.id);
    }
  }, [handleSelect, hasChildren, isExpanded, toggleNode, task.id]);

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

  const nodeClasses = [
    'task-node',
    `task-node--level-${level.toString()}`,
    `task-node--status-${task.status.toLowerCase()}`,
    isSelected ? 'task-node--selected' : '',
    hasChildren ? 'task-node--has-children' : '',
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={nodeClasses} data-testid={testId}>
      <div
        className="task-node__content"
        onClick={handleSelect}
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
        
        <span className="task-node__status" data-testid={`${testId}-status`}>
          {getStatusIcon()}
        </span>
        
        <span className="task-node__description" data-testid={`${testId}-description`}>
          {task.description}
        </span>
        
        <span className="task-node__id" data-testid={`${testId}-id`}>
          #{task.id.slice(-6)}
        </span>
      </div>

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