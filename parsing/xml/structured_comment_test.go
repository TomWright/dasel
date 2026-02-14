package xml_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	daselxml "github.com/tomwright/dasel/v3/parsing/xml"
)

// newTestReaderWriter creates a reader and writer for round-trip testing.
func newTestReaderWriter(t *testing.T) (parsing.Reader, parsing.Writer) {
	t.Helper()
	r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("Unexpected error creating reader: %s", err)
	}
	w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("Unexpected error creating writer: %s", err)
	}
	return r, w
}

// assertRoundTrip performs a read-write round-trip and verifies expected strings are present.
func assertRoundTrip(t *testing.T, r parsing.Reader, w parsing.Writer, input string, expected []string) {
	t.Helper()
	data, err := r.Read([]byte(input))
	if err != nil {
		t.Fatalf("Unexpected error reading XML: %s", err)
	}

	output, err := w.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error writing XML: %s", err)
	}

	outputStr := string(output)
	for _, exp := range expected {
		if !strings.Contains(outputStr, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, outputStr)
		}
	}
}

// assertOrderedContains verifies that the given strings appear in the specified order in output.
func assertOrderedContains(t *testing.T, output string, ordered []string) {
	t.Helper()
	for i := 1; i < len(ordered); i++ {
		prevIdx := strings.Index(output, ordered[i-1])
		currIdx := strings.Index(output, ordered[i])
		if prevIdx < 0 {
			t.Errorf("Expected %q to be present in output, got:\n%s", ordered[i-1], output)
			continue
		}
		if currIdx < 0 {
			t.Errorf("Expected %q to be present in output, got:\n%s", ordered[i], output)
			continue
		}
		if prevIdx >= currIdx {
			t.Errorf("Expected %q before %q in output, got:\n%s", ordered[i-1], ordered[i], output)
		}
	}
}

// TestXmlReader_StructuredModeWithComments tests comment preservation in structured mode
func TestXmlReader_StructuredModeWithComments(t *testing.T) {
	t.Run("basic_comment_structured", func(t *testing.T) {
		options := parsing.DefaultReaderOptions()
		options.Ext = map[string]string{"xml-mode": "structured"}
		r, err := daselxml.XML.NewReader(options)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		input := `<!--comment--><root><child>text</child></root>`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		childrenNode, err := data.GetMapKey("children")
		if err != nil {
			t.Fatalf("Expected 'children' key in structured mode: %s", err)
		}

		childrenLen, err := childrenNode.SliceLen()
		if err != nil {
			t.Fatalf("Expected children to be a slice: %s", err)
		}

		if childrenLen == 0 {
			t.Fatalf("Expected at least one child element")
		}

		rootElement, err := childrenNode.GetSliceIndex(0)
		if err != nil {
			t.Fatalf("Expected to get first child: %s", err)
		}

		comments, ok := rootElement.MetadataValue("xml_comments")
		if !ok {
			t.Errorf("Expected xml_comments metadata to exist")
		}

		if comments == nil {
			t.Errorf("Expected comments to be preserved in structured mode")
		}
	})
}

