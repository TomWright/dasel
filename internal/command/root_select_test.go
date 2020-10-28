package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"io/ioutil"
	"strings"
	"testing"
)

var jsonDataSingle = `{"x": "asd"}`
var yamlDataSingle = `x: asd`
var tomlDataSingle = `x="asd"`
var xmlDataSingle = `<x>asd</x>`

var jsonData = `{
  "id": "1111",
  "details": {
    "name": "Tom",
  	"age": 27,
    "addresses": [
      {
        "street": "101 Some Street",
        "town": "Some Town",
        "county": "Some Country",
        "postcode": "XXX XXX",
        "primary": true
      },
      {
        "street": "34 Another Street",
        "town": "Another Town",
        "county": "Another County",
        "postcode": "YYY YYY"
      }
    ]
  }
}`

var yamlData = `
id: 1111
details:
  name: Tom
  age: 27
  addresses:
  - street: 101 Some Street
    town: Some Town
    county: Some County
    postcode: XXX XXX
    primary: true
  - street: 34 Another Street
    town: Another Town
    county: Another County
    postcode: YYY YYY
`

var tomlData = `id = "1111"
[details]
  name = "Tom"
  age = 27
  [[details.addresses]]
    street =  "101 Some Street"
    town = "Some Town"
    county = "Some County"
    postcode = "XXX XXX"
    primary = true
  [[details.addresses]]
    street = "34 Another Street"
    town = "Another Town"
    county = "Another County"
    postcode = "YYY YYY"
`

var xmlData = `<data>
	<id>1111</id>
	<details>
		<name>Tom</name>
		<age>27</age>
		<addresses primary="true">
			<street>101 Some Street</street>
			<town>Some Town</town>
			<county>Some County</county>
			<postcode>XXX XXX</postcode>
		</addresses>
		<addresses>
			<street>34 Another Street</street>
			<town>Another Town</town>
			<county>Another County</county>
			<postcode>YYY YYY</postcode>
		</addresses>
	</details>
</data>
`

func TestRootCMD_Select(t *testing.T) {
	t.Run("JSON", selectTestForParser("json", jsonData, jsonDataSingle))
	t.Run("YAML", selectTestForParser("yaml", yamlData, yamlDataSingle))
	t.Run("TOML", selectTestForParser("toml", tomlData, tomlDataSingle))
	t.Run("InvalidFile", expectErr(
		[]string{"select", "-f", "bad.json", "-s", "x"},
		"could not open input file",
	))
	t.Run("MissingParser", expectErr(
		[]string{"select", "-s", "x"},
		"parser flag required when reading from stdin",
	))
	t.Run("Stdin", expectOutput(
		`{"name": "Tom"}`,
		[]string{"select", "-f", "stdin", "-p", "json", "-s", ".name"},
		`Tom
`,
	))
	t.Run("StdinAlias", expectOutput(
		`{"name": "Tom"}`,
		[]string{"select", "-f", "-", "-p", "json", "-s", ".name"},
		`Tom
`,
	))
}

func selectTest(in string, parser string, selector string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"select", "-p", parser, selector,
		}

		cmd.SetOut(outputBuffer)
		cmd.SetIn(strings.NewReader(in))
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if out != string(output) {
			t.Errorf("expected result %v, got %v", out, string(output))
		}
	}
}

func selectTestFromFile(inputPath string, selector string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"select", "-f", inputPath, "-s", selector,
		}

		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if out != string(output) {
			t.Errorf("expected result %v, got %v", out, string(output))
		}
	}
}

func TestRootCMD_Select_XML(t *testing.T) {
	t.Run("RootElement", selectTest(xmlDataSingle, "xml", ".", "map[x:asd]\n", nil))
	t.Run("SingleProperty", selectTest(xmlData, "xml", ".data.id", "1111\n", nil))
	t.Run("ObjectProperty", selectTest(xmlData, "xml", ".data.details.name", "Tom\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[0].street", "101 Some Street\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[1].street", "34 Another Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=XXX XXX).street", "101 Some Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=YYY YYY).street", "34 Another Street\n", nil))
	t.Run("Attribute", selectTest(xmlData, "xml", ".data.details.addresses.(-primary=true).street", "101 Some Street\n", nil))
}

func selectTestForParser(parser string, data string, singleData string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("RootElement", selectTest(singleData, parser, ".", "map[x:asd]\n", nil))
		t.Run("SingleProperty", selectTest(data, parser, ".id", "1111\n", nil))
		t.Run("ObjectProperty", selectTest(data, parser, ".details.name", "Tom\n", nil))
		t.Run("Index", selectTest(data, parser, ".details.addresses.[0].street", "101 Some Street\n", nil))
		t.Run("Index", selectTest(data, parser, ".details.addresses.[1].street", "34 Another Street\n", nil))
		t.Run("DynamicString", selectTest(data, parser, ".details.addresses.(postcode=XXX XXX).street", "101 Some Street\n", nil))
		t.Run("DynamicString", selectTest(data, parser, ".details.addresses.(postcode=YYY YYY).street", "34 Another Street\n", nil))

		switch parser {
		case "json":
			t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.json", ".preferences.favouriteColour", "red\n", nil))
		case "yaml":
			t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.yaml", ".preferences.favouriteColour", "red\n", nil))
		}
	}
}
