/**
 * Test utilities for React Testing Library and property-based testing
 * Provides common testing helpers and setup functions
 */

import { render, RenderOptions } from '@testing-library/react';
import { ReactElement } from 'react';

import { TreeNode } from '@/types/tree.js';

/**
 * Custom render function that includes common providers
 */
interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  // Add provider options when needed
}

export const renderWithProviders = (
  ui: ReactElement,
  options: CustomRenderOptions = {}
): ReturnType<typeof render> => {
  return render(ui, { ...options });
};

/**
 * Helper to create mock API responses
 */
export const createMockApiResponse = <T>(data: T, delay = 0): Promise<T> => {
  return new Promise((resolve) => {
    setTimeout(() => resolve(data), delay);
  });
};

/**
 * Helper to create mock API error responses
 */
export const createMockApiError = (message: string, status = 500, delay = 0): Promise<never> => {
  return new Promise((_, reject) => {
    setTimeout(() => {
      const error = new Error(message);
      (error as Error & { status: number }).status = status;
      reject(error);
    }, delay);
  });
};

/**
 * Helper to wait for async operations in tests
 */
export const waitFor = (ms: number): Promise<void> => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

/**
 * Helper to create consistent test IDs
 */
export const createTestId = (component: string, element?: string): string => {
  return element ? `${component}-${element}` : component;
};

/**
 * Helper to validate tree structure consistency
 */
export const validateTreeStructure = (nodes: TreeNode[]): boolean => {
  const visitedIds = new Set<string>();
  
  const validateNode = (node: TreeNode, expectedLevel: number): boolean => {
    // Check for duplicate IDs
    if (visitedIds.has(node.task.id)) {
      return false;
    }
    visitedIds.add(node.task.id);
    
    // Check level consistency
    if (node.level !== expectedLevel) {
      return false;
    }
    
    // Check children
    return node.children.every((child: TreeNode, index: number) => {
      // Check parent-child relationship
      if (child.task.parentId !== node.task.id) {
        return false;
      }
      
      // Check position
      if (child.task.position !== index) {
        return false;
      }
      
      // Recursively validate children
      return validateNode(child, expectedLevel + 1);
    });
  };
  
  return nodes.every(node => validateNode(node, 0));
};

/**
 * Helper to extract all task IDs from a tree structure
 */
export const extractAllTaskIds = (nodes: TreeNode[]): string[] => {
  const ids: string[] = [];
  
  const traverse = (node: TreeNode): void => {
    ids.push(node.task.id);
    if (node.children && node.children.length > 0) {
      node.children.forEach(traverse);
    }
  };
  
  nodes.forEach(traverse);
  return ids;
};

// Re-export common testing utilities
export * from '@testing-library/react';
export * from '@testing-library/user-event';
export { default as userEvent } from '@testing-library/user-event';