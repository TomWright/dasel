package xml_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/xml"
)

type testCase struct {
	in     string
	assert func(t *testing.T, res *model.Value)
}

func (tc testCase) run(t *testing.T) {
	r, err := xml.XML.NewReader()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	res, err := r.Read([]byte(tc.in))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tc.assert(t, res)
}

type rwTestCase struct {
	in  string
	out string
}

func (tc rwTestCase) run(t *testing.T) {
	if tc.out == "" {
		tc.out = tc.in
	}
	r, err := xml.XML.NewReader()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	w, err := xml.XML.NewWriter(parsing.WriterOptions{})
	res, err := r.Read([]byte(tc.in))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	out, err := w.Write(res)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal([]byte(tc.out), out) {
		t.Errorf("unexpected output: %s", cmp.Diff(tc.out, string(out)))
	}
}

func TestYamlValue_UnmarshalXML(t *testing.T) {
	//t.Run("generic", rwTestCase{
	//	in: `<html>
	//<head>
	//    <title>Test</title>
	//</head>
	//<body>
	//    <h1>Test</h1>
	//    <p class="testing">Test</p>
	//    <div>
	//        <a href="test2.html">Test 2</a>
	//    </div>
	//	<div>
	//		<p>Hello</p>
	//		<p>World</p>
	//	</div>
	//</body>
	//</html>`,
	//	}.run)
}
