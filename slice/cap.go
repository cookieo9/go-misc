package slice

import (
	"reflect"
	"unsafe"
)

var (
	_ShrinkCapacityNotPointer = "ShrinkCapacity: not passed a pointer to a slice"
	_ShrinkCapacityIncrease   = "ShrinkCapacity: attempt to increase capacity"
)

// ShrinkCapacity reduces the capacity of the given slice in place.
// If the new capacity is smaller then the exising length, the
// length is also reduced to the new capacity.
//
// Will panic if not given a pointer to a slice, or if capacity
// would be increased by the change.
func ShrinkCapacity(slicePointer interface{}, capacity int) {
	pointerValue := reflect.ValueOf(slicePointer)

	t := pointerValue.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		panic(_ShrinkCapacityNotPointer)
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(pointerValue.Pointer()))

	// Prevent increasing capacity
	if sh.Cap < capacity {
		panic(_ShrinkCapacityIncrease)
	}

	// Enforce output len <= cap
	sh.Cap = capacity
	if sh.Len > sh.Cap {
		sh.Len = sh.Cap
	}
}

// HardSlice performs a slicing operation on the given array or slice
// but sets the capacity of the new slice to it's length instead of the
// remaining extra elements in the slice. The new slice is returned.
//
// Appending to the new slice should always result in a memory copy.
func HardSlice(source interface{}, begin, end int) interface{} {
	sourceValue := reflect.ValueOf(source)

	slice := sourceValue.Slice(begin, end)
	outputPtr := reflect.New(slice.Type())
	output := outputPtr.Elem()
	output.Set(slice)

	ShrinkCapacity(outputPtr.Interface(), output.Len())
	return output.Interface()
}
