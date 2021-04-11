package dasel_test

import (
	"github.com/tomwright/dasel"
	"testing"
)

func testExtractNextSelector(in string, exp string, expRead int) func(t *testing.T) {
	return func(t *testing.T) {
		got, read := dasel.ExtractNextSelector(in)
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
		if read != expRead {
			t.Errorf("expected read of %d, got %d", expRead, read)
		}
	}
}

func TestExtractNextSelector(t *testing.T) {
	t.Run("Simple", testExtractNextSelector(`.metadata.name`, `.metadata`, 9))
	t.Run("EscapedDot", testExtractNextSelector(`.before\.after.name`, `.before.after`, 14))
	t.Run("EscapedSpace", testExtractNextSelector(`.before\ after.name`, `.before after`, 14))
	t.Run("DynamicWithPath", testExtractNextSelector(`.(.before.a=b).after.name`, `.(.before.a=b)`, 14))
	t.Run("EscapedFirstDot", testExtractNextSelector(`\.name`, `.name`, 6))
	t.Run("SimpleProp", testExtractNextSelector(`.name`, `.name`, 5))
	t.Run("SimpleIndex", testExtractNextSelector(`.[123]`, `.[123]`, 6))
	t.Run("SimpleLength", testExtractNextSelector(`.[#]`, `.[#]`, 4))
}