// TestXmlRoundTrip_CommentPreservation tests round-trip preservation of XML comments
func TestXmlRoundTrip_CommentPreservation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "single comment",
			input: `<!--This is a comment-->
<root>
  <child>text</child>
</root>
`,
			expected: []string{"<!--This is a comment-->", "<root>", "<child>text</child>"},
		},
		{
			name: "multiple comments before root",
			input: `<!--First comment-->
<!--Second comment-->
<root>
  <child>text</child>
</root>
`,
			expected: []string{"<!--First comment-->", "<!--Second comment-->"},
		},
		{
			name:     "comment before complex child element",
			input:    `<root><!--Section comment--><section><item>text</item></section></root>`,
			expected: []string{"<!--Section comment-->", "<section>"},
		},
		{
			name: "processing instruction and comments",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!--Document comment-->
<root>
  <child>text</child>
</root>
`,
			expected: []string{`<?xml version="1.0" encoding="UTF-8"?>`, "<!--Document comment-->", "<root>"},
		},
		{
			name: "nested elements and comments",
			input: `<!--Root level comment-->
<Document>
  <!--Section comment-->
  <Section>
    <Item>value1</Item>
    <Item>value2</Item>
  </Section>
</Document>
`,
			expected: []string{"<!--Root level comment-->", "<!--Section comment-->", "<Document>", "<Section>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := newTestReaderWriter(t)
			assertRoundTrip(t, r, w, tt.input, tt.expected)
		})
	}
}

// TestXmlRoundTrip_EdgeCases tests edge cases for comment handling
func TestXmlRoundTrip_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty comment",
			input:    `<!----><root><child>text</child></root>`,
			expected: []string{"<!---->", "<root>"},
		},
		{
			name:     "whitespace-only comment",
			input:    `<!--   --><root><child>text</child></root>`,
			expected: []string{"<!--   -->"},
		},
		{
			name:     "comment between sibling elements",
			input:    `<root><first>one</first><!--between siblings--><second>two</second></root>`,
			expected: []string{"<!--between siblings-->", "<first>one</first>", "<second>two</second>"},
		},
		{
			name:     "multiple comments between siblings",
			input:    `<root><a>1</a><!--comment1--><!--comment2--><b>2</b></root>`,
			expected: []string{"<!--comment1-->", "<!--comment2-->"},
		},
		{
			name:     "trailing comment after last child",
			input:    `<root><child>text</child><!--trailing comment--></root>`,
			expected: []string{"<!--trailing comment-->", "<child>text</child>"},
		},
		{
			name:     "multiple trailing comments",
			input:    `<root><child>text</child><!--trailing1--><!--trailing2--></root>`,
			expected: []string{"<!--trailing1-->", "<!--trailing2-->"},
		},
		{
			name:     "trailing comment in nested structure",
			input:    `<root><parent><child>text</child><!--inner trailing--></parent><!--outer trailing--></root>`,
			expected: []string{"<!--inner trailing-->", "<!--outer trailing-->"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := newTestReaderWriter(t)
			assertRoundTrip(t, r, w, tt.input, tt.expected)
		})
	}
}

// TestXmlRoundTrip_SpecialCommentContent tests comments with special but valid content
func TestXmlRoundTrip_SpecialCommentContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "special characters",
			input:    `<!--Comment with <special> & "characters" 'here'--><root><child>text</child></root>`,
			expected: []string{"<!--Comment with <special>"},
		},
		{
			name:     "single dash",
			input:    `<!--Comment with a-single-dash--><root><child>text</child></root>`,
			expected: []string{"<!--Comment with a-single-dash-->"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := newTestReaderWriter(t)
			assertRoundTrip(t, r, w, tt.input, tt.expected)
		})
	}
}

// TestXmlReader_SecurityLimits tests error handling for security limits
func TestXmlReader_SecurityLimits(t *testing.T) {
	t.Run("reject oversized XML input", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		largeContent := strings.Repeat("x", 10_000_001)
		input := "<root>" + largeContent + "</root>"

		_, err = r.Read([]byte(input))
		if err == nil {
			t.Errorf("Expected error for oversized XML input")
		}
		if !strings.Contains(err.Error(), "exceeds maximum size") {
			t.Errorf("Expected error about maximum size, got: %s", err)
		}
	})

	t.Run("reject oversized comment", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		largeComment := strings.Repeat("x", 10_001)
		input := "<!--" + largeComment + "--><root><child>text</child></root>"

		_, err = r.Read([]byte(input))
		if err == nil {
			t.Errorf("Expected error for oversized comment")
		}
		if !strings.Contains(err.Error(), "exceeds maximum length") {
			t.Errorf("Expected error about maximum length, got: %s", err)
		}
	})

	t.Run("reject document with too many comments", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		var builder strings.Builder
		builder.WriteString("<root>")
		for i := 0; i < 1001; i++ {
			builder.WriteString("<!--comment-->")
		}
		builder.WriteString("<child>text</child></root>")
		input := builder.String()

		_, err = r.Read([]byte(input))
		if err == nil {
			t.Errorf("Expected error for too many comments")
		}
		if !strings.Contains(err.Error(), "exceeds maximum comment count") {
			t.Errorf("Expected error about maximum comment count, got: %s", err)
		}
	})
}

// TestXmlRoundTrip_ProcessingInstructionReset tests that processing instructions are not duplicated across siblings
func TestXmlRoundTrip_ProcessingInstructionReset(t *testing.T) {
	t.Run("processing instructions not duplicated to siblings", func(t *testing.T) {
		r, w := newTestReaderWriter(t)

		input := `<?xml version="1.0"?><root><first>one</first><second>two</second></root>`

		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)

		piCount := strings.Count(outputStr, `<?xml version="1.0"?>`)
		if piCount != 1 {
			t.Errorf("Expected PI to appear exactly once, but found %d occurrences in:\n%s", piCount, outputStr)
		}

		if !strings.Contains(outputStr, "<first>one</first>") {
			t.Errorf("Expected first element to be preserved in:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<second>two</second>") {
			t.Errorf("Expected second element to be preserved in:\n%s", outputStr)
		}
	})

	t.Run("sibling elements do not inherit each others processing instructions", func(t *testing.T) {
		options := parsing.DefaultReaderOptions()
		options.Ext = map[string]string{"xml-mode": "structured"}
		r, err := daselxml.XML.NewReader(options)
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		input := `<root><?target instruction?><first>one</first><second>two</second></root>`

		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		children, err := data.GetMapKey("children")
		if err != nil {
			t.Fatalf("Expected children key: %s", err)
		}

		childLen, _ := children.SliceLen()
		if childLen < 1 {
			t.Fatalf("Expected at least one child")
		}

		rootEl, err := children.GetSliceIndex(0)
		if err != nil {
			t.Fatalf("Expected to get root element: %s", err)
		}

		rootChildren, err := rootEl.GetMapKey("children")
		if err != nil {
			t.Fatalf("Expected root children: %s", err)
		}

		rootChildLen, _ := rootChildren.SliceLen()
		if rootChildLen < 2 {
			t.Fatalf("Expected at least two children in root, got %d", rootChildLen)
		}

		firstChild, _ := rootChildren.GetSliceIndex(0)
		_, hasPI := firstChild.MetadataValue("xml_processing_instructions")
		if !hasPI {
			t.Log("First child does not have PI (which is expected behavior after fix)")
		}

		secondChild, _ := rootChildren.GetSliceIndex(1)
		_, secondHasPI := secondChild.MetadataValue("xml_processing_instructions")
		if secondHasPI {
			t.Errorf("Second child should NOT have processing instructions (PI was incorrectly duplicated)")
		}
	})
}

// TestXmlRoundTrip_ProcessingInstructions tests round-trip preservation of processing instructions
func TestXmlRoundTrip_ProcessingInstructions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "xml declaration",
			input: `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<root>
  <child>text</child>
</root>
`,
			expected: []string{`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`},
		},
		{
			name: "multiple processing instructions",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="style.xsl"?>
<root>
  <child>text</child>
</root>
`,
			expected: []string{`<?xml version="1.0" encoding="UTF-8"?>`, `<?xml-stylesheet type="text/xsl" href="style.xsl"?>`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := newTestReaderWriter(t)
			assertRoundTrip(t, r, w, tt.input, tt.expected)
		})
	}
}
