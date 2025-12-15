package toml

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/pelletier/go-toml/v2/unstable"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Reader = (*tomlReader)(nil)

func newTOMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &tomlReader{}, nil
}

type tomlReader struct{}

// Read reads a value from a byte slice.
func (j *tomlReader) Read(data []byte) (*model.Value, error) {
	p := &unstable.Parser{}
	p.Reset(data)

	root := model.NewMapValue()

	var active *model.Value

	for p.NextExpression() {
		expr := p.Expression()
		switch expr.Kind {
		case unstable.Invalid, unstable.Comment:
			// ignore
			continue
		case unstable.KeyValue:
			keyParts, val, err := j.parseKeyValueNode(p, expr)
			if err != nil {
				return nil, err
			}
			if active != nil {
				if err := setDottedKey(root, active, keyParts, val); err != nil {
					return nil, err
				}
			} else {
				if err := setDottedKey(root, nil, keyParts, val); err != nil {
					return nil, err
				}
			}

		case unstable.Table:
			parts, quoted, err := extractKeyFromTableNode(p, expr)
			if err != nil {
				return nil, err
			}
			m, err := ensureMapAt(root, parts)
			if err != nil {
				return nil, err
			}
			// Record table header parts and quoting info on the map so the writer
			// can reproduce the header exactly if needed.
			m.SetMetadataValue("toml_table_header_parts", parts)
			m.SetMetadataValue("toml_table_header_quoted", quoted)
			active = m

		case unstable.ArrayTable:
			parts, quoted, err := extractKeyFromTableNode(p, expr)
			if err != nil {
				return nil, err
			}
			slice, err := ensureSliceAt(root, parts)
			if err != nil {
				return nil, err
			}
			obj := model.NewMapValue()
			// Mark this object as created via an array table so the writer can
			// emit [[...]] headers for it.
			obj.SetMetadataValue("toml_array_table", true)
			obj.SetMetadataValue("toml_table_header_parts", parts)
			obj.SetMetadataValue("toml_table_header_quoted", quoted)
			if err := slice.Append(obj); err != nil {
				return nil, err
			}
			active = obj

		default:
			// top-level value nodes are unexpected; ignore
		}
	}

	return root, nil
}

// readNode parses a value node (not table/keyvalue headers).
func (j *tomlReader) readNode(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	switch n.Kind {
	// Meta
	case unstable.Invalid:
		return "", nil, nil
	case unstable.Comment:
		return "", nil, nil
	case unstable.Key:
		return "", model.NewStringValue(string(n.Data)), nil

	// Container values
	case unstable.Array:
		v, err := j.readArrayValue(p, n)
		return "", v, err
	case unstable.InlineTable:
		return j.readInlineTable(p, n)

	// Values
	case unstable.String:
		// Create string value and attach TOML-specific style metadata derived
		// from the raw bytes so the writer can reproduce the original form.
		raw := p.Raw(n.Raw)
		v := model.NewStringValue(string(n.Data))
		// Determine style based on raw delimiters
		style := "basic"
		if len(raw) >= 3 && bytes.HasPrefix(raw, []byte("''")) && bytes.HasPrefix(raw, []byte("'''")) {
			style = "multiline_literal"
		} else if len(raw) >= 3 && bytes.HasPrefix(raw, []byte("\"\"\"")) {
			style = "multiline_basic"
		} else if len(raw) >= 1 && raw[0] == '\'' {
			style = "literal"
		} else {
			style = "basic"
		}
		v.SetMetadataValue("toml_style", style)
		return "", v, nil
	case unstable.Bool:
		return "", model.NewBoolValue(string(n.Data) == "true"), nil
	case unstable.Float:
		f, err := strconv.ParseFloat(string(n.Data), 64)
		if err != nil {
			return "", nil, err
		}
		return "", model.NewFloatValue(f), nil
	case unstable.Integer:
		i64, err := strconv.ParseInt(string(n.Data), 10, 64)
		if err != nil {
			return "", nil, err
		}
		return "", model.NewIntValue(int64(i64)), nil
	case unstable.LocalDate:
		return "", model.NewStringValue(string(n.Data)), nil
	case unstable.LocalTime:
		return "", model.NewStringValue(string(n.Data)), nil
	case unstable.LocalDateTime:
		return "", model.NewStringValue(string(n.Data)), nil
	case unstable.DateTime:
		return "", model.NewStringValue(string(n.Data)), nil
	}

	return "", nil, fmt.Errorf("unhandled node kind: %s", n.Kind.String())
}

// parseKeyValueNode extracts the key segments and value from a KeyValue node without consuming parser expressions.
func (j *tomlReader) parseKeyValueNode(p *unstable.Parser, n *unstable.Node) ([]string, *model.Value, error) {
	i := n.Children()
	var keyParts []string
	var val *model.Value

	for i.Next() {
		child := i.Node()
		if child.Kind == unstable.Key {
			keyParts = append(keyParts, string(child.Data))
			continue
		}
		_, v, err := j.readNode(p, child)
		if err != nil {
			return nil, nil, err
		}
		val = v
	}

	if len(keyParts) == 0 {
		return nil, nil, fmt.Errorf("missing key in key/value node")
	}
	if val == nil {
		return nil, nil, fmt.Errorf("missing value in key/value node")
	}

	return keyParts, val, nil
}

