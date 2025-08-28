package solve

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parse_sef"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
)

const EPSILON = 1e-9
const PRECISIONERROR = 1e-2

const solvePath = "/solve"
const textPlain = "text/plain"
const applicationJson = "application/json"
const methodPost = "POST"
const contentType = "Content-Type"

const pageNotFound = "404 PAGE NOT FOUND"
const methodNotAllowed = "405 METHOD NOT ALLOWED"
const unsupportedMediaType = "415 UNSUPPORTED MEDIA TYPE"
const internalServerError = "500 INTERNAL SERVER ERROR"

const enableObjective = true
const disableObjective = false

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

	if r.Header.Get(contentType) != textPlain {
		http.Error(w, unsupportedMediaType, http.StatusUnsupportedMediaType)
		return
	}

	progBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	progStr := string(progBytes)
	prog, idTable, err := parse_sef.ParseSEF(progStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	objectiveConst, objective, err := getExprArr(prog.Objective.Expr, idTable, enableObjective)
	if err != nil {
		http.Error(w, "error converting objective into array: "+err.Error(), http.StatusBadRequest)
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
	numSlack := 0
	for i, constraint := range prog.Constraints {
		_, curConstraintArr, err := getExprArr(constraint.Left, idTable, disableObjective)
		if err != nil {
			http.Error(w, "error converting constraint row "+strconv.Itoa(i)+" into array: "+err.Error(), http.StatusBadRequest)
			return
		}
		constraintsLHS = append(constraintsLHS, curConstraintArr)

		nl, ok := constraint.Right.(*parser.NumberLiteral)
		if !ok {
			http.Error(w, "right hand side is not NumberLiteral on constraint row"+strconv.Itoa(i)+": "+err.Error(), http.StatusBadRequest)
			return
		}
		constraintsRHS = append(constraintsRHS, nl.Value)

		switch constraint.Operator.Type {
		case lexer.TokenLessEqual:
			numSlack++
			// 1 because we're adding on just a ("1 * s") slack variable
			constraintsSlack[i] = 1
		case lexer.TokenEqual:
			continue
		case lexer.TokenGreaterEqual:
			numSlack++
			// -1 because we're subtracting a slack variable ("-1 * s")
			constraintsSlack[i] = -1
		default:
			// shouldn't have any other operator types
			http.Error(w, "invalid comparison operator on constraint row"+strconv.Itoa(i)+": "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	toPositive := allFreeVariables(idTable, make(map[string]struct{}))
	progArrays := SimplexProgramArrays{
		objective:        objective,
		objectiveConst:   objectiveConst,
		constraintsLHS:   constraintsLHS,
		constraintsRHS:   constraintsRHS,
		constraintsSlack: constraintsSlack,
		numSlack:         numSlack,
	}

	idTableInverse := getTableInverse(idTable)
	progStrings, err := simplexInput(progArrays, toPositive, idTable, idTableInverse)
	if err != nil {
		http.Error(w, "error converting arrays into strings: "+err.Error(), http.StatusBadRequest)
		return
	}

	rowSize := strconv.Itoa(len(constraintsLHS))
	// before converted colSize + number of slack variables we added + number of complementary variables we added (complementary := a - b for a, b >= 0)
	colSize := strconv.Itoa(numSlack + len(toPositive) + len(objective))
	output, err := callSimplex(progStrings, rowSize, colSize)
	if err != nil {
		http.Error(w, "error calling simplex method: "+err.Error(), http.StatusBadRequest)
		return
	}

	res, err := parseResult(output, idTableInverse)
	if err != nil {
		http.Error(w, "error parsing simplex method final result: "+err.Error(), http.StatusBadRequest)
		return
	}
	unsubstitutedSolution, err := retrieveOriginalVariables(numSlack, res.Solution, toPositive, idTableInverse)
	if err != nil {
		http.Error(w, "error converting final result variables back to original form: "+err.Error(), http.StatusBadRequest)
		return
	}
	res.Solution = unsubstitutedSolution

	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
