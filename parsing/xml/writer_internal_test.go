package xml

import (
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// TestXmlWriter_CommentValidation tests validation of comment content during write
func TestXmlWriter_CommentValidation(t *testing.T) {
	t.Run("reject comment containing double dash sequence", func(t *testing.T) {
		w, err := newXMLWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// Create a child value with an invalid comment containing -- sequence
		// The comments are attached to child elements, not the root
		childValue := model.NewMapValue()
		childValue.SetMetadataValue("xml_comments", []*xmlComment{{Text: "invalid--comment"}})
		if err := childValue.SetMapKey("item", model.NewStringValue("text")); err != nil {
			t.Fatalf("Unexpected error setting map key: %s", err)
		}

		// Create root value containing the child
		rootValue := model.NewMapValue()
		if err := rootValue.SetMapKey("child", childValue); err != nil {
			t.Fatalf("Unexpected error setting map key: %s", err)
		}

		_, err = w.Write(rootValue)
		if err == nil {
			t.Errorf("Expected error for comment containing double dash")
		}
		if err != nil && !strings.Contains(err.Error(), "cannot contain '--'") {
			t.Errorf("Expected error about double dash sequence, got: %s", err)
		}
	})

	t.Run("accept valid comment without double dash", func(t *testing.T) {
		w, err := newXMLWriter(parsing.DefaultWriterOptions())
		if err != nil {
			t.Fatalf("Unexpected error creating writer: %s", err)
		}

		// Create a child value with a valid comment
		childValue := model.NewMapValue()
		childValue.SetMetadataValue("xml_comments", []*xmlComment{{Text: "valid comment with single - dash"}})
		if err := childValue.SetMapKey("item", model.NewStringValue("text")); err != nil {
			t.Fatalf("Unexpected error setting map key: %s", err)
		}

		// Create root value containing the child
		rootValue := model.NewMapValue()
		if err := rootValue.SetMapKey("child", childValue); err != nil {
			t.Fatalf("Unexpected error setting map key: %s", err)
		}

		output, err := w.Write(rootValue)
		if err != nil {
			t.Fatalf("Unexpected error writing XML: %s", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "<!--valid comment with single - dash-->") {
			t.Errorf("Expected output to contain valid comment, got:\n%s", outputStr)
		}
	})
}

// Test_valueToString tests the valueToString function for all supported types
func Test_valueToString(t *testing.T) {
	t.Run("null value returns empty string", func(t *testing.T) {
		v := model.NewNullValue()
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for null value: %s", err)
		}
		if result != "" {
			t.Errorf("Expected empty string for null value, got: %q", result)
		}
	})

	t.Run("string value", func(t *testing.T) {
		v := model.NewStringValue("hello world")
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for string value: %s", err)
		}
		if result != "hello world" {
			t.Errorf("Expected 'hello world', got: %q", result)
		}
	})

	t.Run("empty string value", func(t *testing.T) {
		v := model.NewStringValue("")
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for empty string value: %s", err)
		}
		if result != "" {
			t.Errorf("Expected empty string, got: %q", result)
		}
	})

	t.Run("int value positive", func(t *testing.T) {
		v := model.NewIntValue(42)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for int value: %s", err)
		}
		if result != "42" {
			t.Errorf("Expected '42', got: %q", result)
		}
	})

	t.Run("int value negative", func(t *testing.T) {
		v := model.NewIntValue(-123)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for negative int value: %s", err)
		}
		if result != "-123" {
			t.Errorf("Expected '-123', got: %q", result)
		}
	})

	t.Run("int value zero", func(t *testing.T) {
		v := model.NewIntValue(0)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for zero int value: %s", err)
		}
		if result != "0" {
			t.Errorf("Expected '0', got: %q", result)
		}
	})

	t.Run("float value", func(t *testing.T) {
		v := model.NewFloatValue(3.14159)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for float value: %s", err)
		}
		if result != "3.14159" {
			t.Errorf("Expected '3.14159', got: %q", result)
		}
	})

	t.Run("float value negative", func(t *testing.T) {
		v := model.NewFloatValue(-2.5)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for negative float value: %s", err)
		}
		if result != "-2.5" {
			t.Errorf("Expected '-2.5', got: %q", result)
		}
	})

	t.Run("float value zero", func(t *testing.T) {
		v := model.NewFloatValue(0.0)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for zero float value: %s", err)
		}
		if result != "0" {
			t.Errorf("Expected '0', got: %q", result)
		}
	})

	t.Run("bool value true", func(t *testing.T) {
		v := model.NewBoolValue(true)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for bool true value: %s", err)
		}
		if result != "true" {
			t.Errorf("Expected 'true', got: %q", result)
		}
	})

	t.Run("bool value false", func(t *testing.T) {
		v := model.NewBoolValue(false)
		result, err := valueToString(v)
		if err != nil {
			t.Errorf("Unexpected error for bool false value: %s", err)
		}
		if result != "false" {
			t.Errorf("Expected 'false', got: %q", result)
		}
	})

	t.Run("map value returns error", func(t *testing.T) {
		v := model.NewMapValue()
		_, err := valueToString(v)
		if err == nil {
			t.Errorf("Expected error for map value")
		}
		if err != nil && !strings.Contains(err.Error(), "cannot format type") {
			t.Errorf("Expected error about formatting type, got: %s", err)
		}
	})

	t.Run("slice value returns error", func(t *testing.T) {
		v := model.NewSliceValue()
		_, err := valueToString(v)
		if err == nil {
			t.Errorf("Expected error for slice value")
		}
		if err != nil && !strings.Contains(err.Error(), "cannot format type") {
			t.Errorf("Expected error about formatting type, got: %s", err)
		}
	})
}
