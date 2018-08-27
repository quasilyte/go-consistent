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
