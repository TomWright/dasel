package xml_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/xml"
)

func TestXmlWriter_Compact(t *testing.T) {
	r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	opts := parsing.DefaultWriterOptions()
	opts.Compact = true
	w, err := xml.XML.NewWriter(opts)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	data, err := r.Read([]byte(`<Document>
  <Sender>Ivanov</Sender>
  <Content>Hello</Content>
</Document>
`))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	xmlBytes, err := w.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expected := "<Document><Sender>Ivanov</Sender><Content>Hello</Content></Document>\n"
	if string(xmlBytes) != expected {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(xmlBytes))
	}
}

func TestXmlWriter_CompactWithPI(t *testing.T) {
	r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	opts := parsing.DefaultWriterOptions()
	opts.Compact = true
	w, err := xml.XML.NewWriter(opts)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	data, err := r.Read([]byte(`<?xml version="1.0" encoding="utf-8"?>
<Root>
  <Name>Test</Name>
</Root>
`))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	xmlBytes, err := w.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expected := `<?xml version="1.0" encoding="utf-8"?><Root><Name>Test</Name></Root>` + "\n"
	if string(xmlBytes) != expected {
		t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(xmlBytes))
	}
}

func TestXmlWriter_CompactWithComments(t *testing.T) {
	r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	opts := parsing.DefaultWriterOptions()
	opts.Compact = true
	w, err := xml.XML.NewWriter(opts)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	data, err := r.Read([]byte(`<!-- top comment -->
<Root>
  <!-- child comment -->
  <Name>Test</Name>
</Root>
`))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	xmlBytes, err := w.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	got := string(xmlBytes)
	// Compact output should contain all elements and comments
	for _, substr := range []string{"<!-- top comment -->", "<!-- child comment -->", "<Root>", "<Name>Test</Name>", "</Root>"} {
		if !strings.Contains(got, substr) {
			t.Fatalf("Expected output to contain %q, got:\n%s", substr, got)
		}
	}
	// Should not contain indentation whitespace (check content before trailing newline)
	trimmed := strings.TrimSuffix(got, "\n")
	if strings.Contains(trimmed, "\n  ") {
		t.Fatalf("Compact output should not contain indentation, got:\n%s", got)
	}
}

func TestXmlReader_Write(t *testing.T) {
	t.Run("nested xml elements", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		w, err := xml.XML.NewWriter(parsing.DefaultWriterOptions())
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

		xmlBytes, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		expected := `<Document>
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
`
		if string(xmlBytes) != expected {
			t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(xmlBytes))
		}
	})

	t.Run("nested xml elements with processing instruction", func(t *testing.T) {
		r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		w, err := xml.XML.NewWriter(parsing.DefaultWriterOptions())
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

		doc, err := data.GetMapKey("Document")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		docProcessingInstructions, ok := doc.MetadataValue("xml_processing_instructions")
		if !ok || docProcessingInstructions == nil {
			t.Fatalf("Expected processing instructions on Document element")
		}

		jsonBytes, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		expected := `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
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
`
		if string(jsonBytes) != expected {
			t.Fatalf("Expected:\n%s\nGot:\n%s", expected, string(jsonBytes))
		}
	})

	t.Run("encode attributes", func(t *testing.T) {
		w, err := xml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		toEncode := model.NewMapValue()
		foo := model.NewMapValue()
		_ = foo.SetMapKey("-fiz", model.NewStringValue("hello"))
		_ = foo.SetMapKey("bar", model.NewStringValue(""))
		_ = toEncode.SetMapKey("foo", foo)

		got, err := w.Write(toEncode)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		exp := []byte(`<foo fiz="hello">
  <bar></bar>
</foo>
`)
		if string(got) != string(exp) {
			t.Errorf("Expected:\n%s\nGot:\n%s", string(exp), string(got))
		}
	})

	t.Run("encode cdata", func(t *testing.T) {
		w, err := xml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		toEncode := model.NewMapValue()
		_ = toEncode.SetMapKey("foo", model.NewStringValue("<bar>baz</bar>"))
		got, err := w.Write(toEncode)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		// TODO : Change this to use CDATA sections rather than escaping.
		exp := []byte(`<foo>&lt;bar&gt;baz&lt;/bar&gt;</foo>
`)
		if string(got) != string(exp) {
			t.Errorf("Expected:\n%s\nGot:\n%s", string(exp), string(got))
		}
	})
}
