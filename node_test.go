package dasel_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
)

// ExampleNode_ReadmeExample tests the code from the readme explanation.
func ExampleNode_ReadmeExample() {
	printNodeValue := func(nodes ...*dasel.Node) {
		for _, n := range nodes {
			fmt.Println(n.InterfaceValue())
		}
	}

	var data interface{}
	_ = json.Unmarshal([]byte(`[{"name": "Tom"}, {"name": "Jim"}]`), &data)

	rootNode := dasel.New(data)

	result, _ := rootNode.Query(".[0].name")
	printNodeValue(result) // Tom

	results, _ := rootNode.QueryMultiple(".[*].name")
	printNodeValue(results...) // Tom \n Jim

	_ = rootNode.Put(".[0].name", "Frank")
	printNodeValue(rootNode) // [ map[name:Frank] map[name:Jim] ]

	_ = rootNode.PutMultiple(".[*].name", "Joe")
	printNodeValue(rootNode) // [ map[name:Joe] map[name:Joe] ]

	outputBytes, _ := json.Marshal(rootNode.InterfaceValue())
	fmt.Println(string(outputBytes)) // [{"name": "Joe"}, {"name": "Joe"}]

	// Output:
	// Tom
	// Tom
	// Jim
	// [map[name:Frank] map[name:Jim]]
	// [map[name:Joe] map[name:Joe]]
	// [{"name":"Joe"},{"name":"Joe"}]
}

// ExampleNode_Query shows how to query data from go code.
func ExampleNode_Query() {
	myData := []byte(`{"name": "Tom"}`)
	var data interface{}
	if err := json.Unmarshal(myData, &data); err != nil {
		panic(err)
	}
	rootNode := dasel.New(data)
	result, err := rootNode.Query(".name")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InterfaceValue())

	// Output:
	// Tom
}

// ExampleNode_Put shows how to update data from go code.
func ExampleNode_Put() {
	myData := []byte(`{"name": "Tom"}`)
	var data interface{}
	if err := json.Unmarshal(myData, &data); err != nil {
		panic(err)
	}
	rootNode := dasel.New(data)
	if err := rootNode.Put(".name", "Jim"); err != nil {
		panic(err)
	}
	fmt.Println(rootNode.InterfaceValue())

	// Output:
	// map[name:Jim]
}

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

