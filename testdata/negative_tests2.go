package ntests2

// In this test suite, (2) option is always used. No warnings should be generated.

import (
	"strconv"
)

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
	_ = &T{}
	_ = &T{}
	_ = &T{}

	// Not a zero value allocation:
	_ = &T{integer: 1}
}

func emptySlice() {
	_ = []rune{}
	_ = []float64{}
	_ = []string{}
}

func nilSliceDecl() {
	_ = []rune(nil)
	_ = []*T(nil)
	_ = []string(nil)
}

var nilSliceVar1 []string
var nilSliceVar2 = []string(nil)

func emptyMap() {
	_ = map[rune]rune{}
	_ = map[*T]*T{}
	_ = map[int]int{}
}

func nilMapDecl() {
	m1 := map[int]int(nil)
	m2 := map[string]*T(nil)
	m3 := map[*T]string(nil)
}

var nilMapVar1 []string
var nilMapVar2 = []string(nil)

func hexLit() {
	_ = 0xFF
	_ = 0xABCDEF
	_ = 0xABCD
}

func rangeCheck(x, low, high int) {
	_ = low <= x && x <= high
	_ = low < x+1 || x+1 < high
	_ = low < x && x < high
}

func andNot(x, y int) {
	_ = x & ^y
	_ = 123 & ^x
	_ = (x + 100) & ^(y + 2)
}

func floatLit() {
	_ = .43
	_ = 1.
	_ = .0
}

func labelCase() {
AllUpper:
Foo:
UpperCamelCase:
LowerCamelCase:
}
