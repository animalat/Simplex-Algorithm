package solve

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"strings"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

// Array version of linear program (to be converted into input for the Simplex calculator)
type SimplexProgramArrays struct {
	objective        []float64
	objectiveConst   float64
	constraintsLHS   [][]float64
	constraintsRHS   []float64
	constraintsSlack []float64
	numSlack         int
}

// Input that is passed into the Simplex calculator
type SimplexProgramStrings struct {
	objectiveOutput      string
	objectiveConstOutput string
	constraintsOutputLHS []string
	constraintsOutputRHS string
}

// API output
type SimplexResult struct {
	Solution    []float64      `json:"solution"`
	ResultType  string         `json:"resultType"`
	Certificate []float64      `json:"certificate"`
	Mapping     map[int]string `json:"mapping"`
}

// Insert element from an Expr to an array version of that Expr
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

// This function takes an Expr and turns it into an array (and possibly an extra constant for the objective function)
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

// This determines all free variables given already positive variables
func allFreeVariables(idTable map[string]int, alreadyPositive map[string]struct{}) map[string]struct{} {
	toPositive := make(map[string]struct{})
	for key := range idTable {
		if _, ok := alreadyPositive[key]; !ok {
			toPositive[key] = struct{}{}
		}
	}

	return toPositive
}

// Uses x := a - b for a, b >= 0 to help turn the LP into standard equality form
func rowInput(row []float64, toPositive map[string]struct{}, idTableInverse map[int]string) (string, error) {
	output := ""
	for i := 0; i < len(row); i++ {
		variable, ok := idTableInverse[i]
		if !ok {
			return "", fmt.Errorf("invalid variable found at index %d", i)
		}

		output += ftos(row[i]) + " "

		if _, ok = toPositive[variable]; ok {
			// we need to add on the subtract too
			output += ftos(-row[i]) + " "
		}
	}
	return output, nil
}

// Add on slack variables to the linear program (help turn it into standard equality form)
func getSlackOutput(numSlackAdded *int, constraintSlack float64, numSlack int) (string, error) {
	if math.Abs(constraintSlack) < EPSILON {
		// no slack variable
		return strings.Repeat("0 ", numSlack), nil
	}

	if *numSlackAdded >= numSlack {
		return "", fmt.Errorf("extra unexpected slack variable: %.2f", constraintSlack)
	}

	rightZeroAmount := numSlack - 1 - (*numSlackAdded)
	slackStr := strings.Repeat("0 ", *numSlackAdded) + ftos(constraintSlack) + " " + strings.Repeat("0 ", rightZeroAmount)
	*numSlackAdded++
	return slackStr, nil
}

// Prepares linear program to be passed into Simplex calculator (gets string to input)
func simplexInput(progArrays SimplexProgramArrays, toPositive map[string]struct{}, idTable map[string]int, idTableInverse map[int]string) (SimplexProgramStrings, error) {
	objective := progArrays.objective
	objectiveConst := progArrays.objectiveConst
	constraintsLHS := progArrays.constraintsLHS
	constraintsRHS := progArrays.constraintsRHS
	constraintsSlack := progArrays.constraintsSlack
	numSlack := progArrays.numSlack

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

// Calls the Simplex calculator (C++)
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
	if err != nil {
		return "", fmt.Errorf("error running solver: %v\noutput: %s", err, string(output))
	}

	return string(output), nil
}

// Gets the raw result from the Simplex calculator
func parseResult(output string, idTableInverse map[int]string) (SimplexResult, error) {
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

	return SimplexResult{
		Solution:    solution,
		ResultType:  resultType,
		Certificate: certificate,
		Mapping:     idTableInverse,
	}, nil
}

// Determines the optimal values of the Linear Program using the output of the Simplex calculator
func retrieveOriginalVariables(numSlack int, arr []float64, toPositive map[string]struct{}, idTableInverse map[int]string) ([]float64, error) {
	curVariableIdx := 0
	var newSolution []float64
	for i := 0; i < len(arr); i++ {
		// stop if it's just slack variables left
		if i >= len(arr)-numSlack {
			break
		}

		variable, ok := idTableInverse[curVariableIdx]
		if !ok {
			return []float64{}, fmt.Errorf("invalid variable index found at index %d, variable number %d", i, curVariableIdx)
		}

		// if we substituted to make a free variable positive, undo the substitution
		if _, ok = toPositive[variable]; ok {
			if i+1 >= len(arr) {
				return []float64{}, fmt.Errorf("non-positive variable found without nonnegative substitutes at index %d, variable number %d", i, curVariableIdx)
			}
			newSolution = append(newSolution, arr[i]-arr[i+1])
			i++
			curVariableIdx++
		} else {
			newSolution = append(newSolution, arr[i])
			curVariableIdx++
		}
	}

	return newSolution, nil
}
