package parser

import "fmt"

type Program struct {
	Decls       []*Decl
	Objective   *Objective
	Constraints []*Constraint
}

type Decl struct {
	ID Token
}

type Objective struct {
	IsMax bool
	Expr  Expr
}

type Constraint struct {
	Left     Expr
	Operator Token
	Right    Expr
	Line     int
}

type Expr interface {
	exprNode()
}

type UnaryExpr struct {
	Operator Token
	Expr     Expr
	Line     int
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
	Line     int
}

type NumberLiteral struct {
	Value float64
	Line  int
}

type Variable struct {
	ID Token
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func (v *Variable) String() string {
	return v.ID.Value
}

func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator.Value, u.Expr)
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator.Value, b.Right)
}

type Parser struct {
	Tokens []Token
	Pos    int
}
