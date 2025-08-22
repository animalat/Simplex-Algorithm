package solve

import (
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
}
