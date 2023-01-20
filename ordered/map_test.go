package ordered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestMap_UnmarshalJSON(t *testing.T) {

	t.Run("UnmarshalMap", func(t *testing.T) {
		dec := json.NewDecoder(bytes.NewReader([]byte(`{"name": "Tom", "age": 29, "things": ["a", 1, false], "settings": {"a": 1, "b": false, "c": 123.456}}`)))
		got, err := UnmarshalJSON(dec)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		m, ok := got.(*Map)
		if !ok {
			t.Errorf("unexpected type: %T", m)
			return
		}

		keyValues := m.KeyValues()

		for _, kv := range keyValues {
			fmt.Printf("%v: %v\n", kv.Key, kv.Value)
		}

	})
}
