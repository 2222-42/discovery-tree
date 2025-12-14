/**
 * End-to-end workflow tests for complete user scenarios
 * Tests full user workflows from task creation to deletion
 */

import { describe, it, expect, beforeAll, afterAll } from 'vitest';
import { DiscoveryTreeApiClient } from '@/services/api/apiClient.js';
import { TaskStatus } from '@/types/task.js';
import { TreeService } from '@/services/tree/treeService.js';

describe('Complete User Workflow Tests', () => {
  let apiClient: DiscoveryTreeApiClient;
  let treeService: TreeService;
  let createdTaskIds: string[] = [];

  beforeAll(() => {
    apiClient = new DiscoveryTreeApiClient({
      baseURL: 'http://localhost:8080/api/v1',
      timeout: 5000,
    });
    treeService = new TreeService();
  });

  afterAll(async () => {
    // Clean up created tasks
    for (const taskId of createdTaskIds) {
      try {
        await apiClient.deleteTask(taskId);
      } catch (error) {
        console.warn(`Failed to cleanup task ${taskId}:`, error);
      }
    }
  });

  it('should complete a full project management workflow', async () => {
    // Step 1: Get the root task (project)
    const rootTask = await apiClient.getRootTask();
    expect(rootTask.status).toBe(TaskStatus.ROOT_WORK_ITEM);

    // Step 2: Create main project phases
    const designPhase = await apiClient.createChildTask('Design Phase', rootTask.id);
    const developmentPhase = await apiClient.createChildTask('Development Phase', rootTask.id);
    const testingPhase = await apiClient.createChildTask('Testing Phase', rootTask.id);
    
    createdTaskIds.push(designPhase.id, developmentPhase.id, testingPhase.id);

    // Step 3: Create sub-tasks for design phase
    const wireframes = await apiClient.createChildTask('Create Wireframes', designPhase.id);
    const mockups = await apiClient.createChildTask('Design Mockups', designPhase.id);
    
    createdTaskIds.push(wireframes.id, mockups.id);

    // Step 4: Create sub-tasks for development phase
    const frontend = await apiClient.createChildTask('Frontend Development', developmentPhase.id);
    const backend = await apiClient.createChildTask('Backend Development', developmentPhase.id);
    
    createdTaskIds.push(frontend.id, backend.id);

    // Step 5: Start working on wireframes
    await apiClient.updateTaskStatus(wireframes.id, TaskStatus.IN_PROGRESS);
    
    // Step 6: Complete wireframes
    await apiClient.updateTaskStatus(wireframes.id, TaskStatus.DONE);

    // Step 7: Start mockups
    await apiClient.updateTaskStatus(mockups.id, TaskStatus.IN_PROGRESS);

    // Step 8: Get all tasks and build tree structure
    const allTasks = await apiClient.getAllTasks();
    const treeNodes = treeService.buildTreeFromTasks(allTasks);

    // Verify tree structure
    expect(treeNodes.length).toBe(1); // Should have one root
    const rootNode = treeNodes[0];
    expect(rootNode.task.id).toBe(rootTask.id);
    expect(rootNode.children.length).toBeGreaterThanOrEqual(3); // At least our 3 phases

    // Find our phases in the tree
    const designNode = rootNode.children.find(child => child.task.id === designPhase.id);
    const developmentNode = rootNode.children.find(child => child.task.id === developmentPhase.id);
    const testingNode = rootNode.children.find(child => child.task.id === testingPhase.id);

    expect(designNode).toBeDefined();
    expect(developmentNode).toBeDefined();
    expect(testingNode).toBeDefined();

    // Verify design phase has sub-tasks
    expect(designNode!.children.length).toBe(2);
    const wireframesNode = designNode!.children.find(child => child.task.id === wireframes.id);
    const mockupsNode = designNode!.children.find(child => child.task.id === mockups.id);

    expect(wireframesNode).toBeDefined();
    expect(mockupsNode).toBeDefined();

    // Verify task statuses
    expect(wireframesNode!.task.status).toBe(TaskStatus.DONE);
    expect(mockupsNode!.task.status).toBe(TaskStatus.IN_PROGRESS);

    // Step 9: Update task descriptions
    await apiClient.updateTask(mockups.id, 'High-Fidelity Design Mockups');
    const updatedMockups = await apiClient.getTask(mockups.id);
    expect(updatedMockups.description).toBe('High-Fidelity Design Mockups');

    // Step 10: Test tree navigation
    const designChildren = await apiClient.getTaskChildren(designPhase.id);
    expect(designChildren.length).toBe(2);
    expect(designChildren.some(child => child.id === wireframes.id)).toBe(true);
    expect(designChildren.some(child => child.id === mockups.id)).toBe(true);

    // Step 11: Complete mockups and start development
    await apiClient.updateTaskStatus(mockups.id, TaskStatus.DONE);
    await apiClient.updateTaskStatus(frontend.id, TaskStatus.IN_PROGRESS);
    await apiClient.updateTaskStatus(backend.id, TaskStatus.IN_PROGRESS);

    // Step 12: Verify final state
    const finalTasks = await apiClient.getAllTasks();
    const finalWireframes = finalTasks.find(task => task.id === wireframes.id);
    const finalMockups = finalTasks.find(task => task.id === mockups.id);
    const finalFrontend = finalTasks.find(task => task.id === frontend.id);
    const finalBackend = finalTasks.find(task => task.id === backend.id);

    expect(finalWireframes?.status).toBe(TaskStatus.DONE);
    expect(finalMockups?.status).toBe(TaskStatus.DONE);
    expect(finalFrontend?.status).toBe(TaskStatus.IN_PROGRESS);
    expect(finalBackend?.status).toBe(TaskStatus.IN_PROGRESS);
  });

  it('should handle task reorganization workflow', async () => {
    // Step 1: Create initial structure
    const rootTask = await apiClient.getRootTask();
    
    const phase1 = await apiClient.createChildTask('Phase 1', rootTask.id);
    const phase2 = await apiClient.createChildTask('Phase 2', rootTask.id);
    
    const task1 = await apiClient.createChildTask('Task 1', phase1.id);
    const task2 = await apiClient.createChildTask('Task 2', phase1.id);
    
    createdTaskIds.push(phase1.id, phase2.id, task1.id, task2.id);

    // Step 2: Verify initial structure
    const phase1Children = await apiClient.getTaskChildren(phase1.id);
    const phase2Children = await apiClient.getTaskChildren(phase2.id);
    
    expect(phase1Children.length).toBe(2);
    expect(phase2Children.length).toBe(0);

    // Step 3: Move task2 to phase2 (simulating drag and drop)
    await apiClient.moveTask(task2.id, phase2.id, 0);

    // Step 4: Verify reorganization
    const newPhase1Children = await apiClient.getTaskChildren(phase1.id);
    const newPhase2Children = await apiClient.getTaskChildren(phase2.id);
    
    expect(newPhase1Children.length).toBe(1);
    expect(newPhase2Children.length).toBe(1);
    expect(newPhase2Children[0].id).toBe(task2.id);

    // Step 5: Verify task2 has correct parent
    const movedTask = await apiClient.getTask(task2.id);
    expect(movedTask.parentId).toBe(phase2.id);
  });

  it('should handle error recovery workflow', async () => {
    // Step 1: Try to create task with invalid parent
    try {
      await apiClient.createChildTask('Invalid Task', 'non-existent-parent');
      expect.fail('Should have thrown error for invalid parent');
    } catch (error) {
      expect(error).toBeDefined();
    }

    // Step 2: Try to get non-existent task
    try {
      await apiClient.getTask('non-existent-task');
      expect.fail('Should have thrown error for non-existent task');
    } catch (error) {
      expect(error).toBeDefined();
    }

    // Step 3: Try to update non-existent task
    try {
      await apiClient.updateTask('non-existent-task', 'New Description');
      expect.fail('Should have thrown error for non-existent task');
    } catch (error) {
      expect(error).toBeDefined();
    }

    // Step 4: Verify system is still functional after errors
    const rootTask = await apiClient.getRootTask();
    expect(rootTask).toBeDefined();
    
    const allTasks = await apiClient.getAllTasks();
    expect(Array.isArray(allTasks)).toBe(true);
  });

  it('should maintain data consistency during concurrent operations', async () => {
    // Step 1: Create test tasks
    const rootTask = await apiClient.getRootTask();
    const parentTask = await apiClient.createChildTask('Concurrent Test Parent', rootTask.id);
    createdTaskIds.push(parentTask.id);

    // Step 2: Create multiple child tasks concurrently
    const childPromises = Array.from({ length: 5 }, (_, i) => 
      apiClient.createChildTask(`Concurrent Child ${i + 1}`, parentTask.id)
    );
    
    const children = await Promise.all(childPromises);
    createdTaskIds.push(...children.map(child => child.id));

    // Step 3: Verify all children were created
    expect(children.length).toBe(5);
    children.forEach((child, index) => {
      expect(child.description).toBe(`Concurrent Child ${index + 1}`);
      expect(child.parentId).toBe(parentTask.id);
    });

    // Step 4: Update all children concurrently
    const updatePromises = children.map((child, index) => 
      apiClient.updateTaskStatus(child.id, index % 2 === 0 ? TaskStatus.IN_PROGRESS : TaskStatus.DONE)
    );
    
    await Promise.all(updatePromises);

    // Step 5: Verify all updates were applied
    const updatedChildren = await apiClient.getTaskChildren(parentTask.id);
    expect(updatedChildren.length).toBe(5);
    
    // Check that each original child has the correct status
    for (let i = 0; i < children.length; i++) {
      const originalChild = children[i];
      const updatedChild = updatedChildren.find(child => child.id === originalChild.id);
      expect(updatedChild).toBeDefined();
      
      const expectedStatus = i % 2 === 0 ? TaskStatus.IN_PROGRESS : TaskStatus.DONE;
      expect(updatedChild!.status).toBe(expectedStatus);
    }

    // Step 6: Build tree and verify structure
    const allTasks = await apiClient.getAllTasks();
    const treeNodes = treeService.buildTreeFromTasks(allTasks);
    
    // Find our parent in the tree
    const findTaskInTree = (nodes: any[], taskId: string): any => {
      for (const node of nodes) {
        if (node.task.id === taskId) return node;
        const found = findTaskInTree(node.children, taskId);
        if (found) return found;
      }
      return null;
    };

    const parentNode = findTaskInTree(treeNodes, parentTask.id);
    expect(parentNode).toBeDefined();
    expect(parentNode.children.length).toBe(5);
  });
});