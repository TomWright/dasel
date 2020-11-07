package command

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

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

func selectTest(in string, parser string, selector string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		outputBuffer := bytes.NewBuffer([]byte{})

		err := runSelectCommand(selectOptions{
			Parser:   parser,
			Selector: selector,
			Reader:   strings.NewReader(in),
			Writer:   outputBuffer,
		}, nil)

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

func newline(x string) string {
	return x + "\n"
}

func TestSelect_JSON(t *testing.T) {
	t.Run("SingleProperty", selectTest(jsonData, "json", ".id", newline(`"1111"`), nil))
	t.Run("ObjectProperty", selectTest(jsonData, "json", ".details.name", newline(`"Tom"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[0].street", newline(`"101 Some Street"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[1].street", newline(`"34 Another Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=XXX XXX).street", newline(`"101 Some Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=YYY YYY).street", newline(`"34 Another Street"`), nil))
}

func TestSelect_YAML(t *testing.T) {
	t.Run("SingleProperty", selectTest(yamlData, "yaml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(yamlData, "yaml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
}

func TestSelect_TOML(t *testing.T) {
	t.Run("SingleProperty", selectTest(tomlData, "toml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(tomlData, "toml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
}
