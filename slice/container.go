package slice

import (
	"errors"
	"fmt"
	"reflect"
)

func check_insert(slice interface{}, index int, item interface{}) (slice_val, item_val reflect.Value) {
	slice_val = reflect.ValueOf(slice)
	if slice_val.Kind() != reflect.Slice {
		panic(errors.New("insert: argument 1 must be slice"))
	}

	item_val = reflect.ValueOf(item)
	if item_val.Type() != slice_val.Type().Elem() {
		panic(fmt.Errorf("insert: %T can't be inserted into %T", item, slice))
	}

	if index < 0 || index > slice_val.Len() {
		panic(fmt.Errorf("insert: index (%d) out of range (0..%d)", index, slice_val.Len()))
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
	slice_val, item_val := check_insert(slice, index, item)

	if index == slice_val.Len() {
		return reflect.Append(slice_val, item_val).Interface()
	}

	begin := slice_val.Slice(0, index+1)
	end := slice_val.Slice(index, slice_val.Len())

	out := reflect.AppendSlice(begin, end)
	out.Index(index).Set(item_val)
	return out.Interface()
}

// Insert adds an element to a slice at a given index. Always allocates
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
	slice_val, item_val := check_insert(slice, index, item)
	begin := slice_val.Slice(0, index)
	end := slice_val.Slice(index, slice_val.Len())

	out := reflect.MakeSlice(slice_val.Type(), 0, slice_val.Len()+1)
	out = reflect.AppendSlice(reflect.Append(reflect.AppendSlice(out, begin), item_val), end)
	return out.Interface()
}

func check_delete(slice interface{}, index int) (slice_val reflect.Value) {
	slice_val = reflect.ValueOf(slice)
	if slice_val.Kind() != reflect.Slice {
		panic(errors.New("delete: argument 1 must be slice"))
	}

	if index < 0 || index > slice_val.Len()-1 {
		panic(fmt.Errorf("delete: idx (%d) out of range (0..%d)", index, slice_val.Len()-1))
	}
	return
}

// Deletes the element of the slice at a given index.
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
	slice_val := check_delete(slice, index)
	begin := slice_val.Slice(0, index)
	end := slice_val.Slice(index+1, slice_val.Len())
	return reflect.AppendSlice(begin, end).Interface()
}

// Deletes the element of the slice at a given index.
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
	slice_val := check_delete(slice, index)
	begin := slice_val.Slice(0, index)
	end := slice_val.Slice(index+1, slice_val.Len())
	tmp := reflect.MakeSlice(slice_val.Type(), 0, slice_val.Len()-1)
	return reflect.AppendSlice(reflect.AppendSlice(tmp, begin), end).Interface()
}
