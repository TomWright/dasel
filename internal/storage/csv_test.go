package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var csvBytes = []byte(`id,name
1,Tom
2,Jim
`)
var csvMap = []map[string]interface{}{
	{
		"id":   "1",
		"name": "Tom",
	},
	{
		"id":   "2",
		"name": "Jim",
	},
}

func TestCSVParser_FromBytes(t *testing.T) {
	got, err := (&storage.CSVParser{}).FromBytes(csvBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(&storage.CSVDocument{
		Value:   csvMap,
		Headers: []string{"id", "name"},
	}, got) {
		t.Errorf("expected %v, got %v", csvMap, got)
	}
}

func TestCSVParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.CSVParser{}).FromBytes(nil)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.CSVParser{}).FromBytes([]byte(`a,b
a,b,c`))
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.CSVParser{}).FromBytes([]byte(`a,b,c
a,b`))
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestCSVParser_ToBytes(t *testing.T) {
	t.Run("CSVDocument", func(t *testing.T) {
		got, err := (&storage.CSVParser{}).ToBytes(&storage.CSVDocument{
			Value:   csvMap,
			Headers: []string{"id", "name"},
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(csvBytes, got) {
			t.Errorf("expected %v, got %v", string(csvBytes), string(got))
		}
	})
	t.Run("SingleDocument", func(t *testing.T) {
		got, err := (&storage.CSVParser{}).ToBytes(&storage.BasicSingleDocument{
			Value: map[string]interface{}{
				"id":   "1",
				"name": "Tom",
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		deepEqualOneOf(t, got, []byte(`id,name
1,Tom
`), []byte(`name,id
Tom,1
`))
	})
	t.Run("SingleDocumentSlice", func(t *testing.T) {
		got, err := (&storage.CSVParser{}).ToBytes(&storage.BasicSingleDocument{
			Value: []interface{}{
				map[string]interface{}{
					"id":   "1",
					"name": "Tom",
				},
				map[string]interface{}{
					"id":   "2",
					"name": "Tommy",
				},
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		deepEqualOneOf(t, got, []byte(`id,name
1,Tom
2,Tommy
`), []byte(`name,id
Tom,1
`))
	})
	t.Run("MultiDocument", func(t *testing.T) {
		got, err := (&storage.CSVParser{}).ToBytes(&storage.BasicMultiDocument{
			Values: []interface{}{
				map[string]interface{}{
					"id":   "1",
					"name": "Tom",
				},
				map[string]interface{}{
					"id":   "2",
					"name": "Jim",
				},
			},
		})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		deepEqualOneOf(t, got, []byte(`id,name
1,Tom
id,name
2,Jim
`), []byte(`name,id
Tom,1
id,name
2,Jim
`), []byte(`id,name
1,Tom
name,id
Jim,2
`), []byte(`name,id
Tom,1
name,id
Jim,2
`))
	})
	t.Run("DefaultDocType", func(t *testing.T) {
		got, err := (&storage.CSVParser{}).ToBytes([]interface{}{"x", "y"})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		deepEqualOneOf(t, got, []byte(`[x y]
`))
	})
}

func deepEqualOneOf(t *testing.T, got []byte, exps ...[]byte) {
	for _, exp := range exps {
		if reflect.DeepEqual(exp, got) {
			return
		}
	}
	t.Errorf("%s did not match any of the expected values", string(got))
}

func TestCSVDocument_Documents(t *testing.T) {
	in := &storage.CSVDocument{
		Value: []map[string]interface{}{
			{
				"id":   1,
				"name": "Tom",
			},
			{
				"id":   2,
				"name": "Jim",
			},
		},
		Headers: []string{"id", "name"},
	}
	exp := []interface{}{
		map[string]interface{}{
			"id":   1,
			"name": "Tom",
		},
		map[string]interface{}{
			"id":   2,
			"name": "Jim",
		},
	}
	got := in.Documents()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestCSVDocument_RealValue(t *testing.T) {
	exp := []map[string]interface{}{
		{
			"id":   1,
			"name": "Tom",
		},
		{
			"id":   2,
			"name": "Jim",
		},
	}
	in := &storage.CSVDocument{
		Value:   exp,
		Headers: []string{"id", "name"},
	}
	got := in.RealValue()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
