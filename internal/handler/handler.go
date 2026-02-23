// Package handler provides HTTP handlers for the calculator REST API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bpospichil/calculator-be/pkg/calculator"
)

// CalculatorHandler handles HTTP requests for calculations.
type CalculatorHandler struct {
	svc *calculator.MathService
}

// NewCalculatorHandler creates a new handler backed by the given MathService.
func NewCalculatorHandler(svc *calculator.MathService) *CalculatorHandler {
	return &CalculatorHandler{svc: svc}
}

// errorResponse is a JSON envelope for error messages.
type errorResponse struct {
	Error string `json:"error"`
}

// calculateResponse is the JSON envelope for successful results.
type calculateResponse struct {
	Result float64 `json:"result"`
}

// Calculate handles POST /calculate requests.
func (h *CalculatorHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req calculator.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload: "+err.Error())
		return
	}

	result, err := h.svc.Calculate(req)
	if err != nil {
		switch {
		case errors.Is(err, calculator.ErrDivisionByZero):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, calculator.ErrUnknownOperation):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{Result: result.Value})
}

// writeJSON serialises v as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a structured JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
