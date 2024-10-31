package xml

import (
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// XML represents the XML file format.
	XML parsing.Format = "xml"
)

var _ parsing.Reader = (*xmlReader)(nil)
var _ parsing.Writer = (*xmlWriter)(nil)

func init() {
	parsing.RegisterReader(XML, newXMLReader)
	// XML writer is not implemented yet
	//parsing.RegisterWriter(XML, newXMLWriter)
}

type xmlAttr struct {
	Name  string
	Value string
}

type xmlElement struct {
	Name     string
	Attrs    []xmlAttr
	Children []*xmlElement
	Content  string
}
