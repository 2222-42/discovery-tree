import React, { useCallback, useEffect } from 'react';

import { useTreeContext } from '../../context/TreeContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import TaskNode from '../TaskNode/TaskNode.js';

import './TreeView.css';

interface TreeViewProps extends BaseComponentProps {
  /** Optional callback when a node is selected */
  onNodeSelect?: (nodeId: string | null) => void;
  /** Whether to show loading state */
  loading?: boolean;
  /** Error message to display */
  error?: string | null;
}

/**
 * TreeView component that renders the hierarchical discovery tree
 * Displays tasks in a tree structure with expand/collapse functionality
 * Requirements: 3.1, 3.2, 3.3, 3.4
 */
export function TreeView({ 
  className = '', 
  'data-testid': testId = 'tree-view',
  onNodeSelect,
  loading = false,
  error = null
}: TreeViewProps): React.JSX.Element {
  const { rootNodes, selectedNodeId, expandedNodes, selectNode } = useTreeContext();

  const handleNodeSelect = useCallback((nodeId: string | null): void => {
    selectNode(nodeId);
    onNodeSelect?.(nodeId);
  }, [selectNode, onNodeSelect]);

  // Handle keyboard navigation for the tree
  const handleKeyDown = useCallback((event: React.KeyboardEvent): void => {
    if (event.key === 'Escape' && selectedNodeId !== null) {
      handleNodeSelect(null);
    }
  }, [selectedNodeId, handleNodeSelect]);

  // Auto-focus the tree when it loads
  useEffect(() => {
    const treeElement = document.querySelector(`[data-testid="${testId}"]`);
    if (treeElement && rootNodes.length > 0) {
      (treeElement as HTMLElement).focus();
    }
  }, [testId, rootNodes.length]);

  // Render loading state
  if (loading) {
    return (
      <div className={`tree-view tree-view--loading ${className}`} data-testid={testId}>
        <div className="tree-view__loading-state">
          <p>Loading tasks...</p>
        </div>
      </div>
    );
  }

  // Render error state
  if (error !== null && error !== '') {
    return (
      <div className={`tree-view tree-view--error ${className}`} data-testid={testId}>
        <div className="tree-view__error-state">
          <p>Error loading tasks: {error}</p>
        </div>
      </div>
    );
  }

  // Render empty state
  if (rootNodes.length === 0) {
    return (
      <div className={`tree-view tree-view--empty ${className}`} data-testid={testId}>
        <div className="tree-view__empty-state">
          <p>No tasks available. Create your first task to get started.</p>
        </div>
      </div>
    );
  }

  // Render tree with hierarchical structure
  return (
    <div 
      className={`tree-view ${className}`} 
      data-testid={testId}
      onKeyDown={handleKeyDown}
      tabIndex={-1}
      role="tree"
      aria-label="Discovery Tree"
    >
      <div className="tree-view__content">
        {rootNodes.map(node => (
          <TaskNode
            key={node.task.id}
            node={node}
            isSelected={selectedNodeId === node.task.id}
            isExpanded={expandedNodes.has(node.task.id)}
            onSelect={handleNodeSelect}
            data-testid={`tree-node-${node.task.id}`}
          />
        ))}
      </div>
    </div>
  );
}

export default TreeView;