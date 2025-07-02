package arglisttests

type Person struct {
	Name string
	Age  int
}

func structLiterals() {
	_ = Person{
		Name: "John",
		Age:  30}

	_ = Person{Name: "Jane",
		Age: 25}

	//= arg list parens: align `)` to a same line with last argument
	_ = Person{
		Name: "Bob",
		Age:  40,
	}
}

func sliceLiterals() {
	_ = []string{
		"first",
		"second",
		"third"}

	_ = []int{1,
		2,
		3}

	//= arg list parens: align `)` to a same line with last argument
	_ = []string{
		"alpha",
		"beta",
		"gamma",
	}

	//= arg list parens: align `)` to a same line with last argument
	_ = []Person{
		{Name: "Alice", Age: 28},
		{Name: "Charlie", Age: 35},
	}
}

func mapLiterals() {
	_ = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3}

	_ = map[int]string{1: "first",
		2: "second",
		3: "third"}

	//= arg list parens: align `)` to a same line with last argument
	_ = map[string]Person{
		"john": {Name: "John", Age: 30},
		"jane": {Name: "Jane", Age: 25},
	}

	//= arg list parens: align `)` to a same line with last argument
	_ = map[string]interface{}{
		"name": "test",
		"age":  42,
	}
}

func functionWithMultipleParams(
	name string,
	age int,
	active bool) {
}

func anotherFunction(param1 string,
	param2 int,
	param3 bool) {
}

// = arg list parens: align `)` to a same line with last argument
func functionWithTrailingComma(
	name string,
	age int,
	active bool,
) {
}

// = arg list parens: align `)` to a same line with last argument
func complexFunction(
	param1 string,
	param2 []int,
	param3 map[string]interface{},
) error {
	return nil
}

type Calculator struct{}

func (c Calculator) Add(
	x int,
	y int) int {
	return x + y
}

func (c Calculator) Multiply(a int,
	b int,
	z int) int {
	return a * b * z
}

// = arg list parens: align `)` to a same line with last argument
func (c Calculator) Divide(
	dividend float64,
	divisor float64,
) float64 {
	return dividend / divisor
}

// = arg list parens: align `)` to a same line with last argument
func (c *Calculator) ComplexOperation(
	param1 string,
	param2 []int,
	param3 map[string]interface{},
) (result interface{}, err error) {
	return nil, nil
}

type MathInterface interface {
	Calculate(
		x int,
		y int) int

	Process(input string,
		options map[string]interface{},
		callback func(string) error) error

	//= arg list parens: align `)` to a same line with last argument
	Transform(
		data []byte,
		format string,
	) ([]byte, error)

	//= arg list parens: align `)` to a same line with last argument
	Validate(
		input interface{},
		rules []string,
		strict bool,
	) bool
}

type AdvancedInterface interface {
	//= arg list parens: align `)` to a same line with last argument
	ComplexMethod(
		param1 string,
		param2 func(int) error,
		param3 chan string,
	) <-chan result

	SimpleMethod(
		a int,
		b string) error
}

type result struct {
	Value interface{}
	Error error
}

func functionLiterals() {
	fn1 := func(
		x int,
		y int) int {
		return x + y
	}

	fn2 := func(a string,
		b string,
		c string) string {
		return a + b + c
	}

	//= arg list parens: align `)` to a same line with last argument
	fn3 := func(
		param1 string,
		param2 int,
	) error {
		return nil
	}

	//= arg list parens: align `)` to a same line with last argument
	fn4 := func(
		data []byte,
		callback func(error),
		options map[string]interface{},
	) {
	}

	_ = fn1
	_ = fn2
	_ = fn3
	_ = fn4
}

func nestedLiterals() {
	_ = []map[string]Person{
		{
			"manager":  {Name: "Alice", Age: 35},
			"employee": {Name: "Bob", Age: 28}},
		{
			"lead":   {Name: "Charlie", Age: 40},
			"junior": {Name: "David", Age: 22}}}

	//= arg list parens: align `)` to a same line with last argument
	_ = map[string][]Person{
		"team1": {
			{Name: "Eve", Age: 30},
			{Name: "Frank", Age: 32},
		},
		"team2": {
			{Name: "Grace", Age: 29},
			{Name: "Henry", Age: 31},
		},
	}
}
