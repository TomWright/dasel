package oflag

import (
	"fmt"
)

// StringList is a flag that can collect multiple strings.
type StringList struct {
	Strings []string
}

// NewStringList returns a new string list flag.
func NewStringList() *StringList {
	return &StringList{
		Strings: []string{},
	}
}

// Type returns the type of flag.
// As far as I can tell this isn't used...
// todo : figure out what this is supposed to return.
func (o *StringList) Type() string {
	panic("implement me")
}

// String returns a string representation of all Override's within the OverrideFlag.
func (o *StringList) String() string {
	return fmt.Sprint(o.Strings)
}

// Set is used to add a new value to the OverrideFlag.
func (o *StringList) Set(value string) error {
	o.Strings = append(o.Strings, value)
	return nil
}
