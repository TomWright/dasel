package storage

// IndentOption returns a write option that sets the given indent.
func IndentOption(indent string) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionIndent,
		Value: indent,
	}
}

// PrettyPrintOption returns an option that enables or disables pretty printing.
func PrettyPrintOption(enabled bool) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionPrettyPrint,
		Value: enabled,
	}
}

// ColouriseOption returns an option that enables or disables colourised output.
func ColouriseOption(enabled bool) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionColourise,
		Value: enabled,
	}
}

// EscapeHTMLOption returns an option that enables or disables HTML escaping.
func EscapeHTMLOption(enabled bool) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionEscapeHTML,
		Value: enabled,
	}
}

// CsvCommaOption returns an option that modifies the separator character for CSV files.
func CsvCommaOption(comma rune) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionCSVComma,
		Value: comma,
	}
}

// CsvCommentOption returns an option that modifies the comment character for CSV files.
func CsvCommentOption(comma rune) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionCSVComment,
		Value: comma,
	}
}

// CsvUseCRLFOption returns an option that modifies the comment character for CSV files.
func CsvUseCRLFOption(enabled bool) ReadWriteOption {
	return ReadWriteOption{
		Key:   OptionCSVUseCRLF,
		Value: enabled,
	}
}

// OptionKey is a defined type for keys within a ReadWriteOption.
type OptionKey string

const (
	// OptionIndent is the key used with IndentOption.
	OptionIndent OptionKey = "indent"
	// OptionPrettyPrint is the key used with PrettyPrintOption.
	OptionPrettyPrint OptionKey = "prettyPrint"
	// OptionColourise is the key used with ColouriseOption.
	OptionColourise OptionKey = "colourise"
	// OptionEscapeHTML is the key used with EscapeHTMLOption.
	OptionEscapeHTML OptionKey = "escapeHtml"
	// OptionCSVComma is the key used with CsvCommaOption.
	OptionCSVComma OptionKey = "csvComma"
	// OptionCSVComment is the key used with CsvCommentOption.
	OptionCSVComment OptionKey = "csvComment"
	// OptionCSVUseCRLF is the key used with CsvUseCRLFOption.
	OptionCSVUseCRLF OptionKey = "csvUseCRLF"
)

// ReadWriteOption is an option to be used when writing.
type ReadWriteOption struct {
	Key   OptionKey
	Value interface{}
}
