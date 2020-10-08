package command_test

import (
	"testing"
)

func TestRootCMD(t *testing.T) {
	t.Run("Select", func(t *testing.T) {
		t.Run("JSON", selectTestForParser("json", jsonData, jsonDataSingle))
		t.Run("YAML", selectTestForParser("yaml", yamlData, yamlDataSingle))
		t.Run("TOML", selectTestForParser("toml", tomlData, tomlDataSingle))
	})
	t.Run("PutString", func(t *testing.T) {
		t.Run("JSON", putStringTestForParserJSON("json"))
		t.Run("YAML", putStringTestForParserYAML("yaml"))
		t.Run("TOML", putStringTestForParserTOML("toml"))
	})
	t.Run("PutInt", func(t *testing.T) {
		t.Run("JSON", putIntTestForParserJSON("json"))
		t.Run("YAML", putIntTestForParserYAML("yaml"))
		t.Run("TOML", putIntTestForParserTOML("toml"))
	})
	t.Run("PutBool", func(t *testing.T) {
		t.Run("JSON", putBoolTestForParserJSON("json"))
		t.Run("YAML", putBoolTestForParserYAML("yaml"))
		t.Run("TOML", putBoolTestForParserTOML("toml"))
	})
}
