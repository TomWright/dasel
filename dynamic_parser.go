package dasel

import "errors"

// ErrDynamicSelectorBracketMismatch is returned when the number of opening brackets doesn't equal that
// of the closing brackets.
var ErrDynamicSelectorBracketMismatch = errors.New("dynamic selector bracket mismatch")

// DynamicSelectorToGroups takes a dynamic selector and splits it into groups.
func DynamicSelectorToGroups(selector string) ([]string, error) {
	i := 0
	tmp := ""
	res := make([]string, 0)
	for _, v := range selector {
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
		} else {
			tmp += string(v)
		}
	}
	if i != 0 {
		return nil, ErrDynamicSelectorBracketMismatch
	}
	return res, nil
}
