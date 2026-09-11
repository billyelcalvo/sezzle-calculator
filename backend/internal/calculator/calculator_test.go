package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(tester *testing.T) {
	tests := []struct {
		name     string
		left     float64
		right    float64
		expected float64
	}{
		{name: "positive numbers", left: 10, right: 5, expected: 15},
		{name: "negative numbers", left: -10, right: -5, expected: -15},
		{name: "mixed signs", left: -10, right: 5, expected: -5},
		{name: "zero", left: 10, right: 0, expected: 10},
		{name: "decimals", left: 0.1, right: 0.2, expected: 0.3},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result := Add(test.left, test.right)
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Add(%v, %v) = %v; want %v", test.left, test.right, result, test.expected)
			}
		})
	}
}

func TestSubtract(tester *testing.T) {
	tests := []struct {
		name     string
		left     float64
		right    float64
		expected float64
	}{
		{name: "positive numbers", left: 10, right: 5, expected: 5},
		{name: "negative result", left: 5, right: 10, expected: -5},
		{name: "negative numbers", left: -10, right: -5, expected: -5},
		{name: "subtract negative", left: 10, right: -5, expected: 15},
		{name: "zero", left: 10, right: 0, expected: 10},
		{name: "decimals", left: 0.3, right: 0.1, expected: 0.2},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result := Subtract(test.left, test.right)
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Subtract(%v, %v) = %v; want %v", test.left, test.right, result, test.expected)
			}
		})
	}
}

func TestMultiply(tester *testing.T) {
	tests := []struct {
		name     string
		left     float64
		right    float64
		expected float64
	}{
		{name: "positive numbers", left: 10, right: 5, expected: 50},
		{name: "mixed signs", left: -10, right: 5, expected: -50},
		{name: "negative numbers", left: -10, right: -5, expected: 50},
		{name: "zero", left: 10, right: 0, expected: 0},
		{name: "decimals", left: 0.1, right: 0.2, expected: 0.02},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result := Multiply(test.left, test.right)
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Multiply(%v, %v) = %v; want %v", test.left, test.right, result, test.expected)
			}
		})
	}
}

func TestDivide(tester *testing.T) {
	tests := []struct {
		name          string
		dividend      float64
		divisor       float64
		expected      float64
		expectedError error
	}{
		{name: "positive numbers", dividend: 10, divisor: 5, expected: 2},
		{name: "fractional result", dividend: 5, divisor: 2, expected: 2.5},
		{name: "negative dividend", dividend: -10, divisor: 5, expected: -2},
		{name: "negative divisor", dividend: 10, divisor: -5, expected: -2},
		{name: "negative numbers", dividend: -10, divisor: -5, expected: 2},
		{name: "zero dividend", dividend: 0, divisor: 5, expected: 0},
		{name: "decimals", dividend: 0.3, divisor: 0.1, expected: 3},
		{name: "division by zero", dividend: 10, divisor: 0, expectedError: ErrDivisionByZero},
		{name: "zero divided by zero", dividend: 0, divisor: 0, expectedError: ErrDivisionByZero},
		{name: "division by negative zero", dividend: 10, divisor: math.Copysign(0, -1), expectedError: ErrDivisionByZero},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result, err := Divide(test.dividend, test.divisor)
			if !errors.Is(err, test.expectedError) {
				tester.Fatalf("Divide(%v, %v) error = %v; want %v", test.dividend, test.divisor, err, test.expectedError)
			}
			if test.expectedError != nil {
				return
			}
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Divide(%v, %v) = %v; want %v", test.dividend, test.divisor, result, test.expected)
			}
		})
	}
}
