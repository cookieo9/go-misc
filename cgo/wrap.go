// Utility functions for CGO.
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

// Identical to cgo.WrapReadWriter, except providing
// read-only services.
func WrapReader(r io.Reader) (unsafe.Pointer, error) {
	return wrapReadWriter(r)
}

// Identical to cgo.WrapReadWriter, except providing
// write-only services.
func WrapWriter(w io.Writer) (unsafe.Pointer, error) {
	return wrapReadWriter(w)
}

// Wraps an io.ReadWriter into a libc *FILE. If the reader
// supports the io.Seeker, or io.Closer interfaces, then
// the extra functionality will be provided through the
// generated *FILE pointer
//
// Even if the reader doesn't implement a Close()
// method, you must still call fclose() on the FILE
// pointer to clean up resources. If you don't want
// the go version to be closed, you'll have to provide
// a version of the object without a Close method, or
// has one which does nothing.
//
// If you write to the *FILE, the write may be buffered
// by libc, so remember to call fflush, or fclose,
// to see changes in go.
//
// This function returns an unsafe.Pointer, because it must
// be cast to a *C.FILE by the code using it. Returning a
// *C.FILE in this package would require a user of this
// package to convert it from a "*cgo.C.FILE" to a
// "mypkg.C.FILE" with an unsafe.Pointer as a stepping stone.
func WrapReadWriter(rw io.ReadWriter) (unsafe.Pointer, error) {
	return wrapReadWriter(rw)
}
