package internal

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// GenerateOptions controls the output format.
type GenerateOptions struct {
	Indent  string
	Compact bool
	Version Version // Output version. Defaults to Version2.
}

// DefaultGenerateOptions returns sensible defaults.
func DefaultGenerateOptions() GenerateOptions {
	return GenerateOptions{
		Indent:  "    ",
		Version: Version2,
	}
}

// Generate writes a Document to the given writer as KDL.
func Generate(w io.Writer, doc *Document, opts GenerateOptions) error {
	if opts.Version == VersionUnknown {
		opts.Version = Version2
	}
	g := &generator{w: w, opts: opts}
	return g.writeDocument(doc, 0)
}

// GenerateString returns a KDL string for the document.
func GenerateString(doc *Document, opts GenerateOptions) (string, error) {
	var sb strings.Builder
	if err := Generate(&sb, doc, opts); err != nil {
		return "", err
	}
	return sb.String(), nil
}

type generator struct {
	w    io.Writer
	opts GenerateOptions
}

func (g *generator) writeDocument(doc *Document, depth int) error {
	for i, node := range doc.Nodes {
		if err := g.writeNode(node, depth); err != nil {
			return err
		}
		if !g.opts.Compact || i < len(doc.Nodes)-1 {
			if err := g.writeNewline(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) writeNode(node *Node, depth int) error {
	if !g.opts.Compact {
		if err := g.writeIndent(depth); err != nil {
			return err
		}
	}

	// Type annotation
	if node.Type != "" {
		if _, err := fmt.Fprintf(g.w, "(%s)", quoteIdentifier(node.Type)); err != nil {
			return err
		}
	}

	// Node name
	if _, err := fmt.Fprint(g.w, quoteIdentifier(node.Name)); err != nil {
		return err
	}

	// Arguments
	for _, arg := range node.Arguments {
		if _, err := fmt.Fprint(g.w, " "); err != nil {
			return err
		}
		if err := g.writeValue(arg); err != nil {
			return err
		}
	}

	// Properties
	for _, prop := range node.Properties {
		if _, err := fmt.Fprint(g.w, " "); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(g.w, "%s=", quoteIdentifier(prop.Key)); err != nil {
			return err
		}
		if err := g.writeValue(prop.Value); err != nil {
			return err
		}
	}

	// Children
	if len(node.Children) > 0 {
		if g.opts.Compact {
			if _, err := fmt.Fprint(g.w, "{"); err != nil {
				return err
			}
			for _, child := range node.Children {
				if err := g.writeNode(child, depth+1); err != nil {
					return err
				}
				if _, err := fmt.Fprint(g.w, ";"); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(g.w, "}"); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(g.w, " {"); err != nil {
				return err
			}
			if err := g.writeNewline(); err != nil {
				return err
			}
			for _, child := range node.Children {
				if err := g.writeNode(child, depth+1); err != nil {
					return err
				}
				if err := g.writeNewline(); err != nil {
					return err
				}
			}
			if err := g.writeIndent(depth); err != nil {
				return err
			}
			if _, err := fmt.Fprint(g.w, "}"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) writeValue(v *Value) error {
	if v.Type != "" {
		if _, err := fmt.Fprintf(g.w, "(%s)", quoteIdentifier(v.Type)); err != nil {
			return err
		}
	}

	switch val := v.Value.(type) {
	case string:
		if _, err := fmt.Fprint(g.w, quoteString(val)); err != nil {
			return err
		}
	case int64:
		if _, err := fmt.Fprintf(g.w, "%d", val); err != nil {
			return err
		}
	case float64:
		if math.IsInf(val, 1) {
			if _, err := fmt.Fprint(g.w, g.keyword("inf")); err != nil {
				return err
			}
		} else if math.IsInf(val, -1) {
			if _, err := fmt.Fprint(g.w, g.keyword("-inf")); err != nil {
				return err
			}
		} else if math.IsNaN(val) {
			if _, err := fmt.Fprint(g.w, g.keyword("nan")); err != nil {
				return err
			}
		} else {
			s := fmt.Sprintf("%g", val)
			// Ensure there's a decimal point so it's clearly a float
			if !strings.Contains(s, ".") && !strings.Contains(s, "e") && !strings.Contains(s, "E") {
				s += ".0"
			}
			if _, err := fmt.Fprint(g.w, s); err != nil {
				return err
			}
		}
	case bool:
		if val {
			if _, err := fmt.Fprint(g.w, g.keyword("true")); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprint(g.w, g.keyword("false")); err != nil {
				return err
			}
		}
	case nil:
		if _, err := fmt.Fprint(g.w, g.keyword("null")); err != nil {
			return err
		}
	default:
		return fmt.Errorf("kdl: unsupported value type %T", val)
	}
	return nil
}

// keyword returns the version-appropriate keyword form.
// In v2: #true, #false, #null, #inf, #-inf, #nan
// In v1: true, false, null (inf/-inf/nan not spec keywords in v1)
func (g *generator) keyword(name string) string {
	if g.opts.Version == Version1 {
		return name
	}
	return "#" + name
}

func (g *generator) writeIndent(depth int) error {
	if g.opts.Indent == "" {
		return nil
	}
	_, err := fmt.Fprint(g.w, strings.Repeat(g.opts.Indent, depth))
	return err
}

func (g *generator) writeNewline() error {
	if g.opts.Compact {
		return nil
	}
	_, err := fmt.Fprint(g.w, "\n")
	return err
}

// quoteIdentifier returns a bare identifier if possible, otherwise quotes it.
func quoteIdentifier(s string) string {
	if s == "" {
		return `""`
	}
	if canBeBareIdent(s) {
		return s
	}
	return quoteString(s)
}

// canBeBareIdent returns true if s is a valid bare KDL identifier.
func canBeBareIdent(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	if !isIdentStart(runes[0]) {
		return false
	}
	for _, r := range runes[1:] {
		if !isIdentChar(r) {
			return false
		}
	}
	// Disallow v1 keywords as bare identifiers (they'd be misinterpreted)
	switch s {
	case "true", "false", "null", "inf", "-inf", "nan":
		return false
	}
	return true
}

// quoteString returns a properly escaped KDL quoted string.
func quoteString(s string) string {
	var sb strings.Builder
	sb.WriteRune('"')
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString(`\"`)
		case '\\':
			sb.WriteString(`\\`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		case '\b':
			sb.WriteString(`\b`)
		case '\f':
			sb.WriteString(`\f`)
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteRune('"')
	return sb.String()
}
