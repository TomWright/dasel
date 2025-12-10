package toml_test

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/toml"
	"testing"
)

func TestTomlReader_Read(t *testing.T) {
	dataBytes := []byte(`title = "TOML Example"
[owner]
name = "Tom Preston-Werner"
`)
	//dob = 1979-05-27T07:32:00-08:00
	dataModel := model.NewMapValue()
	//parsedTime, err := time.Parse(time.RFC3339, "1979-05-27T07:32:00-08:00")
	//if err != nil {
	//	t.Fatalf("unexpected error: %v", err)
	//}
	ownerMap := model.NewMapValue()
	_ = ownerMap.SetMapKey("name", model.NewStringValue("Tom Preston-Werner"))
	//_ = ownerMap.SetMapKey("dob", model.NewValue(parsedTime))
	_ = dataModel.SetMapKey("title", model.NewStringValue("TOML Example"))
	_ = dataModel.SetMapKey("owner", ownerMap)

	r, err := toml.TOML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := r.Read(dataBytes)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	matchResult, err := got.Equal(dataModel)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	matchResultBool, err := matchResult.BoolValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !matchResultBool {
		t.Errorf("expected\n%s\ngot\n%s", dataModel.String(), got.String())
	}
}
