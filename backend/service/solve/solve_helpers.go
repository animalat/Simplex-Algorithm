package solve

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strings"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

type SimplexProgramArrays struct {
	objective        []float64
	objectiveConst   float64
	constraintsLHS   [][]float64
	constraintsRHS   []float64
	constraintsSlack []float64
	numSlack         int
}

type SimplexProgramStrings struct {
	objectiveOutput      string
	objectiveConstOutput string
	constraintsOutputLHS []string
	constraintsOutputRHS string
}

type SimplexResult struct {
	Solution    []float64 `json:"solution"`
	ResultType  string    `json:"resultType"`
	Certificate []float64 `json:"certificate"`
}

func insertElem(objective []float64, objectiveConst *float64, e parser.Expr, idTable map[string]int, isObjective bool) error {
	switch expr := e.(type) {
	case *parser.NumberLiteral:
		if !isObjective {
			return fmt.Errorf("unexpected NumberLiteral in non-objective Expr: %v", expr)
		}

		*objectiveConst = expr.Value
	case *parser.Variable:
		idx, ok := idTable[expr.ID.Value]
		if !ok {
			return fmt.Errorf("undeclared variable: %v", expr)
		}
		// 1 because just a variable would be "1"
		objective[idx] = 1
	case *parser.BinaryExpr:
		if expr.Operator.Type != lexer.TokenAsterisk {
			return fmt.Errorf("unexpected term: %v", expr)
		}

		// this is the last expr
		v, ok := expr.Right.(*parser.Variable)
		if !ok {
			return fmt.Errorf("unexpected type: %T", expr.Right)
		}

		nl, ok := expr.Left.(*parser.NumberLiteral)
		if !ok {
			return fmt.Errorf("unexpected type: %T", expr.Left)
		}

		idx, ok := idTable[v.ID.Value]
		if !ok {
			return fmt.Errorf("undeclared variable: %v", expr)
		}
		objective[idx] = nl.Value
	default:
		return fmt.Errorf("unexpected type: %T", expr)
	}

	return nil
}

// Requires: semantics must be run on the program before passing (this is not a semantics check)
func getExprArr(e parser.Expr, idTable map[string]int, isObjective bool) (float64, []float64, error) {
	objective := make([]float64, len(idTable))
	objectiveConst := 0.0

	curExpr := e

	for {
		switch expr := curExpr.(type) {
		case *parser.NumberLiteral, *parser.Variable:
			if err := insertElem(objective, &objectiveConst, expr, idTable, isObjective); err != nil {
				return 0, nil, err
			}
			return objectiveConst, objective, nil
		case *parser.BinaryExpr:
			if expr.Operator.Type == lexer.TokenAsterisk {
				if err := insertElem(objective, &objectiveConst, expr, idTable, isObjective); err != nil {
					return 0, nil, err
				}
				return objectiveConst, objective, nil
			}

			// otherwise, we grab the expr, and then check left
			if err := insertElem(objective, &objectiveConst, expr.Right, idTable, isObjective); err != nil {
				return 0, nil, err
			}

			curExpr = expr.Left
		default:
			return 0, nil, fmt.Errorf("unexpected type: %T", expr)
		}
	}
}

func allFreeVariables(idTable map[string]int, alreadyPositive map[string]struct{}) map[string]struct{} {
	toPositive := make(map[string]struct{})
	for key := range idTable {
		if _, ok := alreadyPositive[key]; !ok {
			toPositive[key] = struct{}{}
		}
	}

	return toPositive
}

func rowInput(row []float64, toPositive map[string]struct{}, idTableInverse map[int]string) (string, error) {
	output := ""
	for i := 0; i < len(row); i++ {
		variable, ok := idTableInverse[i]
		if !ok {
			return "", fmt.Errorf("invalid variable found at index %d", i)
		}

		output += ftos(row[i]) + " "

		if _, ok = toPositive[variable]; !ok {
			// we need to add on the subtract too
			output += ftos(-row[i])
		}
	}
	return output, nil
}

