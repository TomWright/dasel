package ini_test

import (
	"fmt"
	"github.com/tomwright/dasel/v3/parsing/ini"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
)

func TestIni(t *testing.T) {
	doc := []byte(`app_mode = development

[paths]
data = /home/git/grafana

[server]
protocol = http
http_port = 9999
enforce_domain = true
`)
	reader, err := ini.INI.NewReader(parsing.DefaultReaderOptions())
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

	fmt.Println(string(newDoc))

	//if string(doc) != string(newDoc) {
	//	t.Fatalf("expected %s, got %s...\n%s", string(doc), string(newDoc), cmp.Diff(string(doc), string(newDoc)))
	//}
}
