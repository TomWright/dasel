package internal

// Document represents a KDL document.
type Document struct {
	Nodes []*Node
}

// Node represents a KDL node.
type Node struct {
	Name       string
	Type       string // type annotation, empty if none
	Arguments  []*Value
	Properties []*Property // ordered
	Children   []*Node
}

// Property represents a KDL property (key=value pair on a node).
type Property struct {
	Key   string
	Value *Value
}

// Value represents a KDL value with optional type annotation.
type Value struct {
	Type  string      // type annotation, empty if none
	Value interface{} // string, int64, float64, bool, or nil
}
