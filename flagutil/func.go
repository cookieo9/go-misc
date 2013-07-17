package flagutil

import (
	"errors"
	"flag"
	"strings"
)

// A BoolFunction is a type wrapping a function to convert it into a flag.Value.
// It acts as a boolean flag (meaning in go 1.1 it can be passed as -flag).
//
// The flag generated has no default value (will be blank in usage).
type BoolFunction func(bool)

// Set calls the wrapped function with either true or false when passed the
// strings "true", "t", "1" (for true), or "false", "f", "0" (for false).
// It is called only when the flag is mentioned on the command line.
//
// An error is returned if the string is not one of the accepted
// true or false values.
func (bf BoolFunction) Set(s string) error {
	switch strings.ToUpper(s) {
	case "TRUE", "T", "1":
		bf(true)
	case "FALSE", "F", "0":
		bf(false)
	default:
		return errors.New("must pass true of false to this flag")
	}
	return nil
}

// IsBoolFlag returns true to indicate to the Go 1.1+ version of the flag
// package that this flag represents a boolean flag. This means that it can
// be set via "-flag" in addition to the existing "-flag=true" and
// "-flag=false" options.
func (bf BoolFunction) IsBoolFlag() bool {
	return true
}

// String returns "" at all times.
func (bf BoolFunction) String() string {
	return ""
}

// A StringFunction is a type wrapping a function to convert it into a flag.Value.
// The function is called every time the flag is passed a value, and it can return
// an error if the string is unacceptable.
//
// The flag generated has no default value (will be blank in usage).
type StringFunction func(string) error

// Set calls the wrapped function and returns its error (if any).
func (sf StringFunction) Set(s string) error {
	return sf(s)
}

// String returns "" at all times.
func (sf StringFunction) String() string {
	return ""
}

// BoolFunc creates a BoolFunction flag in the default command
// line FlagSet maintained by the flag package.
func BoolFunc(fn func(bool), name, usage string) {
	flag.Var(BoolFunction(fn), name, usage)
}

// StringFunc creates a StringFunction flag in the default command
// line FlagSet maintained by the flag package.
func StringFunc(fn func(string) error, name, usage string) {
	flag.Var(StringFunction(fn), name, usage)
}
