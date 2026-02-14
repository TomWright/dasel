package xml_test

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	daselxml "github.com/tomwright/dasel/v3/parsing/xml"
)

// TestXmlRoundTrip_ElementOrdering tests that interleaved same-named siblings
// preserve their original document order during round-trips (Issue #196).
func TestXmlRoundTrip_ElementOrdering(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expected        string   // exact expected output (use empty string to skip exact match)
		contains        []string // strings that must appear in output
		orderedContains []string // strings that must appear in this order
	}{
		{
			name: "issue 196 reproducer",
			input: `<message>
  <heading>Look</heading>
  <warning>Hello World</warning>
  <heading>Above</heading>
  <string>Last</string>
</message>`,
			expected: `<message>
  <heading>Look</heading>
  <warning>Hello World</warning>
  <heading>Above</heading>
  <string>Last</string>
</message>
`,
		},
		{
			name: "interleaved a-b-a pattern",
			input: `<root>
  <a>1</a>
  <b>2</b>
  <a>3</a>
</root>`,
			expected: `<root>
  <a>1</a>
  <b>2</b>
  <a>3</a>
</root>
`,
		},
		{
			name: "complex interleaving a-b-c-a-b",
			input: `<root>
  <a>1</a>
  <b>2</b>
  <c>3</c>
  <a>4</a>
  <b>5</b>
</root>`,
			expected: `<root>
  <a>1</a>
  <b>2</b>
  <c>3</c>
  <a>4</a>
  <b>5</b>
</root>
`,
		},
		{
			name: "all same name children",
			input: `<root>
  <item>1</item>
  <item>2</item>
  <item>3</item>
</root>`,
			expected: `<root>
  <item>1</item>
  <item>2</item>
  <item>3</item>
</root>
`,
		},
		{
			name: "no interleaving - unique children",
			input: `<root>
  <a>1</a>
  <b>2</b>
  <c>3</c>
</root>`,
			expected: `<root>
  <a>1</a>
  <b>2</b>
  <c>3</c>
</root>
`,
		},
		{
			name: "deeply nested interleaving",
			input: `<root>
  <parent>
    <x>1</x>
    <y>2</y>
    <x>3</x>
  </parent>
  <other>text</other>
  <parent>
    <a>4</a>
    <b>5</b>
    <a>6</a>
  </parent>
</root>`,
			contains: []string{
				"<x>1</x>",
				"<y>2</y>",
				"<x>3</x>",
				"<a>4</a>",
				"<b>5</b>",
				"<a>6</a>",
			},
			orderedContains: []string{
				"<x>1</x>",
				"<y>2</y>",
				"<x>3</x>",
			},
		},
		{
			name: "elements with attributes interleaved",
			input: `<root>
  <item id="1">first</item>
  <other>middle</other>
  <item id="2">second</item>
</root>`,
			contains: []string{
				`<item id="1">first</item>`,
				"<other>middle</other>",
				`<item id="2">second</item>`,
			},
			orderedContains: []string{
				`<item id="1">first</item>`,
				"<other>middle</other>",
				`<item id="2">second</item>`,
			},
		},
		{
			name: "empty elements interleaved",
			input: `<root>
  <a></a>
  <b>text</b>
  <a></a>
</root>`,
			expected: `<root>
  <a></a>
  <b>text</b>
  <a></a>
</root>
`,
		},
		{
			name: "comments with interleaved elements",
			input: `<root>
  <!-- first comment -->
  <a>1</a>
  <!-- middle comment -->
  <b>2</b>
  <a>3</a>
</root>`,
			contains: []string{
				"<a>1</a>",
				"<b>2</b>",
				"<a>3</a>",
			},
			orderedContains: []string{
				"<a>1</a>",
				"<b>2</b>",
				"<a>3</a>",
			},
		},
		{
			name: "parent with text content and interleaved children",
			input: `<doc>
  <title>Hello</title>
  <para>World</para>
  <title>Goodbye</title>
</doc>`,
			expected: `<doc>
  <title>Hello</title>
  <para>World</para>
  <title>Goodbye</title>
</doc>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w := newTestReaderWriter(t)

			data, err := r.Read([]byte(tt.input))
			if err != nil {
				t.Fatalf("Unexpected error reading XML: %s", err)
			}

			output, err := w.Write(data)
			if err != nil {
				t.Fatalf("Unexpected error writing XML: %s", err)
			}

			outputStr := string(output)

			if tt.expected != "" {
				if outputStr != tt.expected {
					t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, outputStr)
				}
			}

			for _, exp := range tt.contains {
				if !strings.Contains(outputStr, exp) {
					t.Errorf("Expected output to contain %q, got:\n%s", exp, outputStr)
				}
			}

			// Verify element ordering for contains-based tests.
			if len(tt.orderedContains) > 1 {
				for i := 1; i < len(tt.orderedContains); i++ {
					prevIdx := strings.Index(outputStr, tt.orderedContains[i-1])
					currIdx := strings.Index(outputStr, tt.orderedContains[i])
					if prevIdx < 0 || currIdx < 0 {
						continue
					}
					if prevIdx >= currIdx {
						t.Errorf("Expected %q before %q in output, got:\n%s",
							tt.orderedContains[i-1], tt.orderedContains[i], outputStr)
					}
				}
			}
		})
	}
}

// TestXmlRoundTrip_ElementOrderingMetadata tests that xml_child_order metadata
// is correctly populated during read.
func TestXmlRoundTrip_ElementOrderingMetadata(t *testing.T) {
	t.Run("metadata records full child order", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<root>
  <a>1</a>
  <b>2</b>
  <a>3</a>
</root>`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		rootVal, err := data.GetMapKey("root")
		if err != nil {
			t.Fatalf("Expected 'root' key: %s", err)
		}

		orderMeta, ok := rootVal.MetadataValue("xml_child_order")
		if !ok {
			t.Fatal("Expected xml_child_order metadata to exist")
		}

		order, ok := orderMeta.([]string)
		if !ok {
			t.Fatal("Expected xml_child_order to be []string")
		}

		expected := []string{"a", "b", "a"}
		if len(order) != len(expected) {
			t.Fatalf("Expected order length %d, got %d: %v", len(expected), len(order), order)
		}
		for i, name := range expected {
			if order[i] != name {
				t.Errorf("Expected order[%d] = %q, got %q", i, name, order[i])
			}
		}
	})

	t.Run("metadata set on elements with children", func(t *testing.T) {
		r, err := daselxml.XML.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		data, err := r.Read([]byte(`<root><leaf>text</leaf></root>`))
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		rootVal, err := data.GetMapKey("root")
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// root has a child (leaf), so it should have ordering metadata
		_, ok := rootVal.MetadataValue("xml_child_order")
		if !ok {
			t.Error("Expected xml_child_order metadata on root since it has children")
		}
	})
}

