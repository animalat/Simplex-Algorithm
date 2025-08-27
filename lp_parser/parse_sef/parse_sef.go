package parse_sef

import (
	"fmt"
	"strings"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
	"github.com/animalat/Simplex-Algorithm/lp_parser/semantics"
	"github.com/animalat/Simplex-Algorithm/lp_parser/simplify"
)

// Combines everything else and returns a parsed, simplified program.
// Note that converting the objective function from MIN to MAX is not a concern of this function.
func ParseSEF(progStr string) (*parser.Program, map[string]int, error) {
	tokens, err := lexer.Tokenize(strings.NewReader(progStr))
	if err != nil {
		return nil, nil, fmt.Errorf("error tokenizing: %v", err)
	}

	parseProg := parser.ConstructParser(tokens)
	prog, err := parseProg.ParseProgram()
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing: %v", err)
	}

	err = simplify.SimplifyProgram(prog)
	if err != nil {
		return nil, nil, fmt.Errorf("error simplifying expression: %v", err)
	}

	idTable, err := semantics.SemanticCheck(prog)
	if err != nil {
		return nil, nil, fmt.Errorf("semantic check failed: %v", err)
	}

	return prog, idTable, nil
}
