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

type OptionKey string

const (
	// OptionIndent is the key used with IndentOption.
	OptionIndent      OptionKey = "indent"
	OptionPrettyPrint OptionKey = "prettyPrint"
)

// ReadWriteOption is an option to be used when writing.
type ReadWriteOption struct {
	Key   OptionKey
	Value interface{}
}
