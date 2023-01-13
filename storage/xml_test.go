package storage_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
	"io"
	"reflect"
	"testing"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var xmlBytes = []byte(`<user>
  <name>Tom</name>
</user>
`)
var xmlMap = map[string]interface{}{
	"user": map[string]interface{}{
		"name": "Tom",
	},
}
var encodedXmlMap = map[string]interface{}{
	"user": map[string]interface{}{
		"name": "Tõm",
	},
}

func TestXMLParser_FromBytes(t *testing.T) {
	got, err := (&storage.XMLParser{}).FromBytes(xmlBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(xmlMap, got.Interface()) {
		t.Errorf("expected %v, got %v", xmlMap, got)
	}
}

func TestXMLParser_FromBytes_Empty(t *testing.T) {
	got, err := (&storage.XMLParser{}).FromBytes([]byte{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !got.IsEmpty() {
		t.Errorf("expected %v, got %v", nil, got)
	}
}

func TestXMLParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.XMLParser{}).FromBytes(nil)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.XMLParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestXMLParser_ToBytes_Default(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf(xmlMap))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(xmlBytes, got) {
		t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
	}
}
func TestXMLParser_ToBytes_SingleDocument(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf(xmlMap).WithMetadata("isSingleDocument", true))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(xmlBytes, got) {
		t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
	}
}
func TestXMLParser_ToBytes_SingleDocument_Colourise(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf(xmlMap).WithMetadata("isSingleDocument", true), storage.ColouriseOption(true))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	expBuf, _ := storage.Colourise(string(xmlBytes), "xml")
	exp := expBuf.Bytes()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestXMLParser_ToBytes_MultiDocument(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf([]interface{}{xmlMap, xmlMap}).WithMetadata("isMultiDocument", true))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := append([]byte{}, xmlBytes...)
	exp = append(exp, xmlBytes...)
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", string(exp), string(got))
	}
}
func TestXMLParser_ToBytes_DefaultValue(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf("asd"))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := []byte(`asd
`)
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", string(exp), string(got))
	}
}
func TestXMLParser_ToBytes_SingleDocumentValue(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf("asd"))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := []byte(`asd
`)
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", string(exp), string(got))
	}
}
func TestXMLParser_ToBytes_MultiDocumentValue(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(dasel.ValueOf([]interface{}{"asd", "123"}).WithMetadata("isMultiDocument", true))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := []byte(`asd
123
`)
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", string(exp), string(got))
	}
}
func TestXMLParser_ToBytes_Entities(t *testing.T) {
	bytes := []byte(`<systemList>
  <system>
    <command>sudo /home/fozz/RetroPie-Setup/retropie_packages.sh retropiemenu launch %ROM% &lt;/dev/tty &gt;/dev/tty</command>
    <extension>.rp .sh</extension>
    <fullname>RetroPie</fullname>
    <name>retropie</name>
    <path>/home/fozz/RetroPie/retropiemenu</path>
    <platform/>
    <theme>retropie</theme>
  </system>
</systemList>
`)

	p := &storage.XMLParser{}
	var doc interface{}

	t.Run("FromBytes", func(t *testing.T) {
		res, err := p.FromBytes(bytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		doc = res.Interface()
		got := doc.(map[string]interface{})["systemList"].(map[string]interface{})["system"].(map[string]interface{})["command"]
		exp := "sudo /home/fozz/RetroPie-Setup/retropie_packages.sh retropiemenu launch %ROM% &lt;/dev/tty &gt;/dev/tty"
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})

	t.Run("ToBytes", func(t *testing.T) {
		gotBytes, err := p.ToBytes(dasel.ValueOf(doc))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(gotBytes)
		exp := string(bytes)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
}

func TestXMLParser_DifferentEncodings(t *testing.T) {
	newXmlBytes := func(newWriter func(io.Writer) io.Writer, encoding, text string) []byte {
		const encodedXmlBytesFmt = `<?xml version='1.0' encoding='%s'?>`
		const xmlBody = `<user><name>%s</name></user>`

		var buf bytes.Buffer

		w := newWriter(&buf)
		fmt.Fprintf(w, xmlBody, text)

		return []byte(fmt.Sprintf(encodedXmlBytesFmt, encoding) + buf.String())
	}

	testCases := []struct {
		name string
		xml  []byte
	}{
		{
			name: "supports ISO-8859-1",
			xml:  newXmlBytes(charmap.ISO8859_1.NewEncoder().Writer, "ISO-8859-1", "Tõm"),
		},
		{
			name: "supports UTF-8",
			xml:  newXmlBytes(unicode.UTF8.NewEncoder().Writer, "UTF-8", "Tõm"),
		},
		{
			name: "supports latin1",
			xml:  newXmlBytes(charmap.Windows1252.NewEncoder().Writer, "latin1", "Tõm"),
		},
		{
			name: "supports UTF-16",
			xml:  newXmlBytes(unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Writer, "UTF-16", "Tõm"),
		},
		{
			name: "supports UTF-16 (big endian)",
			xml:  newXmlBytes(unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewEncoder().Writer, "UTF-16BE", "Tõm"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := (&storage.XMLParser{}).FromBytes(tc.xml)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !reflect.DeepEqual(encodedXmlMap, got.Interface()) {
				t.Errorf("expected %v, got %v", encodedXmlMap, got)
			}
		})
	}
}
