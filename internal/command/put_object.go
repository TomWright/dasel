package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/oflag"
	"strings"
)

func putObjectCommand() *cobra.Command {
	typeList := oflag.NewStringList()

	cmd := &cobra.Command{
		Use:   "object -f <file> -s <selector> <value>",
		Short: "Update a string property in the given file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flag("file").Value.String()
			outFlag := cmd.Flag("out").Value.String()
			parserFlag := cmd.Flag("parser").Value.String()
			selectorFlag := cmd.Flag("selector").Value.String()

			parser, err := getParser(fileFlag, parserFlag)
			if err != nil {
				return err
			}
			rootNode, err := getRootNode(fileFlag, parser)
			if err != nil {
				return err
			}

			if len(typeList.Strings) != len(args) {
				return fmt.Errorf("exactly %d types are required, got %d", len(args), len(typeList.Strings))
			}

			updateValue := map[string]interface{}{}

			for k, arg := range args {
				splitArg := strings.Split(arg, "=")
				name := splitArg[0]
				value := strings.Join(splitArg[1:], "=")
				parsedValue, err := parseValue(value, typeList.Strings[k])
				if err != nil {
					return fmt.Errorf("could not parse value [%s]: %w", name, err)
				}
				updateValue[name] = parsedValue
			}

			if err := rootNode.Put(selectorFlag, updateValue); err != nil {
				return fmt.Errorf("could not put value: %w", err)
			}

			if err := writeNodeToOutput(rootNode, parser, fileFlag, outFlag); err != nil {
				return fmt.Errorf("could not write output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().VarP(typeList, "type", "t", "Types of the variables in the object.")
	if err := cmd.MarkFlagRequired("type"); err != nil {
		panic("could not mark flag as required: type")
	}

	return cmd
}
