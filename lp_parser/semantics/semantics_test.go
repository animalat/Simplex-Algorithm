package semantics

import (
	"strings"
	"testing"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

func TestSemantics_Term(t *testing.T) {
	numberTest := &parser.NumberLiteral{Value: 5, Line: 0}
	m := make(map[string]bool)
	err := checkTerm(enableObjective, numberTest, m)
	if err != nil {
		t.Errorf("valid parser.NumberLiteral failed: %v", numberTest)
	}

	err = checkTerm(disableObjective, numberTest, m)
	if err == nil {
		t.Errorf("invalid parser.NumberLiteral passed: %v", numberTest)
	}

	variableTest := &parser.Variable{ID: lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 0}}
	err = checkTerm(enableObjective, variableTest, m)
	if err == nil {
		t.Errorf("invalid parser.Variable passed: %v", variableTest)
	}

	m[variableTest.ID.Value] = true
	err = checkTerm(disableObjective, variableTest, m)
	if err != nil {
		t.Errorf("valid parser.Variable failed: %v", variableTest)
	}

	unaryExprTest := &parser.UnaryExpr{Operator: lexer.Token{Type: lexer.TokenMinus, Value: "-", Line: 0}, Expr: variableTest, Line: 0}
	err = checkTerm(disableObjective, unaryExprTest, m)
	if err == nil {
		t.Errorf("invalid parser.UnaryExpr passed: %v", unaryExprTest)
	}

	variableTest2 := &parser.Variable{ID: lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 0}}
	unaryExprTest2 := &parser.UnaryExpr{Operator: lexer.Token{Type: lexer.TokenMinus, Value: "-", Line: 0}, Expr: variableTest2, Line: 0}
	err = checkTerm(disableObjective, unaryExprTest2, m)
	if err == nil {
		t.Errorf("invalid parser.UnaryExpr passed: %v", unaryExprTest2)
	}

	binaryTest := &parser.BinaryExpr{Left: numberTest, Operator: lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 0}, Right: variableTest, Line: 0}
	err = checkTerm(disableObjective, binaryTest, m)
	if err != nil {
		t.Errorf("valid parser.BinaryExpr failed: %v", binaryTest)
	}

	binaryTest2 := &parser.BinaryExpr{Left: numberTest, Operator: lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest2, m)
	if err == nil {
		t.Errorf("invalid parser.BinaryExpr passed: %v", binaryTest2)
	}

	binaryTest3 := &parser.BinaryExpr{Left: variableTest, Operator: lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest3, m)
	if err == nil {
		t.Errorf("invalid parser.BinaryExpr passed: %v", binaryTest3)
	}

	m[variableTest2.ID.Value] = true
	binaryTest4 := &parser.BinaryExpr{Left: variableTest, Operator: lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest4, m)
	if err == nil {
		t.Errorf("invalid parser.BinaryExpr passed: %v", binaryTest4)
	}
}

func assertProg(t *testing.T, s string) error {
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

	err = SemanticCheck(prog)
	if err != nil {
		return err
	}

	if testing.Verbose() {
		parser.PrintParse(prog)
	}

	return nil
}

func TestSemantics_Expr(t *testing.T) {
	err := assertProg(t, "let x1; let x2; let x3; max x1 + x2 + 3; s.t. x1 + x2 <= 3; x1 + x2 + 3 * x3 >= 5;")
	if err != nil {
		t.Errorf("valid program failed (%v)", err)
	}

	err = assertProg(t, "let x1; let x2; let x3; max 4 * x1 + 5 * x2 + 3; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;")
	if err != nil {
		t.Errorf("valid program failed (%v)", err)
	}

	err = assertProg(t, "let x1; let x2; let x3; max x1 + x2 + 3; s.t. 5 * x1 + 3 * x2 * 5 - 5 <= 3; x1 + x2 + 3 * x3 + 1 >= 5;")
	if err == nil {
		t.Errorf("invalid program passed")
	}

	err = assertProg(t, "let x1; let x2; max x1 + 5 + x2; s.t. x1 + x2 = 5;")
	if err == nil {
		t.Errorf("invalid program passed")
	}
}
