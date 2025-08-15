package parser

import "fmt"

const enableObjective = true
const disableObjective = false

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
		return fmt.Errorf("invalid UnaryExpr (should not have UnaryExpr at this stage): %s", expr)
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right

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
		return fmt.Errorf("unknown Expr type: %T", e)
	}
}

func checkExpr(isObjectiveAndFirst bool, e Expr, idTable map[string]bool) error {
	switch expr := e.(type) {
	case *Variable, *NumberLiteral, *UnaryExpr:
		return checkTerm(isObjectiveAndFirst, e, idTable)
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right
		op := expr.Operator

		switch op.Type {
		case TokenPlus:
			// check right, recurse on left
			if err := checkTerm(isObjectiveAndFirst, right, idTable); err != nil {
				return err
			}

			return checkExpr(disableObjective, left, idTable)
		case TokenAsterisk:
			return checkTerm(isObjectiveAndFirst, expr, idTable)
		default:
			return fmt.Errorf("invalid Expr operator: %s", expr)
		}
	default:
		return fmt.Errorf("unknown Expr type: %T", e)
	}
}

func checkNumber(e Expr) error {
	if _, ok := e.(*NumberLiteral); !ok {
		return fmt.Errorf("constant not found at RHS: %v", e)
	}

	return nil
}

func SemanticCheck(p *Program) error {
	idTable := make(map[string]bool)
	for _, decl := range p.Decls {
		idTable[decl.ID.Value] = true
	}

	if err := checkExpr(enableObjective, p.Objective.Expr, idTable); err != nil {
		return err
	}

	for _, constraint := range p.Constraints {
		if err := checkExpr(disableObjective, constraint.Left, idTable); err != nil {
			return err
		}

		if err := checkNumber(constraint.Right); err != nil {
			return err
		}
	}

	return nil
}
