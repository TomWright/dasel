package dasel_test

import (
	"fmt"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
	"testing"
)

var (
	tom = map[string]interface{}{
		"name": "Tom",
		"age":  28,
	}
	amelia = map[string]interface{}{
		"name": "Amelia",
		"age":  26,
	}
	people = []map[string]interface{}{tom, amelia}
	mapC   = map[string]interface{}{
		"thing": "1",
	}
	mapB = map[string]interface{}{
		"c":      mapC,
		"people": people,
	}
	mapA = map[string]interface{}{
		"b": mapB,
	}
	mapRoot = map[string]interface{}{
		"a": mapA,
	}
)

func TestParseSelector(t *testing.T) {
	t.Run("NonIntIndex", func(t *testing.T) {
		_, err := dasel.ParseSelector(".[a]")
		exp := &dasel.InvalidIndexErr{Index: "a"}
		if err == nil || err.Error() != exp.Error() {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("InvalidDynamicComparison", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(x<2)")
		exp := &dasel.UnknownComparisonOperatorErr{Operator: "<"}
		if err == nil || err.Error() != exp.Error() {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
}

func TestNode_Query_File(t *testing.T) {
	tests := []struct {
		Name     string
		Selector string
		Exp      string
	}{
		{Name: "Property", Selector: "name", Exp: "Tom"},
		{Name: "ChildProperty", Selector: "preferences.favouriteColour", Exp: "red"},
		{Name: "Index", Selector: "colours.[0]", Exp: "red"},
		{Name: "Index", Selector: "colours.[1]", Exp: "green"},
		{Name: "Index", Selector: "colours.[2]", Exp: "blue"},
		{Name: "IndexProperty", Selector: "colourCodes.[0].name", Exp: "red"},
		{Name: "IndexProperty", Selector: "colourCodes.[1].name", Exp: "green"},
		{Name: "IndexProperty", Selector: "colourCodes.[2].name", Exp: "blue"},
		{Name: "DynamicProperty", Selector: "colourCodes.(name=red).rgb", Exp: "ff0000"},
		{Name: "DynamicProperty", Selector: "colourCodes.(name=green).rgb", Exp: "00ff00"},
		{Name: "DynamicProperty", Selector: "colourCodes.(name=blue).rgb", Exp: "0000ff"},
		{Name: "MultipleDynamicProperty", Selector: "colourCodes.(name=red)(rgb=ff0000).name", Exp: "red"},
		{Name: "MultipleDynamicProperty", Selector: "colourCodes.(name=green)(rgb=00ff00).name", Exp: "green"},
		{Name: "MultipleDynamicProperty", Selector: "colourCodes.(name=blue)(rgb=0000ff).name", Exp: "blue"},
	}

	fileTest := func(filename string) func(t *testing.T) {
		return func(t *testing.T) {
			parser, err := storage.NewParserFromFilename(filename)
			if err != nil {
				t.Errorf("could not get parser: %s", err)
				return
			}

			value, err := storage.LoadFromFile(filename, parser)
			if err != nil {
				t.Errorf("could not load value from file: %s", err)
				return
			}

			for _, testCase := range tests {
				tc := testCase
				t.Run(tc.Name, func(t *testing.T) {
					node, err := dasel.New(value).Query(tc.Selector)
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}

					if exp, got := tc.Exp, fmt.Sprint(node.Value); exp != got {
						t.Errorf("expected value `%s`, got `%s`", exp, got)
					}
				})
			}
		}
	}

	t.Run("JSON", fileTest("./tests/assets/example.json"))
	t.Run("YAML", fileTest("./tests/assets/example.yaml"))
}

func TestNode_Query_Data(t *testing.T) {
	t.Run("ParentChildPathToProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.c.thing")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "1", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToIndexProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.[1].name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Amelia", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToDynamicProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.(name=Tom).name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Tom", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToMultipleDynamicProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.(name=Tom)(age=28).name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Tom", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
}

func TestNode_Put(t *testing.T) {
	data := map[string]interface{}{
		"people": []map[string]interface{}{
			{
				"id":   1,
				"name": "Tom",
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
		"names": []string{
			"Tom",
			"Jim",
		},
	}
	rootNode := dasel.New(data)

	t.Run("ExistingValue", func(t *testing.T) {
		err := rootNode.Put("people.(id=1).name", "Thomas")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("people.(id=1).name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Thomas", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ExistingIntValue", func(t *testing.T) {
		err := rootNode.Put("people.[0].id", 3)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("people.[0].id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := 3, got.Value.(int); exp != got {
			t.Errorf("expected %d, got %d", exp, got)
		}
	})
	t.Run("NewPropertyOnExistingObject", func(t *testing.T) {
		err := rootNode.Put("people.(id=3).age", 27)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("people.(id=3).age")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := 27, got.Value.(int); exp != got {
			t.Errorf("expected %d, got %d", exp, got)
		}
	})
	t.Run("AppendObjectToList", func(t *testing.T) {
		err := rootNode.Put("people.[]", map[string]interface{}{
			"id":   1,
			"name": "Bob",
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("people.[2].id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if exp, got := 1, got.Value.(int); exp != got {
			t.Errorf("expected %d, got %d", exp, got)
		}
		got, err = rootNode.Query("people.[2].name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if exp, got := "Bob", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("AppendStringToList", func(t *testing.T) {
		err := rootNode.Put("names.[]", "Bob")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("names.[2]")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if exp, got := "Bob", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("NilRootNode", func(t *testing.T) {
		rootNode := dasel.New(nil)
		err := rootNode.Put("name", "Thomas")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Thomas", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("NilChain", func(t *testing.T) {
		rootNode := dasel.New(nil)
		err := rootNode.Put("my.name", "Thomas")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("my.name")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Thomas", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("NilChainToListIndex", func(t *testing.T) {
		rootNode := dasel.New(nil)
		err := rootNode.Put("my.favourite.people.[0]", "Tom")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("my.favourite.people.[0]")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Tom", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("NilChainToListNextAvailableIndex", func(t *testing.T) {
		rootNode := dasel.New(nil)
		err := rootNode.Put("my.favourite.people.[]", "Tom")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := rootNode.Query("my.favourite.people.[0]")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		if exp, got := "Tom", got.Value.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
}
