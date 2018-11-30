package ntests3

// In this test suite, (3) option is always used. No warnings should be generated.

func labelCase() {
allUpper:
foo:
anotherLowerCamelCase:
lowerCamelCase:
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

	_ = len(b) >= 1
	_ = len(m) >= 1
	_ = len(ch) >= 1
	_ = len(ch) >= 1
}
