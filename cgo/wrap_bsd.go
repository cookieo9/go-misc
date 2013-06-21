// +build darwin freebsd netbsd openbsd

package cgo

//
// #include <stdio.h>
// #include <sys/errno.h>
// #include "indirect_bsd.h"
//
// static void seterr ( int e ) { errno = e; }
//
import "C"

import (
	"io"
	"os"
	"unsafe"
)

//export reader
func reader(cookie_ptr unsafe.Pointer, buf *C.char, size C.int) C.int {
	cookie := (*cookie_t)(cookie_ptr)

	rdr, ok := cookie.Reader()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}

	buffer := make_slice(buf, int(size))
	n, err := rdr.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return C.int(n)
		}
		C.seterr(C.EIO)
		return -1
	}

	return C.int(n)
}

//export writer
func writer(cookie_ptr unsafe.Pointer, buf *C.char, size C.int) C.int {
	cookie := (*cookie_t)(cookie_ptr)

	wtr, ok := cookie.Writer()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}

	buffer := make_slice(buf, int(size))
	n, err := wtr.Write(buffer)
	if err != nil {
		C.seterr(C.EIO)
		return -1
	}

	return C.int(n)
}

//export closer
func closer(cookie_ptr unsafe.Pointer) C.int {
	cookie := (*cookie_t)(cookie_ptr)
	defer free_cookie(cookie)

	cls, ok := cookie.Closer()
	if !ok {
		return 0
	}

	if err := cls.Close(); err != nil {
		C.seterr(C.EIO)
		return -1
	}

	return 0
}

//export seeker
func seeker(cookie_ptr unsafe.Pointer, position C.fpos_t, whence C.int) C.fpos_t {
	cookie := (*cookie_t)(cookie_ptr)

	skr, ok := cookie.Seeker()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}

	var w int
	// Not sure if C.SEEK_* matches os.SEEK_* in all cases.
	switch whence {
	case C.SEEK_SET:
		w = os.SEEK_SET
	case C.SEEK_CUR:
		w = os.SEEK_CUR
	case C.SEEK_END:
		w = os.SEEK_END
	default:
		C.seterr(C.EINVAL)
		return -1
	}

	ret, err := skr.Seek(int64(position), w)
	if err != nil {
		C.seterr(C.EINVAL)
		return -1
	}
	return C.fpos_t(ret)
}

func wrapReadWriter(rw interface{}) (unsafe.Pointer, error) {
	cookie := new_cookie(rw)

	rdr := C.c_reader
	wtr := C.c_writer
	cls := C.c_closer
	skr := C.c_seeker

	if _, ok := rw.(io.Reader); !ok {
		rdr = nil
	}

	if _, ok := rw.(io.Writer); !ok {
		wtr = nil
	}

	if _, ok := rw.(io.Seeker); !ok {
		skr = nil
	}

	f, err := C.funopen(unsafe.Pointer(cookie), rdr, wtr, skr, cls)
	return unsafe.Pointer(f), err
}
