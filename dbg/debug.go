// Package dbg provides simple debugging functions.
// These debugging messages which can be easily turned on or off
// using either a variable at runtime, or a build flag at compile time.
package dbg

import (
	"log"
)

// A Debug value of true will print debugging messages
// to the default log through its methods, false will not.
type Debug bool

// On returns the state of this Debug value as a boolean
func (d Debug) On() bool {
	return bool(d)
}

// Println calls log.Println on the arguments if this Debug
// is active.
func (d Debug) Println(v ...interface{}) {
	if d {
		log.Println(v...)
	}
}

// Printf calls log.Printf on the arguments if this Debug
// is active.
func (d Debug) Printf(format string, v ...interface{}) {
	if d {
		log.Printf(format, v...)
	}
}

// Print calls log.Print on the arguments if this Debug
// is active.
func (d Debug) Print(v ...interface{}) {
	if d {
		log.Print(v...)
	}
}
