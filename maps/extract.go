// Utility functions to extract keys, values, or both from maps.
package maps

import "reflect"

// Pulls the keys out of a map into a slice of the
// appropriate type.
func GetKeys(mapval, sliceptr interface{}) {
	mv := reflect.ValueOf(mapval)
	sv := reflect.ValueOf(sliceptr).Elem()

	for _, key := range mv.MapKeys() {
		sv.Set(reflect.Append(sv, key))
	}
}

// Pulls the values out of a map into a slice of the
// appropriate type.
func GetVals(mapval, sliceptr interface{}) {
	mv := reflect.ValueOf(mapval)
	sv := reflect.ValueOf(sliceptr).Elem()

	for _, key := range mv.MapKeys() {
		sv.Set(reflect.Append(sv, mv.MapIndex(key)))
	}
}

// Pull the key/value pairs out of a map into a slice
// of the type struct { Key <KeyType>; Val <ValType> }.
func GetPairs(mapval, sliceptr interface{}) {
	mv := reflect.ValueOf(mapval)
	sv := reflect.ValueOf(sliceptr).Elem()
	et := sv.Type().Elem()
	tmp := reflect.New(et).Elem()
	kf := tmp.Field(0)
	vf := tmp.Field(1)

	for _, key := range mv.MapKeys() {
		kf.Set(key)
		vf.Set(mv.MapIndex(key))
		sv.Set(reflect.Append(sv, tmp))
	}
}
