package solve

import (
	"math"
	"strings"
	"testing"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
	"github.com/animalat/Simplex-Algorithm/lp_parser/semantics"
)

func floatsEqual(a float64, b float64) bool {
	return math.Abs(a-b) < EPSILON
}

func assertProg(t *testing.T, s string, objectiveConstWanted float64, objectiveWanted []float64) error {
	t.Helper()

	tokens, err := lexer.Tokenize(strings.NewReader(s))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	p := &parser.Parser{Tokens: tokens}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	idTable, err := semantics.SemanticCheck(prog)
	if err != nil {
		return err
	}

	if testing.Verbose() {
		parser.PrintParse(prog)
	}

	objectiveConst, objective, err := getExprArr(prog.Objective.Expr, idTable, enableObjective)
	if err != nil {
		return err
	}

	// check const
	if !floatsEqual(objectiveConst, objectiveConstWanted) {
		t.Fatalf("objective constants not equal. Wanted: %.2f, Obtained: %.2f", objectiveConstWanted, objectiveConst)
	}

	// check objective numbers
	if len(objective) != len(objectiveWanted) {
		t.Fatalf("objective functions not equal size: Wanted size %d, obtained size %d", len(objectiveWanted), len(objective))
	}

	for i := range objective {
		if !floatsEqual(objective[i], objectiveWanted[i]) {
			t.Fatalf("objective values not equal. Wanted: %.2f, Obtained: %.2f, Position: %d", objectiveWanted[i], objective[i], i)
		}
	}

	return nil
}

func TestSolve_GetExprArr(t *testing.T) {
	err := assertProg(t, "let x1; let x2; let x3; max x1 + x2 + 3; s.t. x1 + x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 3, []float64{1, 1, 0})
	if err != nil {
		t.Errorf("GetExprArr failed (%v)", err)
	}

	err = assertProg(t, "let x1; let x2; let x3; max 4 * x1 + 5 * x3; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 0, []float64{4, 0, 5})
	if err != nil {
		t.Errorf("GetExprArr failed (%v)", err)
	}

	err = assertProg(t, "let x1; let x2; let x3; let x4; max 4 * x1 + x2 + 0 * x3 + 5 * x4 + 100; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 100, []float64{4, 1, 0, 5})
	if err != nil {
		t.Errorf("GetExprArr failed (%v)", err)
	}

	err = assertProg(t, "let x1; let x2; let x3; let x4; let x5; max 4 * x1 + x2 + 0 * x3 + 5 * x4 + 10 * x5 + 100; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 100, []float64{4, 1, 0, 5, 10})
	if err != nil {
		t.Errorf("GetExprArr failed (%v)", err)
	}
}
