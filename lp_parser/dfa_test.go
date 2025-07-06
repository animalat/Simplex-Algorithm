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

func TestTokenize(t *testing.T) {
	input := "x + y\n42<=100"
	expectedTokens := []string{"x", "+", "y", "42", "<=", "100"}

	tokens, err := Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Tokenize() returned %d tokens, want %d", len(tokens), len(expectedTokens))
	}

	for i, token := range tokens {
		if token.Value != expectedTokens[i] {
			t.Errorf("Token %d = %q; want %q", i, token.Value, expectedTokens[i])
		}
	}

	if testing.Verbose() {
		for _, token := range tokens {
			fmt.Printf("token type: %q, token value: %q\n", token.Type, token.Value)
		}
	}
}

func TestTokenizeLP(t *testing.T) {
	input := "let x1;let x2; max x1 + x2;\ns.t. x1+x2<=5;"
	expectedTokens := []string{"let", "x1", ";", "let", "x2", ";", "max", "x1",
		"+", "x2", ";", "s.t.", "x1", "+", "x2", "<=", "5", ";"}

	tokens, err := Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	if len(tokens) != len(expectedTokens) {
		t.Fatalf("Tokenize() returned %d tokens, want %d", len(tokens), len(expectedTokens))
	}

	for i, token := range tokens {
		if token.Value != expectedTokens[i] {
			t.Errorf("Token %d = %q; want %q", i, token.Value, expectedTokens[i])
		}
	}

	if testing.Verbose() {
		for _, token := range tokens {
			fmt.Printf("token type: %q, token value: %q\n", token.Type, token.Value)
		}
	}
}
