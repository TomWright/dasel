package cli_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
	"github.com/tomwright/dasel/v3/parsing/toml"
	"github.com/tomwright/dasel/v3/parsing/yaml"
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
			args = append(args, "-i", i.format.String(), "-o", o.format.String())
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
	jsonInputData := newStringWithFormat(json.JSON, `{
	"oneTwoThree": 123,
	"oneTwoDotThree": 12.3,
	"hello": "world",
	"boolFalse": false,
	"boolTrue": true,
	"stringFalse": "false",
	"stringTrue": "true",
	"sliceOfNumbers": [1, 2, 3, 4, 5],
	"mapData": {
		"oneTwoThree": 123,
		"oneTwoDotThree": 12.3,
		"hello": "world",
		"boolFalse": false,
		"boolTrue": true,
		"stringFalse": "false",
		"stringTrue": "true",
		"sliceOfNumbers": [1, 2, 3, 4, 5],
		"mapData": {
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
	yamlInputData := newStringWithFormat(yaml.YAML, `oneTwoThree: 123
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
mapData:
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
    mapData:
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

	tomlInputData := newStringWithFormat(toml.TOML, `
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = 'world'
boolFalse = false
boolTrue = true
stringFalse = 'false'
stringTrue = 'true'
sliceOfNumbers = [1, 2, 3, 4, 5]

[mapData]
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = 'world'
boolFalse = false
boolTrue = true
stringFalse = 'false'
stringTrue = 'true'
sliceOfNumbers = [1, 2, 3, 4, 5]

[mapData.mapData]
oneTwoThree = 123
oneTwoDotThree = 12.3
hello = 'world'
boolFalse = false
boolTrue = true
stringFalse = 'false'
stringTrue = 'true'
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
						newStringWithFormat(json.JSON, `"world"`),
						newStringWithFormat(yaml.YAML, `world`),
						newStringWithFormat(toml.TOML, `'world'`),
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
						newStringWithFormat(json.JSON, `123`),
						newStringWithFormat(yaml.YAML, `123`),
						newStringWithFormat(toml.TOML, `123`),
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
						newStringWithFormat(json.JSON, `12.3`),
						newStringWithFormat(yaml.YAML, `12.3`),
						newStringWithFormat(toml.TOML, `12.3`),
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
							newStringWithFormat(json.JSON, `true`),
							newStringWithFormat(yaml.YAML, `true`),
							newStringWithFormat(toml.TOML, `true`),
						},
					}.run)
					t.Run("false", testCases{
						selector: prefix + "boolFalse",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(json.JSON, `false`),
							newStringWithFormat(yaml.YAML, `false`),
							newStringWithFormat(toml.TOML, `false`),
						},
					}.run)
					t.Run("true string", testCases{
						selector: prefix + "stringTrue",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(json.JSON, `"true"`),
							newStringWithFormat(yaml.YAML, `"true"`),
							newStringWithFormat(toml.TOML, `'true'`),
						},
					}.run)
					t.Run("false string", testCases{
						selector: prefix + "stringFalse",
						in: []bytesWithFormat{
							jsonInputData,
							yamlInputData,
						},
						out: []bytesWithFormat{
							newStringWithFormat(json.JSON, `"false"`),
							newStringWithFormat(yaml.YAML, `"false"`),
							newStringWithFormat(toml.TOML, `'false'`),
						},
					}.run)
				})
			}
		}

		t.Run("root", newTestsWithPrefix(""))
		t.Run("nested once", newTestsWithPrefix("mapData."))
		t.Run("nested twice", newTestsWithPrefix("mapData.mapData."))
	})
}
