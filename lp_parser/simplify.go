package parser

import "fmt"

func SimplifyProgram(p *Program) error {
	var err error
	p.Objective.Expr, err = SimplifyExpr(p.Objective.Expr)
	if err != nil {
		return err
	}
	for _, constraint := range p.Constraints {
		constraint.Left, err = SimplifyExpr(constraint.Left)
		if err != nil {
			return err
		}
		constraint.Right, err = SimplifyExpr(constraint.Right)
		if err != nil {
			return err
		}
	}

	return nil
}

func SimplifyExpr(expr Expr) (Expr, error) {
	expr, err := Distribute(expr)
	if err != nil {
		return expr, err
	}
	expr, err = Flatten(expr)
	if err != nil {
		return expr, err
	}
	expr, err = CombineLikeTerms(expr)
	if err != nil {
		return expr, err
	}
	expr, err = ConstantFold(expr)
	if err != nil {
		return expr, err
	}
	return expr, nil
}

func Distribute(expr Expr) (Expr, error) {
	switch e := expr.(type) {
	case *BinaryExpr:
		left, err := Distribute(e.Left)
		if err != nil {
			return expr, err
		}
		right, err := Distribute(e.Right)
		if err != nil {
			return expr, err
		}

		if e.Operator.Type == TokenAsterisk {
			// a * (b + c) => a * b + a * c
			if rBin, ok := right.(*BinaryExpr); ok && (rBin.Operator.Type == TokenPlus || rBin.Operator.Type == TokenMinus) {
				return &BinaryExpr{
					Left: &BinaryExpr{
						Left:     left,
						Operator: e.Operator,
						Right:    rBin.Left,
					},
					Operator: rBin.Operator,
					Right: &BinaryExpr{
						Left:     left,
						Operator: e.Operator,
						Right:    rBin.Right,
					},
				}, nil
			}

			// (a + b) * c => a * c + b * c
			if lBin, ok := left.(*BinaryExpr); ok && (lBin.Operator.Type == TokenPlus || lBin.Operator.Type == TokenMinus) {
				return &BinaryExpr{
					Left: &BinaryExpr{
						Left:     lBin.Left,
						Operator: e.Operator,
						Right:    right,
					},
					Operator: lBin.Operator,
					Right: &BinaryExpr{
						Left:     lBin.Right,
						Operator: e.Operator,
						Right:    right,
					},
				}, nil
			}
		}

		if e.Operator.Type == TokenDivide {
			// (a + b) / c => a / c + b / c
			if lBin, ok := left.(*BinaryExpr); ok && (lBin.Operator.Type == TokenPlus || lBin.Operator.Type == TokenMinus) {
				return &BinaryExpr{
					Left: &BinaryExpr{
						Left:     lBin.Left,
						Operator: e.Operator,
						Right:    right,
					},
					Operator: lBin.Operator,
					Right: &BinaryExpr{
						Left:     lBin.Right,
						Operator: e.Operator,
						Right:    right,
					},
				}, nil
			}
		}

		// otherwise, can't distribute
		return &BinaryExpr{Left: left, Operator: e.Operator, Right: right}, nil
	case *UnaryExpr:
		expr, err := Distribute(e.Expr)
		if err != nil {
			return expr, err
		}
		return &UnaryExpr{Operator: e.Operator, Expr: expr}, nil
	case *NumberLiteral, *Variable:
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected expression type in Distribute: %T", expr)
	}
}

func Flatten(expr Expr) (Expr, error) {
	switch e := expr.(type) {
	case *Variable, *NumberLiteral:
	case *UnaryExpr:
	case *BinaryExpr:
	}
}

func CombineLikeTerms(expr Expr) (Expr, error) {
	return expr, nil
}

func ConstantFold(expr Expr) (Expr, error) {
	return expr, nil
}
