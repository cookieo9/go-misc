package slice

import (
	"reflect"
	. "testing"
)

var insertTests = []struct {
	Slice  interface{}
	Index  int
	Item   interface{}
	Result interface{}
}{
	{[]int{}, 0, 1, []int{1}},
	{[]int{}, 1, 1, nil},
	{[]int{1, 2, 3}, 0, -99, []int{-99, 1, 2, 3}},
	{[]int{1, 2, 3}, 1, 50, []int{1, 50, 2, 3}},
	{[]int{1, 2, 3}, 3, 99, []int{1, 2, 3, 99}},
	{[]int{1, 2, 3}, 2, "q", nil},
	{5, 2, 1, nil},
}

type insertFunc func(interface{}, int, interface{}) interface{}

func runInsertTest(t *T, slice interface{}, index int, item interface{}, result interface{}, f insertFunc, name string) {
	defer func() {
		need_panic := result == nil
		if need_panic {
			if err := recover(); err == nil {
				t.Error("\tDid not get expected panic")
			} else {
				t.Logf("\tGot expected panic: %v", err)
			}
		} else {
			if err := recover(); err != nil {
				t.Errorf("\tUnexpected panic: %v", err)
			}
		}
	}()

	t.Logf("Running test %s(%#v, %d, %#v) = %#v", name, slice, index, item, result)
	r := f(slice, index, item)

	if !reflect.DeepEqual(r, result) {
		t.Errorf("\tUnexpected result: %v", r)
	} else {
		t.Logf("\tGot expected result: %v", r)
	}
}

func TestInsert(t *T) {
	for _, test := range insertTests {
		runInsertTest(t, test.Slice, test.Index, test.Item, test.Result, Insert, "Insert")
	}

	for _, test := range insertTests {
		runInsertTest(t, test.Slice, test.Index, test.Item, test.Result, InsertCopy, "InsertCopy")
	}
}

func BenchmarkInsert(b *B) {
	arr := []int{1, 2, 3, 4}

	for i := 0; i < b.N; i++ {
		_ = Insert(arr, 1, 2).([]int)
	}
}

func BenchmarkInsertCopy(b *B) {
	arr := []int{1, 2, 3, 4}

	for i := 0; i < b.N; i++ {
		_ = InsertCopy(arr, 1, 2).([]int)
	}
}

func BenchmarkInsertBuiltin(b *B) {
	arr := []int{1, 2, 3, 4}

	for i := 0; i < b.N; i++ {
		brr := append(arr[:2], arr[1:]...)
		brr[1] = 2
	}
}
