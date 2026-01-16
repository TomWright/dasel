package json

import (
	// "encoding/json"
	json "github.com/goccy/go-json"
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// JSON represents the JSON file format.
	JSON parsing.Format = "json"

	jsonOpenObject  = json.Delim('{')
	jsonCloseObject = json.Delim('}')
	jsonOpenArray   = json.Delim('[')
	jsonCloseArray  = json.Delim(']')
)

func init() {
	parsing.RegisterReader(JSON, newJSONReader)
	parsing.RegisterWriter(JSON, newJSONWriter)
}
