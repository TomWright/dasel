package hcl

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/zclconf/go-cty/cty"
)

func newHCLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &hclWriter{}, nil
}

type hclWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *hclWriter) Write(value *model.Value) ([]byte, error) {
	f, err := j.valueToFile(value)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := f.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *hclWriter) valueToFile(v *model.Value) (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()

	body := f.Body()

	if err := j.addValueToBody(nil, v, body); err != nil {
		return nil, err
	}

	return f, nil
}

func (j *hclWriter) addValueToBody(previousLabels []string, v *model.Value, body *hclwrite.Body) error {
	if !v.IsMap() {
		return fmt.Errorf("hcl body is expected to be a map, got %s", v.Type())
	}

	kvs, err := v.MapKeyValues()
	if err != nil {
		return err
	}

	blocks := make([]*hclwrite.Block, 0)
	for _, kv := range kvs {
		switch kv.Value.Type() {
		case model.TypeMap:
			block, err := j.valueToBlock(kv.Key, previousLabels, kv.Value)
			if err != nil {
				return fmt.Errorf("failed to encode %q to hcl block: %w", kv.Key, err)
			}
			blocks = append(blocks, block)
		case model.TypeSlice:
			vals := make([]cty.Value, 0)

			allMaps := true

			if err := kv.Value.RangeSlice(func(_ int, value *model.Value) error {
				ctyVal, err := j.valueToCty(value)
				if err != nil {
					return err
				}
				vals = append(vals, ctyVal)

				if !value.IsMap() {
					allMaps = false
				}
				return nil
			}); err != nil {
				return err
			}

			if allMaps {
				if err := kv.Value.RangeSlice(func(_ int, value *model.Value) error {
					block, err := j.valueToBlock(kv.Key, previousLabels, value)
					if err != nil {
						return fmt.Errorf("failed to encode %q to hcl block: %w", kv.Key, err)
					}
					blocks = append(blocks, block)
					return nil
				}); err != nil {
					return err
				}
			} else {
				body.SetAttributeValue(kv.Key, cty.TupleVal(vals))
			}

		default:
			ctyVal, err := j.valueToCty(kv.Value)
			if err != nil {
				return fmt.Errorf("failed to encode attribute %q: %w", kv.Key, err)
			}
			body.SetAttributeValue(kv.Key, ctyVal)
		}
	}

	for _, block := range blocks {
		body.AppendBlock(block)
	}

	return nil
}

func (j *hclWriter) valueToCty(v *model.Value) (cty.Value, error) {
	switch v.Type() {
	case model.TypeString:
		val, err := v.StringValue()
		if err != nil {
			return cty.Value{}, err
		}
		return cty.StringVal(val), nil
	case model.TypeBool:
		val, err := v.BoolValue()
		if err != nil {
			return cty.Value{}, err
		}
		return cty.BoolVal(val), nil
	case model.TypeInt:
		val, err := v.IntValue()
		if err != nil {
			return cty.Value{}, err
		}
		return cty.NumberIntVal(val), nil
	case model.TypeFloat:
		val, err := v.FloatValue()
		if err != nil {
			return cty.Value{}, err
		}
		return cty.NumberFloatVal(val), nil
	case model.TypeNull:
		return cty.NullVal(cty.NilType), nil
	case model.TypeSlice:
		var vals []cty.Value
		if err := v.RangeSlice(func(_ int, value *model.Value) error {
			ctyVal, err := j.valueToCty(value)
			if err != nil {
				return err
			}
			vals = append(vals, ctyVal)
			return nil
		}); err != nil {
			return cty.Value{}, err
		}
		return cty.TupleVal(vals), nil
	case model.TypeMap:
		mapV := map[string]cty.Value{}
		if err := v.RangeMap(func(s string, value *model.Value) error {
			ctyVal, err := j.valueToCty(value)
			if err != nil {
				return err
			}
			mapV[s] = ctyVal
			return nil
		}); err != nil {
			return cty.Value{}, err
		}
		return cty.ObjectVal(mapV), nil
	default:
		return cty.Value{}, fmt.Errorf("unhandled type when converting to cty value %q", v.Type())
	}
}

func (j *hclWriter) valueToBlock(key string, labels []string, v *model.Value) (*hclwrite.Block, error) {
	if !v.IsMap() {
		return nil, fmt.Errorf("must be map")
	}

	b := hclwrite.NewBlock(key, labels)

	if err := j.addValueToBody(labels, v, b.Body()); err != nil {
		return nil, err
	}

	return b, nil
}
