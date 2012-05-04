/*
Package slice provides simplified sorting for generic slices, as well as building heaps from those slices.

*/
package slice

import (
	"container/heap"
	"reflect"
	"sort"
)

type sliceWrap struct {
	reflect.Value
	t reflect.Value
}

func (sw *sliceWrap) Swap(i, j int) {
	ival, jval := sw.Index(i), sw.Index(j)
	sw.t.Set(ival)
	ival.Set(jval)
	jval.Set(sw.t)
}

func (sw *sliceWrap) Push(x interface{}) {
	sw.Set(reflect.Append(sw.Value, reflect.ValueOf(x)))
}

func (sw *sliceWrap) Pop() (x interface{}) {
	l := sw.Len() - 1
	x = sw.Index(l).Interface()
	sw.Set(sw.Slice(0, l))
	return
}

func wrapSlice(slice interface{}) sliceWrap {
	v := reflect.Indirect(reflect.ValueOf(slice))
	t := reflect.New(v.Type().Elem()).Elem()
	return sliceWrap{
		Value: v,
		t:     t,
	}
}

type untyped struct {
	sliceWrap
	f func(a, b interface{}) bool
}

func (u untyped) Less(i, j int) bool {
	ival, jval := u.Index(i).Interface(), u.Index(j).Interface()
	return u.f(ival, jval)
}

type typed struct {
	sliceWrap
	f reflect.Value
}

func (t typed) Less(i, j int) bool {
	args := []reflect.Value{t.Index(i), t.Index(j)}
	return t.f.Call(args)[0].Bool()
}

// A slice that additionally supports slice.Interface can be trivially wrapped
// by the routines in this package and thus be used as a sort.Interface, and
// possibly a heap.Interface
type Interface interface {
	Less(i, j int) bool
}

type wrapped struct {
	Interface
	sliceWrap
}

// WrapTyped wraps a given slice with a typed comparator function to
// make a value that satisfies the sort.Interface and heap.Interface
// interfaces.
//
// The comparator is a function used to compare two elements in
// the slice. Its signature must be func(a,b T) bool where T is
// the element type of the given slice. The comparator should
// return true if the value in a should be sorted before the
// value in b.
//
//	slice := []int8{3,4,1,5,2}
//	ascending := WrapTyped(slice,func(a,b int8)bool {return a<=b})
//	decending := WrapTyped(slice,func(a,b int8)bool {return a>=b})
//	
//	sort.Sort(ascending) // slice will be sorted in ascending order
//	sort.Sort(decending) // slice will be sorted in decending order
//
// The slice parameter can be either a slice, or a pointer to a
// slice. If given a pointer to a slice, the Push and Pop methods
// will work properly according to the heap.Interface interface.
// Otherwise, attempting to call those methods (eg: by the heap package
// functions) on a wrapped non-pointer will result in a panic.
func WrapTyped(slice, comparator interface{}) heap.Interface {
	return &typed{wrapSlice(slice), reflect.ValueOf(comparator)}
}

// WrapUntyped wraps a given slice with an untyped comparator function
// to make a value that satifies the sort.Interface and heap.Interface
// interfaces.
//
// The meanings of the slice and comparator parameters are the same as
// in WrapTyped, only the signature of the comparator function is
// different.
//
// Example:
//
//	slice := []int8{3,4,1,5,2}
//	ascending := WrapUntyped(slice, func(a,b interface{}) bool {return a.(int8) <= b.(int8)})
//	decending := WrapUntyped(slice, func(a,b interface{}) bool {return a.(int8) >= b.(int8)})
//
//	sort.Sort(ascending) // slice will be sorted in ascending order
//	sort.Sort(decending) // slice will be sorted in decending order
func WrapUntyped(slice interface{}, comparator func(a, b interface{}) bool) heap.Interface {
	return &untyped{wrapSlice(slice), comparator}
}

// WrapInterface wraps a slice that satisfies the Interface interface
// to make a value that satifies the sort.Interface and heap.Interface
// interfaces.
//
// If a pointer is passed as the value of slice, then like with WrapTyped
// and WrapUntyped, the resulting value can be used successfully in
// heap operations.
//
// Example:
//
//	type int8asc []int8
//	func (s int8asc) Less(i,j int) bool { return s[i] <= s[j] }
//
//	slice := []int8{3,4,1,5,2}
//	ascending := WrapInterface(int8asc(slice))
//	h := WrapInterface(&int8asc(slice))
//
//	sort.Sort(ascending)
func WrapInterface(slice Interface) heap.Interface {
	return &wrapped{slice, wrapSlice(slice)}
}

// Wrap is a utility method that calls one of the more specific Wrap*
// functions. This utility is less type-safe than the other functions,
// so expect panics.
func Wrap(slice interface{}, args ...interface{}) heap.Interface {
	if x, ok := slice.(Interface); ok {
		return WrapInterface(x)
	}
	if f, ok := args[0].(func(a, b interface{}) bool); ok {
		return WrapUntyped(slice, f)
	}
	return WrapTyped(slice, args[0])
}

// SortTyped passes its arguments unaltered to slice.WrapTyped, and then
// calls sort.Sort on the result.
func SortTyped(slice interface{}, comparator interface{}) {
	sort.Sort(WrapTyped(slice, comparator))
}

// SortUntyped passes its arguments unaltered to slice.WrapUntyped, and then
// calls sort.Sort on the result.
func SortUntyped(slice interface{}, comparator func(a, b interface{}) bool) {
	sort.Sort(WrapUntyped(slice, comparator))
}

// SortInterface passes its arguments unaltered to slice.WrapInterface, and then
// calls sort.Sort on the result.
func SortInterface(slice Interface) {
	sort.Sort(WrapInterface(slice))
}

// Sort passes its arguments unaltered to slice.Wrap, and then
// calls sort.Sort on the result.
func Sort(slice interface{}, args ...interface{}) {
	sort.Sort(Wrap(slice, args...))
}
