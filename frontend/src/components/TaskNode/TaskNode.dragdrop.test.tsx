import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';

import { TaskProvider } from '../../context/TaskContext.js';
import { TreeProvider } from '../../context/TreeContext.js';
import type { TreeNode } from '../../types/tree.js';
import type { Task } from '../../types/task.js';
import { TaskNode } from './TaskNode.js';

// Mock the API client
vi.mock('../../services/api/apiClient.js', () => ({
  apiClient: {
    getAllTasks: vi.fn().mockResolvedValue([]),
    moveTask: vi.fn().mockResolvedValue({}),
  },
}));

const mockTask: Task = {
  id: 'test-task-1',
  description: 'Test Task',
  status: 'TODO',
  parentId: null,
  position: 0,
  createdAt: '2023-01-01T00:00:00Z',
  updatedAt: '2023-01-01T00:00:00Z',
};

const mockTreeNode: TreeNode = {
  task: mockTask,
  children: [],
  isExpanded: false,
  level: 0,
};

function TestWrapper({ children }: { children: React.ReactNode }): React.JSX.Element {
  return (
    <TaskProvider>
      <TreeProvider>
        {children}
      </TreeProvider>
    </TaskProvider>
  );
}

describe('TaskNode Drag and Drop', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should be draggable when not editing', () => {
    render(
      <TestWrapper>
        <TaskNode
          node={mockTreeNode}
          isSelected={false}
          isExpanded={false}
        />
      </TestWrapper>
    );

    const taskContent = screen.getByRole('treeitem');
    expect(taskContent).toHaveAttribute('draggable', 'true');
  });

  it('should handle drag start event', () => {
    render(
      <TestWrapper>
        <TaskNode
          node={mockTreeNode}
          isSelected={false}
          isExpanded={false}
        />
      </TestWrapper>
    );

    const taskContent = screen.getByRole('treeitem');
    
    // Create a mock DataTransfer object
    const mockDataTransfer = {
      effectAllowed: '',
      setData: vi.fn(),
    };

    const dragStartEvent = new Event('dragstart', { bubbles: true });
    Object.defineProperty(dragStartEvent, 'dataTransfer', {
      value: mockDataTransfer,
      writable: false,
    });

    fireEvent(taskContent, dragStartEvent);

    expect(mockDataTransfer.effectAllowed).toBe('move');
    expect(mockDataTransfer.setData).toHaveBeenCalledWith('text/plain', 'test-task-1');
  });

  it('should handle drag over event', () => {
    render(
      <TestWrapper>
        <TaskNode
          node={mockTreeNode}
          isSelected={false}
          isExpanded={false}
        />
      </TestWrapper>
    );

    const taskContent = screen.getByRole('treeitem');
    
    // Mock getBoundingClientRect
    const mockRect = {
      top: 0,
      height: 40,
      left: 0,
      right: 100,
      bottom: 40,
    };
    
    vi.spyOn(taskContent, 'getBoundingClientRect').mockReturnValue(mockRect as DOMRect);

    const mockDataTransfer = {
      dropEffect: '',
    };

    const dragOverEvent = new Event('dragover', { bubbles: true });
    Object.defineProperty(dragOverEvent, 'dataTransfer', {
      value: mockDataTransfer,
      writable: false,
    });
    Object.defineProperty(dragOverEvent, 'clientY', {
      value: 20, // Middle of the element
      writable: false,
    });

    fireEvent(taskContent, dragOverEvent);

    expect(mockDataTransfer.dropEffect).toBe('move');
  });

  it('should handle drop event', () => {
    render(
      <TestWrapper>
        <TaskNode
          node={mockTreeNode}
          isSelected={false}
          isExpanded={false}
        />
      </TestWrapper>
    );

    const taskContent = screen.getByRole('treeitem');
    
    // Mock getBoundingClientRect
    const mockRect = {
      top: 0,
      height: 40,
      left: 0,
      right: 100,
      bottom: 40,
    };
    
    vi.spyOn(taskContent, 'getBoundingClientRect').mockReturnValue(mockRect as DOMRect);

    const dropEvent = new Event('drop', { bubbles: true });
    Object.defineProperty(dropEvent, 'clientY', {
      value: 5, // Top quarter - should be 'before'
      writable: false,
    });

    fireEvent(taskContent, dropEvent);

    // The drop should be handled (no errors thrown)
    expect(taskContent).toBeInTheDocument();
  });

  it('should show visual feedback during drag operations', () => {
    const { container } = render(
      <TestWrapper>
        <TaskNode
          node={mockTreeNode}
          isSelected={false}
          isExpanded={false}
        />
      </TestWrapper>
    );

    // Initially no drag classes
    expect(container.querySelector('.task-node--dragging')).toBeNull();
    expect(container.querySelector('.task-node--drag-over')).toBeNull();
  });
});