// Adds a simple system for inserting debugging messages which can be easily turned on or off
// using either a variable at runtime, or a build flag at compile time.
package dbg

import (
	"log"
)

// A Debug value of true will print debugging messages
// to the default log through its methods, false will not.
type Debug bool

// Returns the active state of this Debug value
// as a boolean
func (d Debug) On() bool {
	return bool(d)
}

// Calls log.Println on the arguments if this Debug
// is active.
func (d Debug) Println(v ...interface{}) {
	if d {
		log.Println(v...)
	}
}

// Calls log.Printf on the arguments if this Debug
// is active.
func (d Debug) Printf(format string, v ...interface{}) {
	if d {
		log.Printf(format, v...)
	}
}

// Calls log.Print on the arguments if this Debug
// is active.
func (d Debug) Print(v ...interface{}) {
	if d {
		log.Print(v...)
	}
}
