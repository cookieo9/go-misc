package maps

import "reflect"

func GetKeys(mapval, sliceptr interface{}) {
	mv := reflect.ValueOf(mapval)
	sv := reflect.ValueOf(sliceptr).Elem()

	for _, key := range mv.MapKeys() {
		sv.Set(reflect.Append(sv, key))
	}
}

func GetVals(mapval, sliceptr interface{}) {
	mv := reflect.ValueOf(mapval)
	sv := reflect.ValueOf(sliceptr).Elem()

	for _, key := range mv.MapKeys() {
		sv.Set(reflect.Append(sv, mv.MapIndex(key)))
	}
}

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
