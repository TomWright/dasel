package dasel

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"reflect"
	"strings"
	"text/template"
)

// FormatNode formats a node with the format template and returns the result.
func FormatNode(node *Node, format string) (*bytes.Buffer, error) {
	tpl, err := formatNodeTemplate(
		&templateNode{
			Node:    node,
			isFirst: true,
			isLast:  true,
		},
	).Parse(format)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = tpl.Execute(buf, node.InterfaceValue())
	return buf, err
}

type templateNode struct {
	*Node

	isFirst bool
	isLast  bool
}

// FormatNodes formats a slice of nodes with the format template and returns the result.
func FormatNodes(nodes []*Node, format string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	nodesLen := len(nodes)
	for k, node := range nodes {
		tpl, err := formatNodeTemplate(
			&templateNode{
				Node:    node,
				isFirst: k == 0,
				isLast:  k == (nodesLen - 1),
			},
		).Parse(format)
		if err != nil {
			return nil, err
		}

		if err := tpl.Execute(buf, node.InterfaceValue()); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

type formatTemplateFuncs struct {
	node *templateNode
}

func (funcs *formatTemplateFuncs) funcMap() template.FuncMap {
	return template.FuncMap{
		"query":          funcs.query,
		"queryMultiple":  funcs.queryMultiple,
		"select":         funcs.query,
		"selectMultiple": funcs.queryMultiple,
		"format":         funcs.format,
		"isFirst":        funcs.isFirst,
		"isLast":         funcs.isLast,
		"newline":        funcs.newline,
		"joinToString":   funcs.joinToString,
	}
}

func (funcs *formatTemplateFuncs) newline() string {
	return "\n"
}

func (funcs *formatTemplateFuncs) isFirst() bool {
	return funcs.node.isFirst
}

func (funcs *formatTemplateFuncs) isLast() bool {
	return funcs.node.isLast
}

func (funcs *formatTemplateFuncs) query(selector string) *Node {
	res, err := funcs.node.Query(selector)
	if err != nil {
		return nil
	}
	return res
}

func (funcs *formatTemplateFuncs) queryMultiple(selector string) []*Node {
	res, err := funcs.node.QueryMultiple(selector)
	if err != nil {
		return nil
	}
	return res
}

func (funcs *formatTemplateFuncs) joinToString(target interface{}, separator string) string {
	targetReflect := reflect.ValueOf(target)

	var elems []string

	switch targetReflect.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetReflect.Len(); i++ {
			el := targetReflect.Index(i)
			elems = append(elems, fmt.Sprintf("%s", el))
		}
	default:
		elems = append(elems, fmt.Sprintf("%s", target))
	}

	return strings.Join(elems, separator)
}

func (funcs *formatTemplateFuncs) format(format string, target interface{}) string {
	switch t := target.(type) {
	case []*Node:
		buf, err := FormatNodes(t, format)
		if err != nil {
			return err.Error()
		}
		res := buf.String()
		return res
	case *Node:
		buf, err := FormatNode(t, format)
		if err != nil {
			return err.Error()
		}
		return buf.String()
	}

	return "<nil>"
}

func formatNodeTemplate(node *templateNode) *template.Template {
	funcs := &formatTemplateFuncs{
		node: node,
	}
	return template.
		New("nodeFormat").
		Funcs(sprig.TxtFuncMap()).
		Funcs(funcs.funcMap())
}
