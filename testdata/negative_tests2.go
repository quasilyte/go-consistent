package ntests2

// In this test suite, (2) option is always used. No warnings should be generated.

// T is an example type.
type T struct {
	integer int
}

func zeroValuePointerAllocation() {
	_ = &T{}
	_ = &T{}
	_ = &T{}
}

func emptySlice() {
	_ = []rune{}
	_ = []float64{}
	_ = []string{}
}

func nilSlice() {
	_ = []rune(nil)
	_ = []*T(nil)
	_ = []string(nil)
}

func emptyMap() {
	_ = map[rune]rune{}
	_ = map[*T]*T{}
	_ = map[int]int{}
}
