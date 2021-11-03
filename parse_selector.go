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
	case nextSel == "[#]":
		sel, err = processParseSelectorLength(nextSel, sel)
	case nextSel == "[@]":
		sel, err = processParseSelectorType(nextSel, sel)
	case strings.HasPrefix(nextSel, "[") && strings.HasSuffix(nextSel, "]"):
		sel, err = processParseSelectorIndex(nextSel, sel)
	default:
		sel, err = processParseSelectorProperty(nextSel, sel)
	}

	return sel, err
}

func getCondition(parts DynamicSelectorParts) (Condition, error) {
	switch parts.Key {
	case "-", "keyValue":
		switch parts.Comparison {
		case "=":
			return &KeyEqualCondition{
				Value: parts.Value,
			}, nil
		case "!=":
			return &KeyEqualCondition{
				Value: parts.Value,
				Not:   true,
			}, nil
		default:
			return nil, &UnknownComparisonOperatorErr{Operator: parts.Comparison}
		}
	default:

		switch parts.Comparison {
		case "=":
			return &EqualCondition{
				Key:   parts.Key,
				Value: parts.Value,
			}, nil
		case "!=":
			return &EqualCondition{
				Key:   parts.Key,
				Value: parts.Value,
				Not:   true,
			}, nil
		case ">=":
			return &SortedComparisonCondition{
				Key:   parts.Key,
				Value: parts.Value,
				Equal: true,
				After: true,
			}, nil
		case ">":
			return &SortedComparisonCondition{
				Key:   parts.Key,
				Value: parts.Value,
				After: true,
			}, nil
		case "<=":
			return &SortedComparisonCondition{
				Key:   parts.Key,
				Value: parts.Value,
				Equal: true,
			}, nil
		case "<":
			return &SortedComparisonCondition{
				Key:   parts.Key,
				Value: parts.Value,
			}, nil
		default:
			return nil, &UnknownComparisonOperatorErr{Operator: parts.Comparison}
		}
	}
}

func processParseSelectorDynamic(selector string, sel Selector) (Selector, error) {
	sel.Type = "DYNAMIC"
	dynamicGroups, err := DynamicSelectorToGroups(selector)
	if err != nil {
		return sel, err
	}

	for _, g := range dynamicGroups {
		parts := FindDynamicSelectorParts(g)
		cond, err := getCondition(parts)
		if err != nil {
			return sel, err
		}
		if cond != nil {
			sel.Conditions = append(sel.Conditions, cond)
		}
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
		parts := FindDynamicSelectorParts(g)
		parts.Key = strings.TrimPrefix(parts.Key, "?:")
		cond, err := getCondition(parts)
		if err != nil {
			return sel, err
		}
		if cond != nil {
			sel.Conditions = append(sel.Conditions, cond)
		}
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

func processParseSelectorLength(selector string, sel Selector) (Selector, error) {
	sel.Type = "LENGTH"
	return sel, nil
}

func processParseSelectorType(selector string, sel Selector) (Selector, error) {
	sel.Type = "TYPE"
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
