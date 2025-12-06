package domain

// ReadinessEvaluator evaluates whether a task is ready to be worked on
type ReadinessEvaluator interface {
	// EvaluateReadiness evaluates the readiness state of a task based on ordering constraints
	EvaluateReadiness(taskID TaskID) (ReadinessState, error)
}

// ReadinessEvaluatorService implements ReadinessEvaluator
type ReadinessEvaluatorService struct {
	repo      TaskRepository
	navigator TreeNavigator
}

// NewReadinessEvaluatorService creates a new ReadinessEvaluatorService
func NewReadinessEvaluatorService(repo TaskRepository, navigator TreeNavigator) *ReadinessEvaluatorService {
	return &ReadinessEvaluatorService{
		repo:      repo,
		navigator: navigator,
	}
}

// EvaluateReadiness evaluates the readiness state of a task
// A task is ready if:
// 1. Its left sibling is DONE (or it has no left sibling)
// 2. All its children are DONE (or it has no children)
func (s *ReadinessEvaluatorService) EvaluateReadiness(taskID TaskID) (ReadinessState, error) {
	// First, verify the task exists
	_, err := s.repo.FindByID(taskID)
	if err != nil {
		return ReadinessState{}, err
	}

	var reasons []string
	leftSiblingComplete := true
	allChildrenComplete := true

	// Check left sibling completion status
	leftSibling, err := s.navigator.GetLeftSibling(taskID)
	if err != nil {
		return ReadinessState{}, err
	}

	if leftSibling != nil {
		// Left sibling exists, check if it's DONE
		if leftSibling.Status() != StatusDONE {
			leftSiblingComplete = false
			reasons = append(reasons, "left sibling is not complete")
		}
	}
	// If leftSibling is nil, leftSiblingComplete remains true (no left sibling)

	// Check all children completion status
	children, err := s.navigator.GetChildren(taskID)
	if err != nil {
		return ReadinessState{}, err
	}

	if len(children) > 0 {
		// Task has children, check if all are DONE
		for _, child := range children {
			if child.Status() != StatusDONE {
				allChildrenComplete = false
				reasons = append(reasons, "not all children are complete")
				break // Only need to find one incomplete child
			}
		}
	}
	// If no children, allChildrenComplete remains true

	// Create and return the readiness state
	return NewReadinessState(leftSiblingComplete, allChildrenComplete, reasons), nil
}
