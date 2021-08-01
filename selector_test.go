package dasel_test

import (
	"github.com/tomwright/dasel"
	"reflect"
	"testing"
)

func testExtractNextSelector(in string, exp string, expRead int) func(t *testing.T) {
	return func(t *testing.T) {
		got, read := dasel.ExtractNextSelector(in)
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if read != expRead {
			t.Errorf("expected read of %d, got %d", expRead, read)
		}
	}
}

func testDynamicSelectorToGroups(in string, exp []string) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := dasel.DynamicSelectorToGroups(in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func testFindDynamicSelectorParts(in string, exp dasel.DynamicSelectorParts) func(t *testing.T) {
	return func(t *testing.T) {
		got := dasel.FindDynamicSelectorParts(in)
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func TestExtractNextSelector(t *testing.T) {
	t.Run("Simple", testExtractNextSelector(`.metadata.name`, `.metadata`, 9))
	t.Run("EscapedDot", testExtractNextSelector(`.before\.after.name`, `.before.after`, 14))
	t.Run("EscapedSpace", testExtractNextSelector(`.before\ after.name`, `.before after`, 14))
	t.Run("DynamicWithPath", testExtractNextSelector(`.(.before.a=b).after.name`, `.(.before.a=b)`, 14))
	t.Run("EscapedFirstDot", testExtractNextSelector(`\.name`, `.name`, 6))
	t.Run("SimpleProp", testExtractNextSelector(`.name`, `.name`, 5))
	t.Run("SimpleIndex", testExtractNextSelector(`.[123]`, `.[123]`, 6))
	t.Run("SimpleLength", testExtractNextSelector(`.[#]`, `.[#]`, 4))
}

func TestDynamicSelectorToGroups(t *testing.T) {
	t.Run("Blank", testDynamicSelectorToGroups("", []string{}))
	t.Run("Single", testDynamicSelectorToGroups("(a=1)", []string{
		"a=1",
	}))
	t.Run("Double", testDynamicSelectorToGroups("(a=1)(b=2)", []string{
		"a=1",
		"b=2",
	}))
	t.Run("Many", testDynamicSelectorToGroups("(a=1)(b=2)(c=3)(d=4)", []string{
		"a=1",
		"b=2",
		"c=3",
		"d=4",
	}))
	t.Run("Nested", testDynamicSelectorToGroups("(a=1.(x=3)(y=4))(b=2)", []string{
		"a=1.(x=3)(y=4)",
		"b=2",
	}))
	t.Run("Dot", testDynamicSelectorToGroups("(a=.)(b=2)", []string{
		"a=.",
		"b=2",
	}))
}

func TestFindDynamicSelectorParts(t *testing.T) {
	t.Run("Blank", testFindDynamicSelectorParts("", dasel.DynamicSelectorParts{
		Key:        "",
		Comparison: "",
		Value:      "",
	}))
	t.Run("Equal", testFindDynamicSelectorParts("a=b", dasel.DynamicSelectorParts{
		Key:        "a",
		Comparison: "=",
		Value:      "b",
	}))
	t.Run("NotEqual", testFindDynamicSelectorParts("a!=b", dasel.DynamicSelectorParts{
		Key:        "a",
		Comparison: "!=",
		Value:      "b",
	}))
	t.Run("MoreThanEqual", testFindDynamicSelectorParts("a>=b", dasel.DynamicSelectorParts{
		Key:        "a",
		Comparison: ">=",
		Value:      "b",
	}))
	t.Run("LessThanEqual", testFindDynamicSelectorParts("a<=b", dasel.DynamicSelectorParts{
		Key:        "a",
		Comparison: "<=",
		Value:      "b",
	}))
	t.Run("NestedGroupIgnored", testFindDynamicSelectorParts("(.(x=y)).x=1", dasel.DynamicSelectorParts{
		Key:        "(.(x=y)).x",
		Comparison: "=",
		Value:      "1",
	}))
	t.Run("NestedGroupIgnored", testFindDynamicSelectorParts("a=(.(x=y))", dasel.DynamicSelectorParts{
		Key:        "a",
		Comparison: "=",
		Value:      "(.(x=y))",
	}))
	t.Run("SelectorOnBothSides", testFindDynamicSelectorParts("(.a.b.c)=(.x.y.z)", dasel.DynamicSelectorParts{
		Key:        "(.a.b.c)",
		Comparison: "=",
		Value:      "(.x.y.z)",
	}))
	t.Run("NestedWithSelectorOnBothSides", testFindDynamicSelectorParts("(.a.b.(1=2).c)=(.x.(3=4).y.z)", dasel.DynamicSelectorParts{
		Key:        "(.a.b.(1=2).c)",
		Comparison: "=",
		Value:      "(.x.(3=4).y.z)",
	}))
}
