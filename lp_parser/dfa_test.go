package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestDFAConstruction(t *testing.T) {
	dfa := NewDFA()

	if len(dfa.Transitions) == 0 {
		t.Error("expected transitions but got none")
	}

	if testing.Verbose() {
		for key, next := range dfa.Transitions {
			fmt.Printf("from %q with %q -> %q\n", key.State, key.Input, next)
		}
	}
}

func assertTokens(t *testing.T, input string, expected []Token) {
	t.Helper()

	tokens, err := Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	if len(tokens) != len(expected) {
		t.Fatalf("Tokenize() returned %d tokens, want %d", len(tokens), len(expected))
	}

	for i, got := range tokens {
		want := expected[i]
		if got.Value != want.Value || got.Type != want.Type || got.Line != want.Line {
			t.Errorf("Token %d = {Type: %q, Value: %q, Line: %d}; want {Type: %q, Value: %q, Line: %d}",
				i, got.Type, got.Value, got.Line, want.Type, want.Value, want.Line)
		}
	}

	if testing.Verbose() {
		for _, token := range tokens {
			fmt.Printf("token type: %q, token value: %q, line number: %d\n", token.Type, token.Value, token.Line)
		}
	}
}

func TestTokenizeSimple(t *testing.T) {
	input := "x + y\n42<=100"
	expected := []Token{
		{Type: TokenId, Value: "x", Line: 1},
		{Type: TokenPlus, Value: "+", Line: 1},
		{Type: TokenId, Value: "y", Line: 1},
		{Type: TokenNumber, Value: "42", Line: 2},
		{Type: TokenLessEqual, Value: "<=", Line: 2},
		{Type: TokenNumber, Value: "100", Line: 2},
	}
	assertTokens(t, input, expected)
}

func TestTokenizeLP(t *testing.T) {
	input := "let x1;let x2;\n max x1 + x2;\ns.t. x1+x2<=5;"
	expected := []Token{
		{Type: TokenLet, Value: "let", Line: 1},
		{Type: TokenId, Value: "x1", Line: 1},
		{Type: TokenSemiColon, Value: ";", Line: 1},
		{Type: TokenLet, Value: "let", Line: 1},
		{Type: TokenId, Value: "x2", Line: 1},
		{Type: TokenSemiColon, Value: ";", Line: 1},
		{Type: TokenMax, Value: "max", Line: 2},
		{Type: TokenId, Value: "x1", Line: 2},
		{Type: TokenPlus, Value: "+", Line: 2},
		{Type: TokenId, Value: "x2", Line: 2},
		{Type: TokenSemiColon, Value: ";", Line: 2},
		{Type: TokenSubjectTo, Value: "s.t.", Line: 3},
		{Type: TokenId, Value: "x1", Line: 3},
		{Type: TokenPlus, Value: "+", Line: 3},
		{Type: TokenId, Value: "x2", Line: 3},
		{Type: TokenLessEqual, Value: "<=", Line: 3},
		{Type: TokenNumber, Value: "5", Line: 3},
		{Type: TokenSemiColon, Value: ";", Line: 3},
	}
	assertTokens(t, input, expected)
}
