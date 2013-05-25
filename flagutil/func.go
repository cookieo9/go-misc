package flagutil

import (
	"errors"
	"flag"
)

type boolFunc func(bool)

func (bf boolFunc) Set(s string) error {
	if s == "true" {
		bf(true)
	}
	if s == "false" {
		bf(false)
	}

	if s != "true" && s != "false" {
		return errors.New("must pass true of false to this flag")
	}
	return nil
}

func (bf boolFunc) IsBoolFlag() bool {
	return true
}

func (bf boolFunc) String() string {
	return ""
}

type stringFunc func(string) error

func (sf stringFunc) Set(s string) error {
	return sf(s)
}

func (sf stringFunc) String() string {
	return ""
}

func BoolFunc(fn func(bool), name, usage string) {
	flag.Var(boolFunc(fn), name, usage)
}

func NewBoolFunc(fn func(bool)) flag.Value {
	return boolFunc(fn)
}

func StringFunc(fn func(string) error, name, usage string) {
	flag.Var(stringFunc(fn), name, usage)
}

func NewStringFunc(fn func(string) error) flag.Value {
	return stringFunc(fn)
}
