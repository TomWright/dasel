package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncReplace(t *testing.T) {
	t.Run("input arg", testCase{
		s:   `replace("hello world", "world", "there")`,
		out: model.NewStringValue("hello there"),
	}.run)

	t.Run("multiple data arg", testCase{
		s:   `replace("hello world", "o", "0", "l", "1")`,
		out: model.NewStringValue("he110 w0r1d"),
	}.run)

	t.Run("data arg", testCase{
		s:   `"hello world".replace("o", "0")`,
		out: model.NewStringValue("hell0 w0rld"),
	}.run)

	t.Run("multiple data arg", testCase{
		s:   `"hello world".replace("o", "0", "h", "H")`,
		out: model.NewStringValue("Hell0 w0rld"),
	}.run)

	t.Run("data arg with input arg ignores data arg", testCase{
		s:   `"bob".replace("hello world", "o", "0", "world", "there")`,
		out: model.NewStringValue("hell0 there"),
	}.run)
}
