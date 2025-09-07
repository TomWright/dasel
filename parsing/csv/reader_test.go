package csv_test

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/csv"
	"testing"
)

func TestCsvReader_Read(t *testing.T) {
	inputBytes := []byte(`name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago`)

	r, err := csv.CSV.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := r.Read(inputBytes)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	exp := model.NewSliceValue()
	row1 := model.NewMapValue()
	if err := row1.SetMapKey("name", model.NewStringValue("Alice")); err != nil {
		t.Fatal(err)
	}
	if err := row1.SetMapKey("age", model.NewStringValue("30")); err != nil {
		t.Fatal(err)
	}
	if err := row1.SetMapKey("city", model.NewStringValue("New York")); err != nil {
		t.Fatal(err)
	}
	if err := exp.Append(row1); err != nil {
		t.Fatal(err)
	}
	row2 := model.NewMapValue()
	if err := row2.SetMapKey("name", model.NewStringValue("Bob")); err != nil {
		t.Fatal(err)
	}
	if err := row2.SetMapKey("age", model.NewStringValue("25")); err != nil {
		t.Fatal(err)
	}
	if err := row2.SetMapKey("city", model.NewStringValue("Los Angeles")); err != nil {
		t.Fatal(err)
	}
	if err := exp.Append(row2); err != nil {
		t.Fatal(err)
	}
	row3 := model.NewMapValue()
	if err := row3.SetMapKey("name", model.NewStringValue("Charlie")); err != nil {
		t.Fatal(err)
	}
	if err := row3.SetMapKey("age", model.NewStringValue("35")); err != nil {
		t.Fatal(err)
	}
	if err := row3.SetMapKey("city", model.NewStringValue("Chicago")); err != nil {
		t.Fatal(err)
	}
	if err := exp.Append(row3); err != nil {
		t.Fatal(err)
	}
	matchRes, err := got.Equal(exp)
	if err != nil {
		t.Fatalf("error comparing values: %v", err)
	}
	match, err := matchRes.BoolValue()
	if err != nil {
		t.Fatalf("error getting bool value: %v", err)
	}
	if !match {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
