// Package math provides utilities to do some math.
package math

import "golang.org/x/exp/constraints"

// Number represents integer or float type.
type Number interface {
	constraints.Integer | constraints.Float
}

// Add add two Number and return a Number.
//
// The number type is using [Number] interface.
//
// It has two parameters: the first Number and the second Number.
// Add returns the added value in Number.
//
// More information on addition can be found at [MathIsFun].
//
// [MathIsFun]: https://www.mathsisfun.com/numbers/addition.html
func Add[T Number](a, b T) T {
	return a + b
}
