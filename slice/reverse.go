package slice

import (
	"reflect"
)

// Reverse reverses the order of the items
// of the given slice (or array).
// Will panic if not passed a slice or array.
func Reverse(slice interface{}) {
	v := reflect.ValueOf(slice)
	l := v.Len()
	for i := 0; i < l/2; i++ {
		a, b := v.Index(i), v.Index(l-1-i)
		t := a.Interface()
		a.Set(b)
		b.Set(reflect.ValueOf(t))
	}
}
