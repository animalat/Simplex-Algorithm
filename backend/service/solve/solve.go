package solve

import (
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

func getObjectiveArr(p *parser.Program, map[string]bool) ([]int, error) {

}

func HandleSolve(w http.ResponseWriter, r *http.Request) {
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
}
