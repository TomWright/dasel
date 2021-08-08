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
		// evaluable, err := gvalInstance.NewEvaluable(g)
		// if err != nil {
		// 	return sel, fmt.Errorf("could not parse dynamic expression: %w", err)
		// }
		//
		// todo : how do we execute sub queries?
		// sel.Conditions = append(sel.Conditions, &GvalCondition{
		// 	Evaluable: evaluable,
		// })

		m := FindDynamicSelectorParts(g)

		var cond Condition

		switch m.Comparison {
		case "=":
			cond = &EqualCondition{
				Key:   m.Key,
				Value: m.Value,
			}
		case "!=":
			cond = &EqualCondition{
				Key:   m.Key,
				Value: m.Value,
				Not:   true,
			}
		case ">=":
			cond = &SortedComparisonCondition{
				Key:   m.Key,
				Value: m.Value,
				Equal: true,
				After: true,
			}
		case ">":
			cond = &SortedComparisonCondition{
				Key:   m.Key,
				Value: m.Value,
				After: true,
			}
		case "<=":
			cond = &SortedComparisonCondition{
				Key:   m.Key,
				Value: m.Value,
				Equal: true,
			}
		case "<":
			cond = &SortedComparisonCondition{
				Key:   m.Key,
				Value: m.Value,
			}
		default:
			return sel, &UnknownComparisonOperatorErr{Operator: m.Comparison}
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
		m := FindDynamicSelectorParts(g)

		m.Key = strings.TrimPrefix(m.Key, "?:")

		var cond Condition
		switch m.Key {
		case "-", "keyValue":
			switch m.Comparison {
			case "=":
				cond = &KeyEqualCondition{
					Value: m.Value,
				}
			case "!=":
				cond = &KeyEqualCondition{
					Value: m.Value,
					Not:   true,
				}
			default:
				return sel, &UnknownComparisonOperatorErr{Operator: m.Comparison}
			}
		default:
			switch m.Comparison {
			case "=":
				cond = &EqualCondition{
					Key:   m.Key,
					Value: m.Value,
				}
			case "!=":
				cond = &EqualCondition{
					Key:   m.Key,
					Value: m.Value,
					Not:   true,
				}
			case ">=":
				cond = &SortedComparisonCondition{
					Key:   m.Key,
					Value: m.Value,
					Equal: true,
					After: true,
				}
			case ">":
				cond = &SortedComparisonCondition{
					Key:   m.Key,
					Value: m.Value,
					After: true,
				}
			case "<=":
				cond = &SortedComparisonCondition{
					Key:   m.Key,
					Value: m.Value,
					Equal: true,
				}
			case "<":
				cond = &SortedComparisonCondition{
					Key:   m.Key,
					Value: m.Value,
				}
			default:
				return sel, &UnknownComparisonOperatorErr{Operator: m.Comparison}
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

func processParseSelectorLength(selector string, sel Selector) (Selector, error) {
	sel.Type = "LENGTH"
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
