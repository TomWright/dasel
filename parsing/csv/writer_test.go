package csv_test

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/csv"
	"testing"
)

func TestCsvWriter_Write(t *testing.T) {
	expBytes := []byte(`name,age,city
Alice,30,New York
Bob,25,Los Angeles
Charlie,35,Chicago
`)

	r, err := csv.CSV.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rows := model.NewSliceValue()
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
	if err := rows.Append(row1); err != nil {
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
	if err := rows.Append(row2); err != nil {
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
	if err := rows.Append(row3); err != nil {
		t.Fatal(err)
	}

	got, err := r.Write(rows)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if string(expBytes) != string(got) {
		t.Errorf("expected:\n%s\ngot:\n%s", string(expBytes), string(got))
	}
}

func TestCsvWriter_RowsWithDifferentKeys(t *testing.T) {
	row := func(kvs ...string) *model.Value {
		v := model.NewMapValue()
		for i := 0; i < len(kvs); i += 2 {
			if err := v.SetMapKey(kvs[i], model.NewStringValue(kvs[i+1])); err != nil {
				t.Fatal(err)
			}
		}
		return v
	}

	rows := model.NewSliceValue()
	for _, r := range []*model.Value{row("a", "1"), row("a", "2", "b", "3"), row("b", "4")} {
		if err := rows.Append(r); err != nil {
			t.Fatal(err)
		}
	}

	w, err := csv.CSV.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatal(err)
	}
	got, err := w.Write(rows)
	if err != nil {
		t.Fatal(err)
	}

	exp := "a,b\n1,\n2,3\n,4\n"
	if string(got) != exp {
		t.Errorf("expected %q, got %q", exp, got)
	}
}
