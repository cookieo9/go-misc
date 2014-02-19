// +build linux
// +build cgo

package cgo

//
// #define _GNU_SOURCE
// #include <stdio.h>
// #include <stdlib.h>
// #include <errno.h>
// #include <sys/errno.h>
// #include "indirect_linux.h"
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
func reader(cookiePtr unsafe.Pointer, buf *C.char, size C.size_t) C.ssize_t {
	cookie := (*cookie)(cookiePtr)

	rdr, ok := cookie.Reader()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}

	buffer := makeSlice(buf, int(size))
	n, err := rdr.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return C.ssize_t(n)
		}
		C.seterr(C.EIO)
		return -1
	}

	return C.ssize_t(n)
}

//export writer
func writer(cookiePtr unsafe.Pointer, buf *C.char, size C.size_t) C.ssize_t {
	cookie := (*cookie)(cookiePtr)

	rdr, ok := cookie.Writer()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}

	buffer := makeSlice(buf, int(size))
	n, err := rdr.Write(buffer)
	if err != nil {
		if err == io.EOF {
			return C.ssize_t(n)
		}
		C.seterr(C.EIO)
		return -1
	}

	return C.ssize_t(n)
}

//export closer
func closer(cookiePtr unsafe.Pointer) C.int {
	cookie := (*cookie)(cookiePtr)
	defer freeCookie(cookie)

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
func seeker(cookiePtr unsafe.Pointer, position *C.off64_t, whence C.int) C.int {
	cookie := (*cookie)(cookiePtr)

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

	ret, err := skr.Seek(int64(*position), w)
	if err != nil {
		C.seterr(C.EINVAL)
		return -1
	}

	*position = C.off64_t(ret)
	return 0
}

func wrapReadWriter(cookie *cookie) (unsafe.Pointer, error) {
	rdr := C.c_reader
	wtr := C.c_writer
	cls := C.c_closer
	skr := C.c_seeker

	if _, ok := cookie.val.(io.Seeker); !ok {
		skr = nil
	}

	if _, ok := cookie.val.(io.Reader); !ok {
		rdr = nil
	}

	if _, ok := cookie.val.(io.Writer); !ok {
		wtr = nil
	}

	fns := C._IO_cookie_io_functions_t{
		read:  rdr,
		write: wtr,
		seek:  skr,
		close: cls,
	}

	mode := "w+"
	if rdr == nil {
		mode = "w"
	}
	if wtr == nil {
		mode = "r"
	}

	cmode := C.CString(mode)
	defer C.free(unsafe.Pointer(cmode))

	f, err := C.fopencookie(unsafe.Pointer(cookie), cmode, fns)
	return unsafe.Pointer(f), err
}
