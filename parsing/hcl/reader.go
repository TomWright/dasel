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
	return &hclReader{}, nil
}

type hclReader struct{}

// Read reads a value from a byte slice.
func (j *hclReader) Read(data []byte) (*model.Value, error) {
	f, _ := hclsyntax.ParseConfig(data, "input", hcl.InitialPos)

	body, ok := f.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("failed to assert file body type")
	}

	return decodeHCLBody(body)
}

func decodeHCLBody(body *hclsyntax.Body) (*model.Value, error) {
	res := model.NewMapValue()

	for _, attr := range body.Attributes {
		val, err := decodeHCLExpr(attr.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode attr %q: %w", attr.Name, err)
		}

		if err := res.SetMapKey(attr.Name, val); err != nil {
			return nil, err
		}
	}

	blockTypeIndexes := make(map[string]int)
	blockValues := make([][]*model.Value, 0)
	for _, block := range body.Blocks {
		if _, ok := blockTypeIndexes[block.Type]; !ok {
			blockValues = append(blockValues, make([]*model.Value, 0))
			blockTypeIndexes[block.Type] = len(blockValues) - 1
		}
		res, err := decodeHCLBlock(block)
		if err != nil {
			return nil, fmt.Errorf("failed to decode block %q: %w", block.Type, err)
		}
		blockValues[blockTypeIndexes[block.Type]] = append(blockValues[blockTypeIndexes[block.Type]], res)
	}

	for t, index := range blockTypeIndexes {
		blocks := blockValues[index]
		switch len(blocks) {
		case 0:
			continue
		case 1:
			if err := res.SetMapKey(t, blocks[0]); err != nil {
				return nil, err
			}
		default:
			val := model.NewSliceValue()
			for _, b := range blocks {
				if err := val.Append(b); err != nil {
					return nil, err
				}
			}
			if err := res.SetMapKey(t, val); err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}

func decodeHCLBlock(block *hclsyntax.Block) (*model.Value, error) {
	res, err := decodeHCLBody(block.Body)
	if err != nil {
		return nil, err
	}

	labels := model.NewSliceValue()
	for _, l := range block.Labels {
		if err := labels.Append(model.NewStringValue(l)); err != nil {
			return nil, err
		}
	}

	if err := res.SetMapKey("labels", labels); err != nil {
		return nil, err
	}

	if err := res.SetMapKey("type", model.NewStringValue(block.Type)); err != nil {
		return nil, err
	}

	return res, nil
}

func decodeHCLExpr(expr hcl.Expression) (*model.Value, error) {
	source := cty.Value{}
	_ = gohcl.DecodeExpression(expr, nil, &source)

	return decodeCtyValue(source)
}

func decodeCtyValue(source cty.Value) (res *model.Value, err error) {
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
	switch {
	case sourceT.IsMapType():
		return nil, fmt.Errorf("map type not implemented")
	case sourceT.IsListType():
		return nil, fmt.Errorf("list type not implemented")
	case sourceT.IsCollectionType():
		return nil, fmt.Errorf("collection type not implemented")
	case sourceT.IsCapsuleType():
		return nil, fmt.Errorf("capsule type not implemented")
	case sourceT.IsTupleType():
		res = model.NewSliceValue()
		it := source.ElementIterator()
		for it.Next() {
			k, v := it.Element()
			// We don't need the index as they should be in order.
			// Just validates the key is correct.
			_, _ = k.AsBigFloat().Float64()

			val, err := decodeCtyValue(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode tuple value: %w", err)
			}

			if err := res.Append(val); err != nil {
				return nil, err
			}
		}
		return res, nil
	case sourceT.IsObjectType():
		return nil, fmt.Errorf("object type not implemented")
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
	case sourceT.IsSetType():
		return nil, fmt.Errorf("set type not implemented")
	default:
		return nil, fmt.Errorf("unhandled type: %s", sourceT.FriendlyName())
	}
}
