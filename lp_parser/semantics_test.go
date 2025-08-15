package parser

import "testing"

func TestSemantics_Term(t *testing.T) {
	numberTest := &NumberLiteral{5, 0}
	m := make(map[string]bool)
	err := checkTerm(enableObjective, numberTest, m)
	if err != nil {
		t.Errorf("valid NumberLiteral failed: %v", numberTest)
	}

	err = checkTerm(disableObjective, numberTest, m)
	if err == nil {
		t.Errorf("invalid NumberLiteral passed: %v", numberTest)
	}

	variableTest := &Variable{ID: Token{Type: TokenId, Value: "x1", Line: 0}}
	err = checkTerm(enableObjective, variableTest, m)
	if err == nil {
		t.Errorf("invalid Variable passed: %v", variableTest)
	}

	m[variableTest.ID.Value] = true
	err = checkTerm(disableObjective, variableTest, m)
	if err != nil {
		t.Errorf("valid Variable failed: %v", variableTest)
	}

	unaryExprTest := &UnaryExpr{Operator: Token{Type: TokenMinus, Value: "-", Line: 0}, Expr: variableTest, Line: 0}
	err = checkTerm(disableObjective, unaryExprTest, m)
	if err != nil {
		t.Errorf("valid UnaryExpr failed: %v", unaryExprTest)
	}
	
	variableTest2 := &Variable{ID: Token{Type: TokenId, Value: "x2", Line: 0}}
	unaryExprTest2 := &UnaryExpr{Operator: Token{Type: TokenMinus, Value: "-", Line: 0}, Expr: variableTest2, Line: 0}
	err = checkTerm(disableObjective, unaryExprTest2, m)
	if err == nil {
		t.Errorf("invalid UnaryExpr passed: %v", unaryExprTest2)
	}

	binaryTest := &BinaryExpr{Left: numberTest, Operator: Token{Type: TokenAsterisk, Value: "*", Line: 0}, Right: variableTest, Line: 0}
	err = checkTerm(disableObjective, binaryTest, m)
	if err != nil {
		t.Errorf("valid BinaryExpr failed: %v", binaryTest)
	}

	binaryTest2 := &BinaryExpr{Left: numberTest, Operator: Token{Type: TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest2, m)
	if err == nil {
		t.Errorf("invalid BinaryExpr passed: %v", binaryTest2)
	}

	binaryTest3 := &BinaryExpr{Left: variableTest, Operator: Token{Type: TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest3, m)
	if err == nil {
		t.Errorf("invalid BinaryExpr passed: %v", binaryTest3)
	}

	m[variableTest2.ID.Value] = true
	binaryTest4 := &BinaryExpr{Left: variableTest, Operator: Token{Type: TokenAsterisk, Value: "*", Line: 0}, Right: variableTest2, Line: 0}
	err = checkTerm(disableObjective, binaryTest4, m)
	if err == nil {
		t.Errorf("invalid BinaryExpr passed: %v", binaryTest4)
	}
}

// func TestSemantics_Expr(t *testing.T) {

// }
