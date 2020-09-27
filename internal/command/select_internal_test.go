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
		})

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

func selectTestForParser(parser string, data string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("SingleProperty", selectTest(data, parser, ".id", "1111", nil))
		t.Run("ObjectProperty", selectTest(data, parser, ".details.name", "Tom", nil))
		t.Run("Index", selectTest(data, parser, ".details.addresses.[0].street", "101 Some Street", nil))
		t.Run("Index", selectTest(data, parser, ".details.addresses.[1].street", "34 Another Street", nil))
		t.Run("DynamicString", selectTest(data, parser, ".details.addresses.(postcode=XXX XXX).street", "101 Some Street", nil))
		t.Run("DynamicString", selectTest(data, parser, ".details.addresses.(postcode=YYY YYY).street", "34 Another Street", nil))
	}
}

func TestSelect(t *testing.T) {
	t.Run("JSON", selectTestForParser("json", jsonData))
	t.Run("YAML", selectTestForParser("yaml", yamlData))
	t.Run("TOML", selectTestForParser("toml", tomlData))
}