func testParseSelector(in string, exp dasel.Selector) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := dasel.ParseSelector(in)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func TestParseSelector(t *testing.T) {
	t.Run("NonIntIndex", func(t *testing.T) {
		_, err := dasel.ParseSelector(".[a]")
		exp := &dasel.InvalidIndexErr{Index: "a"}
		if err == nil || err.Error() != exp.Error() {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("InvalidDynamicBracketCount", func(t *testing.T) {
		_, err := dasel.ParseSelector(".((name=x)")
		exp := dasel.ErrDynamicSelectorBracketMismatch
		if err == nil || !errors.Is(err, exp) {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("InvalidDynamicComparison", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(x<=>2)")
		exp := &dasel.UnknownComparisonOperatorErr{Operator: "<=>"}
		if err == nil || err.Error() != exp.Error() {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("MultipleSearchGroups", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(?:a=b)(a=b)")
		exp := "require exactly 1 group in search selector"
		if err == nil || err.Error() != exp {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("UnknownComparisonOperator", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(a<=>b)")
		exp := "unknown comparison operator: <=>"
		if err == nil || err.Error() != exp {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("UnknownSearchComparisonOperator", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(?:a<=>b)")
		exp := "unknown comparison operator: <=>"
		if err == nil || err.Error() != exp {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("UnknownSearchKeyComparisonOperator", func(t *testing.T) {
		_, err := dasel.ParseSelector(".(?:->b)")
		exp := "unknown comparison operator: >"
		if err == nil || err.Error() != exp {
			t.Errorf("expected error %v, got %v", exp, err)
		}
	})
	t.Run("SearchEqual", testParseSelector(".(?:name=asd)", dasel.Selector{
		Raw:       ".(?:name=asd)",
		Current:   ".(?:name=asd)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.EqualCondition{
				Key:   "name",
				Value: "asd",
			},
		},
	}))
	t.Run("SearchNotEqual", testParseSelector(".(?:name!=asd)", dasel.Selector{
		Raw:       ".(?:name!=asd)",
		Current:   ".(?:name!=asd)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.EqualCondition{
				Key:   "name",
				Value: "asd",
				Not:   true,
			},
		},
	}))
	t.Run("SearchMoreThan", testParseSelector(".(?:name.[#]>3)", dasel.Selector{
		Raw:       ".(?:name.[#]>3)",
		Current:   ".(?:name.[#]>3)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				After: true,
			},
		},
	}))
	t.Run("SearchMoreThanEqual", testParseSelector(".(?:name.[#]>=3)", dasel.Selector{
		Raw:       ".(?:name.[#]>=3)",
		Current:   ".(?:name.[#]>=3)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				After: true,
				Equal: true,
			},
		},
	}))
	t.Run("SearchLessThan", testParseSelector(".(?:name.[#]<3)", dasel.Selector{
		Raw:       ".(?:name.[#]<3)",
		Current:   ".(?:name.[#]<3)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
			},
		},
	}))
	t.Run("SearchLessThanEqual", testParseSelector(".(?:name.[#]<=3)", dasel.Selector{
		Raw:       ".(?:name.[#]<=3)",
		Current:   ".(?:name.[#]<=3)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				Equal: true,
			},
		},
	}))
	t.Run("SearchKey", testParseSelector(".(?:-=asd)", dasel.Selector{
		Raw:       ".(?:-=asd)",
		Current:   ".(?:-=asd)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.KeyEqualCondition{
				Value: "asd",
			},
		},
	}))
	t.Run("SearchKeyNotEqual", testParseSelector(".(?:-!=asd)", dasel.Selector{
		Raw:       ".(?:-!=asd)",
		Current:   ".(?:-!=asd)",
		Remaining: "",
		Type:      "SEARCH",
		Conditions: []dasel.Condition{
			&dasel.KeyEqualCondition{
				Value: "asd",
				Not:   true,
			},
		},
	}))
	t.Run("DynamicKey", testParseSelector(".(-=asd)", dasel.Selector{
		Raw:       ".(-=asd)",
		Current:   ".(-=asd)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.KeyEqualCondition{
				Value: "asd",
			},
		},
	}))
	t.Run("DynamicKeyNotEqual", testParseSelector(".(-!=asd)", dasel.Selector{
		Raw:       ".(-!=asd)",
		Current:   ".(-!=asd)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.KeyEqualCondition{
				Value: "asd",
				Not:   true,
			},
		},
	}))
	t.Run("DynamicEqual", testParseSelector(".(name=asd)", dasel.Selector{
		Raw:       ".(name=asd)",
		Current:   ".(name=asd)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.EqualCondition{
				Key:   "name",
				Value: "asd",
			},
		},
	}))
	t.Run("DynamicNotEqual", testParseSelector(".(name!=asd)", dasel.Selector{
		Raw:       ".(name!=asd)",
		Current:   ".(name!=asd)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.EqualCondition{
				Key:   "name",
				Value: "asd",
				Not:   true,
			},
		},
	}))
	t.Run("DynamicMoreThan", testParseSelector(".(name.[#]>3)", dasel.Selector{
		Raw:       ".(name.[#]>3)",
		Current:   ".(name.[#]>3)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				After: true,
			},
		},
	}))
	t.Run("DynamicMoreThanEqual", testParseSelector(".(name.[#]>=3)", dasel.Selector{
		Raw:       ".(name.[#]>=3)",
		Current:   ".(name.[#]>=3)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				After: true,
				Equal: true,
			},
		},
	}))
	t.Run("DynamicLessThan", testParseSelector(".(name.[#]<3)", dasel.Selector{
		Raw:       ".(name.[#]<3)",
		Current:   ".(name.[#]<3)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
			},
		},
	}))
	t.Run("DynamicLessThanEqual", testParseSelector(".(name.[#]<=3)", dasel.Selector{
		Raw:       ".(name.[#]<=3)",
		Current:   ".(name.[#]<=3)",
		Remaining: "",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.SortedComparisonCondition{
				Key:   "name.[#]",
				Value: "3",
				Equal: true,
			},
		},
	}))
	t.Run("Unicode", testParseSelector(".(name=Ägir).b", dasel.Selector{
		Raw:       ".(name=Ägir).b",
		Current:   ".(name=Ägir)",
		Remaining: ".b",
		Type:      "DYNAMIC",
		Conditions: []dasel.Condition{
			&dasel.EqualCondition{
				Key:   "name",
				Value: "Ägir",
			},
		},
	}))
}

