package domain

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// TestGenerators_ValidDescription verifies that GenValidDescription generates valid descriptions
func TestGenerators_ValidDescription(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenValidDescription generates non-empty descriptions", prop.ForAll(
		func(desc string) bool {
			// Valid descriptions should not be empty after trimming
			return len(desc) > 0 && len(desc) < 1000
		},
		GenValidDescription(),
	))

	properties.TestingRun(t)
}

// TestGenerators_InvalidDescription verifies that GenInvalidDescription generates invalid descriptions
func TestGenerators_InvalidDescription(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("GenInvalidDescription generates empty or whitespace-only descriptions", prop.ForAll(
		func(desc string) bool {
			// Invalid descriptions should be empty or whitespace-only
			_, err := NewTask(desc, nil, 0)
			return err != nil
		},
		GenInvalidDescription(),
	))

	properties.TestingRun(t)
}

// TestGenerators_ValidStatus verifies that GenValidStatus generates valid status values
func TestGenerators_ValidStatus(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenValidStatus generates valid status values", prop.ForAll(
		func(status Status) bool {
			return status.IsValid()
		},
		GenValidStatus(),
	))

	properties.TestingRun(t)
}

// TestGenerators_InvalidStatus verifies that GenInvalidStatus generates invalid status values
func TestGenerators_InvalidStatus(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenInvalidStatus generates invalid status values", prop.ForAll(
		func(status Status) bool {
			return !status.IsValid()
		},
		GenInvalidStatus(),
	))

	properties.TestingRun(t)
}

// TestGenerators_ValidPosition verifies that GenValidPosition generates valid positions
func TestGenerators_ValidPosition(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenValidPosition generates non-negative positions", prop.ForAll(
		func(pos int) bool {
			return pos >= 0
		},
		GenValidPosition(),
	))

	properties.TestingRun(t)
}

// TestGenerators_InvalidPosition verifies that GenInvalidPosition generates invalid positions
func TestGenerators_InvalidPosition(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenInvalidPosition generates negative positions", prop.ForAll(
		func(pos int) bool {
			return pos < 0
		},
		GenInvalidPosition(),
	))

	properties.TestingRun(t)
}

// TestGenerators_TaskID verifies that GenTaskID generates unique task IDs
func TestGenerators_TaskID(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("GenTaskID generates valid task IDs", prop.ForAll(
		func(id TaskID) bool {
			return id.String() != ""
		},
		GenTaskID(),
	))

	properties.TestingRun(t)
}

// TestGenerators_SimpleTree verifies that GenSimpleTree generates valid tree structures
func TestGenerators_SimpleTree(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("GenSimpleTree generates valid tree structures", prop.ForAll(
		func(tree *TreeNode) bool {
			if tree == nil || tree.Task == nil {
				return false
			}
			
			// Root should have no parent
			if tree.Task.ParentID() != nil {
				return false
			}
			
			// All children should have the root as parent
			rootID := tree.Task.ID()
			for i, child := range tree.Children {
				if child == nil || child.Task == nil {
					return false
				}
				if child.Task.ParentID() == nil || !child.Task.ParentID().Equals(rootID) {
					return false
				}
				// Check position matches index
				if child.Task.Position() != i {
					return false
				}
			}
			
			return true
		},
		GenSimpleTree(NewInMemoryTaskRepository()),
	))

	properties.TestingRun(t)
}

// TestGenerators_TaskWithChildren verifies that GenTaskWithChildren generates correct structures
func TestGenerators_TaskWithChildren(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	numChildren := 3
	properties.Property("GenTaskWithChildren generates correct number of children", prop.ForAll(
		func(tree *TreeNode) bool {
			if tree == nil || tree.Task == nil {
				return false
			}
			
			// Should have exactly numChildren children
			if len(tree.Children) != numChildren {
				return false
			}
			
			// All children should have correct parent and positions
			parentID := tree.Task.ID()
			for i, child := range tree.Children {
				if child == nil || child.Task == nil {
					return false
				}
				if child.Task.ParentID() == nil || !child.Task.ParentID().Equals(parentID) {
					return false
				}
				if child.Task.Position() != i {
					return false
				}
			}
			
			return true
		},
		GenTaskWithChildren(NewInMemoryTaskRepository(), numChildren),
	))

	properties.TestingRun(t)
}

// TestGenerators_Siblings verifies that GenSiblings generates correct sibling structures
func TestGenerators_Siblings(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	count := 4
	properties.Property("GenSiblings generates correct number of siblings", prop.ForAll(
		func(siblings []*Task) bool {
			if len(siblings) != count {
				return false
			}
			
			// All siblings should have the same parent
			if len(siblings) == 0 {
				return false
			}
			
			firstParentID := siblings[0].ParentID()
			if firstParentID == nil {
				return false
			}
			
			for i, sibling := range siblings {
				if sibling.ParentID() == nil || !sibling.ParentID().Equals(*firstParentID) {
					return false
				}
				// Check position matches index
				if sibling.Position() != i {
					return false
				}
			}
			
			return true
		},
		GenSiblings(NewInMemoryTaskRepository(), count),
	))

	properties.TestingRun(t)
}

// TestGenerators_TaskWithStatus verifies that GenTaskWithStatus generates tasks with correct status
func TestGenerators_TaskWithStatus(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	targetStatus := StatusDONE
	properties.Property("GenTaskWithStatus generates tasks with correct status", prop.ForAll(
		func(task *Task) bool {
			if task == nil {
				return false
			}
			return task.Status() == targetStatus
		},
		GenTaskWithStatus(NewInMemoryTaskRepository(), targetStatus),
	))

	properties.TestingRun(t)
}

// TestGenerators_TreeWithDepth verifies that GenTreeWithDepth generates trees with correct depth
func TestGenerators_TreeWithDepth(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	maxDepth := 2
	properties.Property("GenTreeWithDepth generates trees with correct depth", prop.ForAll(
		func(tree *TreeNode) bool {
			if tree == nil || tree.Task == nil {
				return false
			}
			
			// Calculate actual depth
			actualDepth := calculateDepth(tree)
			
			// Depth should not exceed maxDepth
			return actualDepth <= maxDepth+1 // +1 because root is at depth 0
		},
		GenTreeWithDepth(NewInMemoryTaskRepository(), maxDepth),
	))

	properties.TestingRun(t)
}

// Helper function to calculate tree depth
func calculateDepth(node *TreeNode) int {
	if node == nil || len(node.Children) == 0 {
		return 1
	}
	
	maxChildDepth := 0
	for _, child := range node.Children {
		childDepth := calculateDepth(child)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}
	
	return 1 + maxChildDepth
}
