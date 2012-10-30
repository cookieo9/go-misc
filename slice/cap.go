package slice

import (
	"reflect"
	"unsafe"
)

var (
	_ShrinkCapacityInvalidType = "slice.ShrinkCapacity: argument 1 not pointer to a slice"
	_ShrinkCapacityIncrease    = "slice.ShrinkCapacity: attempt to increase capacity"
	_ShrinkCapacityNegative    = "slice.ShrinkCapacity: negative target capacity"
)

// ShrinkCapacity reduces the capacity of the given slice in place.
// If the new capacity is smaller then the exising length, the
// length is also reduced to the new capacity.
//
// Will panic if not given a pointer to a slice, or if capacity
// would be increased by the change.
func ShrinkCapacity(slicePointer interface{}, capacity int) {
	pointerValue := reflect.ValueOf(slicePointer)

	if t := pointerValue.Type(); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		panic(_ShrinkCapacityInvalidType)
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(pointerValue.Pointer()))

	// Prevent increasing capacity
	if sh.Cap < capacity {
		panic(_ShrinkCapacityIncrease)
	}

	if capacity < 0 {
		panic(_ShrinkCapacityNegative)
	}

	// Enforce output len <= cap
	sh.Cap = capacity
	if sh.Len > sh.Cap {
		sh.Len = sh.Cap
	}
}

var (
	_HardSliceIndexOutOfBounds = "slice.HardSlice: slice index out of bounds"
	_HardSliceInvalidType      = "slice.HardSlice: argument 1 not a slice"
)

// HardSlice performs a slicing operation on the given slice
// but sets the capacity of the new slice to it's length instead of the
// remaining extra elements in the slice. The new slice is returned.
//
// Appending to the new slice should always result in a memory copy.
func HardSlice(source interface{}, begin, end int) interface{} {
	sourceValue := reflect.ValueOf(source)
	sourceType := sourceValue.Type()

	if sourceType.Kind() != reflect.Slice {
		panic(_HardSliceInvalidType)
	}

	length := sourceValue.Len()
	if valid := 0 <= begin && begin <= end && end <= length; !valid {
		panic(_HardSliceIndexOutOfBounds)
	}

	outputPtr := reflect.New(sourceType)
	output := outputPtr.Elem()

	slice := sourceValue.Slice(begin, end)
	output.Set(slice)

	ShrinkCapacity(outputPtr.Interface(), end-begin)
	return output.Interface()
}
