package simplify

import (
	"fmt"
	"strings"
	"testing"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

func TestSimplify_Distribute(t *testing.T) {
	input := "let x1; let x2; max (3 + 2) * (x1 + x2); s.t. ((3 * 4 * (1 + 9)) * x2 + 15 + (1 + 5 + 2 * 2) * x1) / 5 <= 5 * (3 + 1); -3 * x1 * 4 * 5 + 5 * -3 * -(4 * 1 + 4) <= 3;"
	tokens, err := lexer.Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenizing failed: %v", err)
	}

	parser := &parser.Parser{Tokens: tokens}
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

	if len(prog.Constraints) != 2 {
		t.Fatalf("Expected 1 constraint, got %d", len(prog.Constraints))
	}

	got = fmt.Sprint(prog.Constraints[0].Left)
	want = "(((24 * x2) + 3) + (2 * x1))"
	if got != want {
		t.Errorf("Constraint left side mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	got = fmt.Sprint(prog.Constraints[0].Right)
	want = "20"
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	got = fmt.Sprint(prog.Constraints[1].Left)
	want = "((-60 * x1) + 120)"
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	got = fmt.Sprint(prog.Constraints[1].Right)
	want = "3"
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}
}
