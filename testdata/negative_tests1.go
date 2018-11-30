package ntests1

// In this test suite, (1) option is always used. No warnings should be generated.

import "strconv"
import "errors"
import "fmt"

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
	_ = new(T)
	_ = new(*T)
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	_ = make([]T, 0)
}

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T)
	_ = make(map[int]int)
}

func hexLit() {
	_ = 0xff
	_ = 0xabcdef
	_ = 0xabcd
}

func rangeCheck(x, low, high int) {
	_ = x > low && x <= high
	_ = x+1 >= low && x+1 < high
	_ = x > low || x < high
}

func andNot(x, y int) {
	_ = x &^ y
	_ = 123 &^ x
	_ = (x + 100) &^ (y + 2)
}

func floatLit() {
	_ = 0.0
	_ = 0.123
	_ = 1.0
}

func labelCase() {
ALL_UPPER:
FOO:
UPPER_CAMEL_CASE:
LOWER_CAMEL_CASE:
}

func untypedConstCoerce() {
	const zero = 0

	var _ int = zero
	var _ int32 = 10
	var _ int64 = (zero + 1)
}

func argListParens() {
	threeArgs(
		1,
		2,
		3)
	threeArgs(1,
		2,
		3)
	threeArgs(
		1,
		2,
		3)
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

	_ = len(b) != 0
	_ = len(m) != 0
	_ = len(ch) != 0
	_ = len(ch) != 0
}
