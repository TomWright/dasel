package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// Security limits for XML parsing to prevent DoS attacks.
// These limits are intentionally conservative to balance usability and safety.
const (
	maxCommentLength = 10_000     // Maximum bytes per comment (10KB) - prevents memory exhaustion from single large comments
	maxTotalComments = 1_000      // Maximum comments per document - prevents abuse via comment flooding
	maxXMLSize       = 10_000_000 // Maximum XML input size (10MB) - prevents processing of excessively large files
)

func newXMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &xmlReader{
		structured: options.Ext["xml-mode"] == "structured",
	}, nil
}

type xmlReader struct {
	structured bool
}

// Read reads a value from a byte slice.
func (j *xmlReader) Read(data []byte) (*model.Value, error) {
	if len(data) > maxXMLSize {
		return nil, fmt.Errorf("XML input exceeds maximum size of %d bytes", maxXMLSize)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = true

	totalComments := 0
	el, err := j.parseElement(decoder, xml.StartElement{
		Name: xml.Name{
			Local: "root",
		},
	}, &totalComments)
	if err != nil {
		return nil, err
	}

	if j.structured {
		return el.toStructuredModel()
	}
	return el.toFriendlyModel()
}

func (e *xmlElement) toStructuredModel() (*model.Value, error) {
	attrs := model.NewMapValue()
	for _, attr := range e.Attrs {
		if err := attrs.SetMapKey(attr.Name, model.NewStringValue(attr.Value)); err != nil {
			return nil, err
		}
	}
	res := model.NewMapValue()
	if len(e.ProcessingInstructions) > 0 {
		res.SetMetadataValue("xml_processing_instructions", e.ProcessingInstructions)
	}
	if len(e.Comments) > 0 {
		res.SetMetadataValue("xml_comments", e.Comments)
	}
	if err := res.SetMapKey("name", model.NewStringValue(e.Name)); err != nil {
		return nil, err
	}
	if err := res.SetMapKey("attrs", attrs); err != nil {
		return nil, err
	}

	if err := res.SetMapKey("content", model.NewStringValue(e.Content)); err != nil {
		return nil, err
	}
	children := model.NewSliceValue()
	for _, child := range e.Children {
		childModel, err := child.toStructuredModel()
		if err != nil {
			return nil, err
		}
		if err := children.Append(childModel); err != nil {
			return nil, err
		}
	}
	if err := res.SetMapKey("children", children); err != nil {
		return nil, err
	}
	return res, nil
}

func (e *xmlElement) toFriendlyModel() (*model.Value, error) {
	if len(e.Attrs) == 0 && len(e.Children) == 0 && len(e.Comments) == 0 {
		return model.NewStringValue(e.Content), nil
	}

	res := model.NewMapValue()
	if len(e.ProcessingInstructions) > 0 {
		res.SetMetadataValue("xml_processing_instructions", e.ProcessingInstructions)
	}
	if len(e.Comments) > 0 {
		res.SetMetadataValue("xml_comments", e.Comments)
	}
	for _, attr := range e.Attrs {
		if err := res.SetMapKey("-"+attr.Name, model.NewStringValue(attr.Value)); err != nil {
			return nil, err
		}
	}

	if len(e.Content) > 0 {
		if err := res.SetMapKey("#text", model.NewStringValue(e.Content)); err != nil {
			return nil, err
		}
	}

	if len(e.Children) > 0 {
		childElementKeys := make([]string, 0)
		childElements := make(map[string][]*xmlElement)
		// childOrder records document-order child names for round-trip preservation.
		childOrder := make([]string, 0, len(e.Children))

		for _, child := range e.Children {
			if _, ok := childElements[child.Name]; !ok {
				childElementKeys = append(childElementKeys, child.Name)
			}
			childElements[child.Name] = append(childElements[child.Name], child)
			childOrder = append(childOrder, child.Name)
		}

		for _, key := range childElementKeys {
			cs := childElements[key]
			switch len(cs) {
			case 0:
				continue
			case 1:
				childModel, err := cs[0].toFriendlyModel()
				if err != nil {
					return nil, err
				}
				if err := res.SetMapKey(key, childModel); err != nil {
					return nil, err
				}
			default:
				children := model.NewSliceValue()
				for _, child := range cs {
					childModel, err := child.toFriendlyModel()
					if err != nil {
						return nil, err
					}
					if err := children.Append(childModel); err != nil {
						return nil, err
					}
				}
				if err := res.SetMapKey(key, children); err != nil {
					return nil, err
				}
			}
		}

		res.SetMetadataValue(xmlChildOrderKey, childOrder)
	}

	return res, nil
}

func (j *xmlReader) parseElement(decoder *xml.Decoder, element xml.StartElement, totalComments *int) (*xmlElement, error) {
	el := &xmlElement{
		Name:                   element.Name.Local,
		Attrs:                  make([]xmlAttr, 0),
		Children:               make([]*xmlElement, 0),
		ProcessingInstructions: make([]*xmlProcessingInstruction, 0),
		Comments:               make([]*xmlComment, 0),
	}

	for _, attr := range element.Attr {
		el.Attrs = append(el.Attrs, xmlAttr{
			Name:  attr.Name.Local,
			Value: attr.Value,
		})
	}

	var processingInstructions []*xmlProcessingInstruction
	var comments []*xmlComment

	for {
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			if el.Name == "root" {
				return el, nil
			}
			return nil, fmt.Errorf("unexpected EOF")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read token: %w", err)
		}
		if t == nil {
			if el.Name == "root" {
				return el, nil
			}
			return nil, fmt.Errorf("unexpected nil token")
		}

		switch t := t.(type) {
		case xml.StartElement:
			child, err := j.parseElement(decoder, t, totalComments)
			if err != nil {
				return nil, err
			}
			if len(comments) > 0 {
				child.Comments = append(comments, child.Comments...)
				comments = nil
			}
			if len(processingInstructions) > 0 {
				child.ProcessingInstructions = processingInstructions
				processingInstructions = nil
			}
			el.Children = append(el.Children, child)
		case xml.CharData:
			stringContent := string(t)
			if strings.TrimSpace(stringContent) == "" {
				continue
			}
			el.Content += stringContent
		case xml.EndElement:
			if len(comments) > 0 {
				el.Comments = append(el.Comments, comments...)
			}
			return el, nil
		case xml.Comment:
			commentText := string(t)
			if len(commentText) > maxCommentLength {
				return nil, fmt.Errorf("comment exceeds maximum length of %d bytes", maxCommentLength)
			}
			if *totalComments >= maxTotalComments {
				return nil, fmt.Errorf("document exceeds maximum comment count of %d", maxTotalComments)
			}
			comment := &xmlComment{
				Text: commentText,
			}
			comments = append(comments, comment)
			*totalComments++
			continue
		case xml.ProcInst:
			pi := &xmlProcessingInstruction{
				Target: t.Target,
				Value:  string(t.Inst),
			}
			processingInstructions = append(processingInstructions, pi)
			continue
		case xml.Directive:
			continue
		case xml.Attr:
			continue
		default:
			return nil, fmt.Errorf("unexpected token: %v", t)
		}
	}
}
