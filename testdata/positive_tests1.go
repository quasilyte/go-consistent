package tests1

// In this test suite, (1) option is always preferred.

// T is an example type.
type T struct {
	integer int
}

func zeroValuePointerAllocation() {
	_ = new(T)
	_ = new(T)
	//= zero value pointer allocation: use new instead of address-of-lit
	_ = &T{}
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	//= empty slice: use empty-slice-make instead of empty-slice-lit
	_ = []string{}
}

func nilSlice() {
	var s1 []int
	var s2 []*T
	//= nil slice: use nil-slice-var instead of nil-slice-lit
	s3 := []string(nil)
}

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T)
	//= empty map: use empty-map-make instead of empty-map-lit
	_ = map[int]int{}
}
