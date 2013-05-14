package dbg

import (
	"log"
)

type Debug bool

func (d Debug) On() bool {
	return bool(d)
}

func (d Debug) Println(v ...interface{}) {
	if d {
		log.Println(v...)
	}
}

func (d Debug) Printf(format string, v ...interface{}) {
	if d {
		log.Printf(format, v...)
	}
}

func (d Debug) Print(v ...interface{}) {
	if d {
		log.Print(v...)
	}
}
