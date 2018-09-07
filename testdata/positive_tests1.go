package tests1

// In this test suite, (1) option is always preferred.

// T is an example type.
type T struct {
	integer int
}

func zeroValPtrAlloc() {
	_ = new(T)
	_ = new(T)
	//= zeroValPtrAlloc: use new instead of address-of-literal
	_ = &T{}
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	//= emptySlice: use make instead of literal
	_ = []string{}
}

func nilSliceDecl() {
	var s1 []int
	var s2 []*T
	//= nilSliceDecl: use var instead of literal
	s3 := []string(nil)
}

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T)
	//= emptyMap: use make instead of literal
	_ = map[int]int{}
}

func nilMapDecl() {
	var m1 map[int]int
	var m2 map[string]*T
	//= nilMapDecl: use var instead of literal
	m3 := map[*T]string(nil)
}

func hexLit() {
	_ = 0xff
	_ = 0xabcdef
	//= hexLit: use a-f instead of A-F
	_ = 0xABCD
}

func rangeCheck(x, low, high int) {
	_ = x > low && x <= high
	_ = x+1 >= low && x+1 < high
	//= rangeCheck: use align-left instead of align-center
	_ = low < x || x < high
}

func andNot(x, y int) {
	_ = x &^ y
	_ = 123 &^ x
	//= andNot: use &^ instead of &-plus-^
	_ = (x + 100) & ^(y + 2)
}

func floatLit() {
	_ = 0.0
	_ = 0.123
	//= floatLit: use explicit-int/frac instead of omitted-int/frac
	_ = 0.
}
