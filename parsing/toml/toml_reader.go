package toml

import (
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

type tomlReader struct {
	skipNext bool
}

// Read reads a value from a byte slice.
func (j *tomlReader) Read(data []byte) (*model.Value, error) {
	p := &unstable.Parser{}
	p.Reset(data)

	res := model.NewMapValue()

	for j.skipNext || p.NextExpression() {
		j.skipNext = false
		next := p.Expression()
		k, v, err := j.readNode(p, next)
		if err != nil {
			return nil, err
		}
		if err := res.SetMapKey(k, v); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (j *tomlReader) readNode(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	switch n.Kind {
	// Meta
	case unstable.Invalid:
		return "", nil, nil
	case unstable.Comment:
		return "", nil, nil
	case unstable.Key:
		return "", model.NewStringValue(string(n.Data)), nil

	// Top level structures
	case unstable.Table:
		return j.readTable(p, n)
	case unstable.ArrayTable:
		return j.readArrayTable(p, n)
	case unstable.KeyValue:
		return j.readKeyValue(p, n)

	// Containers values
	case unstable.Array:
		v, err := j.readArrayValue(p, n)
		return "", v, err
	case unstable.InlineTable:
		return j.readInlineTable(p, n)

	// Values
	case unstable.String:
		return "", model.NewStringValue(string(n.Data)), nil
	case unstable.Bool:
		return "", model.NewBoolValue(string(n.Data) == "true"), nil
	case unstable.Float:
		f, err := strconv.ParseFloat(string(n.Data), 64)
		if err != nil {
			return "", nil, err
		}
		return "", model.NewFloatValue(f), nil
	case unstable.Integer:
		i, err := strconv.Atoi(string(n.Data))
		if err != nil {
			return "", nil, err
		}
		return "", model.NewIntValue(int64(i)), nil
	case unstable.LocalDate:
	case unstable.LocalTime:
	case unstable.LocalDateTime:
	case unstable.DateTime:
	}

	return "", nil, fmt.Errorf("unhandled node kind: %s", n.Kind.String())
}

func (j *tomlReader) readKeyValue(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	i := n.Children()

	i.Next()
	valueNode := i.Node()

	_, value, err := j.readNode(p, valueNode)
	if err != nil {
		return "", nil, err
	}

	i.Next()
	keyNode := i.Node()

	return string(keyNode.Data), value, nil
}

func (j *tomlReader) readTable(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	res := model.NewMapValue()

	i := n.Children()

	var k string

	for i.Next() {
		childNode := i.Node()
		if childNode.Kind == unstable.Key {
			k = string(childNode.Data)
			continue
		}
		return k, nil, fmt.Errorf("expected table child node, got %s", childNode.Kind.String())
	}
	if k == "" {
		return k, nil, fmt.Errorf("missing table child key node")
	}

	for p.NextExpression() {
		key, value, err := j.readNode(p, p.Expression())
		if err != nil {
			return k, nil, err
		}
		if err := res.SetMapKey(key, value); err != nil {
			return k, nil, err
		}
	}

	return k, res, nil
}

func (j *tomlReader) readArrayTable(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	i := n.Children()
	var k string

	for i.Next() {
		childNode := i.Node()
		if childNode.Kind == unstable.Key {
			k = string(childNode.Data)
			continue
		}
		return k, nil, fmt.Errorf("expected table child node, got %s", childNode.Kind.String())
	}
	if k == "" {
		return k, nil, fmt.Errorf("missing table child key node")
	}

	obj := model.NewMapValue()

	for p.NextExpression() {
		if p.Expression().Kind != unstable.KeyValue {
			j.skipNext = true
			break
		}
		expr := p.Expression()
		next := expr.Next()
		key, value, err := j.readNode(p, expr)
		fmt.Println(key, value, next, err)
		if err != nil {
			return k, nil, err
		}

		if err := obj.SetMapKey(key, value); err != nil {
			return k, nil, err
		}
	}

	return k, obj, nil
}

func (j *tomlReader) readInlineTable(p *unstable.Parser, n *unstable.Node) (string, *model.Value, error) {
	res := model.NewMapValue()
	res.SetMetadataValue("toml_inline_table", true)

	i := n.Children()
	for i.Next() {
		childNode := i.Node()
		key, val, err := j.readNode(p, childNode)
		if err != nil {
			return "", nil, err
		}
		if err := res.SetMapKey(key, val); err != nil {
			return "", nil, err
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
