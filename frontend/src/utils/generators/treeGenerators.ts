/**
 * Fast-check generators for TreeNode data structures
 * Used for property-based testing of tree-related functionality
 */

import * as fc from 'fast-check';

import { TreeNode } from '@/types/tree.js';

import { taskArb, rootTaskArb, childTaskArb, taskIdArb } from './taskGenerators.js';

/**
 * Generator for tree levels (depth in the tree)
 */
export const treeLevelArb = fc.integer({ min: 0, max: 5 });

/**
 * Generator for expansion state (boolean)
 */
export const expansionStateArb = fc.boolean();

/**
 * Generator for simple TreeNode without children
 * Used as building blocks for more complex tree structures
 */
export const leafTreeNodeArb = fc.record({
  task: taskArb,
  children: fc.constant([]),
  isExpanded: expansionStateArb,
  level: treeLevelArb,
});

/**
 * Generator for TreeNode with potential children
 * Creates realistic tree structures with proper parent-child relationships
 */
export const treeNodeArb: fc.Arbitrary<TreeNode> = fc.letrec(tie => ({
  node: fc.record({
    task: taskArb,
    children: fc.oneof(
      fc.constant([]),
      fc.array(tie('node'), { minLength: 0, maxLength: 3 })
    ),
    isExpanded: expansionStateArb,
    level: treeLevelArb,
  }).map(node => {
    // Ensure children have correct level and parent relationships
    const adjustedChildren = (node.children as TreeNode[]).map((child, index) => ({
      ...child,
      level: node.level + 1,
      task: {
        ...child.task,
        parentId: node.task.id,
        position: index
      }
    }));
    
    return {
      ...node,
      children: adjustedChildren
    };
  })
})).node;

/**
 * Generator for root TreeNode (level 0, no parent)
 */
export const rootTreeNodeArb = fc.record({
  task: rootTaskArb,
  children: fc.array(treeNodeArb, { minLength: 0, maxLength: 4 }),
  isExpanded: expansionStateArb,
  level: fc.constant(0),
}).map(node => ({
  ...node,
  children: (node.children as TreeNode[]).map((child, index) => ({
    ...child,
    level: 1,
    task: {
      ...child.task,
      parentId: node.task.id,
      position: index
    }
  }))
}));

/**
 * Generator for arrays of root TreeNodes (forest structure)
 */
export const treeForestArb = fc.array(rootTreeNodeArb, { minLength: 1, maxLength: 5 });

/**
 * Generator for Set of expanded node IDs
 */
export const expandedNodesSetArb = fc.array(taskIdArb, { minLength: 0, maxLength: 10 })
  .map(ids => new Set(ids));

/**
 * Generator for TreeState objects
 */
export const treeStateArb = fc.record({
  expandedNodes: expandedNodesSetArb,
  selectedNodeId: fc.oneof(fc.constant(null), taskIdArb),
  rootNodes: treeForestArb,
});

/**
 * Generator for TreeDragOperation objects
 */
export const treeDragOperationArb = fc.record({
  draggedTaskId: taskIdArb,
  targetParentId: fc.oneof(fc.constant(null), taskIdArb),
  targetPosition: fc.integer({ min: 0, max: 20 }),
});

/**
 * Generator for valid tree structures (no cycles, proper hierarchy)
 * Creates trees that satisfy all structural constraints
 */
export const validTreeArb = fc.integer({ min: 1, max: 4 }).chain(depth => {
  // Use a counter to ensure unique IDs
  let idCounter = 0;
  
  const createTreeAtLevel = (level: number, parentId: string | null): fc.Arbitrary<TreeNode> => {
    if (level >= depth) {
      // Leaf node
      return fc.record({
        task: parentId ? childTaskArb : rootTaskArb,
        children: fc.constant([]),
        isExpanded: fc.constant(false),
        level: fc.constant(level),
      }).map(node => {
        const uniqueId = `task-${++idCounter}-${Date.now()}-${Math.random().toString(36).substring(2, 11)}`;
        return {
          ...node,
          task: {
            ...node.task,
            id: uniqueId,
            parentId
          }
        };
      });
    }
    
    // Internal node with children
    return fc.record({
      task: parentId ? childTaskArb : rootTaskArb,
      childCount: fc.integer({ min: 0, max: 2 }),
      isExpanded: expansionStateArb,
      level: fc.constant(level),
    }).chain(nodeData => {
      const uniqueId = `task-${++idCounter}-${Date.now()}-${Math.random().toString(36).substring(2, 11)}`;
      const nodeTask = {
        ...nodeData.task,
        id: uniqueId,
        parentId
      };
      
      if (nodeData.childCount === 0) {
        return fc.constant({
          task: nodeTask,
          children: [],
          isExpanded: false,
          level: nodeData.level,
        });
      }
      
      return fc.array(
        createTreeAtLevel(level + 1, nodeTask.id),
        { minLength: nodeData.childCount, maxLength: nodeData.childCount }
      ).map(children => ({
        task: nodeTask,
        children: children.map((child, index) => ({
          ...child,
          task: {
            ...child.task,
            position: index
          }
        })),
        isExpanded: nodeData.isExpanded,
        level: nodeData.level,
      }));
    });
  };
  
  return createTreeAtLevel(0, null);
});

/**
 * Generator for tree arrays with multiple valid trees
 */
export const validTreeArrayArb = fc.array(validTreeArb, { minLength: 1, maxLength: 4 });

/**
 * Generator for empty tree state
 */
export const emptyTreeStateArb = fc.constant({
  expandedNodes: new Set<string>(),
  selectedNodeId: null,
  rootNodes: [],
});

/**
 * Generator for tree paths (arrays of task IDs representing path from root to node)
 */
export const treePathArb = fc.array(taskIdArb, { minLength: 1, maxLength: 6 });

/**
 * Utility function to extract all task IDs from a tree structure
 */
export const extractTaskIds = (nodes: TreeNode[]): string[] => {
  const ids: string[] = [];
  
  const traverse = (node: TreeNode) => {
    ids.push(node.task.id);
    node.children.forEach(traverse);
  };
  
  nodes.forEach(traverse);
  return ids;
};

/**
 * Generator for TreeState with consistent expanded nodes
 * Ensures expanded node IDs actually exist in the tree
 */
export const consistentTreeStateArb = treeForestArb.chain(rootNodes => {
  const allTaskIds = extractTaskIds(rootNodes);
  
  return fc.record({
    expandedNodes: fc.subarray(allTaskIds).map(ids => new Set(ids)),
    selectedNodeId: fc.oneof(
      fc.constant(null),
      allTaskIds.length > 0 ? fc.constantFrom(...allTaskIds) : fc.constant(null)
    ),
    rootNodes: fc.constant(rootNodes),
  });
});