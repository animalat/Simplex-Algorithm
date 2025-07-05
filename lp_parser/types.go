package parser

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

const (
	TokenLet          TokenType = "LET"
	TokenSubjectTo    TokenType = "ST"
	TokenId           TokenType = "ID"
	TokenNumber       TokenType = "NUMBER"
	TokenSemiColon    TokenType = "SEMICOLON"
	TokenEqual        TokenType = "EQ"
	TokenLessEqual    TokenType = "LEQ"
	TokenGreaterEqual TokenType = "GEQ"
	TokenPlus         TokenType = "PLUS"
	TokenMinus        TokenType = "MINUS"
	TokenAsterisk     TokenType = "ASTERISK"
	TokenDivide       TokenType = "SLASH"
	TokenLParen       TokenType = "LPAREN"
	TokenRParen       TokenType = "RPAREN"
	TokenEOF          TokenType = "EOF"
)
