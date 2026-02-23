package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/calculator-api/internal/handler"
	"github.com/calculator-api/pkg/calculator"
)

func newHandler() *handler.CalculatorHandler {
	return handler.NewCalculatorHandler(calculator.NewMathService())
}

func TestCalculate_Success(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected float64
	}{
		{"add", `{"a":2,"b":3,"operation":"add"}`, 5},
		{"subtract", `{"a":10,"b":4,"operation":"subtract"}`, 6},
		{"multiply", `{"a":3,"b":7,"operation":"multiply"}`, 21},
		{"divide", `{"a":20,"b":4,"operation":"divide"}`, 5},
	}

	h := newHandler()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Calculate(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
			}

			var resp struct {
				Result float64 `json:"result"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp.Result != tc.expected {
				t.Errorf("result = %v; want %v", resp.Result, tc.expected)
			}
		})
	}
}

func TestCalculate_Errors(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantMsg:    "only POST is allowed",
		},
		{
			name:       "invalid JSON",
			method:     http.MethodPost,
			body:       `{bad json}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "invalid JSON payload",
		},
		{
			name:       "division by zero",
			method:     http.MethodPost,
			body:       `{"a":1,"b":0,"operation":"divide"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "division by zero",
		},
		{
			name:       "unknown operation",
			method:     http.MethodPost,
			body:       `{"a":1,"b":2,"operation":"power"}`,
			wantStatus: http.StatusBadRequest,
			wantMsg:    "unknown operation",
		},
	}

	h := newHandler()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/calculate", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Calculate(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d; want %d", rec.Code, tc.wantStatus)
			}

			var errResp struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}
			if !contains(errResp.Error, tc.wantMsg) {
				t.Errorf("error = %q; want to contain %q", errResp.Error, tc.wantMsg)
			}
		})
	}
}

func TestCalculate_ContentType(t *testing.T) {
	h := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(`{"a":1,"b":2,"operation":"add"}`))
	rec := httptest.NewRecorder()

	h.Calculate(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", ct)
	}
}

// contains checks if s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