func extractValues(nodes []*dasel.Node) []interface{} {
	gotValues := make([]interface{}, len(nodes))
	for i, n := range nodes {
		gotValues[i] = n.InterfaceValue()
	}
	return gotValues
}

func testNodeQueryMultipleArray(selector string, expValues []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		nodes, err := dasel.New([]map[string]interface{}{
			{
				"name": "Tom",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}).QueryMultiple(selector)
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		gotValues := extractValues(nodes)

		if !valuesAreEqual(expValues, gotValues) {
			t.Errorf("expected %v, got %v", expValues, gotValues)
		}
	}
}

func testNodeQueryMultipleMap(selector string, expValues []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		nodes, err := dasel.New(map[string]interface{}{
			"personA": map[string]interface{}{
				"name": "Tom",
				"age":  27,
			},
			"personB": map[string]interface{}{
				"name": "Jim",
				"age":  27,
			},
			"personC": map[string]interface{}{
				"name": "Amelia",
				"age":  25,
			},
		}).QueryMultiple(selector)
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		gotValues := extractValues(nodes)

		if !valuesAreEqual(expValues, gotValues) {
			t.Errorf("expected %v, got %v", expValues, gotValues)
		}
	}
}

func TestNode_QueryMultiple(t *testing.T) {
	t.Run("SingleResult", testNodeQueryMultipleArray(".[0].name", []interface{}{
		"Tom",
	}))
	t.Run("SingleResultDynamicEqual", testNodeQueryMultipleArray(".(age=25).name", []interface{}{
		"Amelia",
	}))
	t.Run("SingleResultDynamicNotEqual", testNodeQueryMultipleArray(".(age!=27).name", []interface{}{
		"Amelia",
	}))
	t.Run("SingleResultDynamicKeyEqual", testNodeQueryMultipleArray(".(-=0).name", []interface{}{
		"Tom",
	}))
	t.Run("MultipleResultDynamicKeyNotEqual", testNodeQueryMultipleArray(".(-!=0).name", []interface{}{
		"Jim", "Amelia",
	}))
	t.Run("MultipleResultSearchKeyNotEqual", testNodeQueryMultipleArray(".[*].(?:-!=name)", []interface{}{
		"27", "27", "25",
	}))
	t.Run("MultipleResultDynamic", testNodeQueryMultipleArray(".(age=27).name", []interface{}{
		"Tom",
		"Jim",
	}))
	t.Run("MultipleResultAnyIndex", testNodeQueryMultipleArray(".[*].name", []interface{}{
		"Tom",
		"Jim",
		"Amelia",
	}))

	t.Run("MapSingleResult", testNodeQueryMultipleMap(".personA.name", []interface{}{
		"Tom",
	}))
	t.Run("MapSingleResultDynamic", testNodeQueryMultipleMap(".(age=25).name", []interface{}{
		"Amelia",
	}))
	t.Run("MapSingleResultDynamic", testNodeQueryMultipleMap(".(age=27).name", []interface{}{
		"Tom",
		"Jim",
	}))
	t.Run("MapMultipleResultAnyIndex", testNodeQueryMultipleMap(".[*].name", []interface{}{
		"Tom",
		"Jim",
		"Amelia",
	}))
}

