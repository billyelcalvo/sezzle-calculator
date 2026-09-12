package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestPower(tester *testing.T) {
	tests := []struct {
		name          string
		base          float64
		exponent      float64
		expected      float64
		expectedError error
	}{
		{name: "positive exponent", base: 2, exponent: 3, expected: 8},
		{name: "zero exponent", base: 5, exponent: 0, expected: 1},
		{name: "negative exponent", base: 2, exponent: -3, expected: 0.125},
		{name: "fractional exponent", base: 9, exponent: 0.5, expected: 3},
		{name: "negative base", base: -2, exponent: 3, expected: -8},
		{name: "zero base", base: 0, exponent: 3, expected: 0},
		{name: "zero to zero follows math Pow", base: 0, exponent: 0, expected: 1},
		{name: "non-real result", base: -2, exponent: 0.5, expectedError: ErrInvalidPower},
		{name: "zero with negative exponent", base: 0, exponent: -1, expectedError: ErrInvalidPower},
		{name: "overflow", base: 1e308, exponent: 2, expectedError: ErrInvalidPower},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result, err := Power(test.base, test.exponent)
			if !errors.Is(err, test.expectedError) {
				tester.Fatalf("Power(%v, %v) error = %v; want %v", test.base, test.exponent, err, test.expectedError)
			}
			if test.expectedError != nil {
				return
			}
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Power(%v, %v) = %v; want %v", test.base, test.exponent, result, test.expected)
			}
		})
	}
}

func TestSquareRoot(tester *testing.T) {
	tests := []struct {
		name          string
		value         float64
		expected      float64
		expectedError error
	}{
		{name: "perfect square", value: 25, expected: 5},
		{name: "zero", value: 0, expected: 0},
		{name: "decimal", value: 2.25, expected: 1.5},
		{name: "irrational result", value: 2, expected: 1.4142135623730951},
		{name: "negative number", value: -1, expectedError: ErrNegativeSquareRoot},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result, err := SquareRoot(test.value)
			if !errors.Is(err, test.expectedError) {
				tester.Fatalf("SquareRoot(%v) error = %v; want %v", test.value, err, test.expectedError)
			}
			if test.expectedError != nil {
				return
			}
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("SquareRoot(%v) = %v; want %v", test.value, result, test.expected)
			}
		})
	}
}

func TestPercentage(tester *testing.T) {
	tests := []struct {
		name     string
		percent  float64
		value    float64
		expected float64
	}{
		{name: "twenty percent of 150", percent: 20, value: 150, expected: 30},
		{name: "zero percent", percent: 0, value: 150, expected: 0},
		{name: "zero value", percent: 20, value: 0, expected: 0},
		{name: "negative percent", percent: -20, value: 150, expected: -30},
		{name: "decimal percent", percent: 12.5, value: 80, expected: 10},
		{name: "over one hundred percent", percent: 150, value: 20, expected: 30},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			result := Percentage(test.percent, test.value)
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("Percentage(%v, %v) = %v; want %v", test.percent, test.value, result, test.expected)
			}
		})
	}
}

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
