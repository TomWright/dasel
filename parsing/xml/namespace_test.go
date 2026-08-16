package xml_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/xml"
)

// TestXMLNamespaceRoundTrip covers namespace prefixes surviving a read/write
// cycle. encoding/xml resolves `soap:Envelope` into a Name whose Space is the
// namespace URI and whose Local is `Envelope`, dropping the prefix, so a naive
// reader turns `xmlns:soap="..."` into a plain `soap="..."` attribute and
// renames every element in that namespace.
func TestXMLNamespaceRoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		exp  string
	}{
		{
			name: "prefixed elements and namespace declarations",
			in: `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <m:GetPrice xmlns:m="https://example.com/prices">
      <m:Item>Apples</m:Item>
    </m:GetPrice>
  </soap:Body>
</soap:Envelope>
`,
			exp: `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <m:GetPrice xmlns:m="https://example.com/prices">
      <m:Item>Apples</m:Item>
    </m:GetPrice>
  </soap:Body>
</soap:Envelope>
`,
		},
		{
			name: "prefixed attributes",
			in:   `<root xmlns:ns="https://example.com/ns"><item ns:id="a" plain="p">hi</item></root>`,
			exp: `<root xmlns:ns="https://example.com/ns">
  <item ns:id="a" plain="p">hi</item>
</root>
`,
		},
		{
			name: "implicit xml prefix needs no declaration",
			in:   `<root xml:lang="en"><a>hi</a></root>`,
			exp: `<root xml:lang="en">
  <a>hi</a>
</root>
`,
		},
		{
			name: "default namespace leaves names unprefixed",
			in:   `<root xmlns="https://example.com/def"><item>hi</item></root>`,
			exp: `<root xmlns="https://example.com/def">
  <item>hi</item>
</root>
`,
		},
		{
			name: "undeclared prefix is preserved verbatim",
			in:   `<root><un:item un:id="a">hi</un:item></root>`,
			exp: `<root>
  <un:item un:id="a">hi</un:item>
</root>
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := xml.XML.NewReader(parsing.DefaultReaderOptions())
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			w, err := xml.XML.NewWriter(parsing.DefaultWriterOptions())
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			v, err := r.Read([]byte(tc.in))
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			out, err := w.Write(v)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if string(out) != tc.exp {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.exp, string(out))
			}
		})
	}
}
