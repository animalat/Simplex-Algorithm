package main

import (
	"log"
	"net/http"

	"github.com/animalat/Simplex-Algorithm/backend/service/solve"
)

func main() {
	http.HandleFunc("/solve", solve.HandleSolve)

	const port = ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}
