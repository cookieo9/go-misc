package slice

import (
	"reflect"
	. "testing"
)

func ExampleHardSlice() {
	base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	a := base[0:5]    // len(a)=5 cap(a)=15
	a = append(a, -6) // modifies base

	a2 := HardSlice(a, 0, 5).([]int) // len(a2)=5 cap(a2)=5
	a2 = append(a2, -7)              // does not modify base
}

func ExampleShrinkCapacity() {
	base := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a := base[0:5]         // len(a)=5 cap(a)=15
	ShrinkCapacity(&a, 10) // len(a)=5 cap(a)=10
}

func TestShrinkCapacity(t *T) {
	t.Parallel()
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

	t.Log("--- Initial State")
	dumpAll()

	t.Log("--- Using ShrinkCapacity")
	ShrinkCapacity(&a, 5)
	ShrinkCapacity(&b, 5)
	dumpAll()

	if cap(a) != 5 {
		t.Errorf("cap(a) should be %d, but got %d instead", 5, cap(a))
	}
	if cap(b) != 5 {
		t.Errorf("cap(b) should be %d, but got %d instead", 5, cap(b))
	}

	t.Log("--- Appending")
	a = append(a, -999)
	b = append(b, -999)
	dumpAll()

	if !reflect.DeepEqual(base, original) {
		t.Error("Base slice modified!")
	}
}

func TestHardSlice(t *T) {
	t.Parallel()
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
		t.Errorf("cap(a) should be %d, but got %d instead", 5, cap(a))
	}
	if cap(b) != 5 {
		t.Errorf("cap(b) should be %d, but got %d instead", 5, cap(b))
	}
	if &base[0] != &a[0] {
		t.Error("HardSlice created a new array backing a when it shouldn't have")
	}
	if &base[5] != &b[0] {
		t.Error("HardSlice created a new array backing b when it shouldn't have")
	}
}

var ShrinkCapacity_GoodTestCases = []struct {
	Slice    []int
	Cap      int
	Expected []int
}{
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5, []int{1, 2, 3, 4, 5}},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0, []int{}},
	{[]int{}, 0, []int{}},
}

func TestShrinkCapacity_Good(t *T) {
	t.Parallel()
	for _, test := range ShrinkCapacity_GoodTestCases {
		t.Logf("ShrinkCapacity(%v,%d) = %v", test.Slice, test.Cap, test.Expected)

		a := test.Slice
		ShrinkCapacity(&a, test.Cap)
		if cap(a) != test.Cap {
			t.Errorf("Expected capacity of %d, got %d", test.Cap, cap(a))
		}
		if !reflect.DeepEqual(a, test.Expected) {
			t.Errorf("Expected %v, got %v", test.Expected, a)
		}
	}
}

var ShrinkCapacity_PanicTestCases = []struct {
	Input interface{}
	Cap   int
	Panic interface{}
}{
	{&[]int{1, 2, 3}, 5, _ShrinkCapacityIncrease},
	{[]int{1, 2, 3}, 2, _ShrinkCapacityNotPointer},
	{&[]int{}, 0, nil},
}

func TestShrinkCapacity_Panic(t *T) {
	t.Parallel()
	for _, test := range ShrinkCapacity_PanicTestCases {
		t.Logf("ShrinkCapacity(%v,%d) => panic(%v)", test.Input, test.Cap, test.Panic)

		var pval interface{}
		func() {
			defer func() {
				pval = recover()
			}()
			ShrinkCapacity(test.Input, test.Cap)
		}()

		if !reflect.DeepEqual(pval, test.Panic) {
			t.Errorf("Expected panic(%v), got panic(%v)", test.Panic, pval)
		}
	}
}
