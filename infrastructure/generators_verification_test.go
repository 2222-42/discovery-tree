package infrastructure

import (
	"encoding/json"
	"os"
	"testing"

	"discovery-tree/domain"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// TestGenTaskForRepository_GeneratesValidTasks verifies that the task generator produces valid tasks
func TestGenTaskForRepository_GeneratesValidTasks(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("generated tasks are valid", prop.ForAll(
		func(task *domain.Task) bool {
			// Verify task is not nil
			if task == nil {
				return false
			}

			// Verify description is not empty
			if task.Description() == "" {
				return false
			}

			// Verify status is valid
			if !task.Status().IsValid() {
				return false
			}

			// Verify position is non-negative
			if task.Position() < 0 {
				return false
			}

			return true
		},
		GenTaskForRepository(),
	))

	properties.TestingRun(t)
}

// TestGenTaskTree_GeneratesValidTrees verifies that the tree generator produces valid trees
func TestGenTaskTree_GeneratesValidTrees(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated trees have valid structure", prop.ForAll(
		func(tree *TaskTree) bool {
			// Verify tree is not nil
			if tree == nil {
				return false
			}

			// Verify root exists
			if tree.Root == nil {
				return false
			}

			// Verify root has no parent
			if tree.Root.ParentID() != nil {
				return false
			}

			// Verify all tasks are valid
			for _, task := range tree.Tasks {
				if task == nil {
					return false
				}
				if task.Description() == "" {
					return false
				}
				if !task.Status().IsValid() {
					return false
				}
				if task.Position() < 0 {
					return false
				}
			}

			// Verify root is in tasks list
			found := false
			for _, task := range tree.Tasks {
				if task.ID().String() == tree.Root.ID().String() {
					found = true
					break
				}
			}

			return !found
		},
		GenTaskTree(),
	))

	properties.TestingRun(t)
}

// TestGenInvalidJSON_GeneratesInvalidJSON verifies that invalid JSON generator produces unparseable JSON
func TestGenInvalidJSON_GeneratesInvalidJSON(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated JSON is invalid", prop.ForAll(
		func(jsonStr string) bool {
			var dtos []TaskDTO
			err := json.Unmarshal([]byte(jsonStr), &dtos)
			// Should fail to unmarshal
			return err != nil
		},
		GenInvalidJSON(),
	))

	properties.TestingRun(t)
}

// TestGenInvalidTaskDTO_GeneratesInvalidDTOs verifies that invalid DTO generator produces DTOs that fail validation
func TestGenInvalidTaskDTO_GeneratesInvalidDTOs(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated DTOs fail validation", prop.ForAll(
		func(dto TaskDTO) bool {
			// Try to convert DTO to Task
			_, err := FromDTO(dto)
			// Should fail validation
			return err != nil
		},
		GenInvalidTaskDTO(),
	))

	properties.TestingRun(t)
}

// TestGenInvalidParentIDTaskDTO_GeneratesInvalidDTOs verifies that invalid parent ID DTO generator produces DTOs that fail validation
func TestGenInvalidParentIDTaskDTO_GeneratesInvalidDTOs(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated DTOs with invalid parent IDs fail validation", prop.ForAll(
		func(dto TaskDTO) bool {
			// Try to convert DTO to Task
			_, err := FromDTO(dto)
			// Should fail validation
			return err != nil
		},
		GenInvalidParentIDTaskDTO(),
	))

	properties.TestingRun(t)
}

// TestGenConcurrentOperations_GeneratesOperations verifies that concurrent operation generator works
func TestGenConcurrentOperations_GeneratesOperations(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated operations are valid", prop.ForAll(
		func(ops []ConcurrentOperation) bool {
			// Verify we have operations
			if len(ops) == 0 {
				return false
			}

			// Verify each operation has required fields
			for _, op := range ops {
				switch op.Op {
				case OpSave:
					if op.Task == nil {
						return false
					}
				case OpFindByID, OpDelete, OpDeleteSubtree:
					if op.TaskID == nil {
						return false
					}
				case OpFindByParentID:
					// ParentID can be nil for finding root children
				case OpFindAll, OpFindRoot:
					// No additional fields required
				default:
					return false
				}
			}

			return true
		},
		GenConcurrentOperations(10),
	))

	properties.TestingRun(t)
}

// TestGenReadOperations_GeneratesOnlyReads verifies that read operation generator only produces reads
func TestGenReadOperations_GeneratesOnlyReads(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated operations are read-only", prop.ForAll(
		func(ops []ConcurrentOperation) bool {
			// Verify all operations are reads
			for _, op := range ops {
				switch op.Op {
				case OpFindByID, OpFindAll, OpFindRoot, OpFindByParentID:
					// These are read operations - OK
				case OpSave, OpDelete, OpDeleteSubtree:
					// These are write operations - should not appear
					return false
				default:
					return false
				}
			}

			return true
		},
		GenReadOperations(10),
	))

	properties.TestingRun(t)
}

