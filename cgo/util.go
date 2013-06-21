package cgo

import (
	"reflect"
	"unsafe"
)

import "C"

func make_slice(ptr *C.char, size int) (buf []byte) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(unsafe.Pointer(ptr))
	sh.Len, sh.Cap = size, size
	return
}
