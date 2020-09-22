package oflag

import (
	"fmt"
	"strings"
)

// Override represents a single override.
type Override struct {
	// Path is a string that may have dot separators within it.
	Path string
	// Value is the value to set in the config file.
	Value interface{}
}

// String returns a string representation of the override.
func (o Override) String() string {
	return fmt.Sprintf("%s=%v", o.Path, o.Value)
}

// NewOverrideFlag returns an OverrideFlag with the given parse func.
func NewOverrideFlag(parse ParseOverrideValueFn) *OverrideFlag {
	return &OverrideFlag{
		overrides: make([]*Override, 0),
		Parse:     parse,
	}
}

// Combine returns a combined list of overrides from all given OverrideFlag's.
func Combine(flags ...*OverrideFlag) []*Override {
	res := make([]*Override, 0)
	for _, f := range flags {
		res = append(res, f.Overrides()...)
	}
	return res
}

// OverrideFlag allows us to collect a list of overrides from command line flags.
type OverrideFlag struct {
	overrides []*Override
	// Parse allows us to parse override values into specific types.
	Parse ParseOverrideValueFn
}

// Overrides returns the collected overrides.
func (of *OverrideFlag) Overrides() []*Override {
	return of.overrides
}

// Type returns the type of flag.
// As far as I can tell this isn't used...
// todo : figure out what this is supposed to return.
func (of *OverrideFlag) Type() string {
	panic("implement me")
}

// String returns a string representation of all Override's within the OverrideFlag.
func (of *OverrideFlag) String() string {
	val := make([]string, len(of.overrides))
	for k, o := range of.overrides {
		val[k] = o.String()
	}
	return strings.Join(val, " ")
}

// Set is used to add a new value to the OverrideFlag.
func (of *OverrideFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	args := strings.Split(value, "=")

	var override *Override
	switch len(args) {
	case 0:
		// No value was given.
		return nil
	case 1:
		// A blank value was given.
		override = &Override{
			Path:  args[0],
			Value: "",
		}
	case 2:
		// A single value was given.
		override = &Override{
			Path:  args[0],
			Value: args[1],
		}
	default:
		// The value contained more than 1 = sign.
		// Assume the extra = values are within the value.
		override = &Override{
			Path:  args[0],
			Value: strings.Join(args[1:], "="),
		}
	}

	if of.Parse != nil {
		var err error
		// We know that override.Value will be a string here.
		override.Value, err = of.Parse(override.Value.(string))
		if err != nil {
			return fmt.Errorf("invalid value for path `%s`: %w", override.Path, err)
		}
	}

	of.overrides = append(of.overrides, override)

	return nil
}
