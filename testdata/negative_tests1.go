package ntests1

// In this test suite, (1) option is always used. No warnings should be generated.

// T is an example type.
type T struct {
	integer int
}

func zeroValPtrAlloc() {
	_ = new(T)
	_ = new(T)
	_ = new(*T)
}

func emptySlice() {
	_ = make([]int, 0)
	_ = make([]float64, 0)
	_ = make([]T, 0)
}

func nilSliceDecl() {
	var s1 []int
	var s2 []*T
	var s3 []T
}

var nilSliceVar1 []string
var nilSliceVar2 = []string(nil)

func emptyMap() {
	_ = make(map[T]T)
	_ = make(map[*T]*T)
	_ = make(map[int]int)
}

func nilMapDecl() {
	var m1 map[int]int
	var m2 map[string]*T
	var m3 map[*T]string
}

var nilMapVar1 []string
var nilMapVar2 = []string(nil)

func hexLit() {
	_ = 0xff
	_ = 0xabcdef
	_ = 0xabcd
}

func rangeCheck(x, low, high int) {
	_ = x > low && x <= high
	_ = x+1 >= low && x+1 < high
	_ = x > low || x < high
}

func andNot(x, y int) {
	_ = x &^ y
	_ = 123 &^ x
	_ = (x + 100) &^ (y + 2)
}

func floatLit() {
	_ = 0.0
	_ = 0.123
	_ = 1.0
}
