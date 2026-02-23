package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bpospichil/calculator-be/internal/handler"
	"github.com/bpospichil/calculator-be/pkg/calculator"
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

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
