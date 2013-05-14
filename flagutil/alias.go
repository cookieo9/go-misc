package flagutil

import (
	"errors"
	"flag"
)

type boolAlias struct {
	target flag.Value
	value  string
}

func (ba boolAlias) String() string {
	if ba.target.String() == ba.value {
		return "true"
	}
	return "false"
}

func (ba boolAlias) Set(s string) error {
	if s == "true" {
		ba.target.Set(ba.value)
	}
	if s != "true" && s != "false" {
		return errors.New("must pass true or false to this flag")
	}
	return nil
}

func (*boolAlias) IsBoolFlag() bool {
	return true
}

// BoolAlias() makes a boolean flag whose purpose is to
// update an existing flag with a predefined string value
// if the new flag is set to true, as well as return "true"
// from String() if target.String() equals the defined value.
func BoolAlias(target flag.Value, name, value, usage string) {
	ba := &boolAlias{
		target: target,
		value:  value,
	}

	flag.Var(ba, name, usage)
}
