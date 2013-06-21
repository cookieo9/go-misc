// +build !darwin,!freebsd,!openbsd,!netbsd

package cgo

import (
	"errors"
	"unsafe"
)

func wrapReadWriter(rw interface{}) (unsafe.Pointer, error) {
	return nil, errors.New("unsupported")
}
