package calculator_test

import (
	"errors"
	"math"
	"testing"

	"github.com/calculator-api/pkg/calculator"
)

// ---------------------------------------------------------------------------
// Unit tests for individual operations
// ---------------------------------------------------------------------------

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"large numbers", 1e15, 1e15, 2e15},
		{"decimals", 0.1, 0.2, 0.30000000000000004}, // IEEE 754
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Add(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Add(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.expected {
				t.Errorf("Add(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive result", 5, 3, 2},
		{"negative result", 3, 5, -2},
		{"zeros", 0, 0, 0},
		{"subtract negative", 5, -3, 8},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Subtract(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Subtract(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.expected {
				t.Errorf("Subtract(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 3, 4, 12},
		{"by zero", 5, 0, 0},
		{"negative numbers", -3, -4, 12},
		{"mixed signs", -3, 4, -12},
		{"identity", 7, 1, 7},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Multiply(tc.a, tc.b)
			if err != nil {
				t.Fatalf("Multiply(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.expected {
				t.Errorf("Multiply(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		expectErr error
	}{
		{"even division", 10, 2, 5, nil},
		{"fractional result", 7, 2, 3.5, nil},
		{"negative divisor", 10, -2, -5, nil},
		{"divide zero", 0, 5, 0, nil},
		{"division by zero", 10, 0, 0, calculator.ErrDivisionByZero},
		{"zero over zero", 0, 0, 0, calculator.ErrDivisionByZero},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Divide(tc.a, tc.b)
			if tc.expectErr != nil {
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("Divide(%v, %v) error = %v; want %v", tc.a, tc.b, err, tc.expectErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.expected {
				t.Errorf("Divide(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Integration-style tests for MathService.Calculate
// ---------------------------------------------------------------------------

func TestMathService_Calculate(t *testing.T) {
	svc := calculator.NewMathService()

	tests := []struct {
		name      string
		req       calculator.Request
		expected  float64
		expectErr error
	}{
		{
			name:     "add via service",
			req:      calculator.Request{A: 1, B: 2, Op: calculator.OpAdd},
			expected: 3,
		},
		{
			name:     "subtract via service",
			req:      calculator.Request{A: 10, B: 4, Op: calculator.OpSubtract},
			expected: 6,
		},
		{
			name:     "multiply via service",
			req:      calculator.Request{A: 3, B: 7, Op: calculator.OpMultiply},
			expected: 21,
		},
		{
			name:     "divide via service",
			req:      calculator.Request{A: 20, B: 4, Op: calculator.OpDivide},
			expected: 5,
		},
		{
			name:      "divide by zero via service",
			req:       calculator.Request{A: 1, B: 0, Op: calculator.OpDivide},
			expectErr: calculator.ErrDivisionByZero,
		},
		{
			name:      "unknown operation",
			req:       calculator.Request{A: 1, B: 2, Op: "modulus"},
			expectErr: calculator.ErrUnknownOperation,
		},
		{
			name:     "add negative decimals",
			req:      calculator.Request{A: -1.5, B: -2.5, Op: calculator.OpAdd},
			expected: -4,
		},
		{
			name:     "multiply by zero via service",
			req:      calculator.Request{A: 999, B: 0, Op: calculator.OpMultiply},
			expected: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.Calculate(tc.req)

			if tc.expectErr != nil {
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("Calculate(%+v) error = %v; want %v", tc.req, err, tc.expectErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Calculate(%+v) unexpected error: %v", tc.req, err)
			}

			if math.Abs(result.Value-tc.expected) > 1e-9 {
				t.Errorf("Calculate(%+v) = %v; want %v", tc.req, result.Value, tc.expected)
			}
		})
	}
}
