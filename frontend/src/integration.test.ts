/**
 * Integration tests for frontend-backend API integration
 * Tests the complete user workflows from task creation to deletion
 */

import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { DiscoveryTreeApiClient } from '@/services/api/apiClient.js';
import { TaskStatus } from '@/types/task.js';

describe('Frontend-Backend Integration Tests', () => {
  let apiClient: DiscoveryTreeApiClient;
  let createdTaskIds: string[] = [];

  beforeAll(() => {
    // Initialize API client with test configuration
    apiClient = new DiscoveryTreeApiClient({
      baseURL: 'http://localhost:8080/api/v1',
      timeout: 5000,
    });
  });

  afterAll(async () => {
    // Clean up created tasks
    for (const taskId of createdTaskIds) {
      try {
        await apiClient.deleteTask(taskId);
      } catch (error) {
        // Ignore cleanup errors
        console.warn(`Failed to cleanup task ${taskId}:`, error);
      }
    }
  });

  it('should successfully connect to the backend API', async () => {
    // Test basic connectivity by fetching all tasks
    const tasks = await apiClient.getAllTasks();
    expect(Array.isArray(tasks)).toBe(true);
  });

  it('should get the existing root task successfully', async () => {
    const rootTask = await apiClient.getRootTask();
    
    expect(rootTask).toBeDefined();
    expect(rootTask.status).toBe(TaskStatus.ROOT_WORK_ITEM);
    expect(rootTask.parentId).toBeNull();
    expect(rootTask.id).toBeTruthy();
  });

  it('should create a child task successfully', async () => {
    // Get the existing root task
    const parentTask = await apiClient.getRootTask();
    
    // Then create a child task
    const childDescription = 'Integration Test Child Task';
    const childTask = await apiClient.createChildTask(childDescription, parentTask.id);
    
    expect(childTask).toBeDefined();
    expect(childTask.description).toBe(childDescription);
    expect(childTask.parentId).toBe(parentTask.id);
    expect(childTask.status).not.toBe(TaskStatus.ROOT_WORK_ITEM);
    expect(childTask.id).toBeTruthy();
    
    createdTaskIds.push(childTask.id);
  });

  it('should retrieve a specific task by ID', async () => {
    // Get the existing root task
    const rootTask = await apiClient.getRootTask();
    
    // Create a child task first
    const originalTask = await apiClient.createChildTask('Task for Retrieval Test', rootTask.id);
    createdTaskIds.push(originalTask.id);
    
    // Retrieve it by ID
    const retrievedTask = await apiClient.getTask(originalTask.id);
    
    expect(retrievedTask).toBeDefined();
    expect(retrievedTask.id).toBe(originalTask.id);
    expect(retrievedTask.description).toBe(originalTask.description);
    expect(retrievedTask.status).toBe(originalTask.status);
  });

  it('should update a task description successfully', async () => {
    // Get the existing root task
    const rootTask = await apiClient.getRootTask();
    
    // Create a child task first
    const originalTask = await apiClient.createChildTask('Original Description', rootTask.id);
    createdTaskIds.push(originalTask.id);
    
    // Update the description
    const newDescription = 'Updated Description';
    const updatedTask = await apiClient.updateTask(originalTask.id, newDescription);
    
    expect(updatedTask).toBeDefined();
    expect(updatedTask.id).toBe(originalTask.id);
    expect(updatedTask.description).toBe(newDescription);
    expect(new Date(updatedTask.updatedAt).getTime()).toBeGreaterThan(
      new Date(originalTask.updatedAt).getTime()
    );
  });

  it('should update task status successfully', async () => {
    // Get the existing root task
    const parentTask = await apiClient.getRootTask();
    
    const childTask = await apiClient.createChildTask('Child for Status Test', parentTask.id);
    createdTaskIds.push(childTask.id);
    
    // Update the status
    const newStatus = TaskStatus.IN_PROGRESS;
    const updatedTask = await apiClient.updateTaskStatus(childTask.id, newStatus);
    
    expect(updatedTask).toBeDefined();
    expect(updatedTask.id).toBe(childTask.id);
    expect(updatedTask.status).toBe(newStatus);
  });

  it('should get task children successfully', async () => {
    // Get the existing root task
    const parentTask = await apiClient.getRootTask();
    
    // Create multiple child tasks
    const child1 = await apiClient.createChildTask('Child 1', parentTask.id);
    const child2 = await apiClient.createChildTask('Child 2', parentTask.id);
    createdTaskIds.push(child1.id, child2.id);
    
    // Get children
    const children = await apiClient.getTaskChildren(parentTask.id);
    
    expect(Array.isArray(children)).toBe(true);
    expect(children.length).toBeGreaterThanOrEqual(2);
    expect(children.some(child => child.id === child1.id)).toBe(true);
    expect(children.some(child => child.id === child2.id)).toBe(true);
  });

  it('should delete a task successfully', async () => {
    // Get the existing root task
    const rootTask = await apiClient.getRootTask();
    
    // Create a child task first
    const task = await apiClient.createChildTask('Task to Delete', rootTask.id);
    
    // Delete it
    await apiClient.deleteTask(task.id);
    
    // Verify it's deleted by trying to retrieve it
    try {
      await apiClient.getTask(task.id);
      // If we get here, the task wasn't deleted
      expect.fail('Task should have been deleted');
    } catch (error) {
      // Expected - task should not be found
      expect(error).toBeDefined();
    }
  });

  it('should handle error cases gracefully', async () => {
    // Test getting a non-existent task
    try {
      await apiClient.getTask('non-existent-id');
      expect.fail('Should have thrown an error for non-existent task');
    } catch (error) {
      expect(error).toBeDefined();
    }
    
    // Test creating a child task with non-existent parent
    try {
      await apiClient.createChildTask('Orphan Task', 'non-existent-parent-id');
      expect.fail('Should have thrown an error for non-existent parent');
    } catch (error) {
      expect(error).toBeDefined();
    }
  });

  it('should maintain data consistency across operations', async () => {
    // Get the existing root task
    const rootTask = await apiClient.getRootTask();
    
    const child1 = await apiClient.createChildTask('Consistency Child 1', rootTask.id);
    const child2 = await apiClient.createChildTask('Consistency Child 2', rootTask.id);
    createdTaskIds.push(child1.id, child2.id);
    
    // Get all tasks and verify the hierarchy
    const allTasks = await apiClient.getAllTasks();
    
    const foundRoot = allTasks.find(task => task.id === rootTask.id);
    const foundChild1 = allTasks.find(task => task.id === child1.id);
    const foundChild2 = allTasks.find(task => task.id === child2.id);
    
    expect(foundRoot).toBeDefined();
    expect(foundChild1).toBeDefined();
    expect(foundChild2).toBeDefined();
    
    expect(foundChild1?.parentId).toBe(rootTask.id);
    expect(foundChild2?.parentId).toBe(rootTask.id);
    
    // Verify children are returned by getTaskChildren
    const children = await apiClient.getTaskChildren(rootTask.id);
    expect(children.some(child => child.id === child1.id)).toBe(true);
    expect(children.some(child => child.id === child2.id)).toBe(true);
  });
});