func getSlackOutput(numSlackAdded *int, constraintSlack float64, numSlack int) (string, error) {
	if math.Abs(constraintSlack) < EPSILON {
		// no slack variable
		return strings.Repeat("0 ", numSlack), nil
	}

	if *numSlackAdded >= numSlack {
		return "", fmt.Errorf("extra unexpected slack variable: %.2f", constraintSlack)
	}

	rightZeroAmount := numSlack - 1 - (*numSlackAdded)
	*numSlackAdded++
	return strings.Repeat("0 ", *numSlackAdded) + ftos(constraintSlack) + " " + strings.Repeat("0 ", rightZeroAmount), nil
}

func simplexInput(progArrays SimplexProgramArrays, toPositive map[string]struct{}, idTable map[string]int) (SimplexProgramStrings, error) {
	objective := progArrays.objective
	objectiveConst := progArrays.objectiveConst
	constraintsLHS := progArrays.constraintsLHS
	constraintsRHS := progArrays.constraintsRHS
	constraintsSlack := progArrays.constraintsSlack
	numSlack := progArrays.numSlack

	idTableInverse := getTableInverse(idTable)
	objectiveOutput, err := rowInput(objective, toPositive, idTableInverse)
	if err != nil {
		return SimplexProgramStrings{}, err
	}
	objectiveOutput += strings.Repeat("0 ", numSlack)
	objectiveConstOutput := ftos(objectiveConst)

	numSlackAdded := 0
	constraintsOutputLHS := make([]string, 0, len(constraintsLHS))
	for i := range constraintsLHS {
		curRowOutput, err := rowInput(constraintsLHS[i], toPositive, idTableInverse)
		if err != nil {
			return SimplexProgramStrings{}, err
		}
		endingStr, err := getSlackOutput(&numSlackAdded, constraintsSlack[i], numSlack)
		if err != nil {
			return SimplexProgramStrings{}, err
		}
		curRowOutput += endingStr + "\n"
		constraintsOutputLHS = append(constraintsOutputLHS, curRowOutput)
	}

	constraintsOutputRHS := ""
	for _, val := range constraintsRHS {
		constraintsOutputRHS += ftos(val) + " "
	}

	return SimplexProgramStrings{
		objectiveOutput:      objectiveOutput,
		objectiveConstOutput: objectiveConstOutput,
		constraintsOutputLHS: constraintsOutputLHS,
		constraintsOutputRHS: constraintsOutputRHS,
	}, nil
}

func callSimplex(progStrings SimplexProgramStrings, rowSize string, colSize string) (string, error) {
	objectiveOutput := progStrings.objectiveOutput
	objectiveConstOutput := progStrings.objectiveConstOutput
	constraintsOutputLHS := progStrings.constraintsOutputLHS
	constraintsOutputRHS := progStrings.constraintsOutputRHS

	cmd := exec.Command("../../../simplex_core/simplex_solver")

	input := ""
	// main matrix, LHS constraints (A)
	input += rowSize + "\n" + colSize + "\n"
	for _, curStr := range constraintsOutputLHS {
		input += curStr + "\n"
	}
	// RHS constraints (B)
	input += rowSize + "\n1\n"
	input += constraintsOutputRHS + "\n"
	// objective (C)
	input += "1\n" + colSize + "\n"
	input += objectiveOutput + "\n"
	// objective constant (z)
	input += objectiveConstOutput + "\n"

	cmd.Stdin = bytes.NewBufferString(input)

	output, err := cmd.CombinedOutput()
	log.Println(input)
	if err != nil {
		return "", fmt.Errorf("error running solver: %v\noutput: %s", err, string(output))
	}

	return string(output), nil
}

func parseResult(output string) (SimplexResult, error) {
	r := strings.NewReader(output)

	var solution []float64
	for r.Len() > 0 {
		var x float64
		_, err := fmt.Fscan(r, &x)
		if err != nil {
			break
		}

		solution = append(solution, x)
	}

	var resultType string
	_, err := fmt.Fscan(r, &resultType)
	if err != nil {
		return SimplexResult{}, fmt.Errorf("no resultType found")
	}

	var certificate []float64
	for r.Len() > 0 {
		var x float64
		_, err := fmt.Fscan(r, &x)
		if err != nil {
			break
		}

		certificate = append(certificate, x)
	}

	log.Printf("hello, %s", resultType)

	return SimplexResult{
		Solution:    solution,
		ResultType:  resultType,
		Certificate: certificate,
	}, nil
}
