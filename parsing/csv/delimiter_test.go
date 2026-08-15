package csv_test

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/csv"
	"testing"
)

func TestCsvReader_MultiByteDelimiter(t *testing.T) {
	opts := parsing.DefaultReaderOptions()
	opts.Ext["csv-delimiter"] = "§"

	r, err := csv.CSV.NewReader(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := r.Read([]byte("name§age\nAlice§30\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rows, err := got.SliceLen()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rows != 1 {
		t.Fatalf("expected 1 row, got %d", rows)
	}

	row, err := got.GetSliceIndex(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, err := row.MapKeys()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 || keys[0] != "name" || keys[1] != "age" {
		t.Fatalf("expected headers [name age], got %v", keys)
	}
}

func TestCsvWriter_MultiByteDelimiter(t *testing.T) {
	opts := parsing.DefaultWriterOptions()
	opts.Ext["csv-delimiter"] = "§"

	w, err := csv.CSV.NewWriter(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rows := model.NewSliceValue()
	row := model.NewMapValue()
	if err := row.SetMapKey("name", model.NewStringValue("Alice")); err != nil {
		t.Fatal(err)
	}
	if err := row.SetMapKey("age", model.NewStringValue("30")); err != nil {
		t.Fatal(err)
	}
	if err := rows.Append(row); err != nil {
		t.Fatal(err)
	}

	got, err := w.Write(rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exp := "name§age\nAlice§30\n"
	if string(got) != exp {
		t.Errorf("expected %q, got %q", exp, string(got))
	}
}

func TestCsvDelimiter_NotASingleCharacter(t *testing.T) {
	ropts := parsing.DefaultReaderOptions()
	ropts.Ext["csv-delimiter"] = ";;"
	if _, err := csv.CSV.NewReader(ropts); err == nil {
		t.Error("expected an error from the reader, got none")
	}

	wopts := parsing.DefaultWriterOptions()
	wopts.Ext["csv-delimiter"] = ";;"
	if _, err := csv.CSV.NewWriter(wopts); err == nil {
		t.Error("expected an error from the writer, got none")
	}
}
