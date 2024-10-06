package cli_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
)

func newStringWithFormat(format parsing.Format, data string) bytesWithFormat {
	return bytesWithFormat{
		format: format,
		data:   append([]byte(data), []byte("\n")...),
	}
}

type bytesWithFormat struct {
	format parsing.Format
	data   []byte
}

type testCases struct {
	selector string
	in       []bytesWithFormat
	out      []bytesWithFormat
	args     []string
	skip     []string
}

func (tcs testCases) run(t *testing.T) {
	for _, i := range tcs.in {
		for _, o := range tcs.out {
			tcName := fmt.Sprintf("%s to %s", i.format.String(), o.format.String())

			if slices.Contains(tcs.skip, tcName) {
				// Run a test and skip for visibility.
				t.Run(tcName, func(t *testing.T) {
					t.Skip()
				})
				continue
			}

			args := slices.Clone(tcs.args)
			args = append(args, "--input", i.format.String(), "--output", o.format.String())
			if tcs.selector != "" {
				args = append(args, tcs.selector)
			}
			tc := testCase{
				args:   args,
				in:     i.data,
				stdout: o.data,
			}
			t.Run(tcName, runTest(tc))
		}
	}
}

func TestCrossFormatHappyPath(t *testing.T) {
	jsonInputData := newStringWithFormat(parsing.JSON, `{
	"oneTwoThree": 123,
	"oneTwoDotThree": 12.3,
	"hello": "world",
	"boolFalse": false,
	"boolTrue": true,
	"stringFalse": "false",
	"stringTrue": "true",
	"sliceOfNumbers": [1, 2, 3, 4, 5],
	"map": {
		"oneTwoThree": 123,
		"oneTwoDotThree": 12.3,
		"hello": "world",
		"boolFalse": false,
		"boolTrue": true,
		"stringFalse": "false",
		"stringTrue": "true",
		"sliceOfNumbers": [1, 2, 3, 4, 5],
		"map": {
			"oneTwoThree": 123,
			"oneTwoDotThree": 12.3,
			"hello": "world",
			"boolFalse": false,
			"boolTrue": true,
			"stringFalse": "false",
			"stringTrue": "true",
			"sliceOfNumbers": [1, 2, 3, 4, 5]
		}
	}
}`)
	yamlInputData := newStringWithFormat(parsing.YAML, `oneTwoThree: 123
oneTwoDotThree: 12.3
hello: world
boolFalse: false
boolTrue: true
stringFalse: "false"
stringTrue: "true"
sliceOfNumbers:
- 1
- 2
- 3
- 4
- 5
map:
    oneTwoThree: 123
    oneTwoDotThree: 12.3
    hello: world
    boolFalse: false
    boolTrue: true
    stringFalse: "false"
    stringTrue: "true"
    sliceOfNumbers:
    - 1
    - 2
    - 3
    - 4
    - 5
    map:
        oneTwoThree: 123
        oneTwoDotThree: 12.3
        hello: world
        boolFalse: false
        boolTrue: true
        stringFalse: "false"
        stringTrue: "true"
        sliceOfNumbers:
        - 1
        - 2
        - 3
        - 4
        - 5
`)

	tomlInputData := newStringWithFormat(parsing.TOML, `
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = 'world'
boolFalse = false
boolTrue = true
stringFalse = 'false'
stringTrue = 'true'
sliceOfNumbers = [1, 2, 3, 4, 5]

[map]
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = "world"
boolFalse = false
boolTrue = true
stringFalse = "false"
stringTrue = "true"
sliceOfNumbers = [1, 2, 3, 4, 5]

[map.map]
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = "world"
boolFalse = false
boolTrue = true
stringFalse = "false"
stringTrue = "true"
sliceOfNumbers = [1, 2, 3, 4, 5]
`)

	t.Run("select", func(t *testing.T) {
		newTestsWithPrefix := func(prefix string) func(*testing.T) {
			return func(t *testing.T) {
				t.Run("string", testCases{
					selector: prefix + "hello",
					in: []bytesWithFormat{
						jsonInputData,
						yamlInputData,
						tomlInputData,
					},
					out: []bytesWithFormat{
						newStringWithFormat(parsing.JSON, `"world"`),
						newStringWithFormat(parsing.YAML, `world`),
						newStringWithFormat(parsing.TOML, `'world'`),
					},
				}.run)
				t.Run("int", testCases{
					selector: prefix + "oneTwoThree",
					in: []bytesWithFormat{
						jsonInputData,
						yamlInputData,
						tomlInputData,
					},
					out: []bytesWithFormat{
						newStringWithFormat(parsing.JSON, `123`),
						newStringWithFormat(parsing.YAML, `123`),
						newStringWithFormat(parsing.TOML, `123`),
					},
					skip: []string{
						// Skipped because the parser outputs as a float.
						"json to toml",
					},
				}.run)
				t.Run("float", testCases{
					selector: prefix + "oneTwoDotThree",
					in: []bytesWithFormat{
						jsonInputData,
						yamlInputData,
						tomlInputData,
					},
					out: []bytesWithFormat{
						newStringWithFormat(parsing.JSON, `12.3`),
						newStringWithFormat(parsing.YAML, `12.3`),
						newStringWithFormat(parsing.TOML, `12.3`),
					},
				}.run)
				t.Run("bool", func(t *testing.T) {
					t.Run("true", testCases{
						selector: prefix + "boolTrue",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(parsing.JSON, `true`),
							newStringWithFormat(parsing.YAML, `true`),
							newStringWithFormat(parsing.TOML, `true`),
						},
					}.run)
					t.Run("false", testCases{
						selector: prefix + "boolFalse",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(parsing.JSON, `false`),
							newStringWithFormat(parsing.YAML, `false`),
							newStringWithFormat(parsing.TOML, `false`),
						},
					}.run)
					t.Run("true string", testCases{
						selector: prefix + "stringTrue",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(parsing.JSON, `"true"`),
							newStringWithFormat(parsing.YAML, `"true"`),
							newStringWithFormat(parsing.TOML, `'true'`),
						},
					}.run)
					t.Run("false string", testCases{
						selector: prefix + "stringFalse",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(parsing.JSON, `"false"`),
							newStringWithFormat(parsing.YAML, `"false"`),
							newStringWithFormat(parsing.TOML, `'false'`),
						},
					}.run)
				})
			}
		}

		t.Run("root", newTestsWithPrefix(""))
		t.Run("nested once", newTestsWithPrefix("map."))
		t.Run("nested twice", newTestsWithPrefix("map.map."))
	})
}
