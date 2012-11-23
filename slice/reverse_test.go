package slice

import (
	"reflect"
	. "testing"
)

var tests = [][2][]int{
	{{3}, {3}},
	{{3, 1}, {1, 3}},
	{{3, 1, 4}, {4, 1, 3}},
	{{3, 1, 4, 1, 5, 9}, {9, 5, 1, 4, 1, 3}},
	{{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}, {5, 3, 5, 6, 2, 9, 5, 1, 4, 1, 3}},
}

func TestReverse(t *T) {
	for _, test := range tests {
		input, expected := test[0], test[1]
		output := make([]int, len(input))
		copy(output, input)
		Reverse(output)

		t.Log("Reversing", input)
		if !reflect.DeepEqual(output, expected) {
			t.Error("Expected", expected, "got", output)
		}
	}
}

const RevArraySize = 1e4

func BenchmarkReverse(b *B) {
	a := make([]int, RevArraySize)
	for i := 0; i < b.N; i++ {
		Reverse(a)
	}
}

func BenchmarkReverseRaw(b *B) {
	a := make([]int, RevArraySize)
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(a)/2; i++ {
			i2 := len(a) - 1 - i
			a[i], a[i2] = a[i2], a[i]
		}
	}
}
