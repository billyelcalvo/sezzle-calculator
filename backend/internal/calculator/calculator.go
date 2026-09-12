package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero     = errors.New("division by zero")
	ErrNegativeSquareRoot = errors.New("cannot calculate square root of a negative number")
	ErrInvalidPower       = errors.New("power result is not a finite real number")
)

func Add(left, right float64) float64 {
	return left + right
}

func Subtract(left, right float64) float64 {
	return left - right
}

func Multiply(left, right float64) float64 {
	return left * right
}

func Divide(dividend, divisor float64) (float64, error) {
	if divisor == 0 {
		return 0, ErrDivisionByZero
	}

	return dividend / divisor, nil
}

func Power(base, exponent float64) (float64, error) {
	result := math.Pow(base, exponent)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrInvalidPower
	}
	return result, nil
}

func SquareRoot(value float64) (float64, error) {
	if value < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return math.Sqrt(value), nil
}

func Percentage(percent, value float64) float64 {
	return (percent / 100) * value
}
