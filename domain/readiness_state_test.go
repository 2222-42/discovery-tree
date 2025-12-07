package domain

import (
	"testing"
)

func TestNewReadinessState_Ready(t *testing.T) {
	// When both left sibling is complete and all children are complete
	state := NewReadinessState(true, true, []string{})

	if !state.IsReady() {
		t.Error("Expected task to be ready when left sibling and all children are complete")
	}

	if !state.LeftSiblingComplete() {
		t.Error("Expected LeftSiblingComplete to be true")
	}

	if !state.AllChildrenComplete() {
		t.Error("Expected AllChildrenComplete to be true")
	}

	if len(state.Reasons()) != 0 {
		t.Errorf("Expected no reasons when ready, got %d reasons", len(state.Reasons()))
	}
}

func TestNewReadinessState_NotReady_LeftSiblingIncomplete(t *testing.T) {
	// When left sibling is not complete
	reasons := []string{"Left sibling is not complete"}
	state := NewReadinessState(false, true, reasons)

	if state.IsReady() {
		t.Error("Expected task to not be ready when left sibling is incomplete")
	}

	if state.LeftSiblingComplete() {
		t.Error("Expected LeftSiblingComplete to be false")
	}

	if !state.AllChildrenComplete() {
		t.Error("Expected AllChildrenComplete to be true")
	}

	if len(state.Reasons()) != 1 {
		t.Errorf("Expected 1 reason, got %d", len(state.Reasons()))
	}

	if state.Reasons()[0] != "Left sibling is not complete" {
		t.Errorf("Expected reason 'Left sibling is not complete', got '%s'", state.Reasons()[0])
	}
}

func TestNewReadinessState_NotReady_ChildrenIncomplete(t *testing.T) {
	// When children are not complete
	reasons := []string{"Not all children are complete"}
	state := NewReadinessState(true, false, reasons)

	if state.IsReady() {
		t.Error("Expected task to not be ready when children are incomplete")
	}

	if !state.LeftSiblingComplete() {
		t.Error("Expected LeftSiblingComplete to be true")
	}

	if state.AllChildrenComplete() {
		t.Error("Expected AllChildrenComplete to be false")
	}

	if len(state.Reasons()) != 1 {
		t.Errorf("Expected 1 reason, got %d", len(state.Reasons()))
	}
}

func TestNewReadinessState_NotReady_BothIncomplete(t *testing.T) {
	// When both left sibling and children are not complete
	reasons := []string{"Left sibling is not complete", "Not all children are complete"}
	state := NewReadinessState(false, false, reasons)

	if state.IsReady() {
		t.Error("Expected task to not be ready when both conditions are not met")
	}

	if state.LeftSiblingComplete() {
		t.Error("Expected LeftSiblingComplete to be false")
	}

	if state.AllChildrenComplete() {
		t.Error("Expected AllChildrenComplete to be false")
	}

	if len(state.Reasons()) != 2 {
		t.Errorf("Expected 2 reasons, got %d", len(state.Reasons()))
	}
}

func TestReadinessState_ReasonsImmutability(t *testing.T) {
	// Test that modifying the returned reasons slice doesn't affect the internal state
	originalReasons := []string{"Reason 1", "Reason 2"}
	state := NewReadinessState(false, false, originalReasons)

	// Get reasons and modify them (intentionally testing immutability)
	reasons := state.Reasons()
	reasons[0] = "Modified"
	_ = append(reasons, "New reason") // Intentionally not using result to test immutability

	// Get reasons again and verify they haven't changed
	newReasons := state.Reasons()
	if len(newReasons) != 2 {
		t.Errorf("Expected 2 reasons, got %d", len(newReasons))
	}

	if newReasons[0] != "Reason 1" {
		t.Errorf("Expected first reason to be 'Reason 1', got '%s'", newReasons[0])
	}
}

func TestReadinessState_InputReasonsImmutability(t *testing.T) {
	// Test that modifying the input reasons slice doesn't affect the internal state
	inputReasons := []string{"Reason 1"}
	state := NewReadinessState(false, true, inputReasons)

	// Modify the input slice (intentionally testing immutability)
	inputReasons[0] = "Modified"
	_ = append(inputReasons, "New reason") // Intentionally not using result to test immutability

	// Verify the state's reasons haven't changed
	stateReasons := state.Reasons()
	if len(stateReasons) != 1 {
		t.Errorf("Expected 1 reason, got %d", len(stateReasons))
	}

	if stateReasons[0] != "Reason 1" {
		t.Errorf("Expected reason to be 'Reason 1', got '%s'", stateReasons[0])
	}
}
