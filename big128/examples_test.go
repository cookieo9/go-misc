package big128_test

import (
	"fmt"
	"github.com/cookieo9/go-misc/big128"
	"math/big"
)

func ExampleNewInt() {
	x := big128.NewInt(52)
	y := big128.NewInt(10)
	z := big128.NewInt(0)
	z.Sub(x, y)

	fmt.Println(z)
	// Output: 42
}

func ExamplePreallocInts() {
	// Compute first 10 fibonacci numbers
	var fibs [10]big.Int
	big128.PreallocInts(fibs[:])

	fibs[0].SetInt64(1)
	fibs[1].SetInt64(1)

	for i := 2; i < len(fibs); i++ {
		fibs[i].Add(&fibs[i-1], &fibs[i-2])
	}

	for _, v := range fibs {
		fmt.Print(v.String(), " ")
	}
	fmt.Println()
	// Output: 1 1 2 3 5 8 13 21 34 55
}

func ExampleNewInts() {
	// Compute first 10 fibonacci numbers
	fibs := big128.NewInts(10)
	fibs[0].SetInt64(1)
	fibs[1].SetInt64(1)

	for i := 2; i < len(fibs); i++ {
		fibs[i].Add(&fibs[i-1], &fibs[i-2])
	}

	for _, v := range fibs {
		fmt.Print(v.String(), " ")
	}
	fmt.Println()
	// Output: 1 1 2 3 5 8 13 21 34 55
}

func ExampleNewIntPtrs() {
	// Compute first 10 fibonacci numbers
	fibs := big128.NewIntPtrs(10)
	fibs[0].SetInt64(1)
	fibs[1].SetInt64(1)

	for i := 2; i < len(fibs); i++ {
		fibs[i].Add(fibs[i-1], fibs[i-2])
	}

	fmt.Println(fibs)
	// Output: [1 1 2 3 5 8 13 21 34 55]
}
