package json_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
)

func TestJson(t *testing.T) {
	doc := []byte(`{
    "string": "foo",
    "int": 1,
    "float": 1.1,
    "bool": true,
    "null": null,
    "array": [
        1,
        2,
        3
    ],
    "object": {
        "key": "value"
    }
}
`)
	reader, err := json.JSON.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}
	writer, err := json.JSON.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatal(err)
	}

	value, err := reader.Read(doc)
	if err != nil {
		t.Fatal(err)
	}

	newDoc, err := writer.Write(value)
	if err != nil {
		t.Fatal(err)
	}

	if string(doc) != string(newDoc) {
		t.Fatalf("expected %s, got %s...\n%s", string(doc), string(newDoc), cmp.Diff(string(doc), string(newDoc)))
	}
}
