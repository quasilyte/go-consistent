package tests3

// In this test suite, (3) option is always preferred.

func labelCase() {
	//= label case: use lowerCamelCase
ALL_UPPER:
	//= label case: use lowerCamelCase
Foo:
anotherLowerCamelCase:
lowerCamelCase:
	goto ALL_UPPER
	goto Foo
	goto anotherLowerCamelCase
	goto lowerCamelCase
}

func nonZeroLenTestChecker() {
	var (
		s  string
		b  []byte
		m  map[int]int
		ch chan int
	)

	// Strings are ignored.
	_ = len(s) > 0
	_ = len(s) > 0
	_ = len(s) > 0

	//= non-zero length test: use `len(s) >= 1`
	_ = len(b) != 0
	//= non-zero length test: use `len(s) >= 1`
	_ = len(m) > 0
	_ = len(ch) >= 1
	_ = len(ch) >= 1
}
