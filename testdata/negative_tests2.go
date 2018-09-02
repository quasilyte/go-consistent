package ntests2

// In this test suite, (2) option is always used. No warnings should be generated.

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
