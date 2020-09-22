package dasel_test

import (
	"errors"
	"github.com/tomwright/dasel"
	"reflect"
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

type processFn func(t *testing.T, n *dasel.Node) (*dasel.Node, bool)
type checkFn func(t *testing.T, n *dasel.Node) bool

type parseTest struct {
	Selector   string
	Input      interface{}
	ProcessFns []processFn
	CheckFns   []checkFn
	Names      []string
	Err        error
}

func (pt parseTest) Add(name string, p processFn, c checkFn) parseTest {
	if pt.ProcessFns == nil {
		pt.ProcessFns = make([]processFn, 0)
	}
	if pt.CheckFns == nil {
		pt.CheckFns = make([]checkFn, 0)
	}
	if pt.Names == nil {
		pt.Names = make([]string, 0)
	}
	pt.Names = append(pt.Names, name)
	pt.ProcessFns = append(pt.ProcessFns, p)
	pt.CheckFns = append(pt.CheckFns, c)
	return pt
}

func TestNode_Query(t *testing.T) {
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

	t.Run("Traversal", func(t *testing.T) {
		tests := parseTest{
			Selector: ".a.b.c.thing",
			Input:    mapRoot,
		}.Add(
			"StartsAtEndElement",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".thing", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual("1", n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", "1", n.Value)
					return false
				}
				return true
			},
		).Add(
			"Previous1",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Previous, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".c", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual(mapC, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapC, n.Value)
					return false
				}
				return true
			},
		).Add(
			"Previous2",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Previous, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".b", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual(mapB, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapB, n.Value)
					return false
				}
				return true
			},
		).Add(
			"Previous3",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Previous, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".a", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual(mapA, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapA, n.Value)
					return false
				}
				return true
			},
		).Add(
			"Next",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Next, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".b", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}

				if !reflect.DeepEqual(mapB, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapB, n.Value)
					return false
				}
				return true
			},
		).Add(
			"Previous4",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Previous, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".a", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual(mapA, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapA, n.Value)
					return false
				}
				return true
			},
		).Add(
			"Previous5",
			func(t *testing.T, n *dasel.Node) (*dasel.Node, bool) { return n.Previous, true },
			func(t *testing.T, n *dasel.Node) bool {
				if exp, got := ".", n.Selector.Current; exp != got {
					t.Errorf("expected selector of `%s`, got `%s`", exp, got)
					return false
				}
				if !reflect.DeepEqual(mapRoot, n.Value) {
					t.Errorf("expected value of\n%s\ngot\n%s", mapRoot, n.Value)
					return false
				}
				return true
			},
		)

		rootNode := dasel.New(tests.Input)
		got, err := rootNode.Query(tests.Selector)
		switch {
		case err == nil && tests.Err != nil:
			t.Errorf("expected err `%v`, got `%v`", tests.Err, err)
			return
		case err != nil && tests.Err == nil:
			t.Errorf("unexpected err `%v`", err)
			return
		case err != nil && tests.Err != nil && !errors.Is(err, tests.Err):
			t.Errorf("expected err `%v`, got `%v`", tests.Err, err)
			return
		}

		ok := true

		for i := 0; i < len(tests.ProcessFns); i++ {
			if !ok {
				break
			}
			t.Run(tests.Names[i], func(t *testing.T) {
				got, ok = tests.ProcessFns[i](t, got)
				if !ok {
					return
				}
				ok = tests.CheckFns[i](t, got)
				if !ok {
					return
				}
			})
		}
	})
}
