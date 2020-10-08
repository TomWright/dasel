package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"strings"
	"testing"
)

func TestRootCMD(t *testing.T) {
	t.Run("Select", func(t *testing.T) {
		t.Run("JSON", selectTestForParser("json", jsonData, jsonDataSingle))
		t.Run("YAML", selectTestForParser("yaml", yamlData, yamlDataSingle))
		t.Run("TOML", selectTestForParser("toml", tomlData, tomlDataSingle))
		t.Run("InvalidFile", func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			args := []string{
				"select", "-f", "bad.json", "-s", "x",
			}

			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err == nil || !strings.Contains(err.Error(), "could not open input file") {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
		t.Run("MissingParser", func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			args := []string{
				"select", "-s", "x",
			}

			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err == nil || !strings.Contains(err.Error(), "parser flag required when reading from stdin") {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
	})
	t.Run("PutString", func(t *testing.T) {
		t.Run("JSON", putStringTestForParserJSON())
		t.Run("YAML", putStringTestForParserYAML())
		t.Run("TOML", putStringTestForParserTOML())
		t.Run("InvalidFile", func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			args := []string{
				"put", "string", "-f", "bad.json", "-s", "x", "y",
			}

			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err == nil || !strings.Contains(err.Error(), "could not open input file") {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
		t.Run("MissingParser", func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			args := []string{
				"put", "string", "-s", "x", "y",
			}

			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err == nil || !strings.Contains(err.Error(), "parser flag required when reading from stdin") {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
	})
	t.Run("PutInt", func(t *testing.T) {
		t.Run("JSON", putIntTestForParserJSON())
		t.Run("YAML", putIntTestForParserYAML())
		t.Run("TOML", putIntTestForParserTOML())
	})
	t.Run("PutBool", func(t *testing.T) {
		t.Run("JSON", putBoolTestForParserJSON())
		t.Run("YAML", putBoolTestForParserYAML())
		t.Run("TOML", putBoolTestForParserTOML())
	})
	t.Run("PutObject", func(t *testing.T) {
		t.Run("JSON", putObjectTestForParserJSON())
		t.Run("YAML", putObjectTestForParserYAML())
		t.Run("TOML", putObjectTestForParserTOML())
	})
}
