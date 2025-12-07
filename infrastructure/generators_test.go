package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"

	"discovery-tree/domain"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

// GenValidDescription generates valid task descriptions
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
	)
}

// GenValidStatus generates valid Status values
func GenValidStatus() gopter.Gen {
	return gen.OneConstOf(
		domain.StatusTODO,
		domain.StatusInProgress,
		domain.StatusDONE,
		domain.StatusBlocked,
		domain.StatusRootWorkItem,
	)
}

// GenValidPosition generates valid positions (non-negative integers)
func GenValidPosition() gopter.Gen {
	return gen.IntRange(0, 100)
}

// GenTaskID generates a TaskID
func GenTaskID() gopter.Gen {
	return gen.Const(0).Map(func(_ int) domain.TaskID {
		return domain.NewTaskID()
	})
}

// GenOptionalTaskID generates an optional TaskID (can be nil)
func GenOptionalTaskID() gopter.Gen {
	return gen.OneGenOf(
		gen.Const((*domain.TaskID)(nil)),
		GenTaskID().Map(func(id domain.TaskID) *domain.TaskID {
			idCopy := id
			return &idCopy
		}),
	)
}

// GenTaskForRepository generates a random task suitable for repository operations
// This ensures the task has valid fields and can be saved/loaded
func GenTaskForRepository() gopter.Gen {
	return gopter.CombineGens(
		GenValidDescription(),
		GenOptionalTaskID(),
		GenValidPosition(),
		GenValidStatus(),
	).Map(func(values []interface{}) *domain.Task {
		desc := values[0].(string)
		parentID := values[1].(*domain.TaskID)
		position := values[2].(int)
		status := values[3].(domain.Status)
		
		task, err := domain.NewTask(desc, parentID, position)
		if err != nil {
			return nil
		}
		
		// Set the desired status
		if err := task.ChangeStatus(status); err != nil {
			return nil
		}
		
		return task
	}).SuchThat(func(v interface{}) bool {
		return v != nil
	})
}

// GenTaskTree generates a random tree of tasks
// Returns the root task and a slice of all tasks in the tree
func GenTaskTree() gopter.Gen {
	return gopter.CombineGens(
		gen.IntRange(0, 3), // maxDepth
		gen.IntRange(0, 3), // childrenPerNode
	).Map(func(values []interface{}) *TaskTree {
		maxDepth := values[0].(int)
		childrenPerNode := values[1].(int)
		return buildRandomTree(nil, 0, maxDepth, childrenPerNode)
	}).SuchThat(func(v interface{}) bool {
		return v != nil && v.(*TaskTree).Root != nil
	})
}

// TaskTree represents a generated tree structure
type TaskTree struct {
	Root  *domain.Task
	Tasks []*domain.Task
}

// buildRandomTree recursively builds a random task tree
func buildRandomTree(parentID *domain.TaskID, currentDepth int, maxDepth int, childrenPerNode int) *TaskTree {
	// Create the root task
	desc := fmt.Sprintf("Task at depth %d", currentDepth)
	task, err := domain.NewTask(desc, parentID, 0)
	if err != nil {
		return nil
	}
	
	tree := &TaskTree{
		Root:  task,
		Tasks: []*domain.Task{task},
	}
	
	// If we haven't reached max depth, add children
	if currentDepth < maxDepth {
		taskID := task.ID()
		for i := 0; i < childrenPerNode; i++ {
			childTree := buildRandomTree(&taskID, currentDepth+1, maxDepth, childrenPerNode)
			if childTree != nil {
				tree.Tasks = append(tree.Tasks, childTree.Tasks...)
			}
		}
	}
	
	return tree
}

// GenTaskTreeWithSize generates a task tree with approximately the specified number of tasks
func GenTaskTreeWithSize(minSize, maxSize int) gopter.Gen {
	return gen.Const(0).Map(func(_ int) *TaskTree {
		// Use middle value for size
		size := minSize + (maxSize-minSize)/2
		
		// Create root
		root, err := domain.NewTask("Root", nil, 0)
		if err != nil {
			return nil
		}
		
		tasks := []*domain.Task{root}
		
		// Add tasks until we reach target size
		for len(tasks) < size {
			// Pick a random parent from existing tasks
			parentIndex := len(tasks) / 2 // Simple strategy: use middle task
			if parentIndex >= len(tasks) {
				parentIndex = len(tasks) - 1
			}
			
			parent := tasks[parentIndex]
			parentID := parent.ID()
			
			// Create child
			desc := fmt.Sprintf("Task %d", len(tasks))
			child, err := domain.NewTask(desc, &parentID, len(tasks))
			if err != nil {
				break
			}
			
			tasks = append(tasks, child)
		}
		
		return &TaskTree{
			Root:  root,
			Tasks: tasks,
		}
	}).SuchThat(func(v interface{}) bool {
		return v != nil && v.(*TaskTree).Root != nil
	})
}

