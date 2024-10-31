package hcl_test

import (
	"fmt"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/hcl"
	"testing"
)

type readTestCase struct {
	in string
}

func (tc readTestCase) run(t *testing.T) {
	r, err := hcl.HCL.NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	in := []byte(tc.in)

	got, err := r.Read(in)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	fmt.Println(got)
}

func TestHclReader_Read(t *testing.T) {
	t.Run("document a", readTestCase{
		in: `io_mode = "async"

service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"

  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "mgmt" {
    command = ["/usr/local/bin/awesome-app", "mgmt"]
  }
}`,
	}.run)
	t.Run("document b", readTestCase{
		in: `resource "aws_instance" "example" {
  # (resource configuration omitted for brevity)

  provisioner "local-exec" {
    command = "echo 'Hello World' >example.txt"
  }
  provisioner "file" {
    source      = "example.txt"
    destination = "/tmp/example.txt"
  }
  provisioner "remote-exec" {
    inline = [
      "sudo install-something -f /tmp/example.txt",
    ]
  }
}`,
	}.run)
}
