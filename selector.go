package dasel

import (
	"fmt"
	"io"
	"strings"
)

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
}

func (r *standardSelectorResolver) Original() string {
	return r.original
}

// nextPart returns the next part.
// It returns true if there are more parts to the selector, or false if we reached the end.
func (r *standardSelectorResolver) nextPart() (string, bool) {
	b := &strings.Builder{}
	bracketDepth := 0
	for {
		readRune, size, err := r.reader.ReadRune()
		if err == io.EOF {
			return b.String(), false
		}
		if size == 0 {
			continue
		}
		if readRune == r.openFunc {
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

	for {
		nextRune, size, err := nextPartReader.ReadRune()
		if err == io.EOF {
			if funcNameBuilder.Len() > 0 {
				funcName = funcNameBuilder.String()
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read selector: %w", err)
		}

		if size == 0 {
			continue
		}

		switch {
		case nextRune == r.openFunc:
			if !hasOpenedFunc {
				hasOpenedFunc = true
				funcName = funcNameBuilder.String()
				if funcName == "" {
					return nil, fmt.Errorf("syntax error around \"%s\"", nextPart)
				}
			} else {
				argBuilder.WriteRune(nextRune)
			}
			bracketDepth++

		case nextRune == r.closeFunc:
			if bracketDepth > 1 {
				argBuilder.WriteRune(nextRune)
			} else if bracketDepth == 1 {
				hasClosedFunc = true
				arg := argBuilder.String()
				if arg != "" {
					args = append(args, argBuilder.String())
				}
			}
			bracketDepth--

		case hasOpenedFunc && nextRune == r.argSeparator:
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
			argBuilder.WriteRune(nextRune)

		case hasClosedFunc:
			// Do not allow anything after the closeFunc
			return nil, fmt.Errorf("syntax error around \"%s\"", nextPart)

		default:
			funcNameBuilder.WriteRune(nextRune)
		}
	}

	if !hasOpenedFunc {
		return &Selector{
			funcName: "property",
			funcArgs: []string{funcName},
		}, nil
	}

	// Missing func close
	// if hasOpenedFunc && !hasClosedFunc {
	// 	return nil, fmt.Errorf("unclosed function around \"%s\"", nextPart)
	// }

	return &Selector{
		funcName: funcName,
		funcArgs: args,
	}, nil

}
