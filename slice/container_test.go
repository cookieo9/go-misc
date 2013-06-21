package slice

import (
	"reflect"
	. "testing"
)

type insertTest struct {
	Slice  interface{}
	Index  int
	Item   interface{}
	Result interface{}
}

func insertTests() []insertTest {
	return []insertTest{
		{[]int{}, 0, 1, []int{1}},
		{[]int{}, 1, 1, nil},
		{[]int{1, 2, 3}, 0, -99, []int{-99, 1, 2, 3}},
		{[]int{1, 2, 3}, 1, 50, []int{1, 50, 2, 3}},
		{[]int{1, 2, 3}, 3, 99, []int{1, 2, 3, 99}},
		{[]int{1, 2, 3}, 2, "q", nil},
		{5, 2, 1, nil},
	}
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
	for _, test := range insertTests() {
		runInsertTest(t, test.Slice, test.Index, test.Item, test.Result, Insert, "Insert")
	}

	for _, test := range insertTests() {
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

type deleteTest struct {
	InitialSlice   interface{}
	Element        int
	ExpectedResult interface{}
}

var deleteTests = func() []deleteTest {
	return []deleteTest{
		{[]int{1, 4, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4, 5}},
		{[]int{}, 0, nil},
		{[]int{42}, 0, []int{}},
		{[]int{42}, 1, nil},
		{"foo", 0, nil},
	}
}

type deleteFunc func(interface{}, int) interface{}

func runDeleteTest(t *T, slice interface{}, element int, expect interface{}, f func(interface{}, int) interface{}, name string) {
	defer func() {
		need_panic := expect == nil
		if need_panic {
			if msg := recover(); msg == nil {
				t.Error("\tDid not get panic when expected")
			} else {
				t.Logf("\tRecieved panic when expected: %v", msg)
			}
		} else {
			if msg := recover(); msg != nil {
				t.Errorf("\tRecieved unexpected panic: %v", msg)
			}
		}
	}()

	t.Logf("Running test on %s(%v,%d) => expecting %v", name, slice, element, expect)
	result := f(slice, element)
	if !reflect.DeepEqual(expect, result) {
		t.Error("\tUnexpected result: ", result)
	} else {
		t.Log("\tExpected result: ", result)
	}
}

func TestDelete(t *T) {
	for _, test := range deleteTests() {
		runDeleteTest(t, test.InitialSlice, test.Element, test.ExpectedResult, Delete, "Delete")
	}

	for _, test := range deleteTests() {
		runDeleteTest(t, test.InitialSlice, test.Element, test.ExpectedResult, DeleteCopy, "DeleteCopy")
	}
}
