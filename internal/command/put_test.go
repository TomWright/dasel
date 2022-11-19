package command

import (
	"fmt"
	"testing"
)

func TestPutCommand(t *testing.T) {

	t.Run("SetTypeOnExistingProperty", func(t *testing.T) {
		tests := []struct {
			name  string
			t     string
			value string
			exp   string
		}{
			{
				t:     "string",
				value: "some string",
				exp:   `"some string"`,
			},
			{
				t:     "int",
				value: "123",
				exp:   `123`,
			},
			{
				name:  "true bool",
				t:     "bool",
				value: "true",
				exp:   `true`,
			},
			{
				name:  "false bool",
				t:     "bool",
				value: "false",
				exp:   `false`,
			},
			{
				t:     "json",
				value: `{"some":"json"}`,
				exp:   `{"some":"json"}`,
			},
		}

		for _, test := range tests {
			tc := test
			if tc.name == "" {
				tc.name = tc.t
			}
			t.Run(tc.name, runTest(
				[]string{"put", "-r", "json", "-t", tc.t, "--pretty=false", "-v", tc.value, "val"},
				[]byte(`{"val":"oldVal"}`),
				newline([]byte(fmt.Sprintf(`{"val":%s}`, tc.exp))),
				nil,
				nil,
			))
		}
	})

	t.Run("SetStringOnExistingNestedProperty", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "user.name"},
		[]byte(`{"user":{"name":"oldName"}}`),
		newline([]byte(`{"user":{"name":"Tom"}}`)),
		nil,
		nil,
	))

	t.Run("CreateStringProperty", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "name"},
		[]byte(`{}`),
		newline([]byte(`{"name":"Tom"}`)),
		nil,
		nil,
	))

	t.Run("CreateNestedStringPropertyOnExistingParent", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "user.name"},
		[]byte(`{"user":{}}`),
		newline([]byte(`{"user":{"name":"Tom"}}`)),
		nil,
		nil,
	))

	t.Run("CreateNestedStringPropertyOnMissingParent", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "user.name"},
		[]byte(`{}`),
		newline([]byte(`{"user":{"name":"Tom"}}`)),
		nil,
		nil,
	))
}
