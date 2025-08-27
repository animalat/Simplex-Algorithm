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

	if prog.Objective.Expr, err = SimplifyExpr(prog.Objective.Expr); err != nil {
		t.Fatalf("Simplification failed on objective: %v", err)
	}

	for i, constraint := range prog.Constraints {
		if constraint.Left, err = SimplifyExpr(constraint.Left); err != nil {
			t.Fatalf("Simplification failed on constraint %d: %v", i, err)
		}

		if constraint.Right, err = SimplifyExpr(constraint.Right); err != nil {
			t.Fatalf("Simplification failed on constraint %d: %v", i, err)
		}
	}

	got := fmt.Sprint(prog.Objective.Expr)
	want := "((5 * x1) + (5 * x2))"
	if got != want {
		t.Errorf("Objective simplify mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	if len(prog.Constraints) != 2 {
		t.Fatalf("Expected 2 constraint, got %d", len(prog.Constraints))
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
		t.Errorf("Constraint left side mismatch:\nGot:  %v\nWant: %v", got, want)
	}

	got = fmt.Sprint(prog.Constraints[1].Right)
	want = "3"
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}
}

func TestSimplify_CollectLikeTerms(t *testing.T) {
	input := "let x1; let x2; max 3 * x1 + x2 + 10 + x1 + 4 * x2 + 5 + 6 + 3; s.t. x1 + x2 + 4 * x1 + 6 * x2 + 4 + 5 <= 3 + x1 + x2 + 3 * x1 + 4 + 3 * x2 + 5;"
	tokens, err := lexer.Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenizing failed: %v", err)
	}

	p := &parser.Parser{Tokens: tokens}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	beforeObjective := prog.Objective.Expr
	prog.Objective.Expr, _, err = CollectLikeTerms(prog.Objective.Expr, &parser.NumberLiteral{Value: 0}, enableObjective, make(map[string]float64))
	if err != nil {
		t.Fatalf("failed to collect like terms on objective: %v", err)
	}
	for i, constraint := range prog.Constraints {
		constraint.Left, constraint.Right, err = CollectLikeTerms(constraint.Left, constraint.Right, disableObjective, make(map[string]float64))
		if err != nil {
			t.Fatalf("failed to collect like terms on constraint %d: %v", i, err)
		}
	}

	if len(prog.Constraints) != 1 {
		t.Fatalf("Expected 1 constraint, got %d", len(prog.Constraints))
	}

	if testing.Verbose() {
		fmt.Printf("before: %v, after: %v", beforeObjective, prog.Objective.Expr)
	}

	wantA := "((3 * x2) + (1 * x1))"
	wantB := "((1 * x1) + (3 * x2))"
	got := fmt.Sprint(prog.Constraints[0].Left)
	if got != wantA && got != wantB {
		t.Errorf("Constraint left side mismatch:\nGot:  %v\nWant: %v", got, wantA)
	}

	want := "3"
	got = fmt.Sprint(prog.Constraints[0].Right)
	if got != want {
		t.Errorf("Constraint right side mismatch:\nGot:  %v\nWant: %v", got, want)
	}
}
