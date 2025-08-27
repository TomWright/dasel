package cli

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

type variable struct {
	Name  string
	Value *model.Value
}

type variables *[]variable

func variableOptions(vars variables) []execution.ExecuteOptionFn {
	var opts []execution.ExecuteOptionFn
	if vars != nil {
		for _, v := range *vars {
			opts = append(opts, execution.WithVariable(v.Name, v.Value))
		}
	}
	return opts
}

type variableMapper struct {
}

// Decode decodes a variable from a flag.
// E.g. --var foo=bar
// E.g. --var foo=json:{"bar":"baz"}
// E.g. --var foo=json:file:/path/to/file.json
func (vm *variableMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	t := ctx.Scan.Pop()

	strVal, ok := t.Value.(string)
	if !ok {
		return fmt.Errorf("expected string value for variable")
	}

	nameValueSplit := strings.SplitN(strVal, "=", 2)
	if len(nameValueSplit) != 2 {
		return fmt.Errorf("invalid variable format, expect foo=bar, or foo=format:file:path")
	}

	res := variable{
		Name: nameValueSplit[0],
	}

	format := "dasel"
	valueRaw := nameValueSplit[1]

	firstSplit := strings.SplitN(valueRaw, ":", 2)
	if len(firstSplit) == 2 {
		format = firstSplit[0]
		valueRaw = firstSplit[1]
	}
	if strings.HasPrefix(valueRaw, "file:") {
		filePath := valueRaw[len("file:"):]

		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer func() {
			_ = f.Close()
		}()
		contents, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read file contents: %w", err)
		}
		valueRaw = string(contents)
	}

	reader, err := parsing.Format(format).NewReader(parsing.DefaultReaderOptions())
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}
	res.Value, err = reader.Read([]byte(valueRaw))
	if err != nil {
		return fmt.Errorf("failed to read value: %w", err)
	}

	target.Elem().Set(reflect.Append(target.Elem(), reflect.ValueOf(res)))

	return nil
}
