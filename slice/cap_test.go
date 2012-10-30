package slice

import (
	"reflect"
	. "testing"
	"unsafe"
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

	t.Log("### Initial State ###")
	dumpAll()

	t.Log("### Using ShrinkCapacity ###")
	ShrinkCapacity(&a, 5)
	ShrinkCapacity(&b, 5)
	dumpAll()

	if cap(a) != 5 {
		t.Errorf("cap(a) should be %d, but got %d instead", 5, cap(a))
	}
	if cap(b) != 5 {
		t.Errorf("cap(b) should be %d, but got %d instead", 5, cap(b))
	}

	t.Log("### Appending ###")
	a = append(a, -999)
	b = append(b, -999)
	dumpAll()

	if !reflect.DeepEqual(base, original) {
		t.Error("Base slice modified!")
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
	{&[]int{1, 2, 3}, -5, _ShrinkCapacityNegative},
	{[]int{1, 2, 3}, 2, _ShrinkCapacityInvalidType},
	{&[]int{}, 0, nil},
}

func TestShrinkCapacity_Panic(t *T) {
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

var HardSlice_GoodTestCases = []struct {
	Input      []int
	Begin, End int
	Expected   []int
}{
	{[]int{1, 2, 3, 4, 5}, 0, 5, []int{1, 2, 3, 4, 5}},
	{[]int{1, 2, 3, 4, 5}, 2, 4, []int{3, 4}},
	{[]int{1, 2, 3, 4, 5}, 4, 4, []int{}},
	{[]int{1, 2, 3, 4, 5}, 0, 0, []int{}},
	{[]int{}, 0, 0, []int{}},
}

func TestHardSlice_Good(t *T) {
	for _, test := range HardSlice_GoodTestCases {
		t.Logf("HardSlice(%v,%d,%d) = %v", test.Input, test.Begin, test.End, test.Expected)

		expected_len_cap := test.End - test.Begin
		output := HardSlice(test.Input, test.Begin, test.End).([]int)

		if cap(output) != expected_len_cap {
			t.Errorf("Expected capacity of %d, got %d", expected_len_cap, cap(output))
		}
		if len(output) != expected_len_cap {
			t.Errorf("Expected length of %d, got %d", expected_len_cap, len(output))
		}

		if !reflect.DeepEqual(output, test.Expected) {
			t.Errorf("Expected %v, got %v", test.Expected, output)
		}
	}
}

var HardSlice_PanicTestCases = []struct {
	Input      interface{}
	Begin, End int
	Panic      interface{}
}{
	{5, 0, 0, _HardSliceInvalidType},
	{&[]int{1, 2, 3, 4, 5}, 0, 0, _HardSliceInvalidType},
	{[]int{1, 2, 3, 4, 5}, 0, 0, nil},
	{[]int{1, 2, 3, 4, 5}, -1, 0, _HardSliceIndexOutOfBounds},
	{[]int{1, 2, 3, 4, 5}, 0, 6, _HardSliceIndexOutOfBounds},
	{[]int{1, 2, 3, 4, 5}, 3, 2, _HardSliceIndexOutOfBounds},
}

func TestHardSlice_Panic(t *T) {
	for _, test := range HardSlice_PanicTestCases {
		t.Logf("HardSlice(%v,%d,%d) => panic(%v)", test.Input, test.Begin, test.End, test.Panic)
		var panicVal interface{}

		func() {
			defer func() {
				panicVal = recover()
			}()
			HardSlice(test.Input, test.Begin, test.End)
		}()

		if !reflect.DeepEqual(panicVal, test.Panic) {
			t.Errorf("Expected panic(%v), got panic(%v)", test.Panic, panicVal)
		}
	}
}

func BenchmarkShrinkCapacity(b *B) {
	base := make([]int, 10)
	for i := 0; i < b.N; i++ {
		a := base[:5]
		ShrinkCapacity(&a, len(a))
	}
}

func shrink_cap_int(slice *[]int, capacity int) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(slice))
	if capacity > sh.Cap {
		panic(_ShrinkCapacityIncrease)
	}

	sh.Cap = capacity
	if sh.Len > sh.Cap {
		sh.Len = sh.Cap
	}
}

func BenchmarkShrinkCapacity_Fixed(b *B) {
	base := make([]int, 10)
	for i := 0; i < b.N; i++ {
		a := base[:5]
		shrink_cap_int(&a, len(a))
	}
}

func BenchmarkHardSlice(b *B) {
	base := make([]int, 10)
	for i := 0; i < b.N; i++ {
		_ = HardSlice(base, 0, 5).([]int)
	}
}

func hard_slice_int(slice []int, begin, end int) (subslice []int) {
	subslice = slice[begin:end]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&subslice))
	sh.Cap = sh.Len
	return
}

func BenchmarkHardSlice_Fixed(b *B) {
	base := make([]int, 10)
	for i := 0; i < b.N; i++ {
		_ = hard_slice_int(base, 0, 5)
	}
}
