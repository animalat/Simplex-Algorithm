package simplify

import (
	"fmt"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

const negativeMultiplicative = -1
const defaultMultiplicative = 1

type isConstant struct {
	value      float64
	isConstant bool
}

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
	expr, err := DistributeFold(expr, defaultMultiplicative)
	if err != nil {
		return expr, err
	}
	expr, err = CombineLikeTerms(expr)
	if err != nil {
		return expr, err
	}
	return expr, nil
}

func doOperation(a float64, b float64, tt lexer.TokenType) float64 {
	var res float64
	switch tt {
	case lexer.TokenPlus:
		res = a + b
	case lexer.TokenMinus:
		res = a - b
	case lexer.TokenAsterisk:
		res = a * b
	case lexer.TokenDivide:
		res = a / b
	}
	return res
}

func exprIsConstant(expr parser.Expr) (isConstant, error) {
	switch e := expr.(type) {
	case *parser.BinaryExpr:
		leftIsConstant, err := exprIsConstant(e.Left)
		if err != nil {
			return isConstant{isConstant: false}, err
		}
		rightIsConstant, err := exprIsConstant(e.Right)
		if err != nil {
			return isConstant{isConstant: false}, err
		}
		if !leftIsConstant.isConstant || !rightIsConstant.isConstant {
			return isConstant{isConstant: false}, nil
		}

		newConst := doOperation(leftIsConstant.value, rightIsConstant.value, e.Operator.Type)

		return isConstant{value: newConst, isConstant: true}, nil
	case *parser.UnaryExpr:
		innerIsConstant, err := exprIsConstant(e.Expr)
		if err != nil {
			return isConstant{isConstant: false}, err
		}
		var multiplicative float64
		if e.Operator.Type == lexer.TokenMinus {
			multiplicative = negativeMultiplicative
		} else if e.Operator.Type == lexer.TokenPlus {
			multiplicative = defaultMultiplicative
		} else {
			return isConstant{isConstant: false}, fmt.Errorf("invalid UnaryExpr operator %v: %v", e.Operator.Value, e)
		}

		if innerIsConstant.isConstant {
			return isConstant{value: innerIsConstant.value * multiplicative, isConstant: true}, nil
		} else {
			return isConstant{isConstant: false}, nil
		}
	case *parser.NumberLiteral:
		return isConstant{value: e.Value, isConstant: true}, nil
	default:
		return isConstant{isConstant: false}, nil
	}
}

func DistributeFold(expr parser.Expr, multiplicative float64) (parser.Expr, error) {
	switch e := expr.(type) {
	case *parser.BinaryExpr:
		switch e.Operator.Type {
		case lexer.TokenPlus, lexer.TokenMinus:
			newLeft, err := DistributeFold(e.Left, multiplicative)
			if err != nil {
				return nil, err
			}

			var isNegativeLHS float64
			if e.Operator.Type == lexer.TokenMinus {
				isNegativeLHS = negativeMultiplicative
			} else {
				isNegativeLHS = defaultMultiplicative
			}

			newRight, err := DistributeFold(e.Right, isNegativeLHS*multiplicative)
			if err != nil {
				return nil, err
			}

			return &parser.BinaryExpr{
				Left:     newLeft,
				Operator: e.Operator,
				Right:    newRight,
				Line:     e.Line,
			}, nil
		case lexer.TokenAsterisk, lexer.TokenDivide:
			leftIsConstant, err := exprIsConstant(e.Left)
			if err != nil {
				return nil, err
			}
			rightIsConstant, err := exprIsConstant(e.Right)
			if err != nil {
				return nil, err
			}
			if leftIsConstant.isConstant && rightIsConstant.isConstant {
				return &parser.NumberLiteral{Value: doOperation(leftIsConstant.value, rightIsConstant.value, e.Operator.Type), Line: e.Line}, nil
			} else if !leftIsConstant.isConstant && !rightIsConstant.isConstant {
				return nil, fmt.Errorf("nonlinear expression (both sides): %v", e)
			} else if !leftIsConstant.isConstant && rightIsConstant.isConstant {
				return DistributeFold(e.Left, doOperation(multiplicative, rightIsConstant.value, e.Operator.Type))
			} else {
				// leftIsConstant.isConstant && !rightIsConstant.isConstant
				if e.Operator.Type == lexer.TokenAsterisk {
					return DistributeFold(e.Right, doOperation(leftIsConstant.value, multiplicative, e.Operator.Type))
				} else {
					// TokenDivide
					return nil, fmt.Errorf("nonlinear expression (RHS rational): %v", e)
				}
			}
		default:
			return nil, fmt.Errorf("invalid operator \"%v\": %v", e.Operator.Value, e)
		}
	case *parser.UnaryExpr:
		if e.Operator.Type == lexer.TokenMinus {
			multiplicative *= negativeMultiplicative
		}
		return DistributeFold(e.Expr, multiplicative)
	case *parser.Variable:
		return &parser.BinaryExpr{
			Left:     &parser.NumberLiteral{Value: multiplicative, Line: e.ID.Line},
			Operator: lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: e.ID.Line},
			Right:    e,
			Line:     e.ID.Line,
		}, nil
	case *parser.NumberLiteral:
		e.Value *= multiplicative
		return e, nil
	default:
		return nil, fmt.Errorf("invalid Expr type found: %T", e)
	}
}

func CombineLikeTerms(expr parser.Expr) (parser.Expr, error) {
	return expr, nil
}
