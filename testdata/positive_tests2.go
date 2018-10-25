package tests2

// In this test suite, (2) option is always preferred.

//= unit import: wrap single-package import spec into parenthesis
import "strconv"

import (
	"errors"
)

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
	//= zero value ptr alloc: use &T{} for *T allocation
	_ = new(T)
	//= zero value ptr alloc: use &T{} for *T allocation
	_ = new(map[string]bool)
	_ = &T{}
	_ = &map[string]bool{}
	_ = &[]int{}
}

func emptySlice() {
	//= empty slice: use []T{}
	_ = make([]int, 0)
	_ = []float64{}
	_ = []string{}
}

func emptyMap() {
	//= empty map: use map[K]V{}
	_ = make(map[T]T)
	_ = map[*T]*T{}
	_ = map[int]int{}
}

func hexLit() {
	//= hex lit: use A-F (upper case) digits
	_ = 0xff
	_ = 0xABCDEF
	_ = 0xABCD
}

func rangeCheck(x, low, high int) {
	//= range check: use align-center, like in `low < x && x < high`
	_ = x > low && x <= high
	_ = low <= x+1 || x+1 <= high
	_ = low <= x || x < high
	_ = low < x || x < high
}

func andNot(x, y int) {
	//= and-not: put a space between & and ^, like in `x & ^y`
	_ = x &^ y
	_ = 123 & ^x
	_ = (x + 100) & ^(y + 2)
}

func floatLit() {
	//= float lit: use implicit int/frac part, like in `1.` and `.1`
	_ = 1.0
	//= float lit: use implicit int/frac part, like in `1.` and `.1`
	_ = 0.123
	_ = 11.
	_ = 0.
	_ = .0
}

func labelCase() {
	//= label case: use UpperCamelCase
ALL_UPPER:
Foo:
UpperCamelCase:
	//= label case: use UpperCamelCase
lowerCamelCase:
}

func untypedConstCoerce() {
	const zero = 0

	//= untyped const coerce: specity type in RHS, like in `var x = T(const)`
	var _ int = zero
	var _ = int32(10)
	var _ = int64(zero + 1)
}
