package simplify

import (
	"fmt"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

func SimplifyProgram(p *parser.Program) error {
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

func SimplifyExpr(expr parser.Expr) (parser.Expr, error) {
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

func Distribute(expr parser.Expr) (parser.Expr, error) {
	switch e := expr.(type) {
	case *parser.BinaryExpr:
		left, err := Distribute(e.Left)
		if err != nil {
			return expr, err
		}
		right, err := Distribute(e.Right)
		if err != nil {
			return expr, err
		}

		if e.Operator.Type == lexer.TokenAsterisk {
			rBin, okr := right.(*parser.BinaryExpr)
			lBin, okl := left.(*parser.BinaryExpr)
			rightCondition := okr && (rBin.Operator.Type == lexer.TokenPlus || rBin.Operator.Type == lexer.TokenMinus)
			leftCondition := okl && (lBin.Operator.Type == lexer.TokenPlus || lBin.Operator.Type == lexer.TokenMinus)

			// (a + b) * (c + d) => a * c + a * d + b * c + b * d
			if rightCondition && leftCondition {
				/*
					TODO: calculate expressions
					ac := &BinaryExpr{Left: lBin.Left, Operator: Right:}

					return &BinaryExpr{
						Left: &BinaryExpr{

						},
						Operator:
						Right: &BinaryExpr{
						}
					}
				*/
			}

			// a * (b + c) => a * b + a * c
			if rightCondition {
				return &parser.BinaryExpr{
					Left: &parser.BinaryExpr{
						Left:     left,
						Operator: e.Operator,
						Right:    rBin.Left,
					},
					Operator: rBin.Operator,
					Right: &parser.BinaryExpr{
						Left:     left,
						Operator: e.Operator,
						Right:    rBin.Right,
					},
				}, nil
			}

			// (a + b) * c => a * c + b * c
			if leftCondition {
				return &parser.BinaryExpr{
					Left: &parser.BinaryExpr{
						Left:     lBin.Left,
						Operator: e.Operator,
						Right:    right,
					},
					Operator: lBin.Operator,
					Right: &parser.BinaryExpr{
						Left:     lBin.Right,
						Operator: e.Operator,
						Right:    right,
					},
				}, nil
			}
		}

		if e.Operator.Type == lexer.TokenDivide {
			// (a + b) / c => a / c + b / c
			if lBin, ok := left.(*parser.BinaryExpr); ok && (lBin.Operator.Type == lexer.TokenPlus || lBin.Operator.Type == lexer.TokenMinus) {
				return &parser.BinaryExpr{
					Left: &parser.BinaryExpr{
						Left:     lBin.Left,
						Operator: e.Operator,
						Right:    right,
					},
					Operator: lBin.Operator,
					Right: &parser.BinaryExpr{
						Left:     lBin.Right,
						Operator: e.Operator,
						Right:    right,
					},
				}, nil
			}
		}

		// otherwise, can't distribute
		return &parser.BinaryExpr{Left: left, Operator: e.Operator, Right: right}, nil
	case *parser.UnaryExpr:
		expr, err := Distribute(e.Expr)
		if err != nil {
			return expr, err
		}
		return &parser.UnaryExpr{Operator: e.Operator, Expr: expr}, nil
	case *parser.NumberLiteral, *parser.Variable:
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected expression type in Distribute: %T", expr)
	}
}

func Flatten(expr parser.Expr) (parser.Expr, error) {
	/*
		switch e := expr.(type) {
		case *Variable, *NumberLiteral:
			return e, nil
		case *UnaryExpr:
			inner, err := Flatten(e.expr)
			if err != nil {
				return nil, err
			}

		case *BinaryExpr:
		}
	*/
	return expr, nil
}

func CombineLikeTerms(expr parser.Expr) (parser.Expr, error) {
	return expr, nil
}

func ConstantFold(expr parser.Expr) (parser.Expr, error) {
	return expr, nil
}
