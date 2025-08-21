package hcl_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/hcl"
)

type readWriteTestCase struct {
	in string
}

func (tc readWriteTestCase) run(t *testing.T) {
	r, err := hcl.HCL.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w, err := hcl.HCL.NewWriter(parsing.DefaultWriterOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	in := []byte(tc.in)

	data, err := r.Read(in)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	got, err := w.Write(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	gotStr := string(got)

	if !cmp.Equal(tc.in, gotStr) {
		t.Errorf("unexpected output: %s", cmp.Diff(tc.in, gotStr))
	}
}

func TestHclReader_ReadWrite(t *testing.T) {
	t.Run("document a", readWriteTestCase{
		in: `io_mode = "async"
service {
  http {
    listen_addr = "127.0.0.1:8080"
    process {
      main {
        command = ["/usr/local/bin/awesome-app", "server"]
      }
      mgmt {
        command = ["/usr/local/bin/awesome-app", "mgmt"]
      }
      mgmt {
        command = ["/usr/local/bin/awesome-app", "mgmt2"]
      }
    }
  }
}
`,
	}.run)
}
