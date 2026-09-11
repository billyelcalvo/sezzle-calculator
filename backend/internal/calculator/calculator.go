package calculator

import "errors"

var ErrDivisionByZero = errors.New("division by zero")

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
