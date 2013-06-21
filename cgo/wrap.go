package cgo

import (
	"io"
	"reflect"
	"unsafe"
)

import "C"

type cookie_t struct {
	val interface{}
}

func (c cookie_t) Closer() (io.Closer, bool) {
	x, ok := c.val.(io.Closer)
	return x, ok
}

func (c cookie_t) Writer() (io.Writer, bool) {
	x, ok := c.val.(io.Writer)
	return x, ok
}

func (c cookie_t) Seeker() (io.Seeker, bool) {
	x, ok := c.val.(io.Seeker)
	return x, ok
}

func (c cookie_t) Reader() (io.Reader, bool) {
	x, ok := c.val.(io.Reader)
	return x, ok
}

var cookie_registry = map[*cookie_t]struct{}{}

func new_cookie(val interface{}) *cookie_t {
	c := &cookie_t{val: val}
	cookie_registry[c] = struct{}{}
	return c
}

func free_cookie(c *cookie_t) {
	delete(cookie_registry, c)
}

func make_slice(ptr *C.char, size int) (buf []byte) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(unsafe.Pointer(ptr))
	sh.Len, sh.Cap = size, size
	return
}

func WrapReader(r io.Reader) (unsafe.Pointer, error) {
	return wrapReadWriter(r)
}

func WrapWriter(w io.Writer) (unsafe.Pointer, error) {
	return wrapReadWriter(w)
}

func WrapReadWriter(rw io.ReadWriter) (unsafe.Pointer, error) {
	return wrapReadWriter(rw)
}
