package cgo

import (
	"io"
)

type cookie struct {
	val     interface{}
	doClose bool
}

var cookieRegistry = map[*cookie]struct{}{}

func newCookie(val interface{}, doClose bool) *cookie {
	c := &cookie{
		val:     val,
		doClose: doClose,
	}
	cookieRegistry[c] = struct{}{}
	return c
}

func freeCookie(c *cookie) {
	delete(cookieRegistry, c)
}

func (c cookie) Closer() (io.Closer, bool) {
	x, ok := c.val.(io.Closer)
	return x, c.doClose && ok
}

func (c cookie) Writer() (io.Writer, bool) {
	x, ok := c.val.(io.Writer)
	return x, ok
}

func (c cookie) Seeker() (io.Seeker, bool) {
	x, ok := c.val.(io.Seeker)
	return x, ok
}

func (c cookie) Reader() (io.Reader, bool) {
	x, ok := c.val.(io.Reader)
	return x, ok
}
