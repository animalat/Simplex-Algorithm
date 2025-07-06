package parser

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

const (
	TokenLet          TokenType = "LET"
	TokenSubjectTo    TokenType = "S.T."
	TokenMin          TokenType = "MIN"
	TokenMax          TokenType = "MAX"
	TokenId           TokenType = "ID"
	TokenNumber       TokenType = "NUMBER"
	TokenDecimal      TokenType = "DECIMAL"
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

const StartingState string = "start"

type TransitionKey struct {
	State string
	Input rune
}

type DFA struct {
	AlphabetSymbols map[rune]bool
	States          map[string]bool
	FinalStates     map[string]TokenType
	Transitions     map[TransitionKey]string
	StartState      string
}
