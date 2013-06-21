package cgo

import (
	"io"
)

type cookie_t struct {
	val      interface{}
	do_close bool
}

var cookie_registry = map[*cookie_t]struct{}{}

func new_cookie(val interface{}, do_close bool) *cookie_t {
	c := &cookie_t{
		val:      val,
		do_close: do_close,
	}
	cookie_registry[c] = struct{}{}
	return c
}

func free_cookie(c *cookie_t) {
	delete(cookie_registry, c)
}

func (c cookie_t) Closer() (io.Closer, bool) {
	x, ok := c.val.(io.Closer)
	return x, c.do_close && ok
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
