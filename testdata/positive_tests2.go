package tests2

// In this test suite, (2) option is always preferred.

// T is an example type.
type T struct {
	integer int
}

func zeroValPtrAlloc() {
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

func nilSliceDecl() {
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

func nilMapDecl() {
	//= nilMapDecl: use literal instead of var
	var m1 map[int]int
	m2 := map[string]*T(nil)
	m3 := map[*T]string(nil)
}

func hexLit() {
	//= hexLit: use A-F instead of a-f
	_ = 0xff
	_ = 0xABCDEF
	_ = 0xABCD
}

func rangeCheck(x, low, high int) {
	//= rangeCheck: use align-center instead of align-left
	_ = x >= low && x <= high
	_ = low < x+1 || x+1 < high
	_ = low < x && x < high
}

func andNot(x, y int) {
	//= andNot: use &-plus-^ instead of &^
	_ = x &^ y
	_ = 123 & ^x
	_ = (x + 100) & ^(y + 2)
}

func floatLit() {
	//= floatLit: use omitted-int/frac instead of explicit-int/frac
	_ = 0.1
	_ = 1.
	_ = .0
}
