package domain

import (
	"testing"
)

// setupTreeNavigatorTest creates a test tree structure:
//
//	Root (pos 0)
//	├── Child1 (pos 0)
//	│   ├── Grandchild1 (pos 0)
//	│   └── Grandchild2 (pos 1)
//	├── Child2 (pos 1)
//	└── Child3 (pos 2)
func setupTreeNavigatorTest(t *testing.T) (*InMemoryTaskRepository, *TreeNavigatorService, map[string]*Task) {
	repo := NewInMemoryTaskRepository()
	navigator := NewTreeNavigatorService(repo)

	// Create root task
	root, err := NewTask("Root Task", nil, 0)
	if err != nil {
		t.Fatalf("Failed to create root task: %v", err)
	}
	if err := repo.Save(root); err != nil {
		t.Fatalf("Failed to save root task: %v", err)
	}

	// Create child tasks
	rootID := root.ID()
	child1, err := NewTask("Child 1", &rootID, 0)
	if err != nil {
		t.Fatalf("Failed to create child1: %v", err)
	}
	if err := repo.Save(child1); err != nil {
		t.Fatalf("Failed to save child1: %v", err)
	}

	child2, err := NewTask("Child 2", &rootID, 1)
	if err != nil {
		t.Fatalf("Failed to create child2: %v", err)
	}
	if err := repo.Save(child2); err != nil {
		t.Fatalf("Failed to save child2: %v", err)
	}

	child3, err := NewTask("Child 3", &rootID, 2)
	if err != nil {
		t.Fatalf("Failed to create child3: %v", err)
	}
	if err := repo.Save(child3); err != nil {
		t.Fatalf("Failed to save child3: %v", err)
	}

	// Create grandchild tasks under child1
	child1ID := child1.ID()
	grandchild1, err := NewTask("Grandchild 1", &child1ID, 0)
	if err != nil {
		t.Fatalf("Failed to create grandchild1: %v", err)
	}
	if err := repo.Save(grandchild1); err != nil {
		t.Fatalf("Failed to save grandchild1: %v", err)
	}

	grandchild2, err := NewTask("Grandchild 2", &child1ID, 1)
	if err != nil {
		t.Fatalf("Failed to create grandchild2: %v", err)
	}
	if err := repo.Save(grandchild2); err != nil {
		t.Fatalf("Failed to save grandchild2: %v", err)
	}

	tasks := map[string]*Task{
		"root":        root,
		"child1":      child1,
		"child2":      child2,
		"child3":      child3,
		"grandchild1": grandchild1,
		"grandchild2": grandchild2,
	}

	return repo, navigator, tasks
}

func TestTreeNavigator_GetParent(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name           string
		taskID         TaskID
		expectedParent *Task
		expectNil      bool
		expectError    bool
	}{
		{
			name:        "root task has no parent",
			taskID:      tasks["root"].ID(),
			expectNil:   true,
			expectError: false,
		},
		{
			name:           "child1 parent is root",
			taskID:         tasks["child1"].ID(),
			expectedParent: tasks["root"],
			expectError:    false,
		},
		{
			name:           "grandchild1 parent is child1",
			taskID:         tasks["grandchild1"].ID(),
			expectedParent: tasks["child1"],
			expectError:    false,
		},
		{
			name:        "non-existent task returns error",
			taskID:      NewTaskID(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, err := navigator.GetParent(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectNil {
				if parent != nil {
					t.Errorf("Expected nil parent but got %v", parent)
				}
			} else {
				if parent == nil {
					t.Fatal("Expected parent but got nil")
				}
				if !parent.ID().Equals(tt.expectedParent.ID()) {
					t.Errorf("Expected parent %v but got %v", tt.expectedParent.ID(), parent.ID())
				}
			}
		})
	}
}

