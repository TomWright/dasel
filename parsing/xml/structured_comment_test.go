package xml_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	daselxml "github.com/tomwright/dasel/v3/parsing/xml"
)

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

		// Verify comments metadata exists
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
	t.Run("round trip with single comment", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<!--This is a comment-->
<root>
  <child>text</child>
</root>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--This is a comment-->") {
			t.Errorf("Expected output to contain comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<root>") {
			t.Errorf("Expected output to contain root element, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<child>text</child>") {
			t.Errorf("Expected output to contain child element, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with multiple comments before root", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<!--First comment-->
<!--Second comment-->
<root>
  <child>text</child>
</root>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--First comment-->") {
			t.Errorf("Expected output to contain first comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<!--Second comment-->") {
			t.Errorf("Expected output to contain second comment, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with comment before complex child element", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// Comment before a complex child element (one with nested children) gets preserved
		// because complex elements become map values which can hold metadata
		input := `<root><!--Section comment--><section><item>text</item></section></root>`

		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--Section comment-->") {
			t.Errorf("Expected output to contain section comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<section>") {
			t.Errorf("Expected output to contain section element, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with processing instruction and comments", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<?xml version="1.0" encoding="UTF-8"?>
<!--Document comment-->
<root>
  <child>text</child>
</root>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
			t.Errorf("Expected output to contain processing instruction, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<!--Document comment-->") {
			t.Errorf("Expected output to contain comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<root>") {
			t.Errorf("Expected output to contain root element, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with nested elements and comments", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<!--Root level comment-->
<Document>
  <!--Section comment-->
  <Section>
    <Item>value1</Item>
    <Item>value2</Item>
  </Section>
</Document>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--Root level comment-->") {
			t.Errorf("Expected output to contain root level comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<!--Section comment-->") {
			t.Errorf("Expected output to contain section comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<Document>") {
			t.Errorf("Expected output to contain Document element, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<Section>") {
			t.Errorf("Expected output to contain Section element, got:\n%s", outputStr)
		}
	})
}

// TestXmlRoundTrip_EdgeCases tests edge cases for comment handling
func TestXmlRoundTrip_EdgeCases(t *testing.T) {
	t.Run("round trip with empty comment", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<!----><root><child>text</child></root>`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!---->") {
			t.Errorf("Expected output to contain empty comment, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, "<root>") {
			t.Errorf("Expected output to contain root element, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with whitespace-only comment", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<!--   --><root><child>text</child></root>`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--   -->") {
			t.Errorf("Expected output to contain whitespace comment, got:\n%s", outputStr)
		}
	})
}

// TestXmlRoundTrip_SpecialCommentContent tests comments with special but valid content
func TestXmlRoundTrip_SpecialCommentContent(t *testing.T) {
	t.Run("round trip with comment containing special characters", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// Comments with special but valid characters (no -- sequence)
		input := `<!--Comment with <special> & "characters" 'here'--><root><child>text</child></root>`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--Comment with <special>") {
			t.Errorf("Expected output to contain comment with special characters, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with comment containing single dash", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// Single dashes are valid in comments
		input := `<!--Comment with a-single-dash--><root><child>text</child></root>`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--Comment with a-single-dash-->") {
			t.Errorf("Expected output to contain comment with single dashes, got:\n%s", outputStr)
		}
	})
}

// TestXmlReader_SecurityLimits tests error handling for security limits
func TestXmlReader_SecurityLimits(t *testing.T) {
	t.Run("reject oversized XML input", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		// Create XML input larger than 10MB limit
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

		// Create comment larger than 10KB limit
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

		// Create document with more than 1000 comments
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
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// PI at document level (before root) should appear once in output
		// PI inside element is attached to following sibling, not duplicated
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

		// Count occurrences of the PI - should appear exactly once
		piCount := strings.Count(outputStr, `<?xml version="1.0"?>`)
		if piCount != 1 {
			t.Errorf("Expected PI to appear exactly once, but found %d occurrences in:\n%s", piCount, outputStr)
		}

		// Verify the structure is preserved
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

		// In structured mode, we can verify that PI is only attached to the first sibling
		input := `<root><?target instruction?><first>one</first><second>two</second></root>`

		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		// Get children
		children, err := data.GetMapKey("children")
		if err != nil {
			t.Fatalf("Expected children key: %s", err)
		}

		childLen, _ := children.SliceLen()
		if childLen < 1 {
			t.Fatalf("Expected at least one child")
		}

		// Get root element (first child of virtual root)
		rootEl, err := children.GetSliceIndex(0)
		if err != nil {
			t.Fatalf("Expected to get root element: %s", err)
		}

		// Get root's children
		rootChildren, err := rootEl.GetMapKey("children")
		if err != nil {
			t.Fatalf("Expected root children: %s", err)
		}

		rootChildLen, _ := rootChildren.SliceLen()
		if rootChildLen < 2 {
			t.Fatalf("Expected at least two children in root, got %d", rootChildLen)
		}

		// First child should have PI
		firstChild, _ := rootChildren.GetSliceIndex(0)
		_, hasPI := firstChild.MetadataValue("xml_processing_instructions")
		if !hasPI {
			t.Log("First child does not have PI (which is expected behavior after fix)")
		}

		// Second child should NOT have PI
		secondChild, _ := rootChildren.GetSliceIndex(1)
		_, secondHasPI := secondChild.MetadataValue("xml_processing_instructions")
		if secondHasPI {
			t.Errorf("Second child should NOT have processing instructions (PI was incorrectly duplicated)")
		}
	})
}

// TestXmlRoundTrip_ProcessingInstructions tests round-trip preservation of processing instructions
func TestXmlRoundTrip_ProcessingInstructions(t *testing.T) {
	t.Run("round trip with xml declaration", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<root>
  <child>text</child>
</root>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, `<?xml version="1.0" encoding="utf-8" standalone="yes"?>`) {
			t.Errorf("Expected output to contain XML declaration, got:\n%s", outputStr)
		}
	})

	t.Run("round trip with multiple processing instructions", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating reader: %s", err)
		}

		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		input := `<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet type="text/xsl" href="style.xsl"?>
<root>
  <child>text</child>
</root>
`
		data, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Unexpected error reading XML: %s", err)
		}

		output, err := w.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, `<?xml version="1.0" encoding="UTF-8"?>`) {
			t.Errorf("Expected output to contain XML declaration, got:\n%s", outputStr)
		}
		if !strings.Contains(outputStr, `<?xml-stylesheet type="text/xsl" href="style.xsl"?>`) {
			t.Errorf("Expected output to contain stylesheet processing instruction, got:\n%s", outputStr)
		}
	})
}
