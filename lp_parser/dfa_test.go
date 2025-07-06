package parser

import (
	"fmt"
	"testing"
)

func TestDFAConstruction(t *testing.T) {
	dfa := NewDFA()

	if len(dfa.Transitions) == 0 {
		t.Error("Expected some transitions but got none")
	}

	// Print transitions (optional, useful for debugging)
	for key, next := range dfa.Transitions {
		fmt.Printf("From %q with %q -> %q\n", key.State, key.Input, next)
	}
}
