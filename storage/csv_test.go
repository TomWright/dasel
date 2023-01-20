package storage_test

import (
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
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
	exp := dasel.ValueOf(csvMap).WithMetadata("csvHeaders", []string{"id", "name"})
	if !reflect.DeepEqual(exp.Interface(), got.Interface()) {
		t.Errorf("expected %v, got %v", exp, got)
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
	t.Run("SingleDocument", func(t *testing.T) {
		value := dasel.ValueOf(map[string]interface{}{
			"id":   "1",
			"name": "Tom",
		}).
			WithMetadata("isSingleDocument", true)
		got, err := (&storage.CSVParser{}).ToBytes(value)
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
		value := dasel.ValueOf([]interface{}{
			map[string]interface{}{
				"id":   "1",
				"name": "Tom",
			},
			map[string]interface{}{
				"id":   "2",
				"name": "Tommy",
			},
		}).
			WithMetadata("isSingleDocument", true)
		got, err := (&storage.CSVParser{}).ToBytes(value)
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
		value := dasel.ValueOf([]interface{}{
			map[string]interface{}{
				"id":   "1",
				"name": "Tom",
			},
			map[string]interface{}{
				"id":   "2",
				"name": "Jim",
			},
		}).
			WithMetadata("isMultiDocument", true)
		got, err := (&storage.CSVParser{}).ToBytes(value)
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
}

func deepEqualOneOf(t *testing.T, got []byte, exps ...[]byte) {
	for _, exp := range exps {
		if reflect.DeepEqual(exp, got) {
			return
		}
	}
	t.Errorf("%s did not match any of the expected values", string(got))
}
