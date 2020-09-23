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
func (o *StringList) Type() string {
	return "Pass multiple times to add multiple values."
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
