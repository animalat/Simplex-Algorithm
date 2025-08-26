package solve

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/animalat/Simplex-Algorithm/lp_parser/lexer"
	"github.com/animalat/Simplex-Algorithm/lp_parser/parser"
	"github.com/animalat/Simplex-Algorithm/lp_parser/semantics"
)

func assertProg(t *testing.T, s string, objectiveConstWanted float64, objectiveWanted []float64) {
	t.Helper()

	tokens, err := lexer.Tokenize(strings.NewReader(s))
	if err != nil {
		t.Fatalf("Tokenize() error: %v", err)
	}

	p := &parser.Parser{Tokens: tokens}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	idTable, err := semantics.SemanticCheck(prog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if testing.Verbose() {
		parser.PrintParse(prog)
	}

	objectiveConst, objective, err := getExprArr(prog.Objective.Expr, idTable, enableObjective)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// check const
	if !floatsEqual(objectiveConst, objectiveConstWanted) {
		t.Fatalf("objective constants not equal. Wanted: %.2f, Obtained: %.2f", objectiveConstWanted, objectiveConst)
	}

	// check objective numbers
	if len(objective) != len(objectiveWanted) {
		t.Fatalf("objective functions not equal size: Wanted size %d, obtained size %d", len(objectiveWanted), len(objective))
	}

	for i := range objective {
		if !floatsEqual(objective[i], objectiveWanted[i]) {
			t.Fatalf("objective values not equal. Wanted: %.2f, Obtained: %.2f, Position: %d", objectiveWanted[i], objective[i], i)
		}
	}
}

func assertPostRequest(t *testing.T, body []byte, solutionWanted []float64, resultTypeWanted string, certificateWanted []float64) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, solvePath, bytes.NewReader(body))
	req.Header.Set(contentType, textPlain)
	w := httptest.NewRecorder()

	HandleSolve(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errorMsg, err := io.ReadAll(res.Body)
		var errorMsgStr string
		if err != nil {
			errorMsgStr = "unable to read error message"
		} else {
			errorMsgStr = string(errorMsg)
		}
		t.Fatalf("expected status %d, got %d. %s", http.StatusOK, res.StatusCode, errorMsgStr)
	}

	if res.Header.Get(contentType) != applicationJson {
		t.Fatalf("expected %s %s, got %s", contentType, applicationJson, res.Header.Get(contentType))
	}

	var output SimplexResult
	if err := json.NewDecoder(res.Body).Decode(&output); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resultTypeWanted != output.ResultType {
		t.Fatalf("expected resultType %s, received type %v", resultTypeWanted, output.ResultType)
	}

	if len(solutionWanted) != len(output.Solution) {
		t.Fatalf("solution wanted of length %d, received length %d", len(solutionWanted), len(output.Solution))
	}

	for i := range solutionWanted {
		if !floatsEqual(solutionWanted[i], output.Solution[i]) {
			t.Fatalf("solutions not equal at index %d: wanted %.2f, received %.2f", i, solutionWanted[i], output.Solution[i])
		}
	}

	if len(certificateWanted) != len(output.Certificate) {
		t.Fatalf("solution wanted of length %d, received length %d", len(certificateWanted), len(output.Certificate))
	}

	for i := range certificateWanted {
		if !floatsEqual(certificateWanted[i], output.Certificate[i]) {
			t.Fatalf("certificates not equal at index %d: wanted %.2f, received %.2f", i, certificateWanted[i], output.Certificate[i])
		}
	}
}

func TestSolve_GetExprArr(t *testing.T) {
	assertProg(t, "let x1; let x2; let x3; max x1 + x2 + 3; s.t. x1 + x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 3, []float64{1, 1, 0})
	assertProg(t, "let x1; let x2; let x3; max 4 * x1 + 5 * x3; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 0, []float64{4, 0, 5})
	assertProg(t, "let x1; let x2; let x3; let x4; max 4 * x1 + x2 + 0 * x3 + 5 * x4 + 100; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 100, []float64{4, 1, 0, 5})
	assertProg(t, "let x1; let x2; let x3; let x4; let x5; max 4 * x1 + x2 + 0 * x3 + 5 * x4 + 10 * x5 + 100; s.t. 5 * x1 + 3 * x2 <= 3; x1 + x2 + 3 * x3 >= 5;", 100, []float64{4, 1, 0, 5, 10})
}

func TestSolve_PostRequest(t *testing.T) {
	assertPostRequest(t, []byte("let x1; max 4 * x1; s.t. 4 * x1 <= 5; x1 >= 0;"), []float64{1.25, 0.00}, "optimal", []float64{1.00})
}
