// Package calculator provides core arithmetic operations exposed as a MathService.
package calculator

import (
	"errors"
	"fmt"
	"sync"
)

// ErrDivisionByZero is returned when a division by zero is attempted.
var ErrDivisionByZero = errors.New("division by zero")

// ErrUnknownOperation is returned when an unsupported operation is requested.
var ErrUnknownOperation = errors.New("unknown operation")

// Operation represents a supported arithmetic operation.
type Operation string

const (
	OpAdd      Operation = "add"
	OpSubtract Operation = "subtract"
	OpMultiply Operation = "multiply"
	OpDivide   Operation = "divide"
)

// OperationFunc defines the signature for an arithmetic operation.
type OperationFunc func(a, b float64) (float64, error)

// Request holds the operands and the desired operation.
type Request struct {
	A  float64   `json:"a"`
	B  float64   `json:"b"`
	Op Operation `json:"operation"`
}

// Result holds the outcome of a calculation.
type Result struct {
	Value float64 `json:"result"`
}

// MathService performs arithmetic calculations.
type MathService struct {
	mu         sync.RWMutex
	operations map[Operation]OperationFunc
}

// NewMathService creates a new MathService instance with the default operations registered.
func NewMathService() *MathService {
	s := &MathService{
		operations: make(map[Operation]OperationFunc),
	}

	s.Register(OpAdd, Add)
	s.Register(OpSubtract, Subtract)
	s.Register(OpMultiply, Multiply)
	s.Register(OpDivide, Divide)

	return s
}

// Register adds (or replaces) an operation in the service.
func (s *MathService) Register(op Operation, fn OperationFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operations[op] = fn
}

// Calculate performs the arithmetic operation described by req.
func (s *MathService) Calculate(req Request) (Result, error) {
	s.mu.RLock()
	fn, ok := s.operations[req.Op]
	s.mu.RUnlock()

	if !ok {
		return Result{}, fmt.Errorf("%w: %q", ErrUnknownOperation, req.Op)
	}

	v, err := fn(req.A, req.B)
	if err != nil {
		return Result{}, err
	}
	return Result{Value: v}, nil
}

// Add returns the sum of a and b.
func Add(a, b float64) (float64, error) {
	return a + b, nil
}

// Subtract returns the difference of a and b.
func Subtract(a, b float64) (float64, error) {
	return a - b, nil
}

// Multiply returns the product of a and b.
func Multiply(a, b float64) (float64, error) {
	return a * b, nil
}

// Divide returns the quotient of a and b.
// It returns ErrDivisionByZero when b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}
