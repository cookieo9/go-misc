package slice

import (
	"container/heap"
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

var profile = flag.Bool("test.profile", false, "Run profiling tests")

func testSlice() []int {
	return []int{4, 2, 1, 5, 3}
}

type sortInterface []int

func (si sortInterface) Less(i, j int) bool {
	return si[i] > si[j]
}

func TestSortInterface(t *testing.T) {
	a := testSlice()
	b := []int{5, 4, 3, 2, 1}
	c := sortInterface(a)
	Sort(c)

	if err := checkSlice(a, b); err != nil {
		t.Fatal(err)
	}
}

func TestSortTyped(t *testing.T) {
	a := testSlice()
	b := []int{5, 4, 3, 2, 1}
	f := func(a, b int) bool { return a > b }
	Sort(a, f)

	if err := checkSlice(a, b); err != nil {
		t.Fatal(err)
	}
}

func TestSortUntyped(t *testing.T) {
	a := testSlice()
	b := []int{5, 4, 3, 2, 1}
	f := func(a, b interface{}) bool { return a.(int) > b.(int) }

	Sort(a, f)

	if err := checkSlice(a, b); err != nil {
		t.Fatal(err)
	}
}

func checkSlice(given, expected interface{}) error {
	if !reflect.DeepEqual(given, expected) {
		return fmt.Errorf("Expected %v got %v", expected, given)
	}
	return nil
}

func genIntArray(n int) []int {
	rand.Seed(int64(n))
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Int()
	}
	return a
}

func BenchmarkSortInterface(b *testing.B) {
	a := genIntArray(b.N)
	c := sortInterface(a)
	b.ResetTimer()
	Sort(c)
}

func BenchmarkSortTyped(b *testing.B) {
	a := genIntArray(b.N)
	b.ResetTimer()
	Sort(a, func(a, b int) bool { return a < b })
}

func BenchmarkSortUntyped(b *testing.B) {
	a := genIntArray(b.N)
	b.ResetTimer()
	Sort(a, func(a, b interface{}) bool { return a.(int) < b.(int) })
}

func BenchmarkBuiltinSort(b *testing.B) {
	a := sort.IntSlice(genIntArray(b.N))
	b.ResetTimer()
	sort.Sort(a)
}

const ProfileCount = 1e6

func TestProfile(t *testing.T) {
	if *profile {
		a := genIntArray(ProfileCount)
		Sort(a, func(a, b int) bool { return a > b })
	}
}

const HeapCount = 1e4

type heapInterface []int

func (hi heapInterface) Less(i, j int) bool {
	return hi[i] < hi[j]
}

func TestHeap(t *testing.T) {
	a := heapInterface(genIntArray(HeapCount))
	h := Wrap(&a)

	heap.Init(h)

	x := a[0]
	for len(a) > 0 {
		y := heap.Pop(h).(int)
		if x > y {
			t.Fatalf("Wasn't expecting a value (%d) less than %d", y, x)
		}
		x = y
	}
}

func BenchmarkHeap(b *testing.B) {
	a := heapInterface(genIntArray(b.N))

	h := Wrap(&a)
	c := make([]int, b.N)

	b.ResetTimer()
	heap.Init(h)

	for i := range c {
		c[i] = heap.Pop(h).(int)
	}
}
