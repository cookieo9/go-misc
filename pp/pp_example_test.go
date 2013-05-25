package pp_test

import (
	"fmt"
	"github.com/cookieo9/go-misc/pp"
)

func ExamplePP() {
	val := []interface{}{"Hello, World!", 42}
	fmt.Println(pp.PP(val))
	fmt.Println(pp.PP([]byte("Hello, World!")))

	// Output:
	// [
	//     "Hello, World!",
	//     42,
	// ]
	// [
	//     72,
	//     101,
	//     108,
	//     108,
	//     111,
	//     44,
	//     32,
	//     87,
	//     111,
	//     114,
	//     108,
	//     100,
	//     33,
	// ]
}
