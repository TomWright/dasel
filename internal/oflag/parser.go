package oflag

import (
	"errors"
	"strconv"
	"strings"
)

// ParseOverrideValueFn defines a function used to parse config values passed in through command line arguments.
type ParseOverrideValueFn func(value string) (interface{}, error)

// StringParser parses string config values.
func StringParser(value string) (interface{}, error) {
	return value, nil
}

// IntParser parses int config values.
func IntParser(value string) (interface{}, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ErrInvalidBool is returned if the given value does not match an expected bool value.
var ErrInvalidBool = errors.New("unexpected bool value")

// BoolParser parses int config values.
func BoolParser(value string) (interface{}, error) {
	switch strings.ToLower(value) {
	case "yes", "y", "true", "t", "1":
		return true, nil
	case "no", "n", "false", "f", "0":
		return false, nil
	default:
		return false, ErrInvalidBool
	}
}
