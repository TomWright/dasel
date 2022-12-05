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
				name:  "float round number",
				t:     "float",
				value: "123",
				exp:   `123`,
			},
			{
				name:  "float 1 decimal place",
				t:     "float",
				value: "123.4",
				exp:   `123.4`,
			},
			{
				name:  "float 5 decimal place",
				t:     "float",
				value: "123.45678",
				exp:   `123.45678`,
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

	t.Run("SetStringOnExistingIndex", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "z", "[1]"},
		[]byte(`["a","b","c"]`),
		newline([]byte(`["a","z","c"]`)),
		nil,
		nil,
	))

	t.Run("SetStringOnExistingNestedIndex", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "z", "[0].[1]"},
		[]byte(`[["a","b","c"],["d","e","f"]]`),
		newline([]byte(`[["a","z","c"],["d","e","f"]]`)),
		nil,
		nil,
	))

	t.Run("AppendStringIndexToRoot", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "z", "[]"},
		[]byte(`[]`),
		newline([]byte(`["z"]`)),
		nil,
		nil,
	))

	t.Run("AppendStringIndexToNestedSlice", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "z", "[0].[]"},
		[]byte(`[[]]`),
		newline([]byte(`[["z"]]`)),
		nil,
		nil,
	))

	t.Run("AppendToChainOfMissingSlicesAndProperties", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "users.[].name.first"},
		[]byte(`{}`),
		newline([]byte(`{"users":[{"name":{"first":"Tom"}}]}`)),
		nil,
		nil,
	))

	t.Run("AppendToEmptyExistingSlice", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "users.[]"},
		[]byte(`{"users":[]}`),
		newline([]byte(`{"users":["Tom"]}`)),
		nil,
		nil,
	))

	t.Run("AppendToEmptyMissingSlice", runTest(
		[]string{"put", "-r", "json", "-t", "string", "--pretty=false", "-v", "Tom", "users.[]"},
		[]byte(`{}`),
		newline([]byte(`{"users":["Tom"]}`)),
		nil,
		nil,
	))
}
