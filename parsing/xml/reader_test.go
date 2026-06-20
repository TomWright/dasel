package xml_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/json"
	"github.com/tomwright/dasel/v3/parsing/xml"
)

func TestXmlReader_Read(t *testing.T) {
	t.Run("nested xml elements", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		w, err := json.JSON.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<Document>
  <Sender>Ivanov</Sender>
  <In_N_Document>
    <N_Document>1024</N_Document>
    <Date_Reg>2024-06-21T15:07:29.0451517+03:00</Date_Reg>
  </In_N_Document>
  <Out_N_Document>
    <N_Document>2043</N_Document>
    <Date_Reg>2024-05-01T00:00:00</Date_Reg>
  </Out_N_Document>
  <Content>Skzzkz</Content>
  <DSP>true</DSP>
</Document>
`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		jsonBytes, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		expected := `{
    "Document": {
        "Sender": "Ivanov",
        "In_N_Document": {
            "N_Document": "1024",
            "Date_Reg": "2024-06-21T15:07:29.0451517+03:00"
        },
        "Out_N_Document": {
            "N_Document": "2043",
            "Date_Reg": "2024-05-01T00:00:00"
        },
        "Content": "Skzzkz",
        "DSP": "true"
    }
}
`
		if string(jsonBytes) != expected {
			t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(jsonBytes))
		}
	})

	t.Run("nested xml elements with processing instruction", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		w, err := json.JSON.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<Document>
  <Sender>Ivanov</Sender>
  <In_N_Document>
    <N_Document>1024</N_Document>
    <Date_Reg>2024-06-21T15:07:29.0451517+03:00</Date_Reg>
  </In_N_Document>
  <Out_N_Document>
    <N_Document>2043</N_Document>
    <Date_Reg>2024-05-01T00:00:00</Date_Reg>
  </Out_N_Document>
  <Content>Skzzkz</Content>
  <DSP>true</DSP>
</Document>
`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		jsonBytes, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		expected := `{
    "Document": {
        "Sender": "Ivanov",
        "In_N_Document": {
            "N_Document": "1024",
            "Date_Reg": "2024-06-21T15:07:29.0451517+03:00"
        },
        "Out_N_Document": {
            "N_Document": "2043",
            "Date_Reg": "2024-05-01T00:00:00"
        },
        "Content": "Skzzkz",
        "DSP": "true"
    }
}
`
		if string(jsonBytes) != expected {
			t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(jsonBytes))
		}
	})

	t.Run("cdata tag", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<foo>
	<![CDATA[<bar>baz</bar>]]>
</foo>
`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		x, err := data.GetMapKey("foo")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		got, err := x.StringValue()
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		exp := "<bar>baz</bar>"
		if exp != got {
			t.Fatalf("Expected value %q but got %q", exp, got)
		}
	})

	t.Run("empty cdata tag", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<foo>
	<![CDATA[]]>
</foo>
`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		x, err := data.GetMapKey("foo")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		got, err := x.StringValue()
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		exp := ""
		if exp != got {
			t.Fatalf("Expected value %q but got %q", exp, got)
		}
	})
}

func TestXMLReader_DepthLimit(t *testing.T) {
	reader, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("deeply nested elements exceeding limit return error", func(t *testing.T) {
		// 10_001 levels — one past the 10_000 limit. Each open tag is 3 bytes (<a>),
		// so total ~60 KB, well under the 10 MB size cap.
		const depth = 10_001
		open := strings.Repeat("<a>", depth)
		close := strings.Repeat("</a>", depth)
		_, err := reader.Read([]byte(open + close))
		if !errors.Is(err, xml.ErrXMLMaxDepthExceeded) {
			t.Fatalf("expected ErrXMLMaxDepthExceeded, got %v", err)
		}
	})

	t.Run("normal nesting within limit succeeds", func(t *testing.T) {
		data := []byte(`<root><a><b><c>hello</c></b></a></root>`)
		_, err := reader.Read(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
