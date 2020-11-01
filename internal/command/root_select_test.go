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

func newline(x string) string {
	return x + "\n"
}

func TestRootCMD_Select(t *testing.T) {
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
		`"Tom"
`,
	))
	t.Run("StdinAlias", expectOutput(
		`{"name": "Tom"}`,
		[]string{"select", "-f", "-", "-p", "json", "-s", ".name"},
		`"Tom"
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

func TestRootCmd_Select_JSON(t *testing.T) {
	t.Run("RootElement", selectTest(jsonDataSingle, "json", ".", newline(`{
  "x": "asd"
}`), nil))
	t.Run("SingleProperty", selectTest(jsonData, "json", ".id", newline(`"1111"`), nil))
	t.Run("ObjectProperty", selectTest(jsonData, "json", ".details.name", newline(`"Tom"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[0].street", newline(`"101 Some Street"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[1].street", newline(`"34 Another Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=XXX XXX).street", newline(`"101 Some Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=YYY YYY).street", newline(`"34 Another Street"`), nil))
	t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.json", ".preferences.favouriteColour", newline(`"red"`), nil))

	t.Run("SubSelector", selectTest(`{
  "users": [
	{
	  "primary": true,
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  }
	},
	{
	  "primary": false,
	  "name": {
		"first": "Jim",
		"last": "Wright"
	  }
	}
  ]
}`, "json", ".users.(name.first=Tom).primary", newline(`true`), nil))

	t.Run("SubSubSelector", selectTest(`{
  "users": [
	{
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  },
      "addresses": [
        {
          "primary": true,
          "number": 123
        },
        {
          "primary": false,
          "number": 456
        }
      ]
	}
  ]
}`, "json", ".users.(.addresses.(.primary=true).number=123).name.first", newline(`"Tom"`), nil))

	t.Run("SubSubAndSelector", selectTest(`{
  "users": [
	{
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  },
      "addresses": [
        {
          "primary": true,
          "number": 123
        },
        {
          "primary": false,
          "number": 456
        }
      ]
	}
  ]
}`, "json", ".users.(.addresses.(.primary=true).number=123)(.name.last=Wright).name.first", newline(`"Tom"`), nil))
}

func TestRootCmd_Select_YAML(t *testing.T) {
	t.Run("RootElement", selectTest(yamlDataSingle, "yaml", ".", newline(`x: asd`), nil))
	t.Run("SingleProperty", selectTest(yamlData, "yaml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(yamlData, "yaml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
	t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.yaml", ".preferences.favouriteColour", newline(`red`), nil))
}

func TestRootCmd_Select_TOML(t *testing.T) {
	t.Run("RootElement", selectTest(tomlDataSingle, "toml", ".", newline(`x = "asd"`), nil))
	t.Run("SingleProperty", selectTest(tomlData, "toml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(tomlData, "toml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
}

func TestRootCMD_Select_XML(t *testing.T) {
	t.Run("RootElement", selectTest(xmlDataSingle, "xml", ".", newline(`<x>asd</x>`), nil))
	t.Run("SingleProperty", selectTest(xmlData, "xml", ".data.id", "1111\n", nil))
	t.Run("ObjectProperty", selectTest(xmlData, "xml", ".data.details.name", "Tom\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[0].street", "101 Some Street\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[1].street", "34 Another Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=XXX XXX).street", "101 Some Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=YYY YYY).street", "34 Another Street\n", nil))
	t.Run("Attribute", selectTest(xmlData, "xml", ".data.details.addresses.(-primary=true).street", "101 Some Street\n", nil))
}