func valuesAreEqual(exp []interface{}, got []interface{}) bool {
	if len(exp) != len(got) {
		return false
	}
	for _, g := range got {
		found := false
		for _, e := range exp {
			if reflect.DeepEqual(g, e) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestNode_PutMultiple(t *testing.T) {
	t.Run("SingleResult", func(t *testing.T) {
		value := []map[string]interface{}{
			{
				"name": "Tom",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}
		err := dasel.New(value).PutMultiple(".[0].name", "Frank")
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		exp := []map[string]interface{}{
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}

		if !reflect.DeepEqual(exp, value) {
			t.Errorf("expected %v, got %v", exp, value)
		}
	})
	t.Run("SingleResultDynamic", func(t *testing.T) {
		value := []map[string]interface{}{
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}
		err := dasel.New(value).PutMultiple(".(age=25).name", "Frank")
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		exp := []map[string]interface{}{
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Frank",
				"age":  "25",
			},
		}

		if !reflect.DeepEqual(exp, value) {
			t.Errorf("expected %v, got %v", exp, value)
		}
	})
	t.Run("MultipleResultDynamic", func(t *testing.T) {
		value := []map[string]interface{}{
			{
				"name": "Tom",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}
		err := dasel.New(value).PutMultiple(".(age=27).name", "Frank")
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		exp := []map[string]interface{}{
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}

		if !reflect.DeepEqual(exp, value) {
			t.Errorf("expected %v, got %v", exp, value)
		}
	})
	t.Run("MultipleResultAnyIndex", func(t *testing.T) {
		value := []map[string]interface{}{
			{
				"name": "Tom",
				"age":  "27",
			},
			{
				"name": "Jim",
				"age":  "27",
			},
			{
				"name": "Amelia",
				"age":  "25",
			},
		}
		err := dasel.New(value).PutMultiple(".[*].name", "Frank")
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}

		exp := []map[string]interface{}{
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Frank",
				"age":  "27",
			},
			{
				"name": "Frank",
				"age":  "25",
			},
		}

		if !reflect.DeepEqual(exp, value) {
			t.Errorf("expected %v, got %v", exp, value)
		}
	})
}

func TestNode_Query(t *testing.T) {
	parser, err := storage.NewReadParserFromFilename("./tests/assets/example.json")
	if err != nil {
		t.Errorf("could not get parser: %s", err)
		return
	}

	value, err := storage.LoadFromFile("./tests/assets/example.json", parser)
	if err != nil {
		t.Errorf("could not load value from file: %s", err)
		return
	}

	t.Run("Valid", func(t *testing.T) {
		node, err := dasel.New(value).Query("preferences.favouriteColour")
		if err != nil {
			t.Errorf("unexpected query error: %s", err)
			return
		}
		if exp, got := "red", fmt.Sprint(node.InterfaceValue()); exp != got {
			t.Errorf("expected value `%s`, got `%s`", exp, got)
		}
	})
	t.Run("NotFound", func(t *testing.T) {
		_, err := dasel.New(value).Query(".colours.[0].a")
		expErr := fmt.Errorf("could not find value: selector [type:PROPERTY selector:.a] does not support value: [kind:string type:string] red")
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
	})
	t.Run("InvalidSelector", func(t *testing.T) {
		_, err := dasel.New(value).Query(".colours.[a]")
		expErr := fmt.Errorf("failed to parse selector: %w", &dasel.InvalidIndexErr{Index: "a"})
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
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
			parser, err := storage.NewReadParserFromFilename(filename)
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
						t.Errorf("unexpected query error: %s", err)
						return
					}

					if exp, got := tc.Exp, fmt.Sprint(node.InterfaceValue()); exp != got {
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
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if exp, got := "1", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToIndexProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.[1].name")
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if exp, got := "Amelia", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToDynamicProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.(name=Tom).name")
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if exp, got := "Tom", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("ParentChildPathToMultipleDynamicProperty", func(t *testing.T) {
		rootNode := dasel.New(mapRoot)

		got, err := rootNode.Query(".a.b.people.(name=Tom)(age=28).name")
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if exp, got := "Tom", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("DynamicPropertyOnMap", func(t *testing.T) {
		rootNode := dasel.New(map[string]interface{}{
			"personA": map[string]interface{}{
				"name": "Tom",
			},
			"personB": map[string]interface{}{
				"name": "Jim",
			},
		})

		got, err := rootNode.Query(".(name=Tom).name")
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if exp, got := "Tom", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
}

func putQueryTest(rootNode *dasel.Node, putSelector string, newValue interface{}, querySelector string) func(t *testing.T) {
	return func(t *testing.T) {
		err := rootNode.Put(putSelector, newValue)
		if err != nil {
			t.Errorf("unexpected put error: %v", err)
			return
		}

		got, err := rootNode.Query(querySelector)
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		if !reflect.DeepEqual(newValue, got.InterfaceValue()) {
			t.Errorf("expected %v, got %v", newValue, got.InterfaceValue())
		}
	}
}

func deleteTest(rootNode *dasel.Node, deleteSelector string, exp interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		err := rootNode.Delete(deleteSelector)
		if err != nil {
			t.Errorf("unexpected delete error: %v", err)
			return
		}

		gotVal := rootNode.InterfaceValue()
		if !reflect.DeepEqual(exp, gotVal) && exp != gotVal {
			t.Errorf("expected [%T] %v, got [%T] %v", exp, exp, gotVal, gotVal)
		}
	}
}

func deleteMultipleTest(rootNode *dasel.Node, deleteSelector string, exp interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		err := rootNode.DeleteMultiple(deleteSelector)
		if err != nil {
			t.Errorf("unexpected delete error: %v", err)
			return
		}

		if !reflect.DeepEqual(exp, rootNode.InterfaceValue()) {
			t.Errorf("expected %v, got %v", exp, rootNode.InterfaceValue())
		}
	}
}

func putQueryMultipleTest(rootNode *dasel.Node, putSelector string, newValue interface{}, querySelector string) func(t *testing.T) {
	return func(t *testing.T) {
		err := rootNode.PutMultiple(putSelector, newValue)
		if err != nil {
			t.Errorf("unexpected put error: %v", err)
			return
		}

		got, err := rootNode.QueryMultiple(querySelector)
		if err != nil {
			t.Errorf("unexpected query error: %v", err)
			return
		}

		for _, n := range got {
			if !reflect.DeepEqual(newValue, n.InterfaceValue()) {
				t.Errorf("expected %v, got %v", newValue, n.InterfaceValue())
			}
		}
	}
}

func TestNode_Put_Query(t *testing.T) {
	data := map[string]interface{}{
		"id": "123",
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

	t.Run("InvalidSelector", func(t *testing.T) {
		err := rootNode.Put("people.[a].name", "Thomas")
		expErr := fmt.Errorf("failed to parse selector: %w", &dasel.InvalidIndexErr{Index: "a"})
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
	})
	t.Run("ExistingSingleString", putQueryTest(rootNode, "id", "456", "id"))
	t.Run("ExistingStringValue", putQueryTest(rootNode, "people.[0].name", "Thomas", "people.(id=1).name"))
	t.Run("ExistingIntValue", putQueryTest(rootNode, "people.(id=1).id", 3, "people.(id=3).id"))
	t.Run("NewPropertyOnExistingObject", putQueryTest(rootNode, "people.(id=3).age", 27, "people.[0].age"))
	t.Run("AppendObjectToList", func(t *testing.T) {
		err := rootNode.Put("people.[]", map[string]interface{}{
			"id":   1,
			"name": "Bob",
		})
		if err != nil {
			t.Errorf("unexpected put error: %v", err)
			return
		}

		got, err := rootNode.Query("people.[2].id")
		if err != nil {
			t.Errorf("unexpected query [1] error: %v", err)
			return
		}
		if exp, got := 1, got.InterfaceValue().(int); exp != got {
			t.Errorf("expected %d, got %d", exp, got)
		}
		got, err = rootNode.Query("people.[2].name")
		if err != nil {
			t.Errorf("unexpected query [2] error: %v", err)
			return
		}
		if exp, got := "Bob", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("AppendStringToList", putQueryTest(rootNode, "names.[]", "Bob", "names.[2]"))
	t.Run("NilRootNode", putQueryTest(dasel.New(nil), "name", "Thomas", "name"))
	t.Run("NilChain", putQueryTest(dasel.New(nil), "my.name", "Thomas", "my.name"))
	t.Run("NilChainToListIndex", putQueryTest(dasel.New(nil), "my.favourite.people.[0]", "Tom", "my.favourite.people.[0]"))
	t.Run("NilChainToListNextAvailableIndex", putQueryTest(dasel.New(nil), "my.favourite.people.[]", "Tom", "my.favourite.people.[0]"))
	t.Run("NilChainToDynamic", putQueryTest(dasel.New(nil), "(name=Jim).name", "Tom", "[0].name"))
}

func TestNode_PutMultiple_Query(t *testing.T) {
	data := map[string]interface{}{
		"id": "123",
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

	t.Run("InvalidSelector", func(t *testing.T) {
		err := rootNode.PutMultiple("people.[a].name", "Thomas")
		expErr := fmt.Errorf("failed to parse selector: %w", &dasel.InvalidIndexErr{Index: "a"})
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
	})
	t.Run("ExistingSingleString", putQueryMultipleTest(rootNode, "id", "456", "id"))
	t.Run("ExistingStringValue", putQueryMultipleTest(rootNode, "people.[0].name", "Thomas", "people.(id=1).name"))
	t.Run("ExistingIntValue", putQueryMultipleTest(rootNode, "people.(id=1).id", 3, "people.(id=3).id"))
	t.Run("NewPropertyOnExistingObject", putQueryMultipleTest(rootNode, "people.(id=3).age", 27, "people.[0].age"))
	t.Run("AppendObjectToList", func(t *testing.T) {
		err := rootNode.PutMultiple("people.[]", map[string]interface{}{
			"id":   1,
			"name": "Bob",
		})
		if err != nil {
			t.Errorf("unexpected put error: %v", err)
			return
		}

		got, err := rootNode.Query("people.[2].id")
		if err != nil {
			t.Errorf("unexpected query [1] error: %v", err)
			return
		}
		if exp, got := 1, got.InterfaceValue().(int); exp != got {
			t.Errorf("expected %d, got %d", exp, got)
		}
		got, err = rootNode.Query("people.[2].name")
		if err != nil {
			t.Errorf("unexpected query [2] error: %v", err)
			return
		}
		if exp, got := "Bob", got.InterfaceValue().(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
	t.Run("AppendStringToList", putQueryMultipleTest(rootNode, "names.[]", "Bob", "names.[2]"))
	t.Run("NilRootNode", putQueryMultipleTest(dasel.New(nil), "name", "Thomas", "name"))
	t.Run("NilChain", putQueryMultipleTest(dasel.New(nil), "my.name", "Thomas", "my.name"))
	t.Run("NilChainToListIndex", putQueryMultipleTest(dasel.New(nil), "my.favourite.people.[0]", "Tom", "my.favourite.people.[0]"))
	t.Run("NilChainToListNextAvailableIndex", putQueryMultipleTest(dasel.New(nil), "my.favourite.people.[]", "Tom", "my.favourite.people.[0]"))
}

func TestNode_Delete(t *testing.T) {
	data := func() map[string]interface{} {
		return map[string]interface{}{
			"id": "123",
			"names": []string{
				"Tom",
				"Jim",
			},
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
		}
	}

	t.Run("InvalidSelector", func(t *testing.T) {
		err := dasel.New(data()).Delete("people.[a].name")
		expErr := fmt.Errorf("failed to parse selector: %w", &dasel.InvalidIndexErr{Index: "a"})
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
	})
	t.Run("ExistingSingleString", deleteTest(dasel.New(data()), "id", map[string]interface{}{
		"names": []string{
			"Tom",
			"Jim",
		},
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
	}))
	t.Run("ExistingStringValue", deleteTest(dasel.New(data()), "people.[0].name", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"id": 1,
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
	}))
	t.Run("ExistingObjectInArray", deleteTest(dasel.New(data()), "people.[0]", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"id":   2,
				"name": "Jim",
			},
		},
	}))
	t.Run("ExistingIntValue", deleteTest(dasel.New(data()), "people.(id=1).id", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"name": "Tom",
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
	}))
	t.Run("ExistingArray", deleteMultipleTest(dasel.New(data()), "names", map[string]interface{}{
		"id": "123",
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
	}))
	t.Run("RootNodeObject", deleteTest(dasel.New(map[string]interface{}{
		"name": "Tom",
	}), ".", map[string]interface{}{}))
	t.Run("RootNodeArray", deleteTest(dasel.New([]interface{}{
		"name", "Tom",
	}), ".", []interface{}{}))
	t.Run("RootNodeUnknown", deleteTest(dasel.New(false), ".", map[string]interface{}{}))
}

func TestNode_DeleteMultiple_Query(t *testing.T) {
	data := func() map[string]interface{} {
		return map[string]interface{}{
			"id": "123",
			"names": []string{
				"Tom",
				"Jim",
			},
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
		}
	}

	t.Run("InvalidSelector", func(t *testing.T) {
		err := dasel.New(data()).DeleteMultiple("people.[a].name")
		expErr := fmt.Errorf("failed to parse selector: %w", &dasel.InvalidIndexErr{Index: "a"})
		if err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
	})
	t.Run("ExistingSingleString", deleteMultipleTest(dasel.New(data()), "id", map[string]interface{}{
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
	}))
	t.Run("ExistingStringValue", deleteMultipleTest(dasel.New(data()), "people.[0].name", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"id": 1,
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
	}))
	t.Run("WildcardProperty", deleteMultipleTest(dasel.New(data()), "people.[*].name", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"id": 1,
			},
			{
				"id": 2,
			},
		},
	}))
	t.Run("ExistingIntValue", deleteMultipleTest(dasel.New(data()), "people.(id=1).id", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{
			{
				"name": "Tom",
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
	}))
	t.Run("AllObjectsInArray", deleteMultipleTest(dasel.New(data()), "people.[*]", map[string]interface{}{
		"id": "123",
		"names": []string{
			"Tom",
			"Jim",
		},
		"people": []map[string]interface{}{},
	}))
	t.Run("ExistingArray", deleteMultipleTest(dasel.New(data()), "names", map[string]interface{}{
		"id": "123",
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
	}))
	t.Run("AllItemsInObject", deleteMultipleTest(dasel.New(data()), ".[*]", map[string]interface{}{}))
	t.Run("RootNodeObject", deleteMultipleTest(dasel.New(map[string]interface{}{
		"name": "Tom",
	}), ".", map[string]interface{}{}))
	t.Run("RootNodeArray", deleteMultipleTest(dasel.New([]interface{}{
		"name", "Tom",
	}), ".", []interface{}{}))
	t.Run("RootNodeUnknown", deleteMultipleTest(dasel.New(false), ".", map[string]interface{}{}))
}

// TestNode_NewFromFile tests parsing from file.
func TestNode_NewFromFile(t *testing.T) {
	noSuchFileName := "no_such_file"
	noSuchFileErrMsg := "could not open file: open " + noSuchFileName
	unmarshalFailMsg := "could not unmarshal data"

	tests := []struct {
		name   string
		file   string
		parser string
		err    error
	}{
		{
			name:   "Existing JSON file and JSON Read parser specified",
			file:   "./tests/assets/example.json",
			parser: "json",
			err:    nil,
		},
		{
			name:   "Existing JSON file and YAML Read parser specified",
			file:   "./tests/assets/example.json",
			parser: "yaml",
			// JSON is a subset of YAML v1.2
			// https://stackoverflow.com/questions/21584985/what-valid-json-files-are-not-valid-yaml-1-1-files
			err: nil,
		},
		{
			name:   "Existing YAML file and YAML Read parser specified",
			file:   "./tests/assets/example.yaml",
			parser: "yaml",
			err:    nil,
		},
		{
			name:   "Existing YAML file and XML Read parser specified",
			file:   "./tests/assets/example.yaml",
			parser: "xml",
			err:    errors.New(unmarshalFailMsg),
		},
		{
			name:   "Existing XML file and XML Read parser specified",
			file:   "./tests/assets/example.xml",
			parser: "xml",
			err:    nil,
		},
		{
			name:   "Existing XML file and JSON Read parser specified",
			file:   "./tests/assets/example.xml",
			parser: "json",
			err:    errors.New(unmarshalFailMsg),
		},
		{
			name:   "Existing JSON file and invalid Read parser specified",
			file:   "./tests/assets/example.json",
			parser: "bad",
			err:    storage.UnknownParserErr{Parser: "bad"},
		},
		{
			name:   "File doesn't exist and JSON Read parser specified",
			file:   noSuchFileName,
			parser: "json",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and YAML Read parser specified",
			file:   noSuchFileName,
			parser: "yaml",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and YML Read parser specified",
			file:   noSuchFileName,
			parser: "yml",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and TOML Read parser specified",
			file:   noSuchFileName,
			parser: "toml",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and XML Read parser specified",
			file:   noSuchFileName,
			parser: "xml",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and CSV Read parser specified",
			file:   noSuchFileName,
			parser: "csv",
			err:    errors.New(noSuchFileErrMsg),
		},
		{
			name:   "File doesn't exist and invalid ('bad') Read parser specified",
			file:   noSuchFileName,
			parser: "bad",
			err:    storage.UnknownParserErr{Parser: "bad"},
		},
		{
			name:   "File doesn't exist and invalid ('-') Read parser specified",
			file:   noSuchFileName,
			parser: "-",
			err:    storage.UnknownParserErr{Parser: "-"},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			_, err := dasel.NewFromFile(tc.file, tc.parser)

			if tc.err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.err, err)
				return
			}
			if tc.err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.err, err)
				return
			}
			if tc.err != nil && err != nil && !strings.Contains(err.Error(), tc.err.Error()) {
				t.Errorf("expected err %v, got %v", tc.err, err)
				return
			}
		})
	}
}

