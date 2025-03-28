package hcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/zclconf/go-cty/cty"
)

func newHCLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &hclReader{
		alwaysReadLabelsToSlice: options.Ext["hcl-block-format"] == "array",
	}, nil
}

type hclReader struct {
	alwaysReadLabelsToSlice bool
}

// Read reads a value from a byte slice.
// Reads the HCL data into a model that follows the HCL JSON spec.
// See https://github.com/hashicorp/hcl/blob/main/json%2Fspec.md
func (r *hclReader) Read(data []byte) (*model.Value, error) {
	f, _ := hclsyntax.ParseConfig(data, "input", hcl.InitialPos)

	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("failed to assert file body type")
	}

	return r.decodeHCLBody(body)
}

func (r *hclReader) decodeHCLBody(body *hclsyntax.Body) (*model.Value, error) {
	res := model.NewMapValue()
	var err error

	for _, attr := range body.Attributes {
		val, err := r.decodeHCLExpr(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode attr %q: %w", attr.Name, err)
		}

		if err := res.SetMapKey(attr.Name, val); err != nil {
			return nil, err
		}
	}

	res, err = r.decodeHCLBodyBlocks(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *hclReader) decodeHCLBodyBlocks(body *hclsyntax.Body, res *model.Value) (*model.Value, error) {
	for _, block := range body.Blocks {
		if err := r.decodeHCLBlock(block, res); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *hclReader) decodeHCLBlock(block *hclsyntax.Block, res *model.Value) error {
	key := block.Type
	v := res
	for _, label := range block.Labels {
		exists, err := v.MapKeyExists(key)
		if err != nil {
			return err
		}

		if exists {
			keyV, err := v.GetMapKey(key)
			if err != nil {
				return err
			}
			v = keyV
		} else {
			keyV := model.NewMapValue()
			if err := v.SetMapKey(key, keyV); err != nil {
				return err
			}
			v = keyV
		}

		key = label
	}

	body, err := r.decodeHCLBody(block.Body)
	if err != nil {
		return err
	}

	exists, err := v.MapKeyExists(key)
	if err != nil {
		return err
	}
	if exists {
		keyV, err := v.GetMapKey(key)
		if err != nil {
			return err
		}

		switch keyV.Type() {
		case model.TypeSlice:
			if err := keyV.Append(body); err != nil {
				return err
			}
		case model.TypeMap:
			// Previous value was a map.
			// Create a new slice containing the previous map and the new map.
			newKeyV := model.NewSliceValue()
			previousKeyV, err := keyV.Copy()
			if err != nil {
				return err
			}
			if err := newKeyV.Append(previousKeyV); err != nil {
				return err
			}
			if err := newKeyV.Append(body); err != nil {
				return err
			}
			if err := keyV.Set(newKeyV); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected type: %s", keyV.Type())
		}
	} else {
		if r.alwaysReadLabelsToSlice {
			slice := model.NewSliceValue()
			if err := slice.Append(body); err != nil {
				return err
			}
			if err := v.SetMapKey(key, slice); err != nil {
				return err
			}
		} else {
			if err := v.SetMapKey(key, body); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *hclReader) decodeHCLExpr(expr hcl.Expression) (*model.Value, error) {
	source := cty.Value{}
	_ = gohcl.DecodeExpression(expr, nil, &source)

	return r.decodeCtyValue(source)
}

func (r *hclReader) decodeCtyValue(source cty.Value) (res *model.Value, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("failed to decode: %v", r)
			return
		}
	}()
	if source.IsNull() {
		return model.NewNullValue(), nil
	}

	sourceT := source.Type()
	if sourceT.HasDynamicTypes() {
		// TODO : Handle DynamicPseudoType.
		// I haben't found a clear way to do this.
		return model.NewNullValue(), nil
	}
	switch {
	case sourceT.IsListType(), sourceT.IsTupleType():
		res = model.NewSliceValue()
		it := source.ElementIterator()
		for it.Next() {
			k, v := it.Element()
			// We don't need the index as they should be in order.
			// Just validates the key is correct.
			_, _ = k.AsBigFloat().Float64()

			val, err := r.decodeCtyValue(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode tuple value: %w", err)
			}

			if err := res.Append(val); err != nil {
				return nil, err
			}
		}
		return res, nil
	case sourceT.IsMapType(), sourceT.IsObjectType(), sourceT.IsSetType():
		v := model.NewMapValue()
		it := source.ElementIterator()
		for it.Next() {
			k, el := it.Element()
			if k.Type() != cty.String {
				return nil, fmt.Errorf("object key must be a string")
			}
			kStr := k.AsString()

			elVal, err := r.decodeCtyValue(el)
			if err != nil {
				return nil, fmt.Errorf("failed to decode object value: %w", err)
			}

			if err := v.SetMapKey(kStr, elVal); err != nil {
				return nil, err
			}
		}
		return v, nil
	case sourceT.IsPrimitiveType():
		switch sourceT {
		case cty.String:
			v := source.AsString()
			return model.NewStringValue(v), nil
		case cty.Bool:
			v := source.True()
			return model.NewBoolValue(v), nil
		case cty.Number:
			v := source.AsBigFloat()
			f64, _ := v.Float64()
			if v.IsInt() {
				return model.NewIntValue(int64(f64)), nil
			}
			return model.NewFloatValue(f64), nil
		default:
			return nil, fmt.Errorf("unhandled primitive type %q", source.Type())
		}
	default:
		return nil, fmt.Errorf("unsupported type: %s", sourceT.FriendlyName())
	}
}