func TestTreeNavigator_GetChildren(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name          string
		taskID        TaskID
		expectedCount int
		expectedOrder []string
		expectError   bool
	}{
		{
			name:          "root has 3 children",
			taskID:        tasks["root"].ID(),
			expectedCount: 3,
			expectedOrder: []string{"Child 1", "Child 2", "Child 3"},
			expectError:   false,
		},
		{
			name:          "child1 has 2 children",
			taskID:        tasks["child1"].ID(),
			expectedCount: 2,
			expectedOrder: []string{"Grandchild 1", "Grandchild 2"},
			expectError:   false,
		},
		{
			name:          "child2 has no children",
			taskID:        tasks["child2"].ID(),
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:        "non-existent task returns error",
			taskID:      NewTaskID(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			children, err := navigator.GetChildren(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(children) != tt.expectedCount {
				t.Errorf("Expected %d children but got %d", tt.expectedCount, len(children))
			}

			// Verify ordering
			for i, expectedDesc := range tt.expectedOrder {
				if i >= len(children) {
					break
				}
				if children[i].Description() != expectedDesc {
					t.Errorf("Expected child at position %d to be '%s' but got '%s'",
						i, expectedDesc, children[i].Description())
				}
			}
		})
	}
}

func TestTreeNavigator_GetSiblings(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name          string
		taskID        TaskID
		expectedCount int
		expectError   bool
	}{
		{
			name:          "child1 has 3 siblings (including itself)",
			taskID:        tasks["child1"].ID(),
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "child2 has 3 siblings (including itself)",
			taskID:        tasks["child2"].ID(),
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "grandchild1 has 2 siblings (including itself)",
			taskID:        tasks["grandchild1"].ID(),
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "root has 1 sibling (itself, no other roots)",
			taskID:        tasks["root"].ID(),
			expectedCount: 1,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siblings, err := navigator.GetSiblings(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(siblings) != tt.expectedCount {
				t.Errorf("Expected %d siblings but got %d", tt.expectedCount, len(siblings))
			}

			// Verify ordering by position
			for i := 0; i < len(siblings)-1; i++ {
				if siblings[i].Position() >= siblings[i+1].Position() {
					t.Error("Siblings are not properly ordered by position")
				}
			}
		})
	}
}

func TestTreeNavigator_GetLeftSibling(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name           string
		taskID         TaskID
		expectedSibling *Task
		expectNil      bool
		expectError    bool
	}{
		{
			name:        "child1 (pos 0) has no left sibling",
			taskID:      tasks["child1"].ID(),
			expectNil:   true,
			expectError: false,
		},
		{
			name:           "child2 (pos 1) left sibling is child1",
			taskID:         tasks["child2"].ID(),
			expectedSibling: tasks["child1"],
			expectError:    false,
		},
		{
			name:           "child3 (pos 2) left sibling is child2",
			taskID:         tasks["child3"].ID(),
			expectedSibling: tasks["child2"],
			expectError:    false,
		},
		{
			name:        "grandchild1 (pos 0) has no left sibling",
			taskID:      tasks["grandchild1"].ID(),
			expectNil:   true,
			expectError: false,
		},
		{
			name:           "grandchild2 (pos 1) left sibling is grandchild1",
			taskID:         tasks["grandchild2"].ID(),
			expectedSibling: tasks["grandchild1"],
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftSibling, err := navigator.GetLeftSibling(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectNil {
				if leftSibling != nil {
					t.Errorf("Expected nil left sibling but got %v", leftSibling)
				}
			} else {
				if leftSibling == nil {
					t.Fatal("Expected left sibling but got nil")
				}
				if !leftSibling.ID().Equals(tt.expectedSibling.ID()) {
					t.Errorf("Expected left sibling %v but got %v",
						tt.expectedSibling.ID(), leftSibling.ID())
				}
			}
		})
	}
}

