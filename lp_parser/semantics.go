package parser

func checkTerm(e *Expr, idTable map[string]bool, isObjectiveAndFirst bool) error {
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
		return fmt.Errorf("Unexpected UnaryExpr: %s", expr)
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right

		if _, ok = left.(*NumberLiteral); !ok {
			return fmt.Errorf("Expected NumberLiteral, received: %s", left)
		}

		if _, ok := right.(*Variable); !ok {
			return fmt.Errorf("Expected Variable, received: %s", right)
		}

		return nil
	}
}

func checkExpr(e *Expr, idTable map[string]bool, isObjectiveAndFirst bool) error {
	switch expr := e.(type) {
	case *Variable:
		return nil
	case *NumberLiteral:
		return fmt.Errorf("Expected no NumberLiteral, received NumberLiteral %s", expr)
	case *UnaryExpr:
		inner := expr.Expr
		if v, ok := inner.(*Variable); ok {
			if _, ok = idTable[v.ID.Value]; ok {
				return nil
			} else {
				return fmt.Errorf("Undeclared variable: %s", v)
			}
		} else if _, ok = inner.(*NumberLiteral); ok && isObjectiveAndFirst {
			return nil
		} else {
			return fmt.Errorf("UnaryExpr with nested Expr: %s", expr)
		}
	case *BinaryExpr:
		left := expr.Left
		right := expr.Right
		

	}
	return nil
}

func SemanticCheck(p *Program) error {
	idTable := make(map[string]bool)
	for _, decl := range p.Decls {
		idTable[decl.ID.Value] = true
	}

	return nil
}
