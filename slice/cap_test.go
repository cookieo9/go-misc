package slice

import (
	"reflect"
	. "testing"
)

func TestShrinkCapacity(t *T) {
	dumpSlice := func(name string, slice []int) {
		t.Logf("%s: addr(%p) len(%d) cap(%d) - %v", name, &slice[0], len(slice), cap(slice), slice)
	}

	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	base := append([]int{}, original...)

	a := base[0:5]
	b := base[5:10]

	dumpAll := func() {
		dumpSlice("original", original)
		dumpSlice("base", base)
		dumpSlice("a", a)
		dumpSlice("b", b)
	}

	t.Log("--- Initial State ---")
	dumpAll()

	t.Log("--- Using ShrinkCapacity ---")
	ShrinkCapacity(&a, 5)
	ShrinkCapacity(&b, 5)
	dumpAll()

	if cap(a) != 5 {
		t.Error("cap(a) should be %d, but got %d instead", 5, cap(a))
	}
	if cap(b) != 5 {
		t.Error("cap(b) should be %d, but got %d instead", 5, cap(b))
	}

	t.Log("--- Appending ---")
	a = append(a, -999)
	b = append(b, -999)
	dumpAll()

	if !reflect.DeepEqual(base, original) {
		t.Error("base slice modified!")
	}
}

func TestHardSlice(t *T) {
	dumpSlice := func(name string, slice []int) {
		t.Logf("%s: addr(%p) len(%d) cap(%d) - %v", name, &slice[0], len(slice), cap(slice), slice)
	}

	base := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a := HardSlice(base, 0, 5).([]int)
	b := HardSlice(base, 5, 10).([]int)

	dumpSlice("base", base)
	dumpSlice("a", a)
	dumpSlice("b", b)

	if cap(a) != 5 {
		t.Error("cap(a) should be %d, but got %d instead", 5, cap(a))
	}
	if cap(b) != 5 {
		t.Error("cap(b) should be %d, but got %d instead", 5, cap(b))
	}
	if &base[0] != &a[0] {
		t.Error("HardSlice created a new array backing a when it shouldn't have")
	}
	if &base[5] != &b[0] {
		t.Error("HardSlice created a new array backing b when it shouldn't have")
	}
}

func ExampleHardSlice() {
	base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	a := base[0:5]    // len(a)=5 cap(a)=15
	a = append(a, -6) // modifies base

	a2 := HardSlice(a, 0, 5).([]int) // len(a2)=5 cap(a2)=5
	a2 = append(a2, -7)              // does not modify base
}
