package dasel

import (
	"fmt"
	"io"
	"strings"
)

type ErrBadSelectorSyntax struct {
	Part    string
	Message string
}

func (e ErrBadSelectorSyntax) Error() string {
	return fmt.Sprintf("bad syntax: %s, around %s", e.Message, e.Part)
}

func (e ErrBadSelectorSyntax) Is(other error) bool {
	o, ok := other.(*ErrBadSelectorSyntax)
	if !ok {
		return false
	}
	if o.Part != "" && o.Part != e.Part {
		return false
	}
	if o.Message != "" && o.Message != e.Message {
		return false
	}
	return true
}

type Selector struct {
	funcName string
	funcArgs []string
}

type SelectorResolver interface {
	Original() string
	Next() (*Selector, error)
}

func NewSelectorResolver(selector string, functions *FunctionCollection) SelectorResolver {
	return &standardSelectorResolver{
		functions:    functions,
		original:     selector,
		reader:       strings.NewReader(selector),
		separator:    '.',
		openFunc:     '(',
		closeFunc:    ')',
		argSeparator: ',',
		escapeChar:   '\\',
	}
}

type standardSelectorResolver struct {
	functions    *FunctionCollection
	original     string
	reader       *strings.Reader
	separator    rune
	openFunc     rune
	closeFunc    rune
	argSeparator rune
	escapeChar   rune
}

func (r *standardSelectorResolver) Original() string {
	return r.original
}

// nextPart returns the next part.
// It returns true if there are more parts to the selector, or false if we reached the end.
func (r *standardSelectorResolver) nextPart() (string, bool) {
	b := &strings.Builder{}
	bracketDepth := 0
	escaped := false
	for {
		readRune, _, err := r.reader.ReadRune()
		if err == io.EOF {
			return b.String(), false
		}
		if escaped {
			b.WriteRune(readRune)
			escaped = false
			continue
		} else if readRune == r.escapeChar {
			b.WriteRune(readRune)
			escaped = true
			continue
		} else if readRune == r.openFunc {
			bracketDepth++
		} else if readRune == r.closeFunc {
			bracketDepth--
		}
		if readRune == r.separator && bracketDepth == 0 {
			return b.String(), true
		}
		b.WriteRune(readRune)
	}
}

func (r *standardSelectorResolver) Next() (*Selector, error) {
	nextPart, moreParts := r.nextPart()
	if nextPart == "" && !moreParts {
		return nil, nil
	}
	if nextPart == "" && moreParts {
		return &Selector{
			funcName: "this",
			funcArgs: []string{},
		}, nil
	}

	if r.functions != nil {
		if s := r.functions.ParseSelector(nextPart); s != nil {
			return s, nil
		}
	}

	var hasOpenedFunc, hasClosedFunc = false, false
	bracketDepth := 0

	var funcNameBuilder = &strings.Builder{}
	var argBuilder = &strings.Builder{}

	nextPartReader := strings.NewReader(nextPart)

	funcName := ""
	args := make([]string, 0)

	escaped := false
	for {
		nextRune, _, err := nextPartReader.ReadRune()
		if err == io.EOF {
			if funcNameBuilder.Len() > 0 {
				funcName = funcNameBuilder.String()
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read selector: %w", err)
		}

		switch {
		case nextRune == r.escapeChar && !escaped:
			escaped = true
			continue

		case nextRune == r.openFunc && !escaped:
			if !hasOpenedFunc {
				hasOpenedFunc = true
				funcName = funcNameBuilder.String()
				if funcName == "" {
					return nil, &ErrBadSelectorSyntax{
						Part:    nextPart,
						Message: "function name required before open bracket",
					}
				}
			} else {
				argBuilder.WriteRune(nextRune)
			}
			bracketDepth++

		case nextRune == r.closeFunc && !escaped:
			if bracketDepth > 1 {
				argBuilder.WriteRune(nextRune)
			} else if bracketDepth == 1 {
				hasClosedFunc = true
				arg := argBuilder.String()
				if arg != "" {
					args = append(args, argBuilder.String())
				}
			} else if bracketDepth < 1 {
				return nil, &ErrBadSelectorSyntax{
					Part:    nextPart,
					Message: "too many closing brackets",
				}
			}
			bracketDepth--

		case hasOpenedFunc && nextRune == r.argSeparator && !escaped:
			if bracketDepth > 1 {
				argBuilder.WriteRune(nextRune)
			} else if bracketDepth == 1 {
				arg := argBuilder.String()
				argBuilder.Reset()
				if arg != "" {
					args = append(args, arg)
				}
			}

		case hasOpenedFunc:
			if escaped {
				escaped = false
			}
			argBuilder.WriteRune(nextRune)

		case hasClosedFunc:
			// Do not allow anything after the closeFunc
			return nil, &ErrBadSelectorSyntax{
				Part:    nextPart,
				Message: "selector function must end after closing bracket",
			}

		default:
			if escaped {
				escaped = false
			}
			funcNameBuilder.WriteRune(nextRune)
		}
	}

	if !hasOpenedFunc {
		return &Selector{
			funcName: "property",
			funcArgs: []string{funcName},
		}, nil
	}

	return &Selector{
		funcName: funcName,
		funcArgs: args,
	}, nil

}
