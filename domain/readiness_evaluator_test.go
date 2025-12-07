package domain

import (
	"testing"
)

func TestReadinessEvaluatorService_EvaluateReadiness(t *testing.T) {
	tests := []struct {
		name                string
		setupTree           func(repo TaskRepository) TaskID // Returns the task ID to evaluate
		expectedReady       bool
		expectedLeftSibling bool
		expectedChildren    bool
		expectedReasons     int // Number of reasons
	}{
		{
			name: "root task with no children is ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)
				return root.ID()
			},
			expectedReady:       true,
			expectedLeftSibling: true,
			expectedChildren:    true,
			expectedReasons:     0,
		},
		{
			name: "root task with all children DONE is ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusDONE)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = child2.ChangeStatus(StatusDONE)
				_ = repo.Save(child2)

				return root.ID()
			},
			expectedReady:       true,
			expectedLeftSibling: true,
			expectedChildren:    true,
			expectedReasons:     0,
		},
		{
			name: "root task with incomplete children is not ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusDONE)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = child2.ChangeStatus(StatusTODO)
				_ = repo.Save(child2)

				return root.ID()
			},
			expectedReady:       false,
			expectedLeftSibling: true,
			expectedChildren:    false,
			expectedReasons:     1,
		},
		{
			name: "leftmost child with no children is ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = repo.Save(child1)

				return child1.ID()
			},
			expectedReady:       true,
			expectedLeftSibling: true,
			expectedChildren:    true,
			expectedReasons:     0,
		},
		{
			name: "middle child with left sibling DONE is ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusDONE)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = repo.Save(child2)

				return child2.ID()
			},
			expectedReady:       true,
			expectedLeftSibling: true,
			expectedChildren:    true,
			expectedReasons:     0,
		},
		{
			name: "middle child with left sibling incomplete is not ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusTODO)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = repo.Save(child2)

				return child2.ID()
			},
			expectedReady:       false,
			expectedLeftSibling: false,
			expectedChildren:    true,
			expectedReasons:     1,
		},
		{
			name: "task with incomplete left sibling and incomplete children is not ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusTODO)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = repo.Save(child2)

				grandchild, _ := NewTask("Grandchild", &child2.id, 0)
				_ = grandchild.ChangeStatus(StatusTODO)
				_ = repo.Save(grandchild)

				return child2.ID()
			},
			expectedReady:       false,
			expectedLeftSibling: false,
			expectedChildren:    false,
			expectedReasons:     2,
		},
		{
			name: "task with left sibling DONE but incomplete children is not ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusDONE)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = repo.Save(child2)

				grandchild, _ := NewTask("Grandchild", &child2.id, 0)
				_ = grandchild.ChangeStatus(StatusInProgress)
				_ = repo.Save(grandchild)

				return child2.ID()
			},
			expectedReady:       false,
			expectedLeftSibling: true,
			expectedChildren:    false,
			expectedReasons:     1,
		},
		{
			name: "rightmost child with left sibling DONE is ready",
			setupTree: func(repo TaskRepository) TaskID {
				root, _ := NewTask("Root", nil, 0)
				_ = repo.Save(root)

				child1, _ := NewTask("Child 1", &root.id, 0)
				_ = child1.ChangeStatus(StatusDONE)
				_ = repo.Save(child1)

				child2, _ := NewTask("Child 2", &root.id, 1)
				_ = child2.ChangeStatus(StatusDONE)
				_ = repo.Save(child2)

				child3, _ := NewTask("Child 3", &root.id, 2)
				_ = repo.Save(child3)

				return child3.ID()
			},
			expectedReady:       true,
			expectedLeftSibling: true,
			expectedChildren:    true,
			expectedReasons:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			repo := NewInMemoryTaskRepository()
			navigator := NewTreeNavigatorService(repo)
			evaluator := NewReadinessEvaluatorService(repo, navigator)

			// Create the tree and get the task ID to evaluate
			taskID := tt.setupTree(repo)

			// Execute
			state, err := evaluator.EvaluateReadiness(taskID)

			// Verify
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if state.IsReady() != tt.expectedReady {
				t.Errorf("expected IsReady=%v, got %v", tt.expectedReady, state.IsReady())
			}

			if state.LeftSiblingComplete() != tt.expectedLeftSibling {
				t.Errorf("expected LeftSiblingComplete=%v, got %v", tt.expectedLeftSibling, state.LeftSiblingComplete())
			}

			if state.AllChildrenComplete() != tt.expectedChildren {
				t.Errorf("expected AllChildrenComplete=%v, got %v", tt.expectedChildren, state.AllChildrenComplete())
			}

			if len(state.Reasons()) != tt.expectedReasons {
				t.Errorf("expected %d reasons, got %d: %v", tt.expectedReasons, len(state.Reasons()), state.Reasons())
			}
		})
	}
}