// TestXmlRoundTrip_FallbackWithoutMetadata tests that the writer correctly
// falls back to insertion-order iteration when xml_child_order is absent.
func TestXmlRoundTrip_FallbackWithoutMetadata(t *testing.T) {
	t.Run("programmatically constructed value without metadata", func(t *testing.T) {
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		// Build a map without xml_child_order metadata (simulating non-XML source).
		root := model.NewMapValue()
		inner := model.NewMapValue()
		_ = inner.SetMapKey("b", model.NewStringValue("2"))
		_ = inner.SetMapKey("a", model.NewStringValue("1"))
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		// Without metadata, should use insertion order: b before a
		bIdx := strings.Index(outputStr, "<b>")
		aIdx := strings.Index(outputStr, "<a>")
		if bIdx < 0 || aIdx < 0 {
			t.Fatalf("Expected both <b> and <a> in output:\n%s", outputStr)
		}
		if bIdx > aIdx {
			t.Errorf("Expected <b> before <a> (insertion order), got:\n%s", outputStr)
		}
	})
}

// TestXmlWriter_OrderingEdgeCases tests defensive code paths in the writer's
// ordering reconstruction logic.
func TestXmlWriter_OrderingEdgeCases(t *testing.T) {
	t.Run("stale metadata - deleted key", func(t *testing.T) {
		// Build a Value with xml_child_order referencing a key not in the map.
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		root := model.NewMapValue()
		inner := model.NewMapValue()
		_ = inner.SetMapKey("a", model.NewStringValue("1"))
		// "b" is NOT in the map, but IS in the ordering metadata.
		inner.SetMetadataValue("xml_child_order", []string{"a", "b", "a"})
		// Only "a" as scalar, so second "a" in metadata will also be skipped (scalar emitted once).
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<a>1</a>") {
			t.Errorf("Expected <a>1</a> in output, got:\n%s", outputStr)
		}
		// Should not panic or produce corrupt output.
	})

	t.Run("counter overflow - more metadata entries than slice values", func(t *testing.T) {
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		root := model.NewMapValue()
		inner := model.NewMapValue()
		items := model.NewSliceValue()
		_ = items.Append(model.NewStringValue("first"))
		// Slice has 1 element, but metadata says 3 occurrences.
		_ = inner.SetMapKey("item", items)
		inner.SetMetadataValue("xml_child_order", []string{"item", "item", "item"})
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		// Should produce only 1 <item>, not panic.
		itemCount := strings.Count(outputStr, "<item>")
		if itemCount != 1 {
			t.Errorf("Expected exactly 1 <item>, got %d in:\n%s", itemCount, outputStr)
		}
	})

	t.Run("new keys not in ordering - appended at end", func(t *testing.T) {
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		root := model.NewMapValue()
		inner := model.NewMapValue()
		_ = inner.SetMapKey("a", model.NewStringValue("1"))
		_ = inner.SetMapKey("b", model.NewStringValue("2"))
		_ = inner.SetMapKey("c", model.NewStringValue("new"))
		// Metadata only knows about "a" and "b", not "c".
		inner.SetMetadataValue("xml_child_order", []string{"b", "a"})
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		// "b" should come before "a" (per metadata), "c" appended at end.
		bIdx := strings.Index(outputStr, "<b>")
		aIdx := strings.Index(outputStr, "<a>")
		cIdx := strings.Index(outputStr, "<c>")
		if bIdx < 0 || aIdx < 0 || cIdx < 0 {
			t.Fatalf("Expected all three elements in output:\n%s", outputStr)
		}
		if bIdx > aIdx {
			t.Errorf("Expected <b> before <a> per metadata ordering, got:\n%s", outputStr)
		}
		if cIdx < aIdx {
			t.Errorf("Expected <c> after <a> (appended at end), got:\n%s", outputStr)
		}
	})

	t.Run("invalid metadata type - fallback to insertion order", func(t *testing.T) {
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		root := model.NewMapValue()
		inner := model.NewMapValue()
		_ = inner.SetMapKey("b", model.NewStringValue("2"))
		_ = inner.SetMapKey("a", model.NewStringValue("1"))
		// Set invalid metadata type (int instead of []string).
		inner.SetMetadataValue("xml_child_order", 42)
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		// Should fall back to insertion order: b before a.
		bIdx := strings.Index(outputStr, "<b>")
		aIdx := strings.Index(outputStr, "<a>")
		if bIdx < 0 || aIdx < 0 {
			t.Fatalf("Expected both elements in output:\n%s", outputStr)
		}
		if bIdx > aIdx {
			t.Errorf("Expected fallback to insertion order (<b> before <a>), got:\n%s", outputStr)
		}
	})

	t.Run("scalar duplicate guard - emit only once", func(t *testing.T) {
		w, err := daselxml.XML.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		root := model.NewMapValue()
		inner := model.NewMapValue()
		_ = inner.SetMapKey("a", model.NewStringValue("only-once"))
		// Metadata says "a" twice, but value is scalar (not slice).
		inner.SetMetadataValue("xml_child_order", []string{"a", "a"})
		_ = root.SetMapKey("root", inner)

		output, err := w.Write(root)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		outputStr := string(output)
		aCount := strings.Count(outputStr, "<a>")
		if aCount != 1 {
			t.Errorf("Expected exactly 1 <a> element, got %d in:\n%s", aCount, outputStr)
		}
	})

	t.Run("idempotent round-trip", func(t *testing.T) {
		r, w := newTestReaderWriter(t)

		input := `<root>
  <a>1</a>
  <b>2</b>
  <a>3</a>
</root>`

		// First round-trip
		data1, err := r.Read([]byte(input))
		if err != nil {
			t.Fatalf("Read 1 error: %s", err)
		}
		output1, err := w.Write(data1)
		if err != nil {
			t.Fatalf("Write 1 error: %s", err)
		}

		// Second round-trip
		data2, err := r.Read(output1)
		if err != nil {
			t.Fatalf("Read 2 error: %s", err)
		}
		output2, err := w.Write(data2)
		if err != nil {
			t.Fatalf("Write 2 error: %s", err)
		}

		if string(output1) != string(output2) {
			t.Errorf("Round-trip not idempotent.\nFirst:\n%s\nSecond:\n%s", string(output1), string(output2))
		}
	})
}
