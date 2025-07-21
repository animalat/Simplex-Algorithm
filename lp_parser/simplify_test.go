package parser

import (
	"strings"
	"testing"
)

func TestSimplify_Distribute(t *testing.T) {
	input := "let x1; let x2; max (3 + 2) * (x1 + x2); s.t. (3 + x1) / x2 <= 5 * (3 + 1)"
	tokens, err := Tokenize(strings.NewReader(input))
	if err != nil {
		t.Fatalf("TestSimplify_Distribute fatal error tokenizing: %v", err)
	}

    parser := &Parser{Tokens: tokens}
    prog, err := parser.ParseProgram()
	if err != nil {
        t.Fatalf("TestSimplify_Distribute fatal error parsing: %v", err)
    }
    
}
