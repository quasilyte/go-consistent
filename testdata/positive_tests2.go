package tests2

// In this test suite, (2) option is always preferred.

// T is an example type.
type T struct {
	integer int
}

func zeroValuePointerAllocation() {
	//= zeroValPtrAlloc: use address-of-literal instead of new
	_ = new(T)
	_ = &T{}
	_ = &T{}
}

func emptySlice() {
	//= emptySlice: use literal instead of make
	_ = make([]int, 0)
	_ = []float64{}
	_ = []string{}
}

func nilSlice() {
	//= nilSliceDecl: use literal instead of var
	var s1 []int
	s2 := []*T(nil)
	s3 := []string(nil)
}

func emptyMap() {
	//= emptyMap: use literal instead of make
	_ = make(map[T]T)
	_ = map[*T]*T{}
	_ = map[int]int{}
}
