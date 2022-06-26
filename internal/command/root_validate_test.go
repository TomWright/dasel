package command_test

import (
	"testing"
)

func TestRootCMD_Validate(t *testing.T) {
	t.Run("AllPass", expectOutputAndErr(
		[]string{"validate", "./../../tests/assets/example.json", "./../../tests/assets/example.yaml"},
		"",
		`pass ./../../tests/assets/example.json
pass ./../../tests/assets/example.yaml
`,
	))

	t.Run("PartialFail", expectOutputAndErr(
		[]string{"validate", "./../../tests/assets/example.json", "./../../tests/assets/example.yaml", "./../../tests/assets/broken.json", "./../../tests/assets/broken.xml"},
		"2 files failed validation",
		`pass ./../../tests/assets/example.json
pass ./../../tests/assets/example.yaml
fail ./../../tests/assets/broken.json could not load input: could not unmarshal data: invalid character '}' after array element
fail ./../../tests/assets/broken.xml could not load input: could not unmarshal data: xml.Decoder.Token() - XML syntax error on line 1: element <a> closed by </b>
`,
	))

	t.Run("AllFail", expectOutputAndErr(
		[]string{"validate", "./../../tests/assets/broken.json", "./../../tests/assets/broken.xml"},
		"2 files failed validation",
		`fail ./../../tests/assets/broken.json could not load input: could not unmarshal data: invalid character '}' after array element
fail ./../../tests/assets/broken.xml could not load input: could not unmarshal data: xml.Decoder.Token() - XML syntax error on line 1: element <a> closed by </b>
`,
	))

	t.Run("NoFilesPass", expectOutputAndErr(
		[]string{"validate"},
		"",
		``,
	))
}
