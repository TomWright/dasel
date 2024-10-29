package cli

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/parsing"
)

type extReadWriteFlag struct {
	Name  string
	Value string
}

type extReadWriteFlags *[]extReadWriteFlag

func applyReaderFlags(readerOptions *parsing.ReaderOptions, f extReadWriteFlags) {
	if f != nil {
		for _, flag := range *f {
			readerOptions.Ext[flag.Name] = flag.Value
		}
	}
}

func applyWriterFlags(writerOptions *parsing.WriterOptions, f extReadWriteFlags) {
	if f != nil {
		for _, flag := range *f {
			writerOptions.Ext[flag.Name] = flag.Value
		}
	}
}

type extReadWriteFlagMapper struct {
}

func (vm *extReadWriteFlagMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	t := ctx.Scan.Pop()

	strVal, ok := t.Value.(string)
	if !ok {
		return fmt.Errorf("expected string value for variable")
	}

	nameValueSplit := strings.SplitN(strVal, "=", 2)
	if len(nameValueSplit) != 2 {
		return fmt.Errorf("invalid read/write flag format, expect foo=bar")
	}

	res := extReadWriteFlag{
		Name:  nameValueSplit[0],
		Value: nameValueSplit[1],
	}

	target.Elem().Set(reflect.Append(target.Elem(), reflect.ValueOf(res)))

	return nil
}
