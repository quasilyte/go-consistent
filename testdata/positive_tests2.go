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
	goto ALL_UPPER
	goto Foo
	goto UpperCamelCase
	goto lowerCamelCase
}

func untypedConstCoerce() {
	const zero = 0

	//= untyped const coerce: specity type in RHS, like in `var x = T(const)`
	var _ int = zero
	var _ = int32(10)
	var _ = int64(zero + 1)
}

func threeArgs(a, b, c int) {}

func argListParens() {
	//= arg list parens: move `)` to the next line and put `,` after the last argument
	threeArgs(
		1,
		2,
		3)
	threeArgs(1,
		2,
		3,
	)
	threeArgs(
		1,
		2,
		3,
	)
}

func nonZeroLenTestChecker() {
	var (
		s  string
		b  []byte
		m  map[int]int
		ch chan int
	)

	// Strings are ignored.
	_ = len(s) >= 1
	_ = len(s) >= 1
	_ = len(s) >= 1

	//= non-zero length test: use `len(s) > 0`
	_ = len(b) != 0
	_ = len(m) > 0
	_ = len(ch) > 0
	//= non-zero length test: use `len(s) > 0`
	_ = len(ch) >= 1
}

func defaultCaseOrder(x int, v interface{}) {
	//= default case order: default case should be the last case
	switch x {
	default:
	case 10:
	}

	switch v.(type) {
	case int:
	case string:
	default:
	}

	switch {
	case x > 20:
	default:
	}
}

func omitTypes(a, b, c int) {
	return
}

func omitTypes2(a, b, c int) {
	return
}

func allTypes(a int, b int, c int) { //= use only one type declaration after several arguments of the same type
	return
}
