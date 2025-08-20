package parser

import (
	"strings"
	"testing"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
)

func tokens(ts ...lexer.Token) []lexer.Token {
	return append(ts, lexer.Token{Type: lexer.TokenEOF, Line: 999})
}

func TestParseProgram_ValidSingleDeclAndObjective(t *testing.T) {
	toks := tokens(
		lexer.Token{Type: lexer.TokenLet, Value: "let", Line: 1},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 1},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 1},

		lexer.Token{Type: lexer.TokenMax, Line: 2},
		lexer.Token{Type: lexer.TokenNumber, Value: "1", Line: 2},
		lexer.Token{Type: lexer.TokenPlus, Value: "+", Line: 2},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 2},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 2},
		lexer.Token{Type: lexer.TokenSubjectTo, Value: "s.t.", Line: 4},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 4},
		lexer.Token{Type: lexer.TokenLessEqual, Value: "<=", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "10", Line: 4},
		lexer.Token{Type: lexer.TokenSemiColon, Value: ";", Line: 4},
	)

	parser := &Parser{Tokens: toks}
	prog, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prog.Decls) != 1 {
		t.Errorf("expected 1 decl, got %d", len(prog.Decls))
	}
	if prog.Objective == nil {
		t.Errorf("expected objective, got nil")
	}

	if testing.Verbose() {
		PrintParse(prog)
	}
}

func TestParseProgram_MultipleDeclsAndConstraints(t *testing.T) {
	toks := tokens(
		// let x1;
		lexer.Token{Type: lexer.TokenLet, Line: 1},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 1},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 1},

		// let x2;
		lexer.Token{Type: lexer.TokenLet, Line: 2},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 2},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 2},

		// max x1 + x2;
		lexer.Token{Type: lexer.TokenMax, Line: 3},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 3},
		lexer.Token{Type: lexer.TokenPlus, Value: "+", Line: 3},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 3},
		lexer.Token{Type: lexer.TokenSemiColon, Value: ";", Line: 3},
		lexer.Token{Type: lexer.TokenSubjectTo, Value: "s.t.", Line: 4},

		// x1 + x2 <= 10;
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 4},
		lexer.Token{Type: lexer.TokenPlus, Value: "+", Line: 4},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 4},
		lexer.Token{Type: lexer.TokenLessEqual, Value: "<=", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "10", Line: 4},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 4},
	)

	parser := &Parser{Tokens: toks}
	prog, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prog.Decls) != 2 {
		t.Errorf("expected 2 decls, got %d", len(prog.Decls))
	}
	if len(prog.Constraints) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(prog.Constraints))
	}

	t.Logf("Parsed objective expression: %v", prog.Objective.Expr)

	if testing.Verbose() {
		PrintParse(prog)
	}
}

func TestParseProgram_MissingSemicolonFails(t *testing.T) {
	toks := tokens(
		// let x1    (missing semicolon)
		lexer.Token{Type: lexer.TokenLet, Line: 1},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 1},

		// max x1;
		lexer.Token{Type: lexer.TokenMax, Line: 2},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 2},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 2},
		lexer.Token{Type: lexer.TokenSubjectTo, Value: "s.t.", Line: 4},
	)

	parser := &Parser{Tokens: toks}
	_, err := parser.ParseProgram()
	if err == nil {
		t.Fatalf("expected error due to missing semicolon, got nil")
	}
}

func TestParseProgram_ComplexTest(t *testing.T) {
	toks := tokens(
		// let x1;
		lexer.Token{Type: lexer.TokenLet, Line: 1},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 1},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 1},

		// let x2;
		lexer.Token{Type: lexer.TokenLet, Line: 2},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 2},
		lexer.Token{Type: lexer.TokenSemiColon, Line: 2},

		// max x1 + (2 - 3 / 4) * x2;
		lexer.Token{Type: lexer.TokenMax, Line: 3},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 3},
		lexer.Token{Type: lexer.TokenPlus, Value: "+", Line: 3},
		lexer.Token{Type: lexer.TokenLParen, Value: "(", Line: 3},
		lexer.Token{Type: lexer.TokenNumber, Value: "2", Line: 3},
		lexer.Token{Type: lexer.TokenMinus, Value: "-", Line: 3},
		lexer.Token{Type: lexer.TokenNumber, Value: "3", Line: 3},
		lexer.Token{Type: lexer.TokenDivide, Value: "/", Line: 3},
		lexer.Token{Type: lexer.TokenNumber, Value: "4", Line: 3},
		lexer.Token{Type: lexer.TokenRParen, Value: ")", Line: 3},
		lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 3},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 3},
		lexer.Token{Type: lexer.TokenSemiColon, Value: ";", Line: 3},
		lexer.Token{Type: lexer.TokenSubjectTo, Value: "s.t.", Line: 4},

		// (2 - 3 / 4) * x2 >= 5;
		lexer.Token{Type: lexer.TokenLParen, Value: "(", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "2", Line: 4},
		lexer.Token{Type: lexer.TokenMinus, Value: "-", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "3", Line: 4},
		lexer.Token{Type: lexer.TokenDivide, Value: "/", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "4", Line: 4},
		lexer.Token{Type: lexer.TokenRParen, Value: ")", Line: 4},
		lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 4},
		lexer.Token{Type: lexer.TokenId, Value: "x2", Line: 4},
		lexer.Token{Type: lexer.TokenGreaterEqual, Value: ">=", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "5", Line: 4},
		lexer.Token{Type: lexer.TokenSemiColon, Value: ";", Line: 4},
		// (2 - 3 / 4) * x2 >= 5;
		lexer.Token{Type: lexer.TokenLParen, Value: "(", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "2", Line: 4},
		lexer.Token{Type: lexer.TokenMinus, Value: "-", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "3", Line: 4},
		lexer.Token{Type: lexer.TokenDivide, Value: "/", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "4", Line: 4},
		lexer.Token{Type: lexer.TokenRParen, Value: ")", Line: 4},
		lexer.Token{Type: lexer.TokenAsterisk, Value: "*", Line: 4},
		lexer.Token{Type: lexer.TokenId, Value: "x1", Line: 4},
		lexer.Token{Type: lexer.TokenGreaterEqual, Value: ">=", Line: 4},
		lexer.Token{Type: lexer.TokenNumber, Value: "5", Line: 4},
		lexer.Token{Type: lexer.TokenSemiColon, Value: ";", Line: 4},
	)

	parser := &Parser{Tokens: toks}
	prog, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prog.Decls) != 2 {
		t.Errorf("expected 2 decls, got %d", len(prog.Decls))
	}
	if len(prog.Constraints) != 2 {
		t.Errorf("expected 2 constraints, got %d", len(prog.Constraints))
	}

	t.Logf("Parsed objective expression: %v", prog.Objective.Expr)

	if testing.Verbose() {
		PrintParse(prog)
	}
}

func TestParseProgram_DefaultTest(t *testing.T) {
	tokens, err := lexer.Tokenize(strings.NewReader("let x1; let x2; let x3; max x1 + x2 + 3; s.t. x1 + x2 <= 3; x1 + x2 + 3 * x3 >= 5;"))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	parser := &Parser{Tokens: tokens}
	prog, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(prog.Decls) != 3 {
		t.Errorf("expected 3 decls, got %d", len(prog.Decls))
	}

	if len(prog.Constraints) != 2 {
		t.Errorf("expected 2 constraints, got %d", len(prog.Constraints))
	}

	t.Logf("Parsed objective expression: %v", prog.Objective.Expr)

	if testing.Verbose() {
		PrintParse(prog)
	}
}
