package hclgo

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// HCLFileToGo takes a HCL file and converts it to a map[string]interface{}.
func HCLFileToGo(file *hcl.File) (map[string]interface{}, error) {
	target := map[string]interface{}{}

	diags := gohcl.DecodeBody(file.Body, nil, &target)
	if diags.HasErrors() {
		return nil, fmt.Errorf("could not decode hcl body: %s", diags.Error())
	}

	res, err := HCLToGo(target)
	if err != nil {
		return nil, err
	}
	return res.(map[string]interface{}), nil
}

// HCLToGo takes HCL data types and converts them to go data types.
func HCLToGo(data interface{}) (interface{}, error) {
	switch val := data.(type) {
	case map[string]interface{}:
		var err error
		for k, v := range val {
			val[k], err = HCLToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return val, nil

	case []interface{}:
		var err error
		for k, v := range val {
			val[k], err = HCLToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return val, nil

	case *hcl.Attribute:
		x, _ := val.Expr.Value(nil)
		switch x.Type() {
		case cty.Bool:
			return x.True(), nil
		case cty.Number:
			floatVal, _ := x.AsBigFloat().Float64()
			return floatVal, nil
		case cty.String:
			return x.AsString(), nil
		default:
			return nil, fmt.Errorf("unhandled hcl attribute type [%s]: %s", val.Name, x.Type().GoString())
		}
	}
	return data, nil
}

// GoToHCLFile takes an interface{} converts it to a HCL file.
func GoToHCLFile(data interface{}) (*hclwrite.File, error) {
	data, err := GoToHCL(data)
	if err != nil {
		return nil, err
	}

	file := hclwrite.NewEmptyFile()

	// todo : catch panics

	gohcl.EncodeIntoBody(data, file.Body())

	return file, nil
}

// GoToHCL takes go data types data types and converts them to HCL data types.
func GoToHCL(data interface{}) (interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
