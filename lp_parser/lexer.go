package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Tokenize(reader io.Reader) ([]Token, error) {
	dfa := NewDFA()
	var tokens []Token

	scanner := bufio.NewScanner(reader)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		words := strings.Fields(line)
		for _, word := range words {
			wordRunes := []rune(word)
			for len(wordRunes) > 0 {
				currentToken, lettersRead, err := dfa.Run(wordRunes, lineNum)

				if err != nil {
					return tokens, err
				}

				tokens = append(tokens, currentToken)
				wordRunes = wordRunes[lettersRead:]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return tokens, fmt.Errorf("failed to read file: %w", err)
	}

	return tokens, nil
}