// extractKeyFromTableNode returns the key segments from a Table/ArrayTable node.
func extractKeyFromTableNode(p *unstable.Parser, n *unstable.Node) ([]string, []bool, error) {
	i := n.Children()
	var parts []string
	var quoted []bool
	for i.Next() {
		child := i.Node()
		if child.Kind == unstable.Key {
			parts = append(parts, string(child.Data))
			raw := p.Raw(child.Raw)
			isQuoted := false
			if len(raw) > 0 && (raw[0] == '"' || raw[0] == '\'') {
				isQuoted = true
			}
			quoted = append(quoted, isQuoted)
			continue
		}
		return nil, nil, fmt.Errorf("expected table child node, got %s", child.Kind.String())
	}
	if len(parts) == 0 {
		return nil, nil, fmt.Errorf("missing table child key node")
	}
	return parts, quoted, nil
}

// ensureMapAt ensures a map exists at the dotted path under root and returns it.
func ensureMapAt(root *model.Value, path []string) (*model.Value, error) {
	if len(path) == 0 {
		return root, nil
	}
	cur := root
	for _, seg := range path {
		exists, err := cur.MapKeyExists(seg)
		if err != nil {
			return nil, err
		}
		if !exists {
			m := model.NewMapValue()
			if err := cur.SetMapKey(seg, m); err != nil {
				return nil, err
			}
			cur = m
			continue
		}
		next, err := cur.GetMapKey(seg)
		if err != nil {
			return nil, err
		}
		if !next.IsMap() {
			return nil, fmt.Errorf("conflicting types at path '%s': expected map", seg)
		}
		cur = next
	}
	return cur, nil
}

// ensureSliceAt ensures a slice exists at the dotted path under root and returns it.
func ensureSliceAt(root *model.Value, path []string) (*model.Value, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path for array table")
	}
	parentPath := path[:len(path)-1]
	finalSeg := path[len(path)-1]
	parent, err := ensureMapAt(root, parentPath)
	if err != nil {
		return nil, err
	}
	exists, err := parent.MapKeyExists(finalSeg)
	if err != nil {
		return nil, err
	}
	if !exists {
		s := model.NewSliceValue()
		if err := parent.SetMapKey(finalSeg, s); err != nil {
			return nil, err
		}
		return s, nil
	}
	v, err := parent.GetMapKey(finalSeg)
	if err != nil {
		return nil, err
	}
	if !v.IsSlice() {
		return nil, fmt.Errorf("conflicting types at path '%s': expected slice", finalSeg)
	}
	return v, nil
}

// setDottedKey sets a value at a (possibly dotted) key within the given container (creating intermediate maps).
// If active is non-nil, it is the current table object to set keys relative to; root is only used when active is nil to
// create implicit parent maps on the root.
func setDottedKey(root, active *model.Value, parts []string, val *model.Value) error {
	if len(parts) == 0 {
		return fmt.Errorf("empty key")
	}
	// If active table provided, we should set relative to it. But parts may be dotted (i.e., multiple segments).
	if active != nil {
		if len(parts) == 1 {
			return active.SetMapKey(parts[0], val)
		}
		parent, err := ensureMapAt(active, parts[:len(parts)-1])
		if err != nil {
			return err
		}
		return parent.SetMapKey(parts[len(parts)-1], val)
	}
	// No active table: set on root
	if len(parts) == 1 {
		return root.SetMapKey(parts[0], val)
	}
	parent, err := ensureMapAt(root, parts[:len(parts)-1])
	if err != nil {
		return err
	}
	return parent.SetMapKey(parts[len(parts)-1], val)
}

func (j *tomlReader) readInlineTable(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	res := model.NewMapValue()
	res.SetMetadataValue("toml_inline_table", true)

	i := n.Children()
	for i.Next() {
		childNode := i.Node()
		// Inline table children are key/value pairs. Handle KeyValue specially.
		switch childNode.Kind {
		case unstable.KeyValue:
			kparts, v, err := j.parseKeyValueNode(p, childNode)
			if err != nil {
				return "", nil, err
			}
			if err := setDottedKey(res, nil, kparts, v); err != nil {
				return "", nil, err
			}
		default:
			// fallback to readNode for other kinds (e.g., Key)
			key, val, err := j.readNode(p, childNode)
			if err != nil {
				return "", nil, err
			}
			if key == "" {
				return "", nil, fmt.Errorf("missing key in inline table child")
			}
			if err := res.SetMapKey(key, val); err != nil {
				return "", nil, err
			}
		}
	}

	return "", res, nil
}

func (j *tomlReader) readArrayValue(p *unstable.Parser, n *unstable.Node) (*model.Value, error) {
	res := model.NewSliceValue()

	i := n.Children()

	for i.Next() {
		childNode := i.Node()

		_, val, err := j.readNode(p, childNode)
		if err != nil {
			return nil, err
		}

		if err := res.Append(val); err != nil {
			return nil, err
		}
	}

	return res, nil
}