// GenInvalidJSON generates various forms of invalid JSON
func GenInvalidJSON() gopter.Gen {
	return gen.OneConstOf(
		// Malformed JSON
		"{invalid json}",
		"[{incomplete",
		"{\"key\": }",
		"[1, 2, 3,]",
		"{\"key\": \"value\",}",
		
		// Not an array
		"{\"id\": \"123\"}",
		"\"just a string\"",
		"123",
		"true",
		
		// Empty/whitespace
		"",
		"   ",
		"\n\n",
		
		// Invalid characters
		"{\"key\": 'single quotes'}",
		"[{\"key\": undefined}]",
		
		// Truncated JSON
		"[{\"id\": \"550e8400-e29b-41d4-a716-446655440000\", \"description\": \"Task\"",
	)
}

// GenInvalidTaskDTO generates TaskDTOs with invalid data
func GenInvalidTaskDTO() gopter.Gen {
	now := time.Now()
	
	// Create DTOs with timestamps set
	emptyDesc := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "",
		Status:      "TODO",
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	invalidID := TaskDTO{
		ID:          "not-a-uuid",
		Description: "Valid description",
		Status:      "TODO",
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	invalidStatus := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "INVALID_STATUS",
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	negativePosition := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		Position:    -1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	return gen.OneConstOf(emptyDesc, invalidID, invalidStatus, negativePosition)
}

// GenInvalidParentIDTaskDTO generates TaskDTOs with invalid parent IDs
func GenInvalidParentIDTaskDTO() gopter.Gen {
	now := time.Now()
	
	invalidParent1 := "not-a-uuid"
	dto1 := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    &invalidParent1,
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	invalidParent2 := ""
	dto2 := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    &invalidParent2,
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	invalidParent3 := "invalid"
	dto3 := TaskDTO{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Description: "Valid description",
		Status:      "TODO",
		ParentID:    &invalidParent3,
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	return gen.OneConstOf(dto1, dto2, dto3)
}

// GenInvalidTaskJSON generates JSON strings with invalid task data
func GenInvalidTaskJSON() gopter.Gen {
	return gen.OneGenOf(
		GenInvalidTaskDTO(),
		GenInvalidParentIDTaskDTO(),
	).Map(func(dto interface{}) string {
		invalidDTO := dto.(TaskDTO)
		data, _ := json.Marshal([]TaskDTO{invalidDTO})
		return string(data)
	})
}

// RepositoryOperation represents an operation that can be performed on a repository
type RepositoryOperation int

const (
	OpSave RepositoryOperation = iota
	OpFindByID
	OpFindAll
	OpFindRoot
	OpFindByParentID
	OpDelete
	OpDeleteSubtree
)

// ConcurrentOperation represents a single operation to be performed concurrently
type ConcurrentOperation struct {
	Op     RepositoryOperation
	TaskID *domain.TaskID // Used for FindByID, Delete, DeleteSubtree, FindByParentID
	Task   *domain.Task   // Used for Save
}

// GenConcurrentOperations generates a sequence of concurrent operations
func GenConcurrentOperations(numOps int) gopter.Gen {
	return gen.SliceOfN(numOps, GenSingleOperation()).SuchThat(func(v interface{}) bool {
		ops := v.([]ConcurrentOperation)
		return len(ops) > 0
	})
}

// Helper generators for each operation type
func genSaveOp() gopter.Gen {
	return gen.Const(0).Map(func(_ int) ConcurrentOperation {
		task, _ := domain.NewTask("Generated task", nil, 0)
		return ConcurrentOperation{
			Op:   OpSave,
			Task: task,
		}
	})
}

func genFindByIDOp() gopter.Gen {
	return gen.Const(0).Map(func(_ int) ConcurrentOperation {
		taskID := domain.NewTaskID()
		return ConcurrentOperation{
			Op:     OpFindByID,
			TaskID: &taskID,
		}
	})
}

func genFindAllOp() gopter.Gen {
	return gen.Const(ConcurrentOperation{
		Op: OpFindAll,
	})
}

func genFindRootOp() gopter.Gen {
	return gen.Const(ConcurrentOperation{
		Op: OpFindRoot,
	})
}

func genFindByParentIDOp() gopter.Gen {
	return gen.Bool().Map(func(hasParent bool) ConcurrentOperation {
		var parentID *domain.TaskID
		if hasParent {
			id := domain.NewTaskID()
			parentID = &id
		}
		return ConcurrentOperation{
			Op:     OpFindByParentID,
			TaskID: parentID,
		}
	})
}

func genDeleteOp() gopter.Gen {
	return gen.Const(0).Map(func(_ int) ConcurrentOperation {
		taskID := domain.NewTaskID()
		return ConcurrentOperation{
			Op:     OpDelete,
			TaskID: &taskID,
		}
	})
}

func genDeleteSubtreeOp() gopter.Gen {
	return gen.Const(0).Map(func(_ int) ConcurrentOperation {
		taskID := domain.NewTaskID()
		return ConcurrentOperation{
			Op:     OpDeleteSubtree,
			TaskID: &taskID,
		}
	})
}

// GenSingleOperation generates a single repository operation
func GenSingleOperation() gopter.Gen {
	return gen.OneGenOf(
		genSaveOp(),
		genFindByIDOp(),
		genFindAllOp(),
		genFindRootOp(),
		genFindByParentIDOp(),
		genDeleteOp(),
		genDeleteSubtreeOp(),
	)
}

// GenReadOperations generates a sequence of read-only operations
func GenReadOperations(numOps int) gopter.Gen {
	return gen.SliceOfN(numOps, gen.OneGenOf(
		genFindByIDOp(),
		genFindAllOp(),
		genFindRootOp(),
		genFindByParentIDOp(),
	)).SuchThat(func(v interface{}) bool {
		ops := v.([]ConcurrentOperation)
		return len(ops) > 0
	})
}

// GenWriteOperations generates a sequence of write operations
func GenWriteOperations(numOps int) gopter.Gen {
	return gen.SliceOfN(numOps, gen.OneGenOf(
		genSaveOp(),
		genDeleteOp(),
		genDeleteSubtreeOp(),
	)).SuchThat(func(v interface{}) bool {
		ops := v.([]ConcurrentOperation)
		return len(ops) > 0
	})
}

// GenMixedOperations generates a mix of read and write operations
func GenMixedOperations(numOps int) gopter.Gen {
	return gen.SliceOfN(numOps, GenSingleOperation()).SuchThat(func(v interface{}) bool {
		ops := v.([]ConcurrentOperation)
		return len(ops) > 0
	})
}

// GenValidTaskDTO generates a valid TaskDTO
func GenValidTaskDTO() gopter.Gen {
	return gopter.CombineGens(
		GenValidDescription(),
		GenOptionalTaskID(),
		GenValidPosition(),
		GenValidStatus(),
	).Map(func(values []interface{}) TaskDTO {
		desc := values[0].(string)
		parentID := values[1].(*domain.TaskID)
		position := values[2].(int)
		status := values[3].(domain.Status)
		
		task, _ := domain.NewTask(desc, parentID, position)
		task.ChangeStatus(status)
		
		return ToDTO(task)
	})
}

// GenTaskDTOSlice generates a slice of valid TaskDTOs
func GenTaskDTOSlice(minSize, maxSize int) gopter.Gen {
	return gen.Const(0).Map(func(_ int) []TaskDTO {
		// Generate a random size between minSize and maxSize
		n := minSize + (maxSize-minSize)/2 // Use middle value for simplicity
		dtos := make([]TaskDTO, n)
		for i := 0; i < n; i++ {
			desc := fmt.Sprintf("Task %d", i)
			task, _ := domain.NewTask(desc, nil, i)
			dtos[i] = ToDTO(task)
		}
		return dtos
	})
}

// GenValidJSON generates valid JSON for task arrays
func GenValidJSON() gopter.Gen {
	return gen.Const(0).Map(func(_ int) string {
		// Generate a small array of tasks
		n := 5
		dtos := make([]TaskDTO, n)
		for i := 0; i < n; i++ {
			desc := fmt.Sprintf("Task %d", i)
			task, _ := domain.NewTask(desc, nil, i)
			dtos[i] = ToDTO(task)
		}
		data, _ := json.MarshalIndent(dtos, "", "  ")
		return string(data)
	})
}
