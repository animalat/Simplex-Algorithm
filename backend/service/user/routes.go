package user

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/solve", h.handleSolve).Methods("POST")
}

func (h *Handler) handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
}
