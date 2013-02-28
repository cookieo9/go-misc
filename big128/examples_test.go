package big128_test

import (
	big128 "."
	"fmt"
	"math/big"
)

func ExamplePreallocInt() {
	// Useless: adds extra code, and doesn't even result
	// in fewer allocations.
	var x, y, z big.Int

	big128.PreallocInt(&x)
	big128.PreallocInt(&y)
	big128.PreallocInt(&z)

	x.SetInt64(52)
	y.SetInt64(10)
	z.Sub(&x, &y)

	fmt.Println(z.String())
	// Output: 42
}
func ExampleNewInt() {
	// Use as drop-in replacement for math/big.NewInt()
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
