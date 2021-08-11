package dasel_test

import (
	"github.com/tomwright/dasel"
	"testing"
)

func testFormatNode(value interface{}, format string, exp string) func(t *testing.T) {
	return func(t *testing.T) {
		node := dasel.New(value)
		buf, err := dasel.FormatNode(node, format)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		got := buf.String()
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func testFormatNodes(values []interface{}, format string, exp string) func(t *testing.T) {
	return func(t *testing.T) {
		nodes := make([]*dasel.Node, len(values))
		for k, v := range values {
			nodes[k] = dasel.New(v)
		}
		buf, err := dasel.FormatNodes(nodes, format)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		got := buf.String()
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func TestFormatNode(t *testing.T) {
	t.Run("InvalidFormatTemplate", func(t *testing.T) {
		_, err := dasel.FormatNode(nil, "{{")
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})
	t.Run("PropertyAccess", testFormatNode(
		map[string]interface{}{
			"name":  "Tom",
			"email": "contact@tomwright.me",
		},
		`{{ .name }}, {{ .email }}`,
		`Tom, contact@tomwright.me`,
	))
	t.Run("QueryAccess", testFormatNode(
		map[string]interface{}{
			"name":  "Tom",
			"email": "contact@tomwright.me",
		},
		`{{ query ".name" }}, {{ query ".email" }}`,
		`Tom, contact@tomwright.me`,
	))
	t.Run("Format", testFormatNode(
		map[string]interface{}{
			"name":  "Tom",
			"email": "contact@tomwright.me",
		},
		`{{ query ".name" | format "{{ . }}" }}, {{ query ".email" | format "{{ . }}" }}`,
		`Tom, contact@tomwright.me`,
	))
	t.Run("QueryAccessInvalidSelector", testFormatNode(
		map[string]interface{}{
			"name":  "Tom",
			"email": "contact@tomwright.me",
		},
		`{{ query ".bad" }}`,
		`<nil>`,
	))
	t.Run("QueryMultipleCommaSeparated", testFormatNode(
		map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name": "Tom",
				},
				{
					"name": "Jim",
				},
				{
					"name": "Frank",
				},
			},
		},
		`{{ queryMultiple ".users.[*]" | format "{{ .name }}{{ if not isLast }},{{ end }}" }}`,
		`Tom,Jim,Frank`,
	))
	t.Run("QueryMultipleLineSeparated", testFormatNode(
		map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name": "Tom",
				},
				{
					"name": "Jim",
				},
				{
					"name": "Frank",
				},
			},
		},
		`{{ queryMultiple ".users.[*]" | format "{{ .name }}{{ if not isLast }}{{ newline }}{{ end }}" }}`,
		`Tom
Jim
Frank`,
	))
	t.Run("QueryMultipleDashSeparated", testFormatNode(
		map[string]interface{}{
			"users": []map[string]interface{}{
				{
					"name": "Tom",
				},
				{
					"name": "Jim",
				},
				{
					"name": "Frank",
				},
			},
		},
		`{{ queryMultiple ".users.[*]" | format "{{ if not isFirst }}---{{ newline }}{{ end }}{{ .name }}{{ if not isLast }}{{ newline }}{{ end }}" }}`,
		`Tom
---
Jim
---
Frank`,
	))
}

func TestFormatNodes(t *testing.T) {
	t.Run("InvalidFormatTemplate", func(t *testing.T) {
		_, err := dasel.FormatNodes([]*dasel.Node{dasel.New("")}, "{{")
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})
	t.Run("PropertyAccess", testFormatNodes(
		[]interface{}{
			map[string]interface{}{
				"name":  "Tom",
				"email": "contact@tomwright.me",
			},
			map[string]interface{}{
				"name":  "Jim",
				"email": "jim@gmail.com",
			},
		},
		"{{ .name }}, {{ .email }}{{ if not isLast }}{{ newline }}{{ end }}",
		`Tom, contact@tomwright.me
Jim, jim@gmail.com`,
	))
	t.Run("QueryAccess", testFormatNodes(
		[]interface{}{
			map[string]interface{}{
				"name":  "Tom",
				"email": "contact@tomwright.me",
			},
			map[string]interface{}{
				"name":  "Jim",
				"email": "jim@gmail.com",
			},
		},
		`{{ query ".name" }}, {{ query ".email" }}{{ if not isLast }}{{ newline }}{{ end }}`,
		`Tom, contact@tomwright.me
Jim, jim@gmail.com`))
	t.Run("QueryAccessInvalidSelector", testFormatNodes(
		[]interface{}{
			map[string]interface{}{
				"name":  "Tom",
				"email": "contact@tomwright.me",
			},
			map[string]interface{}{
				"name":  "Jim",
				"email": "jim@gmail.com",
			},
		},
		`{{ query ".bad" }}{{ newline }}`,
		`<nil>
<nil>
`,
	))
}
