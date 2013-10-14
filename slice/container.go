package slice

import (
	"errors"
	"fmt"
	"reflect"
)

func checkInsert(slice interface{}, index int, item interface{}) (sliceVal, itemVal reflect.Value) {
	sliceVal = reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		panic(errors.New("insert: argument 1 must be slice"))
	}

	itemVal = reflect.ValueOf(item)
	if itemVal.Type() != sliceVal.Type().Elem() {
		panic(fmt.Errorf("insert: %T can't be inserted into %T", item, slice))
	}

	if index < 0 || index > sliceVal.Len() {
		panic(fmt.Errorf("insert: index (%d) out of range (0..%d)", index, sliceVal.Len()))
	}

	return
}

// Insert adds an element to a slice at a given index. If len(slice) < cap(slice),
// then the original slice is expanded, otherwise a new slice is allocated.
// This behaviour is analogous to builtin.append().
//
// The returned slice will reference the array where the insert occured, and
// will have apropriate values for cap and len. The returned slice will be
// the same type as the given slice, but stored in an interface{}.
//
// Uses reflection to perform this action on slices of any type.
// Will panic if:
//	- slice argument is not a slice type
//	- slice argument's element type doesn't match item's type
//	- idx is not in the range 0 .. len(slice)
//
// Equivalent to:
//	slice = append(slice[:idx+1],slice[idx:]...)
//	slice[idx] = item
//
// NOTE: Other slices referencing the same underlying array (including the argument
// slice) may have their contents altered, if the change is made in place. To
// guarantee that a new array will always be created use the InsertCopy function.
func Insert(slice interface{}, index int, item interface{}) interface{} {
	sliceVal, itemVal := checkInsert(slice, index, item)

	if index == sliceVal.Len() {
		return reflect.Append(sliceVal, itemVal).Interface()
	}

	begin := sliceVal.Slice(0, index+1)
	end := sliceVal.Slice(index, sliceVal.Len())

	out := reflect.AppendSlice(begin, end)
	out.Index(index).Set(itemVal)
	return out.Interface()
}

// InsertCopy adds an element to a slice at a given index. Always allocates
// a new slice, and never modifies the original memory.
//
// The returned slice will be the same type as the given slice, but
// stored in an interface{}.
//
// Uses reflection to perform this action on slices of any type.
// Will panic if:
//	- slice argument is not a slice type
//	- slice argument's element type doesn't match item's type
//	- idx is not in the range 0 .. len(slice)
//
// Equivalent to:
// 	begin, end := slice[:idx], slice[idx:]
// 	slice = append(append(append(make([]T,0,len(slice)+1), begin...),item),end...)
func InsertCopy(slice interface{}, index int, item interface{}) interface{} {
	sliceVal, itemVal := checkInsert(slice, index, item)
	begin := sliceVal.Slice(0, index)
	end := sliceVal.Slice(index, sliceVal.Len())

	out := reflect.MakeSlice(sliceVal.Type(), 0, sliceVal.Len()+1)
	out = reflect.AppendSlice(reflect.Append(reflect.AppendSlice(out, begin), itemVal), end)
	return out.Interface()
}

func checkDelete(slice interface{}, index int) (sliceVal reflect.Value) {
	sliceVal = reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		panic(errors.New("delete: argument 1 must be slice"))
	}

	if index < 0 || index > sliceVal.Len()-1 {
		panic(fmt.Errorf("delete: idx (%d) out of range (0..%d)", index, sliceVal.Len()-1))
	}
	return
}

// Delete removes the element of the slice at a given index.
// The returned slice is a reference to the modified original array.
//
// The returned slice will be the same type as the given slice, but
// stored in an interface{}.
//
// Uses reflection to perform this action on slices of any type.
// Will panic if:
//	- slice argument is not a slice type
//	- idx is not in the range 0 .. len(slice)-1
//
// Equivalent to:
//	slice = append(slice[:idx],slice[idx+1:]...)
//
// NOTE: Any other slices (including the argument slice)
// referencing the original underlying array may have their contents
// altered by this call. This behaviour is consistent with append().
func Delete(slice interface{}, index int) interface{} {
	sliceVal := checkDelete(slice, index)
	begin := sliceVal.Slice(0, index)
	end := sliceVal.Slice(index+1, sliceVal.Len())
	return reflect.AppendSlice(begin, end).Interface()
}

// DeleteCopy removes the element of the slice at a given index.
// The returned slice references a new array, the original
// array and all referencing slices are un-altered.
//
// The returned slice will be the same type as the given slice, but
// stored in an interface{}.
//
// Uses reflection to perform this action on slices of any type.
// Will panic if:
//	- slice argument is not a slice type
//	- idx is not in the range 0 .. len(slice)-1
//
// Equivalent to:
//	slice = append(append(make([]T,0,len(slice)-1),slice[:idx]...),slice[idx+1]...)
func DeleteCopy(slice interface{}, index int) interface{} {
	sliceVal := checkDelete(slice, index)
	begin := sliceVal.Slice(0, index)
	end := sliceVal.Slice(index+1, sliceVal.Len())
	tmp := reflect.MakeSlice(sliceVal.Type(), 0, sliceVal.Len()-1)
	return reflect.AppendSlice(reflect.AppendSlice(tmp, begin), end).Interface()
}
