import React from 'react';

import { useTreeContext } from '../../context/TreeContext.js';
import type { BaseComponentProps } from '../../types/app.js';
import TaskNode from '../TaskNode/TaskNode.js';

interface TreeViewProps extends BaseComponentProps {
  /** Optional callback when a node is selected */
  onNodeSelect?: (nodeId: string | null) => void;
}

/**
 * TreeView component that renders the hierarchical discovery tree
 * Displays tasks in a tree structure with expand/collapse functionality
 */
export function TreeView({ 
  className = '', 
  'data-testid': testId = 'tree-view',
  onNodeSelect 
}: TreeViewProps): React.JSX.Element {
  const { rootNodes, selectedNodeId, expandedNodes } = useTreeContext();

  const handleNodeSelect = (nodeId: string | null): void => {
    onNodeSelect?.(nodeId);
  };

  if (rootNodes.length === 0) {
    return (
      <div className={`tree-view tree-view--empty ${className}`} data-testid={testId}>
        <div className="tree-view__empty-state">
          <p>No tasks available. Create your first task to get started.</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`tree-view ${className}`} data-testid={testId}>
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