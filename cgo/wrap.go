// Utility functions for CGO.
package cgo

import (
	"io"
	"unsafe"
)

// Identical to cgo.WrapReadWriter, except providing
// read-only services.
func WrapReader(r io.Reader, do_close bool) (unsafe.Pointer, error) {
	return wrapReadWriter(new_cookie(r, do_close))
}

// Identical to cgo.WrapReadWriter, except providing
// write-only services.
func WrapWriter(w io.Writer, do_close bool) (unsafe.Pointer, error) {
	return wrapReadWriter(new_cookie(w, do_close))
}

// Wraps an io.ReadWriter into a libc *FILE. If the reader
// supports the io.Seeker, or io.Closer interfaces, then
// the extra functionality will be provided through the
// generated *FILE pointer
//
// Even if the reader doesn't implement a Close()
// method, you must still call fclose() on the FILE
// pointer to clean up resources. If you don't want
// the go version to be closed, pass false for do_close.
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
func WrapReadWriter(rw io.ReadWriter, do_close bool) (unsafe.Pointer, error) {
	return wrapReadWriter(new_cookie(rw, do_close))
}