// TestNode_Write tests writing to Writer.
func TestNode_Write(t *testing.T) {
	type testData = map[string]interface{}
	type args struct {
		parser     string
		compact    bool
		escapeHTML bool
	}

	commonData := testData{">": "&"}
	xmlData := testData{
		"key": testData{
			"value": "<&>",
		},
	}
	xmlEscapedData := testData{
		"key": testData{
			"value": "&lt;&amp;&gt;",
		},
	}

	tests := []struct {
		name       string
		data       testData
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "JSON, compact, no escape",
			data: commonData,
			args: args{
				parser:     "json",
				compact:    true,
				escapeHTML: false,
			},
			wantWriter: "{\">\":\"&\"}\n",
			wantErr:    false,
		},
		{
			name: "JSON, pretty, escape",
			data: commonData,
			args: args{
				parser:     "json",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "{\n  \"\\u003e\": \"\\u0026\"\n}\n",
			wantErr:    false,
		},

		{
			name: "YAML, compact, no escape",
			data: commonData,
			args: args{
				parser:     "yaml",
				compact:    true,
				escapeHTML: false,
			},
			wantWriter: "'>': '&'\n",
			wantErr:    false,
		},
		{
			name: "YAML, pretty, escape",
			data: commonData,
			args: args{
				parser:     "yaml",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "'>': '&'\n",
			wantErr:    false,
		},

		{
			name: "YML, compact, no escape",
			data: commonData,
			args: args{
				parser:     "yml",
				compact:    true,
				escapeHTML: false,
			},
			wantWriter: "'>': '&'\n",
			wantErr:    false,
		},
		{
			name: "YML, pretty, escape",
			data: commonData,
			args: args{
				parser:     "yml",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "'>': '&'\n",
			wantErr:    false,
		},

		{
			name: "TOML, compact, no escape",
			data: commonData,
			args: args{
				parser:     "toml",
				compact:    true,
				escapeHTML: false,
			},
			wantWriter: "\">\" = \"&\"\n",
			wantErr:    false,
		},
		{
			name: "TOML, pretty, escape",
			data: commonData,
			args: args{
				parser:     "toml",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "\">\" = \"&\"\n",
			wantErr:    false,
		},

		{
			name: "PLAIN, pretty, escape",
			data: commonData,
			args: args{
				parser:     "plain",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "map[>:&]\n",
			wantErr:    false,
		},
		{
			name: "PLAIN, pretty, no escape",
			data: commonData,
			args: args{
				parser:     "plain",
				compact:    false,
				escapeHTML: false,
			},
			wantWriter: "map[>:&]\n",
			wantErr:    false,
		},

		{
			name: "-, pretty, escape",
			data: commonData,
			args: args{
				parser:     "-",
				compact:    false,
				escapeHTML: true,
			},
			wantWriter: "map[>:&]\n",
			wantErr:    false,
		},
		{
			name: "-, pretty, no escape",
			data: commonData,
			args: args{
				parser:     "-",
				compact:    false,
				escapeHTML: false,
			},
			wantWriter: "map[>:&]\n",
			wantErr:    false,
		},

		{
			name: "XML, pretty, no escape",
			data: xmlData,
			args: args{
				parser:     "xml",
				compact:    false,
				escapeHTML: false,
			},
			// Invalid XML, but since we were asked not to escape,
			// this is probably the best we should do
			wantWriter: "<key>\n  <value><&></value>\n</key>\n",
			wantErr:    false,
		},
		{
			name: "Pre-escaped XML, pretty, no escape",
			data: xmlEscapedData,
			args: args{
				parser:     "xml",
				compact:    false,
				escapeHTML: false,
			},
			wantWriter: "<key>\n  <value>&lt;&amp;&gt;</value>\n</key>\n",
			wantErr:    false,
		},

		{
			name: "Invalid Write parser: bad",
			data: commonData,
			args: args{
				parser:     "bad",
				compact:    false,
				escapeHTML: false,
			},
			wantWriter: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			node := dasel.New(tt.data)
			if err := node.Write(writer, tt.args.parser, tt.args.compact, tt.args.escapeHTML); (err != nil) != tt.wantErr {
				t.Errorf("Node.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Node.Write() = %v want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