func TestReadinessEvaluatorService_EvaluateReadiness_NonExistentTask(t *testing.T) {
	// Setup
	repo := NewInMemoryTaskRepository()
	navigator := NewTreeNavigatorService(repo)
	evaluator := NewReadinessEvaluatorService(repo, navigator)

	// Create a task ID that doesn't exist
	nonExistentID := NewTaskID()

	// Execute
	_, err := evaluator.EvaluateReadiness(nonExistentID)

	// Verify
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}

	if _, ok := err.(NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}

func TestReadinessEvaluatorService_EvaluateReadiness_ComplexTree(t *testing.T) {
	// Setup a more complex tree structure
	repo := NewInMemoryTaskRepository()
	navigator := NewTreeNavigatorService(repo)
	evaluator := NewReadinessEvaluatorService(repo, navigator)

	// Create tree:
	//        Root
	//       /  |  \
	//      A   B   C
	//     /|   |
	//    D E   F
	//
	// Status: A=DONE, D=DONE, E=DONE, B=TODO, F=TODO, C=TODO

	root, _ := NewTask("Root", nil, 0)
	_ = repo.Save(root)

	taskA, _ := NewTask("A", &root.id, 0)
	_ = taskA.ChangeStatus(StatusDONE)
	_ = repo.Save(taskA)

	taskB, _ := NewTask("B", &root.id, 1)
	_ = taskB.ChangeStatus(StatusTODO)
	_ = repo.Save(taskB)

	taskC, _ := NewTask("C", &root.id, 2)
	_ = taskC.ChangeStatus(StatusTODO)
	_ = repo.Save(taskC)

	taskD, _ := NewTask("D", &taskA.id, 0)
	_ = taskD.ChangeStatus(StatusDONE)
	_ = repo.Save(taskD)

	taskE, _ := NewTask("E", &taskA.id, 1)
	_ = taskE.ChangeStatus(StatusDONE)
	_ = repo.Save(taskE)

	taskF, _ := NewTask("F", &taskB.id, 0)
	_ = taskF.ChangeStatus(StatusTODO)
	_ = repo.Save(taskF)

	// Test evaluations
	tests := []struct {
		name          string
		taskID        TaskID
		expectedReady bool
	}{
		{"Root not ready (children incomplete)", root.ID(), false},
		{"A ready (all children DONE, no left sibling)", taskA.ID(), true},
		{"B not ready (children incomplete, left sibling DONE)", taskB.ID(), false},
		{"C not ready (left sibling incomplete)", taskC.ID(), false},
		{"D ready (no children, no left sibling)", taskD.ID(), true},
		{"E ready (no children, left sibling DONE)", taskE.ID(), true},
		{"F ready (no children, no left sibling)", taskF.ID(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, err := evaluator.EvaluateReadiness(tt.taskID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if state.IsReady() != tt.expectedReady {
				t.Errorf("expected IsReady=%v, got %v (reasons: %v)", 
					tt.expectedReady, state.IsReady(), state.Reasons())
			}
		})
	}
}
