/**
 * Tests for the test generators to ensure they produce valid data
 * This validates that our property-based testing infrastructure is working
 */

import { describe, it, expect } from 'vitest';
import * as fc from 'fast-check';
import {
  taskArb,
  rootTaskArb,
  childTaskArb,
  taskArrayArb,
  treeNodeArb,
  validTreeArb,
  consistentTreeStateArb,
} from './index.js';
import { TaskStatus } from '@/types/task.js';
import { validateTreeStructure } from '@/utils/testUtils.js';

describe('Task Generators', () => {
  it('should generate valid task objects', () => {
    fc.assert(
      fc.property(taskArb, (task) => {
        // Basic structure validation
        expect(task).toHaveProperty('id');
        expect(task).toHaveProperty('description');
        expect(task).toHaveProperty('status');
        expect(task).toHaveProperty('parentId');
        expect(task).toHaveProperty('position');
        expect(task).toHaveProperty('createdAt');
        expect(task).toHaveProperty('updatedAt');

        // Type validation
        expect(typeof task.id).toBe('string');
        expect(typeof task.description).toBe('string');
        expect(Object.values(TaskStatus)).toContain(task.status);
        expect(typeof task.position).toBe('number');
        expect(task.position).toBeGreaterThanOrEqual(0);

        // Date validation
        expect(() => new Date(task.createdAt)).not.toThrow();
        expect(() => new Date(task.updatedAt)).not.toThrow();
        expect(new Date(task.updatedAt).getTime()).toBeGreaterThanOrEqual(
          new Date(task.createdAt).getTime()
        );

        // Description should not be empty
        expect(task.description.trim().length).toBeGreaterThan(0);
      }),
      { numRuns: 50 }
    );
  });

  it('should generate valid root tasks', () => {
    fc.assert(
      fc.property(rootTaskArb, (task) => {
        expect(task.parentId).toBeNull();
        expect(task.status).toBe(TaskStatus.ROOT_WORK_ITEM);
      }),
      { numRuns: 50 }
    );
  });

  it('should generate valid child tasks', () => {
    fc.assert(
      fc.property(childTaskArb, (task) => {
        expect(task.parentId).not.toBeNull();
        expect(task.parentId).toBeTruthy();
        expect([TaskStatus.TODO, TaskStatus.IN_PROGRESS, TaskStatus.DONE]).toContain(task.status);
      }),
      { numRuns: 50 }
    );
  });

  it('should generate valid task arrays with at least one root', () => {
    fc.assert(
      fc.property(taskArrayArb, (tasks) => {
        expect(tasks.length).toBeGreaterThan(0);
        
        // Should have at least one root task
        const rootTasks = tasks.filter(task => task.parentId === null);
        expect(rootTasks.length).toBeGreaterThan(0);
        
        // All root tasks should have ROOT_WORK_ITEM status
        rootTasks.forEach(task => {
          expect(task.status).toBe(TaskStatus.ROOT_WORK_ITEM);
        });
        
        // All task IDs should be unique
        const ids = tasks.map(task => task.id);
        const uniqueIds = new Set(ids);
        expect(uniqueIds.size).toBe(ids.length);
      }),
      { numRuns: 30 }
    );
  });
});

describe('Tree Generators', () => {
  it('should generate valid tree nodes', () => {
    fc.assert(
      fc.property(treeNodeArb, (node) => {
        // Basic structure validation
        expect(node).toHaveProperty('task');
        expect(node).toHaveProperty('children');
        expect(node).toHaveProperty('isExpanded');
        expect(node).toHaveProperty('level');

        // Type validation
        expect(Array.isArray(node.children)).toBe(true);
        expect(typeof node.isExpanded).toBe('boolean');
        expect(typeof node.level).toBe('number');
        expect(node.level).toBeGreaterThanOrEqual(0);

        // Children should have correct parent relationship
        node.children.forEach((child, index) => {
          expect(child.task.parentId).toBe(node.task.id);
          expect(child.task.position).toBe(index);
          expect(child.level).toBe(node.level + 1);
        });
      }),
      { numRuns: 30 }
    );
  });

  it('should generate valid tree structures', () => {
    fc.assert(
      fc.property(validTreeArb, (tree) => {
        // Root should have level 0 and no parent
        expect(tree.level).toBe(0);
        expect(tree.task.parentId).toBeNull();

        // Validate entire tree structure
        expect(validateTreeStructure([tree])).toBe(true);
      }),
      { numRuns: 20 }
    );
  });

  it('should generate consistent tree states', () => {
    fc.assert(
      fc.property(consistentTreeStateArb, (treeState) => {
        // Basic structure validation
        expect(treeState).toHaveProperty('expandedNodes');
        expect(treeState).toHaveProperty('selectedNodeId');
        expect(treeState).toHaveProperty('rootNodes');

        // Expanded nodes should be a Set
        expect(treeState.expandedNodes instanceof Set).toBe(true);

        // If there's a selected node, it should exist in the tree
        if (treeState.selectedNodeId) {
          const allIds = new Set<string>();
          const collectIds = (node: any) => {
            allIds.add(node.task.id);
            node.children.forEach(collectIds);
          };
          treeState.rootNodes.forEach(collectIds);
          
          expect(allIds.has(treeState.selectedNodeId)).toBe(true);
        }

        // All expanded nodes should exist in the tree
        const allIds = new Set<string>();
        const collectIds = (node: any) => {
          allIds.add(node.task.id);
          node.children.forEach(collectIds);
        };
        treeState.rootNodes.forEach(collectIds);

        treeState.expandedNodes.forEach(expandedId => {
          expect(allIds.has(expandedId)).toBe(true);
        });
      }),
      { numRuns: 20 }
    );
  });
});