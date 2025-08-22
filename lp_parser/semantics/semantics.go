package semantics

import (
	"fmt"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

const enableObjective = true
const disableObjective = false

func checkTerm(isObjectiveAndFirst bool, e parser.Expr, idTable map[string]int) error {
	switch expr := e.(type) {
	case *parser.Variable:
		if _, ok := idTable[expr.ID.Value]; !ok {
			return fmt.Errorf("undeclared Variable: %s", expr)
		}

		return nil
	case *parser.NumberLiteral:
		if isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("expected no NumberLiteral, received NumberLiteral %s", expr)
		}
	case *parser.UnaryExpr:
		return fmt.Errorf("invalid UnaryExpr (should not have UnaryExpr at this stage): %s", expr)
	case *parser.BinaryExpr:
		left := expr.Left
		right := expr.Right

		if _, ok := left.(*parser.NumberLiteral); !ok {
			return fmt.Errorf("expected NumberLiteral, received: %s", left)
		}

		v, ok := right.(*parser.Variable)
		if !ok {
			return fmt.Errorf("expected Variable, received: %s", right)
		}

		if _, ok = idTable[v.ID.Value]; !ok {
			return fmt.Errorf("undeclared Variable: %s", v)
		}

		return nil
	default:
		return fmt.Errorf("unknown Expr type: %T", e)
	}
}

func checkExpr(isObjectiveAndFirst bool, e parser.Expr, idTable map[string]int) error {
	switch expr := e.(type) {
	case *parser.Variable, *parser.NumberLiteral, *parser.UnaryExpr:
		return checkTerm(isObjectiveAndFirst, e, idTable)
	case *parser.BinaryExpr:
		left := expr.Left
		right := expr.Right
		op := expr.Operator

		switch op.Type {
		case lexer.TokenPlus:
			// check right, recurse on left
			if err := checkTerm(isObjectiveAndFirst, right, idTable); err != nil {
				return err
			}

			return checkExpr(disableObjective, left, idTable)
		case lexer.TokenAsterisk:
			return checkTerm(isObjectiveAndFirst, expr, idTable)
		default:
			return fmt.Errorf("invalid Expr operator: %s", expr)
		}
	default:
		return fmt.Errorf("unknown Expr type: %T", e)
	}
}

func checkNumber(e parser.Expr) error {
	if _, ok := e.(*parser.NumberLiteral); !ok {
		return fmt.Errorf("constant not found at RHS: %v", e)
	}

	return nil
}

func SemanticCheck(p *parser.Program) (map[string]int, error) {
	idTable := make(map[string]int)
	for i, decl := range p.Decls {
		if _, ok := idTable[decl.ID.Value]; ok {
			return nil, fmt.Errorf("duplicate variable: %v", decl.ID.Value)
		}

		idTable[decl.ID.Value] = i
	}

	if err := checkExpr(enableObjective, p.Objective.Expr, idTable); err != nil {
		return nil, err
	}

	for _, constraint := range p.Constraints {
		if err := checkExpr(disableObjective, constraint.Left, idTable); err != nil {
			return nil, err
		}

		if err := checkNumber(constraint.Right); err != nil {
			return nil, err
		}
	}

	return idTable, nil
}
