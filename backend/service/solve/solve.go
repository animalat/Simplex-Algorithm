package solve

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
	"github.com/animalat/Simplex-Algorithm/lp_parser/semantics"
)

const solvePath = "/solve"
const methodPost = "POST"

const badRequest = "400 BAD REQUEST"
const pageNotFound = "404 PAGE NOT FOUND"
const methodNotAllowed = "405 METHOD NOT ALLOWED"
const internalServerError = "500 INTERNAL SERVER ERROR"

const enableObjective = true
const disableObjective = false

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

func HandleSolve(w http.ResponseWriter, r *http.Request) {
	// TODO: handle malicious input
	defer r.Body.Close()

	if r.URL.Path != solvePath {
		http.Error(w, pageNotFound, http.StatusNotFound)
		return
	}

	if r.Method != methodPost {
		http.Error(w, methodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	progBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	progStr := string(progBytes)

	tokens, err := lexer.Tokenize(strings.NewReader(progStr))
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	parseProg := parser.ConstructParser(tokens)
	prog, err := parseProg.ParseProgram()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	idTable, err := semantics.SemanticCheck(prog)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	// TODO: pass into solver
	objectiveConst, objective, err := getExprArr(prog.Objective.Expr, idTable, enableObjective)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	if !prog.Objective.IsMax {
		// flip sign to make it a maximization problem
		for i := range objective {
			objective[i] *= -1
		}
		objectiveConst *= -1
	}

	constraintsLHS := make([][]float64, 0, len(prog.Constraints))
	constraintsRHS := make([]float64, 0, len(prog.Constraints))
	constraintsSlack := make([]float64, len(prog.Constraints))
	for i, constraint := range prog.Constraints {
		_, curConstraintArr, err := getExprArr(constraint.Left, idTable, disableObjective)
		if err != nil {
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
		constraintsLHS = append(constraintsLHS, curConstraintArr)

		nl, ok := constraint.Right.(*parser.NumberLiteral)
		if !ok {
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
		constraintsRHS = append(constraintsRHS, nl.Value)

		switch constraint.Operator.Type {
		case lexer.TokenLessEqual:
			// 1 because we're adding on just a ("1 * s") slack variable
			constraintsSlack[i] = 1
		case lexer.TokenGreaterEqual:
			// -1 because we're subtracting a slack variable ("-1 * s")
			constraintsSlack[i] = -1
		default:
			continue
		}
	}
	// TODO: combine constraintsLHS and constraintsSlack
	// TODO: add zeroes onto objective to correspond to slack variables
	// TODO: force positive variables (maybe even recognize variables already >= 0 and also flip variables <= 0? later...)
}
