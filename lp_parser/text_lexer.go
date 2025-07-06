package parser

import (
	"io"
)

func Tokenize(reader io.Reader) ([]Token, error) {
	// scanner := bufio.NewScanner(reader)
	// lineNum := 0
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	lineNum++

	// 	words := strings.Fields(line)
	// 	for _, word := range words {
	// 		//somethinghere
	// 	}
	// }

	// if err := scanner.Err(); err != nil {
	// 	return nil, fmt.Errorf("failed to read file: %w", err)
	// }

	return nil, nil
}
