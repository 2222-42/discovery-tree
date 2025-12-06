package domain

import (
	"fmt"
	"strings"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

// Generator for valid task descriptions (non-empty, non-whitespace-only strings)
func GenValidDescription() gopter.Gen {
	return gen.OneGenOf(
		gen.Const("Valid task description"),
		gen.Const("Root task"),
		gen.Const("Child task"),
		gen.Const("Important work item"),
		gen.Const("Feature implementation"),
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) > 0 && len(s) < 100
		}),
		gen.Identifier().SuchThat(func(s string) bool {
			return len(s) > 0
		}),
	).SuchThat(func(v interface{}) bool {
		s, ok := v.(string)
		return ok && strings.TrimSpace(s) != ""
	})
}

// Generator for invalid task descriptions (empty or whitespace-only strings)
func GenInvalidDescription() gopter.Gen {
	return gen.OneConstOf(
		"",
		"   ",
		"\t\t",
		"\n\n",
		" \t\n ",
		"     ",
	)
}

// Generator for valid Status values
func GenValidStatus() gopter.Gen {
	return gen.OneConstOf(
		StatusTODO,
		StatusInProgress,
		StatusDONE,
		StatusBlocked,
		StatusRootWorkItem,
	)
}

// Generator for invalid Status values
func GenInvalidStatus() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(Status(-1)),
		gen.Const(Status(-100)),
		gen.Const(Status(100)),
		gen.Const(Status(999)),
		gen.IntRange(-1000, -1).Map(func(i int) Status { return Status(i) }),
		gen.IntRange(6, 1000).Map(func(i int) Status { return Status(i) }),
	)
}

// Generator for valid positions (non-negative integers)
func GenValidPosition() gopter.Gen {
	return gen.IntRange(0, 100)
}

// Generator for invalid positions (negative integers)
func GenInvalidPosition() gopter.Gen {
	return gen.IntRange(-100, -1)
}

// Generator for TaskID
func GenTaskID() gopter.Gen {
	return gen.Const(0).Map(func(_ int) TaskID {
		return NewTaskID()
	})
}

// Generator for optional TaskID (can be nil)
func GenOptionalTaskID() gopter.Gen {
	return gen.OneGenOf(
		gen.Const((*TaskID)(nil)),
		GenTaskID().Map(func(id TaskID) *TaskID { 
			idCopy := id
			return &idCopy 
		}),
	)
}

// TreeNode represents a node in a generated tree structure
type TreeNode struct {
	Task     *Task
	Children []*TreeNode
}

