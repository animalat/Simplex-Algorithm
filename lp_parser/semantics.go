package parser

func checkTerm(bool isObjectiveAndFirst, e *Expr, idTable map[string]bool) error {
	switch expr := e.(type) {
	case *Variable:
		if _, ok := idTable[expr.ID.Type]; !ok {
			return fmt.Errorf("Undeclared Variable: %s", expr)
		}

		return nil
	case *NumberLiteral:
		if isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("Expected no NumberLiteral, received NumberLiteral %s", expr)
		}
	case *UnaryExpr:
		inner := expr.Expr
		if v, ok := inner.(*Variable); ok {
			if _, ok = idTable[v.ID.Value]; ok {
				return nil
			} else {
				return fmt.Errorf("Undeclared Variable: %s", inner)
			}
		} else if _, ok := inner.(*NumberLiteral); ok && isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("Invalid UnaryExpr: %s", expr)
		}
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right

		if _, ok := left.(*NumberLiteral); !ok {
			return fmt.Errorf("Expected NumberLiteral, received: %s", left)
		}

		v, ok := right.(*Variable);
		if !ok {
			return fmt.Errorf("Expected Variable, received: %s", right)
		}

		if _, ok = idTable[v.ID.Value]; !ok {
			return fmt.Errorf("Undeclared Variable: %s", v)
		}

		return nil
	case default:
		return fmt.Errorf("Unknown type Expr")
	}
}

func checkExpr(isObjectiveAndFirst bool, e *Expr, idTable map[string]bool) error {
	switch expr := e.(type) {
	case *Variable, *NumberLiteral, *UnaryExpr:
		return checkTerm(true, e, idTable)
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right
		op := expr.Operator

		switch op.Type {
		case TokenPlus:
			// recurse on left, check right
			err := checkTerm(false, right, idTable)
			if err != nil {
				return err
			}

			return checkExpr(true, left, idTable)
		case TokenAsterisk:
			return checkTerm(true, expr, idTable)
		case default:
			return fmt.Errorf("Invalid Expr operator: %s", expr)
		}
	case default:
		return fmt.Errorf("Unknown type Expr")
	}
}

func SemanticCheck(p *Program) error {
	idTable := make(map[string]bool)
	for _, decl := range p.Decls {
		idTable[decl.ID.Value] = true
	}

	return nil
}
