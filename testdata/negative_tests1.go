package ntests1

// In this test suite, (1) option is always used. No warnings should be generated.

// T is an example type.
type T struct {
	integer int
}

func zeroValuePointerAllocation() {
	_ = new(T)
	_ = new(T)
	_ = new(*T)
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	_ = make([]T, 0)
}

func nilSlice() {
	var s1 []int
	var s2 []*T
	var s3 []T
}

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T)
	_ = make(map[int]int)
}
