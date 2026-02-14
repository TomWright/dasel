package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newXMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &xmlWriter{
		options: options,
	}, nil
}

type xmlWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *xmlWriter) Write(value *model.Value) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := xml.NewEncoder(buf)
	defer func() {
		_ = writer.Close()
	}()
	writer.Indent("", "  ")

	element, err := j.toElement("root", value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to element: %w", err)
	}
	for _, c := range element.Children {
		if err := writer.Encode(c); err != nil {
			return nil, err
		}
		if err := writer.EncodeToken(xml.CharData("\n")); err != nil {
			return nil, err
		}
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}
	outBytes := buf.Bytes()
	if !bytes.HasSuffix(outBytes, []byte("\n")) {
		outBytes = append(outBytes, '\n')
	}
	return outBytes, nil
}

func (j *xmlWriter) toElement(key string, value *model.Value) (*xmlElement, error) {
	readProcessingInstructions := func() []*xmlProcessingInstruction {
		if piMeta, ok := value.MetadataValue("xml_processing_instructions"); ok && piMeta != nil {
			if pis, ok := piMeta.([]*xmlProcessingInstruction); ok {
				return pis
			}
		}
		return nil
	}
	readComments := func() []*xmlComment {
		if commentMeta, ok := value.MetadataValue("xml_comments"); ok && commentMeta != nil {
			if comments, ok := commentMeta.([]*xmlComment); ok {
				return comments
			}
		}
		return nil
	}
	switch value.Type() {

	case model.TypeString:
		strVal, err := valueToString(value)
		return &xmlElement{
			Name:                   key,
			Content:                strVal,
			ProcessingInstructions: readProcessingInstructions(),
			Comments:               readComments(),
		}, err

	case model.TypeMap:
		kvs, err := value.MapKeyValues()
		if err != nil {
			return nil, err
		}

		el := &xmlElement{
			Name:                   key,
			ProcessingInstructions: readProcessingInstructions(),
			Comments:               readComments(),
		}

		if err := extractAttrsAndText(kvs, el); err != nil {
			return nil, err
		}

		orderMeta, hasOrder := value.MetadataValue(xmlChildOrderKey)
		if hasOrder {
			childOrder, ok := orderMeta.([]string)
			if !ok {
				hasOrder = false
			} else {
				if err := j.buildChildrenOrdered(kvs, el, childOrder); err != nil {
					return nil, err
				}
			}
		}

		if !hasOrder {
			if err := j.buildChildrenUnordered(kvs, el); err != nil {
				return nil, err
			}
		}

		return el, nil
	case model.TypeSlice:
		el := &xmlElement{
			Name:                   "root",
			ProcessingInstructions: readProcessingInstructions(),
			Comments:               readComments(),
			useChildrenOnly:        true,
		}
		if err := value.RangeSlice(func(i int, value *model.Value) error {
			childEl, err := j.toElement(key, value)
			if err != nil {
				return err
			}
			el.appendChild(childEl)

			return nil
		}); err != nil {
			return nil, err
		}
		return el, nil
	default:
		return nil, fmt.Errorf("xml writer does not support value type: %s", value.Type())
	}
}

// extractAttrsAndText iterates kvs and extracts "-" prefixed attributes into
// el.Attrs and "#text" into el.Content.
func extractAttrsAndText(kvs []model.KeyValue, el *xmlElement) error {
	for _, kv := range kvs {
		if strings.HasPrefix(kv.Key, "-") {
			attr := xmlAttr{
				Name: kv.Key[1:],
			}
			var err error
			attr.Value, err = valueToString(kv.Value)
			if err != nil {
				return fmt.Errorf("failed to convert attribute %q to string: %w", attr.Name, err)
			}
			el.Attrs = append(el.Attrs, attr)
			continue
		}

		if kv.Key == "#text" {
			var err error
			el.Content, err = valueToString(kv.Value)
			if err != nil {
				return fmt.Errorf("failed to convert content to string: %w", err)
			}
			continue
		}
	}
	return nil
}