func TestTreeNavigator_GetRightSibling(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name            string
		taskID          TaskID
		expectedSibling *Task
		expectNil       bool
		expectError     bool
	}{
		{
			name:            "child1 (pos 0) right sibling is child2",
			taskID:          tasks["child1"].ID(),
			expectedSibling: tasks["child2"],
			expectError:     false,
		},
		{
			name:            "child2 (pos 1) right sibling is child3",
			taskID:          tasks["child2"].ID(),
			expectedSibling: tasks["child3"],
			expectError:     false,
		},
		{
			name:        "child3 (pos 2) has no right sibling",
			taskID:      tasks["child3"].ID(),
			expectNil:   true,
			expectError: false,
		},
		{
			name:            "grandchild1 (pos 0) right sibling is grandchild2",
			taskID:          tasks["grandchild1"].ID(),
			expectedSibling: tasks["grandchild2"],
			expectError:     false,
		},
		{
			name:        "grandchild2 (pos 1) has no right sibling",
			taskID:      tasks["grandchild2"].ID(),
			expectNil:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rightSibling, err := navigator.GetRightSibling(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectNil {
				if rightSibling != nil {
					t.Errorf("Expected nil right sibling but got %v", rightSibling)
				}
			} else {
				if rightSibling == nil {
					t.Fatal("Expected right sibling but got nil")
				}
				if !rightSibling.ID().Equals(tt.expectedSibling.ID()) {
					t.Errorf("Expected right sibling %v but got %v",
						tt.expectedSibling.ID(), rightSibling.ID())
				}
			}
		})
	}
}

func TestTreeNavigator_GetRoot(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	root, err := navigator.GetRoot()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if root == nil {
		t.Fatal("Expected root task but got nil")
	}

	if !root.ID().Equals(tasks["root"].ID()) {
		t.Errorf("Expected root %v but got %v", tasks["root"].ID(), root.ID())
	}

	if root.Status() != StatusRootWorkItem {
		t.Errorf("Expected root status to be RootWorkItem but got %v", root.Status())
	}
}

func TestTreeNavigator_GetSubtree(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	tests := []struct {
		name          string
		taskID        TaskID
		expectedCount int
		expectError   bool
	}{
		{
			name:          "root subtree includes all 6 tasks",
			taskID:        tasks["root"].ID(),
			expectedCount: 6,
			expectError:   false,
		},
		{
			name:          "child1 subtree includes child1 and 2 grandchildren",
			taskID:        tasks["child1"].ID(),
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "child2 subtree includes only child2 (no children)",
			taskID:        tasks["child2"].ID(),
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "grandchild1 subtree includes only itself",
			taskID:        tasks["grandchild1"].ID(),
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:        "non-existent task returns error",
			taskID:      NewTaskID(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subtree, err := navigator.GetSubtree(tt.taskID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(subtree) != tt.expectedCount {
				t.Errorf("Expected subtree with %d tasks but got %d", tt.expectedCount, len(subtree))
			}

			// Verify the first task is the requested task
			if !subtree[0].ID().Equals(tt.taskID) {
				t.Error("First task in subtree should be the requested task")
			}
		})
	}
}

func TestTreeNavigator_GetSubtree_PreservesStructure(t *testing.T) {
	_, navigator, tasks := setupTreeNavigatorTest(t)

	// Get the subtree starting from child1
	subtree, err := navigator.GetSubtree(tasks["child1"].ID())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify that child1 is included
	foundChild1 := false
	foundGrandchild1 := false
	foundGrandchild2 := false

	for _, task := range subtree {
		if task.ID().Equals(tasks["child1"].ID()) {
			foundChild1 = true
		}
		if task.ID().Equals(tasks["grandchild1"].ID()) {
			foundGrandchild1 = true
		}
		if task.ID().Equals(tasks["grandchild2"].ID()) {
			foundGrandchild2 = true
		}
	}

	if !foundChild1 {
		t.Error("Subtree should include child1")
	}
	if !foundGrandchild1 {
		t.Error("Subtree should include grandchild1")
	}
	if !foundGrandchild2 {
		t.Error("Subtree should include grandchild2")
	}

	// Verify that child2 and child3 are NOT included
	for _, task := range subtree {
		if task.ID().Equals(tasks["child2"].ID()) || task.ID().Equals(tasks["child3"].ID()) {
			t.Error("Subtree should not include siblings of child1")
		}
	}
}
