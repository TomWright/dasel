package cli_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
)

type inputProvider interface {
	Format() parsing.Format
	UserObject() []byte
	ListOfNumbers() []byte
	ListOfStrings() []byte
	UserName() []byte
}

type jsonInputProvider struct{}

func (j jsonInputProvider) Format() parsing.Format {
	return parsing.JSON
}

func (j jsonInputProvider) UserObject() []byte {
	return []byte(`{"name":"Tom"}`)
}

func (j jsonInputProvider) ListOfNumbers() []byte {
	return []byte(`[1,2,3]`)
}

func (j jsonInputProvider) ListOfStrings() []byte {
	return []byte(`["a","b","c"]`)
}

func (j jsonInputProvider) UserName() []byte {
	return []byte(`"Tom"`)
}

type yamlInputProvider struct{}

func (y yamlInputProvider) Format() parsing.Format {
	return parsing.YAML
}

func (y yamlInputProvider) UserObject() []byte {
	return []byte(`name: Tom`)
}

func (y yamlInputProvider) ListOfNumbers() []byte {
	return []byte(`- 1
- 2
- 3`)
}

func (y yamlInputProvider) ListOfStrings() []byte {
	return []byte(`- a
- b
- c`)
}

func (y yamlInputProvider) UserName() []byte {
	return []byte(`Tom`)
}

type tomlInputProvider struct{}

func (t tomlInputProvider) Format() parsing.Format {
	return parsing.TOML
}

func (t tomlInputProvider) UserObject() []byte {
	return []byte(`name = "Tom"`)
}

func (t tomlInputProvider) ListOfNumbers() []byte {
	return []byte(`[1, 2, 3]`)
}

func (t tomlInputProvider) ListOfStrings() []byte {
	return []byte(`["a", "b", "c"]`)
}

func (t tomlInputProvider) UserName() []byte {
	return []byte(`Tom`)
}

func TestGeneric(t *testing.T) {
	t.Run("json", runGenericTests(jsonInputProvider{}))
	//t.Run("yaml", runGenericTests(yamlInputProvider{}))
	//t.Run("toml", runGenericTests(tomlInputProvider{}))
}

func runGenericTests(i inputProvider) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("root", runTest(testCase{
			args:   []string{"-i", i.Format().String(), ``},
			in:     i.UserObject(),
			stdout: i.UserObject(),
		}))
		t.Run("top level string", runTest(testCase{
			args:   []string{"-i", i.Format().String(), `name`},
			in:     i.UserObject(),
			stdout: i.UserName(),
		}))
	}
}