// GenSimpleTree generates a simple tree with a root and a few children
// This is a helper that creates a tree structure in the repository
func GenSimpleTree(repo TaskRepository) gopter.Gen {
	return gen.IntRange(0, 5).Map(func(numChildren int) *TreeNode {
		// Create root task
		root, err := NewTask("Root task", nil, 0)
		if err != nil {
			return nil
		}
		if err := repo.Save(root); err != nil {
			return nil
		}
		
		rootNode := &TreeNode{
			Task:     root,
			Children: make([]*TreeNode, 0, numChildren),
		}
		
		// Create children
		parentID := root.ID()
		for i := 0; i < numChildren; i++ {
			childDesc := fmt.Sprintf("Child %d", i)
			child, err := NewTask(childDesc, &parentID, i)
			if err != nil {
				continue
			}
			if err := repo.Save(child); err != nil {
				continue
			}
			
			childNode := &TreeNode{
				Task:     child,
				Children: []*TreeNode{},
			}
			rootNode.Children = append(rootNode.Children, childNode)
		}
		
		return rootNode
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

// GenTreeWithDepth generates a tree with specified depth
func GenTreeWithDepth(repo TaskRepository, maxDepth int) gopter.Gen {
	return gen.IntRange(0, 3).Map(func(childrenPerNode int) *TreeNode {
		return buildTreeWithDepth(repo, nil, 0, 0, maxDepth, childrenPerNode)
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

func buildTreeWithDepth(repo TaskRepository, parentID *TaskID, position int, currentDepth int, maxDepth int, childrenPerNode int) *TreeNode {
	// Create the task
	desc := fmt.Sprintf("Task at depth %d, position %d", currentDepth, position)
	task, err := NewTask(desc, parentID, position)
	if err != nil {
		return nil
	}
	if err := repo.Save(task); err != nil {
		return nil
	}
	
	node := &TreeNode{
		Task:     task,
		Children: []*TreeNode{},
	}
	
	// If we haven't reached max depth, add children
	if currentDepth < maxDepth {
		taskID := task.ID()
		for i := 0; i < childrenPerNode; i++ {
			childNode := buildTreeWithDepth(repo, &taskID, i, currentDepth+1, maxDepth, childrenPerNode)
			if childNode != nil {
				node.Children = append(node.Children, childNode)
			}
		}
	}
	
	return node
}

// GenTaskWithChildren generates a task with a specified number of children
func GenTaskWithChildren(repo TaskRepository, numChildren int) gopter.Gen {
	return GenValidDescription().Map(func(desc string) *TreeNode {
		// Create parent task
		parent, err := NewTask(desc, nil, 0)
		if err != nil {
			return nil
		}
		if err := repo.Save(parent); err != nil {
			return nil
		}
		
		parentNode := &TreeNode{
			Task:     parent,
			Children: make([]*TreeNode, 0, numChildren),
		}
		
		// Create children
		parentID := parent.ID()
		for i := 0; i < numChildren; i++ {
			childDesc := fmt.Sprintf("Child %d of %s", i, desc)
			child, err := NewTask(childDesc, &parentID, i)
			if err != nil {
				continue
			}
			if err := repo.Save(child); err != nil {
				continue
			}
			
			childNode := &TreeNode{
				Task:     child,
				Children: []*TreeNode{},
			}
			parentNode.Children = append(parentNode.Children, childNode)
		}
		
		return parentNode
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

// GenSiblings generates a list of sibling tasks under the same parent
func GenSiblings(repo TaskRepository, count int) gopter.Gen {
	return GenValidDescription().Map(func(parentDesc string) []*Task {
		// Create parent task
		parent, err := NewTask(parentDesc, nil, 0)
		if err != nil {
			return nil
		}
		if err := repo.Save(parent); err != nil {
			return nil
		}
		
		parentID := parent.ID()
		siblings := make([]*Task, 0, count)
		
		// Create sibling tasks
		for i := 0; i < count; i++ {
			siblingDesc := fmt.Sprintf("Sibling %d", i)
			sibling, err := NewTask(siblingDesc, &parentID, i)
			if err != nil {
				continue
			}
			if err := repo.Save(sibling); err != nil {
				continue
			}
			siblings = append(siblings, sibling)
		}
		
		return siblings
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

// GenTaskWithStatus generates a task with a specific status
func GenTaskWithStatus(repo TaskRepository, status Status) gopter.Gen {
	return GenValidDescription().Map(func(desc string) *Task {
		// Create task
		task, err := NewTask(desc, nil, 0)
		if err != nil {
			return nil
		}
		
		// Set the desired status
		if err := task.ChangeStatus(status); err != nil {
			return nil
		}
		
		if err := repo.Save(task); err != nil {
			return nil
		}
		
		return task
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

// FlattenTree converts a TreeNode into a flat list of all tasks in the tree
func FlattenTree(node *TreeNode) []*Task {
	if node == nil {
		return []*Task{}
	}
	
	tasks := []*Task{node.Task}
	for _, child := range node.Children {
		tasks = append(tasks, FlattenTree(child)...)
	}
	return tasks
}

// GetAllTaskIDs extracts all task IDs from a tree
func GetAllTaskIDs(node *TreeNode) []TaskID {
	tasks := FlattenTree(node)
	ids := make([]TaskID, len(tasks))
	for i, task := range tasks {
		ids[i] = task.ID()
	}
	return ids
}
