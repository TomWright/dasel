package dasel

import (
	"errors"
	"strings"
)

// ErrDynamicSelectorBracketMismatch is returned when the number of opening brackets doesn't equal that
// of the closing brackets.
var ErrDynamicSelectorBracketMismatch = errors.New("dynamic selector bracket mismatch")

// ExtractNextSelector returns the next selector from the given input.
func ExtractNextSelector(input string) (string, int) {
	escapedIndex := -1
	res := ""
	i := 0
	read := 0
	for k, v := range input {
		curRuneStr := string(v)
		curRuneLength := len(curRuneStr)
		if escapedIndex == k-1 && k != 0 {
			// last character was escape character
			res += curRuneStr
			read += curRuneLength
			continue
		}

		if v == '(' || v == '[' {
			i++
		} else if v == ')' || v == ']' {
			i--
		}

		if v == '\\' {
			escapedIndex = k
			read += curRuneLength
			continue
		}

		if i == 0 && v == '.' && k != 0 {
			break
		}
		res += curRuneStr
		read += curRuneLength
	}
	return res, read
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

// DynamicSelectorParts contains the parts for a dynamic selector.
type DynamicSelectorParts struct {
	Key        string
	Comparison string
	Value      string
}

var comparisons = []string{
	"=",
	"!=",
	"<",
	"<=",
	">",
	">=",
}

func isBuildingComparison(comparison string) bool {
	for _, c := range comparisons {
		if strings.HasPrefix(c, comparison) {
			return true
		}
	}
	return false
}

func isValidComparison(comparison string) bool {
	for _, c := range comparisons {
		if comparison == c {
			return true
		}
	}
	return false
}

// FindDynamicSelectorParts extracts the parts from the dynamic selector given.
func FindDynamicSelectorParts(selector string) DynamicSelectorParts {
	i := 0
	parts := DynamicSelectorParts{}
	for _, v := range selector {
		switch {

		// Start of a group
		case v == '(':
			if parts.Comparison == "" {
				parts.Key += string(v)
			} else {
				parts.Value += string(v)
			}
			i++

		// End of a group
		case v == ')':
			i--
			if parts.Comparison == "" {
				parts.Key += string(v)
			} else {
				parts.Value += string(v)
			}

		// Matches a comparison operator
		case i == 0 && isValidComparison(parts.Comparison+string(v)):
			parts.Comparison += string(v)

		// Is building a comparison character
		case i == 0 && isBuildingComparison(parts.Comparison+string(v)):
			parts.Comparison += string(v)

		// Add to key or value based on comparison existence
		default:
			if parts.Comparison == "" {
				parts.Key += string(v)
			} else {
				parts.Value += string(v)
			}
		}
	}
	return parts
}
