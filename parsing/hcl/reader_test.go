package hcl_test

import (
	"fmt"
	"testing"

	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/hcl"
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
	t.Run("document c", readTestCase{
		in: `image_id = "ami-123"
cluster_min_nodes = 2
cluster_decimal_nodes = 2.2
cluster_max_nodes = true
availability_zone_names = [
"us-east-1a",
"us-west-1c",
]
docker_ports = [{
internal = 8300
external = 8300
protocol = "tcp"
},
{
internal = 8301
external = 8301
protocol = "tcp"
}
]`,
	}.run)
}
