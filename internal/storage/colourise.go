package storage

import (
	"bytes"
	"github.com/alecthomas/chroma/quick"
)

// ColouriseStyle is the style used when colourising output.
const ColouriseStyle = "solarized-dark256"

// ColouriseFormatter is the formatter used when colourising output.
const ColouriseFormatter = "terminal"

// ColouriseBuffer colourises the given buffer in-place.
func ColouriseBuffer(content *bytes.Buffer, lexer string) error {
	contentString := content.String()
	content.Reset()
	return quick.Highlight(content, contentString, lexer, ColouriseFormatter, ColouriseStyle)
}

// Colourise colourises the given string.
func Colourise(content string, lexer string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	return buf, quick.Highlight(buf, content, lexer, ColouriseFormatter, ColouriseStyle)
}
