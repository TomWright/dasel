package toml

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Reader = (*tomlReader)(nil)

func newTOMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &tomlReader{}, nil
}

type tomlReader struct{}

//var _ unstable.Unmarshaler = (*tomlValue)(nil)

// Read reads a value from a byte slice.
func (j *tomlReader) Read(data []byte) (*model.Value, error) {
	var unmarshalled any
	if err := toml.Unmarshal(data, &unmarshalled); err != nil {
		return nil, err
	}
	return model.NewValue(&unmarshalled), nil
}

// Read reads a value from a byte slice.
//func (j *tomlReader) Read(data []byte) (*model.Value, error) {
//	decoder := toml.NewDecoder(bytes.NewReader(data)).
//		EnableUnmarshalerInterface()
//	var unmarshalled tomlValue
//	if err := decoder.Decode(&unmarshalled); err != nil {
//		return nil, err
//	}
//	return unmarshalled.value, nil
//}

//type tomlValue struct {
//	value *model.Value
//	node  *unstable.Node
//}

//func (t *tomlValue) UnmarshalTOML(value *unstable.Node) error {
//	t.node = value
//	fmt.Println(value.Kind.String())
//	fmt.Println(string(value.Data))
//	fmt.Println(string(value.Next().Data))
//	fmt.Println(string(value.Next().Kind.String()))
//	fmt.Println(string(value.Next().Data))
//	fmt.Println(string(value.Next().Kind.String()))
//	fmt.Println(value.Child())
//	//panic("asd")
//	err := t.parseNodeValue()
//	if err != nil {
//		return err
//	}
//	fmt.Println("here")
//	return nil
//}

//func (t *tomlValue) parseNodeValue() error {
//	switch t.node.Kind {
//	case unstable.String:
//		t.value = model.NewStringValue(string(t.node.Data))
//	case unstable.Integer:
//		t.value = model.NewIntValue(0)
//	case unstable.Float:
//		t.value = model.NewFloatValue(0)
//	case unstable.Bool:
//		t.value = model.NewBoolValue(false)
//	case unstable.DateTime:
//		t.value = model.NewStringValue("")
//	case unstable.Array:
//		slice := model.NewSliceValue()
//		//for _, item := range value.Elements {
//		//	v, err := t.readNode(item)
//		//	if err != nil {
//		//		return nil, err
//		//	}
//		//	slice.Append(v)
//		//}
//		t.value = slice
//	case unstable.Table:
//		m := model.NewMapValue()
//		//for k, v := range value.Fields {
//		//	mv, err := t.readNode(v)
//		//	if err != nil {
//		//		return nil, err
//		//	}
//		//	m.Set(k, mv)
//		//}
//		t.value = m
//	case unstable.InlineTable:
//		m := model.NewMapValue()
//		//for k, v := range value.Fields {
//		//	mv, err := t.readNode(v)
//		//	if err != nil {
//		//		return nil, err
//		//	}
//		//	m.Set(k, mv)
//		//}
//		t.value = m
//	default:
//		t.value = model.NewNullValue()
//		return fmt.Errorf("unhandled TOML node kind: %v", t.node.Kind)
//	}
//	return nil
//}
