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
		if got.Value != want.Value || got.Type != want.Type {
			t.Errorf("Token %d = {Type: %q, Value: %q}; want {Type: %q, Value: %q}",
				i, got.Type, got.Value, want.Type, want.Value)
		}
	}

	if testing.Verbose() {
		for _, token := range tokens {
			fmt.Printf("token type: %q, token value: %q\n", token.Type, token.Value)
		}
	}
}

func TestTokenizeSimple(t *testing.T) {
	input := "x + y\n42<=100"
	expected := []Token{
		{Type: TokenId, Value: "x"},
		{Type: TokenPlus, Value: "+"},
		{Type: TokenId, Value: "y"},
		{Type: TokenNumber, Value: "42"},
		{Type: TokenLessEqual, Value: "<="},
		{Type: TokenNumber, Value: "100"},
	}
	assertTokens(t, input, expected)
}

func TestTokenizeLP(t *testing.T) {
	input := "let x1;let x2; max x1 + x2;\ns.t. x1+x2<=5;"
	expected := []Token{
		{Type: TokenLet, Value: "let"},
		{Type: TokenId, Value: "x1"},
		{Type: TokenSemiColon, Value: ";"},
		{Type: TokenLet, Value: "let"},
		{Type: TokenId, Value: "x2"},
		{Type: TokenSemiColon, Value: ";"},
		{Type: TokenMax, Value: "max"},
		{Type: TokenId, Value: "x1"},
		{Type: TokenPlus, Value: "+"},
		{Type: TokenId, Value: "x2"},
		{Type: TokenSemiColon, Value: ";"},
		{Type: TokenSubjectTo, Value: "s.t."},
		{Type: TokenId, Value: "x1"},
		{Type: TokenPlus, Value: "+"},
		{Type: TokenId, Value: "x2"},
		{Type: TokenLessEqual, Value: "<="},
		{Type: TokenNumber, Value: "5"},
		{Type: TokenSemiColon, Value: ";"},
	}
	assertTokens(t, input, expected)
}
