package xml

import (
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// XML represents the XML file format.
	XML parsing.Format = "xml"
)

// xmlChildOrderKey is the metadata key for preserving child element
// document order during XML round-trips. Value type: []string.
const xmlChildOrderKey = "xml_child_order"

var _ parsing.Reader = (*xmlReader)(nil)
var _ parsing.Writer = (*xmlWriter)(nil)

func init() {
	parsing.RegisterReader(XML, newXMLReader)
	parsing.RegisterWriter(XML, newXMLWriter)
}

type xmlAttr struct {
	Name  string
	Value string
}

type xmlProcessingInstruction struct {
	Target string
	Value  string
}

type xmlComment struct {
	Text string
}

type xmlElement struct {
	Name                   string
	Attrs                  []xmlAttr
	Children               []*xmlElement
	Content                string
	ProcessingInstructions []*xmlProcessingInstruction
	Comments               []*xmlComment
	useChildrenOnly        bool
	depth                  int // Tracks nesting depth for proper indentation
}

// appendChild appends child's children (if useChildrenOnly) or the child itself.
func (el *xmlElement) appendChild(child *xmlElement) {
	if child.useChildrenOnly {
		el.Children = append(el.Children, child.Children...)
	} else {
		el.Children = append(el.Children, child)
	}
}