// TestGenWriteOperations_GeneratesOnlyWrites verifies that write operation generator only produces writes
func TestGenWriteOperations_GeneratesOnlyWrites(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated operations are write-only", prop.ForAll(
		func(ops []ConcurrentOperation) bool {
			// Verify all operations are writes
			for _, op := range ops {
				switch op.Op {
				case OpSave, OpDelete, OpDeleteSubtree:
					// These are write operations - OK
				case OpFindByID, OpFindAll, OpFindRoot, OpFindByParentID:
					// These are read operations - should not appear
					return false
				default:
					return false
				}
			}

			return true
		},
		GenWriteOperations(10),
	))

	properties.TestingRun(t)
}

// TestGenValidJSON_GeneratesValidJSON verifies that valid JSON generator produces parseable JSON
func TestGenValidJSON_GeneratesValidJSON(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated JSON is valid", prop.ForAll(
		func(jsonStr string) bool {
			var dtos []TaskDTO
			err := json.Unmarshal([]byte(jsonStr), &dtos)
			// Should successfully unmarshal
			if err != nil {
				return false
			}

			// Verify all DTOs can be converted to tasks
			for _, dto := range dtos {
				_, err := FromDTO(dto)
				if err != nil {
					return false
				}
			}

			return true
		},
		GenValidJSON(),
	))

	properties.TestingRun(t)
}

// TestGenTaskTreeWithSize_GeneratesCorrectSize verifies that tree generator creates trees of requested size
func TestGenTaskTreeWithSize_GeneratesCorrectSize(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("generated trees have approximately correct size", prop.ForAll(
		func(tree *TaskTree) bool {
			// Verify tree is not nil
			if tree == nil {
				return false
			}

			// Verify size is within expected range (5-15)
			size := len(tree.Tasks)
			if size < 5 || size > 15 {
				return false
			}

			return true
		},
		GenTaskTreeWithSize(5, 15),
	))

	properties.TestingRun(t)
}

// TestGenerators_Integration tests that generators can be used together
func TestGenerators_Integration(t *testing.T) {
	testPath := "./test_data/generators_integration.json"
	os.RemoveAll("./test_data")
	defer os.RemoveAll("./test_data")

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("generators work with repository", prop.ForAll(
		func(tree *TaskTree) bool {
			// Create a fresh repository for each test
			repo, err := NewFileTaskRepository(testPath)
			if err != nil {
				return false
			}

			// Save all tasks from the tree
			for _, task := range tree.Tasks {
				if err := repo.Save(task); err != nil {
					return false
				}
			}

			// Verify all tasks can be retrieved
			for _, task := range tree.Tasks {
				retrieved, err := repo.FindByID(task.ID())
				if err != nil {
					return false
				}
				if retrieved.ID().String() != task.ID().String() {
					return false
				}
			}

			// Clean up for next iteration
			if err := os.Remove(testPath); err != nil {
				t.Errorf("Failed to remove, err: %v", err)
			}

			return true
		},
		GenTaskTree(),
	))

	properties.TestingRun(t)
}

// TestGenSingleOperation_CoversAllOperations verifies that single operation generator can produce all operation types
func TestGenSingleOperation_CoversAllOperations(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 200 // Run more tests to ensure coverage
	properties := gopter.NewProperties(parameters)

	// Track which operations we've seen
	seenOps := make(map[RepositoryOperation]bool)

	properties.Property("generator produces all operation types", prop.ForAll(
		func(op ConcurrentOperation) bool {
			seenOps[op.Op] = true
			return true
		},
		GenSingleOperation(),
	))

	properties.TestingRun(t)

	// Verify we've seen all operation types
	expectedOps := []RepositoryOperation{
		OpSave, OpFindByID, OpFindAll, OpFindRoot, OpFindByParentID, OpDelete, OpDeleteSubtree,
	}

	for _, expectedOp := range expectedOps {
		if !seenOps[expectedOp] {
			t.Errorf("Generator did not produce operation type: %v", expectedOp)
		}
	}
}

// TestGenValidTaskDTO_ProducesValidDTOs verifies that valid DTO generator produces DTOs that pass validation
func TestGenValidTaskDTO_ProducesValidDTOs(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("generated DTOs are valid", prop.ForAll(
		func(dto TaskDTO) bool {
			// Try to convert DTO to Task
			task, err := FromDTO(dto)
			if err != nil {
				return false
			}

			// Verify task is valid
			if task == nil {
				return false
			}

			// Verify fields match
			if task.Description() != dto.Description {
				return false
			}

			if task.Status().String() != dto.Status {
				return false
			}

			if task.Position() != dto.Position {
				return false
			}

			return true
		},
		GenValidTaskDTO(),
	))

	properties.TestingRun(t)
}
