package slice

import (
	"reflect"
	"unsafe"
)

// ShrinkCapacity reduces the capacity of the given slice in place.
// If the new capacity is smaller then the exising length, the
// length is also reduced to the new capacity.
//
// Will panic if not given a pointer to a slice, or if capacity
// would be increased by the change.
func ShrinkCapacity(slicePointer interface{}, capacity int) {
	pointerValue := reflect.ValueOf(slicePointer)

	if pointerValue.Kind() != reflect.Ptr || pointerValue.Elem().Kind() != reflect.Slice {
		panic("ShrinkCapacity: not passed a pointer to a slice")
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(pointerValue.Pointer()))
	l, c := sh.Len, sh.Cap

	// Prevent increasing capacity
	if c < capacity {
		panic("ShrinkCapacity: attempt to increase capacity")
	}

	// Enforce output len <= cap
	c = capacity
	if l > c {
		l = c
	}

	sh.Len, sh.Cap = l, c
}

// HardSlice performs a slicing operation on the given array or slice
// but sets the capacity of the new slice to it's length instead of the
// remaining extra elements in the slice. The new slice is returned.
//
// Appending to the new slice should always result in a memory copy.
func HardSlice(source interface{}, begin, end int) interface{} {
	sourceValue := reflect.ValueOf(source)
	slice := sourceValue.Slice(begin, end)
	output := reflect.New(slice.Type()).Elem()
	output.Set(slice)
	ShrinkCapacity(output.Addr().Interface(), output.Len())
	return output.Interface()
}
