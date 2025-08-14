package parser

import "fmt"

func checkTerm(isObjectiveAndFirst bool, e Expr, idTable map[string]bool) error {
	switch expr := e.(type) {
	case *Variable:
		if _, ok := idTable[expr.ID.Value]; !ok {
			return fmt.Errorf("undeclared Variable: %s", expr)
		}

		return nil
	case *NumberLiteral:
		if isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("expected no NumberLiteral, received NumberLiteral %s", expr)
		}
	case *UnaryExpr:
		inner := expr.Expr
		if v, ok := inner.(*Variable); ok {
			if _, ok = idTable[v.ID.Value]; ok {
				return nil
			} else {
				return fmt.Errorf("undeclared Variable: %s", inner)
			}
		} else if _, ok := inner.(*NumberLiteral); ok && isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("invalid UnaryExpr: %s", expr)
		}
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right

		// TODO: implement UnaryExpr with nested NumberLiteral check
		if _, ok := left.(*NumberLiteral); !ok {
			return fmt.Errorf("expected NumberLiteral, received: %s", left)
		}

		v, ok := right.(*Variable)
		if !ok {
			return fmt.Errorf("expected Variable, received: %s", right)
		}

		if _, ok = idTable[v.ID.Value]; !ok {
			return fmt.Errorf("undeclared Variable: %s", v)
		}

		return nil
	default:
		return fmt.Errorf("unknown type Expr")
	}
}

func checkExpr(isObjectiveAndFirst bool, e Expr, idTable map[string]bool) error {
	switch expr := e.(type) {
	case *Variable, *NumberLiteral, *UnaryExpr:
		return checkTerm(true, e, idTable)
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right
		op := expr.Operator

		switch op.Type {
		case TokenPlus:
			// check right, recurse on left
			if err := checkTerm(false, right, idTable); err != nil {
				return err
			}

			return checkExpr(true, left, idTable)
		case TokenAsterisk:
			return checkTerm(true, expr, idTable)
		default:
			return fmt.Errorf("invalid Expr operator: %s", expr)
		}
	default:
		return fmt.Errorf("unknown type Expr")
	}
}

func SemanticCheck(p *Program) error {
	idTable := make(map[string]bool)
	for _, decl := range p.Decls {
		idTable[decl.ID.Value] = true
	}

	if err := checkExpr(true, p.Objective.Expr, idTable); err != nil {
		return err
	}

	return nil
}
