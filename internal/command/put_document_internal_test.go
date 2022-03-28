package command

import (
	"errors"
	"github.com/tomwright/dasel/storage"
	"strings"
	"testing"
)

func TestPut_Document(t *testing.T) {
	t.Run("SingleFailingWriter", func(t *testing.T) {
		err := runPutDocumentCommand(putDocumentOpts{
			Parser:         "json",
			Selector:       ".[0]",
			Reader:         strings.NewReader(`[{"name": "Tom"}]`),
			DocumentString: `{"name": "Frank"}`,
			Writer:         &failingWriter{},
		}, nil)

		if err == nil || !errors.Is(err, errFailingWriterErr) {
			t.Errorf("expected error %v, got %v", errFailingWriterErr, err)
			return
		}
	})
	t.Run("MultiFailingWriter", func(t *testing.T) {
		err := runPutDocumentCommand(putDocumentOpts{
			Parser:         "json",
			Selector:       ".[*]",
			Reader:         strings.NewReader(`[{"name": "Tom"}]`),
			DocumentString: `{"name": "Frank"}`,
			Writer:         &failingWriter{},
			Multi:          true,
		}, nil)

		if err == nil || !errors.Is(err, errFailingWriterErr) {
			t.Errorf("expected error %v, got %v", errFailingWriterErr, err)
			return
		}
	})
	t.Run("InvalidDocumentParser", func(t *testing.T) {
		err := runPutDocumentCommand(putDocumentOpts{
			Parser:         "json",
			Selector:       ".[*]",
			Reader:         strings.NewReader(`[{"name": "Tom"}]`),
			DocumentString: `{"name": "Frank"}`,
			DocumentParser: "bad",
		}, nil)

		exp := &storage.UnknownParserErr{Parser: "bad"}

		if err == nil || !strings.HasSuffix(err.Error(), exp.Error()) {
			t.Errorf("expected error %v, got %v", exp, err)
			return
		}
	})
	t.Run("InvalidDocument", func(t *testing.T) {
		err := runPutDocumentCommand(putDocumentOpts{
			Parser:         "json",
			Selector:       ".[*]",
			Reader:         strings.NewReader(`[{"name": "Tom"}]`),
			DocumentString: `{"name": "Frank}`,
			DocumentParser: "json",
		}, nil)

		exp := "could not parse document: could not unmarshal data: unexpected EOF"

		if err == nil || err.Error() != exp {
			t.Errorf("expected error %v, got %v", exp, err)
			return
		}
	})
}