// buildChildrenOrdered reconstructs child elements using counter-based ordering
// from the childOrder metadata slice.
func (j *xmlWriter) buildChildrenOrdered(kvs []model.KeyValue, el *xmlElement, childOrder []string) error {
	// Build local map for fast lookups without GetMapKey overhead.
	childValues := make(map[string]*model.Value, len(kvs))
	for _, kv := range kvs {
		if !strings.HasPrefix(kv.Key, "-") && kv.Key != "#text" {
			childValues[kv.Key] = kv.Value
		}
	}

	counters := make(map[string]int, len(childValues))
	seen := make(map[string]bool, len(childValues))

	for _, name := range childOrder {
		seen[name] = true
		childVal, exists := childValues[name]
		if !exists {
			// Key not in map (stale metadata after delete), skip.
			counters[name]++
			continue
		}

		index := counters[name]
		counters[name]++

		if childVal.Type() == model.TypeSlice {
			// SliceLen error is impossible after TypeSlice type check above.
			sliceLen, _ := childVal.SliceLen()
			if index >= sliceLen {
				// Counter overflow (more metadata entries than actual values), skip.
				continue
			}
			item, sliceErr := childVal.GetSliceIndex(index)
			if sliceErr != nil {
				// Should not happen after bounds check; skip gracefully if model is inconsistent.
				continue
			}
			childEl, childErr := j.toElement(name, item)
			if childErr != nil {
				return fmt.Errorf("failed to convert child element %q to element: %w", name, childErr)
			}
			el.appendChild(childEl)
		} else {
			if index >= 1 {
				// Scalar value, can only be emitted once.
				continue
			}
			childEl, childErr := j.toElement(name, childVal)
			if childErr != nil {
				return fmt.Errorf("failed to convert child element %q to element: %w", name, childErr)
			}
			el.appendChild(childEl)
		}
	}

	// Append any map keys not in the ordering (new keys from mutations).
	for _, kv := range kvs {
		if strings.HasPrefix(kv.Key, "-") || kv.Key == "#text" || seen[kv.Key] {
			continue
		}
		childEl, childErr := j.toElement(kv.Key, kv.Value)
		if childErr != nil {
			return fmt.Errorf("failed to convert child element %q to element: %w", kv.Key, childErr)
		}
		el.appendChild(childEl)
	}

	return nil
}

// buildChildrenUnordered iterates map keys in insertion order, skipping
// attributes and #text (backward-compatible fallback).
func (j *xmlWriter) buildChildrenUnordered(kvs []model.KeyValue, el *xmlElement) error {
	for _, kv := range kvs {
		if strings.HasPrefix(kv.Key, "-") || kv.Key == "#text" {
			continue
		}
		childEl, childErr := j.toElement(kv.Key, kv.Value)
		if childErr != nil {
			return fmt.Errorf("failed to convert child element %q to element: %w", kv.Key, childErr)
		}
		el.appendChild(childEl)
	}
	return nil
}

func valueToString(v *model.Value) (string, error) {
	if v.IsNull() {
		return "", nil
	}

	switch v.Type() {
	case model.TypeString:
		stringValue, err := v.StringValue()
		if err != nil {
			return "", err
		}
		return stringValue, nil
	case model.TypeInt:
		i, err := v.IntValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case model.TypeFloat:
		i, err := v.FloatValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%g", i), nil
	case model.TypeBool:
		i, err := v.BoolValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%t", i), nil
	default:
		return "", fmt.Errorf("xml writer cannot format type %s to string", v.Type())
	}
}

// indentString returns the indentation for a given depth level.
func indentString(depth int) string {
	return strings.Repeat("  ", depth)
}

func (e *xmlElement) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// Write processing instructions before the element (document-level)
	if len(e.ProcessingInstructions) > 0 {
		for _, pi := range e.ProcessingInstructions {
			if err := enc.EncodeToken(xml.ProcInst{
				Target: pi.Target,
				Inst:   []byte(pi.Value),
			}); err != nil {
				return err
			}
			if err := enc.EncodeToken(xml.CharData("\n")); err != nil {
				return err
			}
		}
	}
	// Write comments before the element (document-level comments only)
	// Child-level comments are written by the parent inside its body
	if e.depth == 0 && len(e.Comments) > 0 {
		for _, comment := range e.Comments {
			if strings.Contains(comment.Text, "--") {
				return fmt.Errorf("comment text cannot contain '--' sequence (invalid XML comment)")
			}
			if err := enc.EncodeToken(xml.Comment(comment.Text)); err != nil {
				return fmt.Errorf("failed to encode comment: %w", err)
			}
			if err := enc.EncodeToken(xml.CharData("\n")); err != nil {
				return err
			}
		}
	}
	start.Name = xml.Name{Local: e.Name}

	if len(e.Attrs) > 0 {
		for _, attr := range e.Attrs {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  xml.Name{Local: attr.Name},
				Value: attr.Value,
			})
		}
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	// TODO : Handle CDATA sections on write.

	if len(e.Content) > 0 {
		if err := enc.EncodeToken(xml.CharData(e.Content)); err != nil {
			return err
		}
	}

	// Write children with their preceding comments
	childDepth := e.depth + 1
	for _, child := range e.Children {
		// Write child's comments inside parent, before the child element
		if len(child.Comments) > 0 {
			for _, comment := range child.Comments {
				if strings.Contains(comment.Text, "--") {
					return fmt.Errorf("comment text cannot contain '--' sequence (invalid XML comment)")
				}
				// Add newline + indentation before comment
				if err := enc.EncodeToken(xml.CharData("\n" + indentString(childDepth))); err != nil {
					return err
				}
				if err := enc.EncodeToken(xml.Comment(comment.Text)); err != nil {
					return fmt.Errorf("failed to encode comment: %w", err)
				}
			}
			// Clear comments so child doesn't write them again
			child.Comments = nil
		}
		// Set child depth for recursive calls
		child.depth = childDepth
		if err := enc.Encode(child); err != nil {
			return err
		}
	}

	return enc.EncodeToken(start.End())
}
