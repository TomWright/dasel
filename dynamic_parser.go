package dasel

import (
	"errors"
)

// ErrDynamicSelectorBracketMismatch is returned when the number of opening brackets doesn't equal that
// of the closing brackets.
var ErrDynamicSelectorBracketMismatch = errors.New("dynamic selector bracket mismatch")

// ExtractNextSelector returns the next selector from the given input.
func ExtractNextSelector(input string) string {
	res := ""
	i := 0
	for k, v := range input {
		if v == '(' || v == '[' {
			i++
		} else if v == ')' || v == ']' {
			i--
		}

		if i == 0 && v == '.' && k != 0 {
			break
		}
		res += string(v)
	}
	return res
}

// DynamicSelectorToGroups takes a dynamic selector and splits it into groups.
func DynamicSelectorToGroups(selector string) ([]string, error) {
	i := 0
	tmp := ""
	res := make([]string, 0)
	for k, v := range selector {
		if v == '(' {
			if i > 0 {
				tmp += string(v)
			} else {
				tmp = ""
			}
			i++
		} else if v == ')' {
			i--
			if i == 0 {
				res = append(res, tmp)
				tmp = ""
			} else {
				tmp += string(v)
			}
		} else if v == '.' && i == 0 && k != 0 {
			return res, nil
		} else {
			tmp += string(v)
		}
	}
	if i != 0 {
		return nil, ErrDynamicSelectorBracketMismatch
	}
	return res, nil
}
