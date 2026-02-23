package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/calculator-api/internal/handler"
	"github.com/calculator-api/pkg/calculator"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	svc := calculator.NewMathService()
	h := handler.NewCalculatorHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", h.Calculate)

	// Simple health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	addr := ":" + port
	log.Printf("calculator-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
