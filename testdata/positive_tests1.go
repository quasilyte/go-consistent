package tests1

// In this test suite, (1) option is always preferred.

import "strconv"

import "errors"

//= unit import: omit parenthesis in a single-package import
import (
	"fmt"
)

var (
	_ = fmt.Printf
	_ = errors.New
	_ = strconv.Atoi
)

// T is an example type.
type T struct {
	integer int
}

func zeroValPtrAlloc() {
	_ = new(T)
	_ = new(map[string]bool)
	_ = new([]int)
	//= zero value ptr alloc: use new(T) for *T allocation
	_ = &T{}
	//= zero value ptr alloc: use new(T) for *T allocation
	_ = &[]int{}
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	//= empty slice: use make([]T, 0)
	_ = []string{}
}

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T, 0)
	//= empty map: use make(map[K]V)
	_ = map[int]int{}
}

func hexLit() {
	_ = 0xff
	_ = 0xabcdef
	//= hex lit: use a-f (lower case) digits
	_ = 0xABCD
}

func rangeCheck(x, low, high int) {
	_ = x > low && x <= high
	_ = x+1 >= low && x+1 < high
	_ = x >= low && x <= high
	//= range check: use align-left, like in `x >= low && x <= high`
	_ = low < x || x < high
}

func andNot(x, y int) {
	_ = x &^ y
	_ = 123 &^ x
	//= and-not: remove a space between & and ^, like in `x &^ y`
	_ = (x + 100) & ^(y + 2)
}

func floatLit() {
	_ = 0.0
	_ = 0.123
	_ = 1.0
	//= float lit: use explicit int/frac part, like in `1.0` and `0.1`
	_ = 0.
	//= float lit: use explicit int/frac part, like in `1.0` and `0.1`
	_ = .0
}

func labelCase() {
ALL_UPPER:
FOO:
	//= label case: use ALL_UPPER
UpperCamelCase:
	//= label case: use ALL_UPPER
lowerCamelCase:
}

func untypedConstCoerce() {
	const zero = 0

	var _ int = zero
	var _ int32 = 10
	//= untyped const coerce: specify type in LHS, like in `var x T = const`
	var _ = int64(zero + 1)
}
