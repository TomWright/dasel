package yaml

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"go.yaml.in/yaml/v4"
)

// YAML represents the YAML file format.
const YAML parsing.Format = "yaml"

func init() {
	parsing.RegisterReader(YAML, newYAMLReader)
	parsing.RegisterWriter(YAML, newYAMLWriter)
}

type yamlValue struct {
	node  *yaml.Node
	value *model.Value
}
