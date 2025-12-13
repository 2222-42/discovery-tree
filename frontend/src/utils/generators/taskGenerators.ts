/**
 * Fast-check generators for Task data structures
 * Used for property-based testing of task-related functionality
 */

import * as fc from 'fast-check';
import { Task, TaskStatus } from '@/types/task.js';

/**
 * Generator for valid task IDs
 * Generates UUIDs and simple alphanumeric IDs
 */
export const taskIdArb = fc.oneof(
  // UUID-like format
  fc.uuid(),
  // Simple alphanumeric IDs
  fc.stringMatching(/^[a-zA-Z0-9]{8,16}$/)
);

/**
 * Generator for task descriptions
 * Creates realistic task descriptions with various lengths and content
 */
export const taskDescriptionArb = fc.oneof(
  // Short descriptions
  fc.stringMatching(/^[A-Z][a-z\s]{10,50}$/),
  // Medium descriptions with punctuation
  fc.stringMatching(/^[A-Z][a-zA-Z0-9\s\.,!?-]{20,100}$/),
  // Longer descriptions
  fc.string({ minLength: 5, maxLength: 200 }).filter(s => s.trim().length > 0)
);

/**
 * Generator for task status values
 */
export const taskStatusArb = fc.constantFrom(
  TaskStatus.TODO,
  TaskStatus.IN_PROGRESS,
  TaskStatus.DONE,
  TaskStatus.ROOT_WORK_ITEM
);

/**
 * Generator for ISO timestamp strings
 */
export const isoTimestampArb = fc.date({ min: new Date('2020-01-01T00:00:00.000Z'), max: new Date('2030-12-31T23:59:59.999Z') })
  .map(date => {
    // Ensure the date is valid before converting to ISO string
    if (isNaN(date.getTime())) {
      return new Date().toISOString();
    }
    return date.toISOString();
  });

/**
 * Generator for task positions (non-negative integers)
 */
export const taskPositionArb = fc.integer({ min: 0, max: 1000 });

/**
 * Generator for parent task IDs (can be null for root tasks)
 */
export const parentIdArb = fc.oneof(
  fc.constant(null),
  taskIdArb
);

/**
 * Generator for complete Task objects
 * Creates valid task instances with all required properties
 */
export const taskArb = fc.record({
  id: taskIdArb,
  description: taskDescriptionArb,
  status: taskStatusArb,
  parentId: parentIdArb,
  position: taskPositionArb,
  createdAt: isoTimestampArb,
  updatedAt: isoTimestampArb,
}).map(task => ({
  ...task,
  // Ensure updatedAt is after or equal to createdAt
  updatedAt: new Date(Math.max(
    new Date(task.createdAt).getTime(),
    new Date(task.updatedAt).getTime()
  )).toISOString()
}));

/**
 * Generator for root tasks (parentId is always null)
 */
export const rootTaskArb = taskArb.map(task => ({
  ...task,
  parentId: null,
  status: TaskStatus.ROOT_WORK_ITEM
}));

/**
 * Generator for child tasks (parentId is never null)
 */
export const childTaskArb = fc.record({
  id: taskIdArb,
  description: taskDescriptionArb,
  status: fc.constantFrom(TaskStatus.TODO, TaskStatus.IN_PROGRESS, TaskStatus.DONE),
  parentId: taskIdArb,
  position: taskPositionArb,
  createdAt: isoTimestampArb,
  updatedAt: isoTimestampArb,
}).map(task => ({
  ...task,
  updatedAt: new Date(Math.max(
    new Date(task.createdAt).getTime(),
    new Date(task.updatedAt).getTime()
  )).toISOString()
}));

/**
 * Generator for arrays of tasks with valid parent-child relationships
 * Creates hierarchical task structures that form valid trees
 */
export const taskArrayArb = fc.integer({ min: 1, max: 10 }).chain(size => {
  return fc.array(taskArb, { minLength: size, maxLength: size }).map(tasks => {
    // Ensure we have at least one root task
    const rootTask = { ...tasks[0], parentId: null, status: TaskStatus.ROOT_WORK_ITEM };
    const remainingTasks = tasks.slice(1);
    
    // Make some tasks children of existing tasks
    const validTasks = [rootTask];
    
    remainingTasks.forEach((task, index) => {
      if (index < remainingTasks.length / 2 && validTasks.length > 0) {
        // Make this task a child of a random existing task
        const parentIndex = Math.floor(Math.random() * validTasks.length);
        const parent = validTasks[parentIndex];
        validTasks.push({
          ...task,
          parentId: parent.id,
          status: fc.sample(fc.constantFrom(TaskStatus.TODO, TaskStatus.IN_PROGRESS, TaskStatus.DONE), 1)[0]
        });
      } else {
        // Keep as root or make it a child with some probability
        const shouldBeChild = Math.random() > 0.3 && validTasks.length > 0;
        if (shouldBeChild) {
          const parentIndex = Math.floor(Math.random() * validTasks.length);
          const parent = validTasks[parentIndex];
          validTasks.push({
            ...task,
            parentId: parent.id,
            status: fc.sample(fc.constantFrom(TaskStatus.TODO, TaskStatus.IN_PROGRESS, TaskStatus.DONE), 1)[0]
          });
        } else {
          validTasks.push({
            ...task,
            parentId: null,
            status: TaskStatus.ROOT_WORK_ITEM
          });
        }
      }
    });
    
    return validTasks;
  });
});

/**
 * Generator for task arrays that form a single tree (one root with children)
 */
export const singleTreeTaskArrayArb = fc.integer({ min: 1, max: 15 }).chain(size => {
  return fc.record({
    rootTask: rootTaskArb,
    childTasks: fc.array(childTaskArb, { minLength: 0, maxLength: size - 1 })
  }).map(({ rootTask, childTasks }) => {
    // Make all child tasks reference the root or other children
    const allTasks = [rootTask];
    
    childTasks.forEach(child => {
      // Randomly assign parent from existing tasks
      const possibleParents = allTasks;
      const parentIndex = Math.floor(Math.random() * possibleParents.length);
      const parent = possibleParents[parentIndex];
      
      allTasks.push({
        ...child,
        parentId: parent.id
      });
    });
    
    return allTasks;
  });
});