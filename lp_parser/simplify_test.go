package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimplify_Distribute(t *testing.T) {
	input := "let x1; let x2; max (3 + 2) * (x1 + x2); s.t. (3 + x1) / x2 <= 5 * (3 + 1);"
	tokens, err := Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenizing failed: %v", err)
	}

	parser := &Parser{Tokens: tokens}
	prog, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if err := SimplifyProgram(prog); err != nil {
		t.Fatalf("Simplification failed: %v", err)
	}

	got := fmt.Sprint(prog.Objective.Expr)
	want := "((5 * x1) + (5 * x2))"
	if got != want {
		t.Errorf("Objective simplify mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	if len(prog.Constraints) != 1 {
		t.Fatalf("Expected 1 constraint, got %d", len(prog.Constraints))
	}

	got = fmt.Sprint(prog.Constraints[0].Left)
	want = "((3 / x2) + (x1 / x2))"
	if got != want {
		t.Errorf("Constraint left side mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	got = fmt.Sprint(prog.Constraints[0].Right)
	want = "(5 * 4)"
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}
}
