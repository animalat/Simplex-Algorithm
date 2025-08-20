package parser

import (
	"fmt"
	"strconv"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
)

func (u *UnaryExpr) exprNode()     {}
func (b *BinaryExpr) exprNode()    {}
func (n *NumberLiteral) exprNode() {}
func (v *Variable) exprNode()      {}

func (p *Parser) Peek() (lexer.Token, error) {
	if p.Pos >= len(p.Tokens) {
		// Return EOF token
		const endLine = -1
		return lexer.Token{Type: lexer.TokenEOF, Value: "", Line: endLine}, nil
	}
	return p.Tokens[p.Pos], nil
}

func (p *Parser) Advance() (lexer.Token, error) {
	token, err := p.Peek()
	if err != nil {
		return token, err
	}

	p.Pos++
	return token, nil
}

func (p *Parser) Expect(tt lexer.TokenType) (lexer.Token, error) {
	token, err := p.Advance()
	if err != nil {
		return token, err
	}

	if tt != token.Type {
		return token, fmt.Errorf("token type does not match at line %d: Expected %s but got %s", token.Line, tt, token.Type)
	}

	return token, nil
}

func (p *Parser) ParseDecl() (*Decl, error) {
	if _, err := p.Expect(lexer.TokenLet); err != nil {
		return nil, err
	}

	token, err := p.Expect(lexer.TokenId)
	if err != nil {
		return nil, err
	}

	if _, err = p.Expect(lexer.TokenSemiColon); err != nil {
		return nil, err
	}

	return &Decl{ID: token}, nil
}

func (p *Parser) ParseObjective() (*Objective, error) {
	token, err := p.Advance()
	if err != nil {
		return nil, err
	}

	if token.Type != lexer.TokenMax && token.Type != lexer.TokenMin {
		return nil, fmt.Errorf("token min or max not found at line %d", token.Line)
	}
	isMax := token.Type == lexer.TokenMax

	expr, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	token, err = p.Expect(lexer.TokenSemiColon)
	if err != nil {
		return nil, err
	}

	return &Objective{IsMax: isMax, Expr: expr}, nil
}

func (p *Parser) ParseConstraint() (*Constraint, error) {
	left, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	op, err := p.Advance()
	if err != nil {
		return nil, err
	}
	if op.Type != lexer.TokenLessEqual && op.Type != lexer.TokenEqual && op.Type != lexer.TokenGreaterEqual {
		return nil, fmt.Errorf("operator not found at line %d", op.Line)
	}

	right, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	if _, err = p.Expect(lexer.TokenSemiColon); err != nil {
		return nil, err
	}

	return &Constraint{Left: left, Operator: op, Right: right}, nil
}

func (p *Parser) ParseExpr() (Expr, error) {
	left, err := p.ParseTerm()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.Peek()
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.TokenPlus && token.Type != lexer.TokenMinus {
			break
		}

		op, err := p.Advance()
		if err != nil {
			return nil, err
		}

		right, err := p.ParseTerm()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}

	return left, nil
}

func (p *Parser) ParseTerm() (Expr, error) {
	left, err := p.ParseFactor()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.Peek()
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.TokenAsterisk && token.Type != lexer.TokenDivide {
			break
		}

		op, err := p.Advance()
		if err != nil {
			return nil, err
		}

		right, err := p.ParseFactor()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left, nil
}

func (p *Parser) ParseFactor() (Expr, error) {
	// unary check
	token, err := p.Peek()
	if err != nil {
		return nil, err
	}

	if token.Type == lexer.TokenMinus || token.Type == lexer.TokenPlus {
		op, err := p.Advance()
		if err != nil {
			return nil, err
		}

		expr, err := p.ParseFactor()
		if err != nil {
			return nil, err
		}

		return &UnaryExpr{Operator: op, Expr: expr}, nil
	}

	// remaining cases
	token, err = p.Advance()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.TokenNumber:
		const doubleSize = 64
		value, err := strconv.ParseFloat(token.Value, doubleSize)
		if err != nil {
			return nil, fmt.Errorf("invalid number token at line %d", token.Line)
		}
		return &NumberLiteral{Value: value, Line: token.Line}, nil
	case lexer.TokenId:
		return &Variable{ID: token}, nil
	case lexer.TokenLParen:
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}

		if _, err = p.Expect(lexer.TokenRParen); err != nil {
			return nil, err
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token with value %s at line %d", token.Value, token.Line)
	}
}

func (p *Parser) ParseProgram() (*Program, error) {
	var decls []*Decl
	for {
		token, err := p.Peek()
		if err != nil {
			return nil, err
		}
		if token.Type != lexer.TokenLet {
			break
		}

		decl, err := p.ParseDecl()
		if err != nil {
			return nil, err
		}

		decls = append(decls, decl)
	}

	objective, err := p.ParseObjective()
	if err != nil {
		return nil, err
	}

	if _, err := p.Expect(lexer.TokenSubjectTo); err != nil {
		return nil, err
	}

	var constraints []*Constraint
	for {
		token, err := p.Peek()
		if err != nil {
			return nil, err
		}
		if token.Type == lexer.TokenEOF {
			break
		}

		constraint, err := p.ParseConstraint()
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, constraint)
	}

	return &Program{Decls: decls, Objective: objective, Constraints: constraints}, nil
}

func PrintParse(p *Program) error {
	for _, decl := range p.Decls {
		fmt.Printf("let %s;\n", decl.ID.Value)
	}

	if p.Objective.IsMax {
		fmt.Print("max ")
	} else {
		fmt.Print("min ")
	}

	fmt.Printf("%s;\n", p.Objective.Expr)

	for _, constraint := range p.Constraints {
		fmt.Printf("%s %s %s;\n", constraint.Left, constraint.Operator.Value, constraint.Right)
	}

	return nil
}
