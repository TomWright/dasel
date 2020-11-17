package dasel

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSelector parses the given selector string and returns a Selector.
func ParseSelector(selector string) (Selector, error) {
	sel := Selector{
		Raw:        selector,
		Current:    "",
		Remaining:  "",
		Type:       "",
		Property:   "",
		Conditions: make([]Condition, 0),
	}

	{
		nextSelector, read := ExtractNextSelector(sel.Raw)
		sel.Current = nextSelector
		sel.Remaining = sel.Raw[read:]
	}

	nextSel := strings.TrimPrefix(sel.Current, ".")
	var err error

	switch {
	case strings.HasPrefix(nextSel, "(?:") && strings.HasSuffix(nextSel, ")"):
		sel, err = processParseSelectorSearch(nextSel, sel)
	case strings.HasPrefix(nextSel, "(") && strings.HasSuffix(nextSel, ")"):
		sel, err = processParseSelectorDynamic(nextSel, sel)
	case nextSel == "[]":
		sel, err = processParseSelectorNextAvailableIndex(nextSel, sel)
	case nextSel == "[*]":
		sel, err = processParseSelectorIndexAny(nextSel, sel)
	case strings.HasPrefix(nextSel, "[") && strings.HasSuffix(nextSel, "]"):
		sel, err = processParseSelectorIndex(nextSel, sel)
	default:
		sel, err = processParseSelectorProperty(nextSel, sel)
	}

	return sel, err
}

func processParseSelectorDynamic(selector string, sel Selector) (Selector, error) {
	sel.Type = "DYNAMIC"
	dynamicGroups, err := DynamicSelectorToGroups(selector)
	if err != nil {
		return sel, err
	}

	for _, g := range dynamicGroups {
		m := dynamicSelectorRegexp.FindStringSubmatch(g)
		if m == nil {
			return sel, fmt.Errorf("invalid search format")
		}

		var cond Condition
		switch m[2] {
		case "=":
			cond = &EqualCondition{
				Key:   m[1],
				Value: m[3],
			}
		default:
			return sel, &UnknownComparisonOperatorErr{Operator: m[2]}
		}

		sel.Conditions = append(sel.Conditions, cond)
	}

	return sel, nil
}

func processParseSelectorSearch(selector string, sel Selector) (Selector, error) {
	sel.Type = "SEARCH"

	dynamicGroups, err := DynamicSelectorToGroups(selector)
	if err != nil {
		return sel, err
	}
	if len(dynamicGroups) != 1 {
		return sel, fmt.Errorf("require exactly 1 group in search selector")
	}

	for _, g := range dynamicGroups {
		m := dynamicSelectorRegexp.FindStringSubmatch(g)
		if m == nil {
			return sel, fmt.Errorf("invalid search format")
		}

		m[1] = strings.TrimPrefix(m[1], "?:")

		var cond Condition
		switch m[1] {
		case "-", "keyValue":
			switch m[2] {
			case "=":
				cond = &KeyEqualCondition{
					Value: m[3],
				}
			default:
				return sel, &UnknownComparisonOperatorErr{Operator: m[2]}
			}
		default:
			switch m[2] {
			case "=":
				cond = &EqualCondition{
					Key:   strings.TrimPrefix(m[1], "?:"),
					Value: m[3],
				}
			default:
				return sel, &UnknownComparisonOperatorErr{Operator: m[2]}
			}
		}

		sel.Conditions = append(sel.Conditions, cond)
	}

	return sel, nil
}

func processParseSelectorNextAvailableIndex(selector string, sel Selector) (Selector, error) {
	sel.Type = "NEXT_AVAILABLE_INDEX"
	return sel, nil
}

func processParseSelectorIndexAny(selector string, sel Selector) (Selector, error) {
	sel.Type = "INDEX_ANY"
	return sel, nil
}

func processParseSelectorIndex(selector string, sel Selector) (Selector, error) {
	sel.Type = "INDEX"
	indexStr := selector[1 : len(selector)-1]
	index, err := strconv.ParseInt(indexStr, 10, 32)
	if err != nil {
		return sel, &InvalidIndexErr{Index: indexStr}
	}
	sel.Index = int(index)
	return sel, nil
}

func processParseSelectorProperty(selector string, sel Selector) (Selector, error) {
	sel.Type = "PROPERTY"
	sel.Property = selector
	return sel, nil
}
