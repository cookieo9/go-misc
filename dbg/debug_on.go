// +build !nodebug

package dbg

// A default Debug value controlled by the nodebug build flag.
var Default = Debug(true)
