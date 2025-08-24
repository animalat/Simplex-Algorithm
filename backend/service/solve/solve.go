package solve

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
	"github.com/animalat/Simplex-Algorithm/lp_parser/semantics"
)

const EPSILON = 1e-9

const solvePath = "/solve"
const methodPost = "POST"

const badRequest = "400 BAD REQUEST"
const pageNotFound = "404 PAGE NOT FOUND"
const methodNotAllowed = "405 METHOD NOT ALLOWED"
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
	numSlack := 0
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
			http.Error(w, badRequest, http.StatusBadRequest)
			return
		}
	}

	toPositive := allFreeVariables(idTable, make(map[string]struct{}))
	progStrings, err := simplexInput(SimplexProgramArrays{objective, objectiveConst, constraintsLHS, constraintsRHS, constraintsSlack, numSlack}, toPositive, idTable)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	rowSize := strconv.Itoa(len(constraintsLHS))
	// before converted colSize + number of slack variables we added + number of complementary variables we added (complementary := a - b for a, b >= 0)
	colSize := strconv.Itoa(numSlack + len(toPositive) + len(objective))
	output, err := callSimplex(progStrings, rowSize, colSize)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}
	// TODO: combine constraintsLHS and constraintsSlack
	// TODO: add zeroes onto objective to correspond to slack variables
	// TODO: force positive variables (maybe even recognize variables already >= 0 and also flip variables <= 0? later...)
}
