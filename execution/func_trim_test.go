package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncTrim(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `"  hello  ".trim()`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("arg input", testCase{
		s:   `trim("  hello  ")`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("tabs and newlines", testCase{
		s:   `trim("\n\thello\t\n")`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("no whitespace", testCase{
		s:   `"hello".trim()`,
		out: model.NewStringValue("hello"),
	}.run)
}

func TestFuncTrimPrefix(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `"hello world".trimPrefix("hello ")`,
		out: model.NewStringValue("world"),
	}.run)
	t.Run("arg input", testCase{
		s:   `trimPrefix("hello world", "hello ")`,
		out: model.NewStringValue("world"),
	}.run)
	t.Run("no match", testCase{
		s:   `"hello world".trimPrefix("xyz")`,
		out: model.NewStringValue("hello world"),
	}.run)
}

func TestFuncTrimSuffix(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `"hello world".trimSuffix(" world")`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("arg input", testCase{
		s:   `trimSuffix("hello world", " world")`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("no match", testCase{
		s:   `"hello world".trimSuffix("xyz")`,
		out: model.NewStringValue("hello world"),
	}.run)
}
