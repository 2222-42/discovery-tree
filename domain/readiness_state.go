package domain

// ReadinessState represents whether a task is ready to be worked on based on ordering constraints
type ReadinessState struct {
	isReady             bool
	leftSiblingComplete bool
	allChildrenComplete bool
	reasons             []string
}

// NewReadinessState creates a new ReadinessState
func NewReadinessState(leftSiblingComplete, allChildrenComplete bool, reasons []string) ReadinessState {
	// A task is ready if:
	// 1. Its left sibling is complete (or it has no left sibling)
	// 2. All its children are complete (or it has no children)
	isReady := leftSiblingComplete && allChildrenComplete

	// Make a copy of reasons to avoid external mutation
	reasonsCopy := make([]string, len(reasons))
	copy(reasonsCopy, reasons)

	return ReadinessState{
		isReady:             isReady,
		leftSiblingComplete: leftSiblingComplete,
		allChildrenComplete: allChildrenComplete,
		reasons:             reasonsCopy,
	}
}

// IsReady returns whether the task is ready to be worked on
func (r ReadinessState) IsReady() bool {
	return r.isReady
}

// LeftSiblingComplete returns whether the left sibling is complete or doesn't exist
func (r ReadinessState) LeftSiblingComplete() bool {
	return r.leftSiblingComplete
}

// AllChildrenComplete returns whether all children are complete or no children exist
func (r ReadinessState) AllChildrenComplete() bool {
	return r.allChildrenComplete
}

// Reasons returns the reasons why the task is not ready (empty if ready)
func (r ReadinessState) Reasons() []string {
	// Return a copy to prevent external mutation
	reasonsCopy := make([]string, len(r.reasons))
	copy(reasonsCopy, r.reasons)
	return reasonsCopy
}
