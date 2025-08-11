package parser

import (
	"strings"
	"testing"
)

func tokens(ts ...Token) []Token {
	return append(ts, Token{Type: TokenEOF, Line: 999})
}

func TestParseProgram_ValidSingleDeclAndObjective(t *testing.T) {
	toks := tokens(
		Token{Type: TokenLet, Value: "let", Line: 1},
		Token{Type: TokenId, Value: "x1", Line: 1},
		Token{Type: TokenSemiColon, Line: 1},

		Token{Type: TokenMax, Line: 2},
		Token{Type: TokenNumber, Value: "1", Line: 2},
		Token{Type: TokenPlus, Value: "+", Line: 2},
		Token{Type: TokenId, Value: "x1", Line: 2},
		Token{Type: TokenSemiColon, Line: 2},
		Token{Type: TokenSubjectTo, Value: "s.t.", Line: 4},
		Token{Type: TokenId, Value: "x1", Line: 4},
		Token{Type: TokenLessEqual, Value: "<=", Line: 4},
		Token{Type: TokenNumber, Value: "10", Line: 4},
		Token{Type: TokenSemiColon, Value: ";", Line: 4},
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
		Token{Type: TokenLet, Line: 1},
		Token{Type: TokenId, Value: "x1", Line: 1},
		Token{Type: TokenSemiColon, Line: 1},

		// let x2;
		Token{Type: TokenLet, Line: 2},
		Token{Type: TokenId, Value: "x2", Line: 2},
		Token{Type: TokenSemiColon, Line: 2},

		// max x1 + x2;
		Token{Type: TokenMax, Line: 3},
		Token{Type: TokenId, Value: "x1", Line: 3},
		Token{Type: TokenPlus, Value: "+", Line: 3},
		Token{Type: TokenId, Value: "x2", Line: 3},
		Token{Type: TokenSemiColon, Value: ";", Line: 3},
		Token{Type: TokenSubjectTo, Value: "s.t.", Line: 4},

		// x1 + x2 <= 10;
		Token{Type: TokenId, Value: "x1", Line: 4},
		Token{Type: TokenPlus, Value: "+", Line: 4},
		Token{Type: TokenId, Value: "x2", Line: 4},
		Token{Type: TokenLessEqual, Value: "<=", Line: 4},
		Token{Type: TokenNumber, Value: "10", Line: 4},
		Token{Type: TokenSemiColon, Line: 4},
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
		Token{Type: TokenLet, Line: 1},
		Token{Type: TokenId, Value: "x1", Line: 1},

		// max x1;
		Token{Type: TokenMax, Line: 2},
		Token{Type: TokenId, Value: "x1", Line: 2},
		Token{Type: TokenSemiColon, Line: 2},
		Token{Type: TokenSubjectTo, Value: "s.t.", Line: 4},
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
		Token{Type: TokenLet, Line: 1},
		Token{Type: TokenId, Value: "x1", Line: 1},
		Token{Type: TokenSemiColon, Line: 1},

		// let x2;
		Token{Type: TokenLet, Line: 2},
		Token{Type: TokenId, Value: "x2", Line: 2},
		Token{Type: TokenSemiColon, Line: 2},

		// max x1 + (2 - 3 / 4) * x2;
		Token{Type: TokenMax, Line: 3},
		Token{Type: TokenId, Value: "x1", Line: 3},
		Token{Type: TokenPlus, Value: "+", Line: 3},
		Token{Type: TokenLParen, Value: "(", Line: 3},
		Token{Type: TokenNumber, Value: "2", Line: 3},
		Token{Type: TokenMinus, Value: "-", Line: 3},
		Token{Type: TokenNumber, Value: "3", Line: 3},
		Token{Type: TokenDivide, Value: "/", Line: 3},
		Token{Type: TokenNumber, Value: "4", Line: 3},
		Token{Type: TokenRParen, Value: ")", Line: 3},
		Token{Type: TokenAsterisk, Value: "*", Line: 3},
		Token{Type: TokenId, Value: "x2", Line: 3},
		Token{Type: TokenSemiColon, Value: ";", Line: 3},
		Token{Type: TokenSubjectTo, Value: "s.t.", Line: 4},

		// (2 - 3 / 4) * x2 >= 5;
		Token{Type: TokenLParen, Value: "(", Line: 4},
		Token{Type: TokenNumber, Value: "2", Line: 4},
		Token{Type: TokenMinus, Value: "-", Line: 4},
		Token{Type: TokenNumber, Value: "3", Line: 4},
		Token{Type: TokenDivide, Value: "/", Line: 4},
		Token{Type: TokenNumber, Value: "4", Line: 4},
		Token{Type: TokenRParen, Value: ")", Line: 4},
		Token{Type: TokenAsterisk, Value: "*", Line: 4},
		Token{Type: TokenId, Value: "x2", Line: 4},
		Token{Type: TokenGreaterEqual, Value: ">=", Line: 4},
		Token{Type: TokenNumber, Value: "5", Line: 4},
		Token{Type: TokenSemiColon, Value: ";", Line: 4},
		// (2 - 3 / 4) * x2 >= 5;
		Token{Type: TokenLParen, Value: "(", Line: 4},
		Token{Type: TokenNumber, Value: "2", Line: 4},
		Token{Type: TokenMinus, Value: "-", Line: 4},
		Token{Type: TokenNumber, Value: "3", Line: 4},
		Token{Type: TokenDivide, Value: "/", Line: 4},
		Token{Type: TokenNumber, Value: "4", Line: 4},
		Token{Type: TokenRParen, Value: ")", Line: 4},
		Token{Type: TokenAsterisk, Value: "*", Line: 4},
		Token{Type: TokenId, Value: "x1", Line: 4},
		Token{Type: TokenGreaterEqual, Value: ">=", Line: 4},
		Token{Type: TokenNumber, Value: "5", Line: 4},
		Token{Type: TokenSemiColon, Value: ";", Line: 4},
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
	tokens, err := Tokenize(strings.NewReader("let x1; let x2; let x3; max x1 + x2 + 3; s.t. x1 + x2 <= 3; x1 + x2 + 3 * x3 >= 5;"))
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