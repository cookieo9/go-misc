package slice

import (
	"fmt"
	"os"
)

func ExampleSort() {
	a := []int{3, 1, 4, 1, 5, 1, 9, 2, 6, 5, 3, 5}

	Sort(a, func(a, b int) bool { return a <= b })
	fmt.Fprintln(os.Stdout, a)
}
