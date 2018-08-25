package tests2

// In this test suite, (2) option is always preffered.

// T is an example type.
type T struct {
	integer int
}

func zeroValuePointerAllocation() {
	//= zero value pointer allocation: use address-of-lit instead of new
	_ = new(T)
	_ = &T{}
	_ = &T{}
}

func emptySlice() {
	//= empty slice: use empty-slice-lit instead of empty-slice-make
	_ = make([]int, 0)
	_ = []float64{}
	_ = []string{}
}

func nilSlice() {
	//= nil slice: use nil-slice-lit instead of nil-slice-var
	var s1 []int
	_ = []*T(nil)
	_ = []string(nil)
}

func emptyMap() {
	//= empty map: use empty-map-lit instead of empty-map-make
	_ = make(map[T]T)
	_ = map[*T]*T{}
	_ = map[int]int{}
}
