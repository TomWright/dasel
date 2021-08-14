package storage_test

import (
	"bytes"
	"github.com/alecthomas/chroma/quick"
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
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

func TestXMLParser_FromBytes(t *testing.T) {
	got, err := (&storage.XMLParser{}).FromBytes(xmlBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(&storage.BasicSingleDocument{Value: xmlMap}, got) {
		t.Errorf("expected %v, got %v", xmlMap, got)
	}
}

func TestXMLParser_FromBytes_Empty(t *testing.T) {
	got, err := (&storage.XMLParser{}).FromBytes([]byte{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if got != nil {
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
	got, err := (&storage.XMLParser{}).ToBytes(xmlMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(xmlBytes, got) {
		t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
	}
}
func TestXMLParser_ToBytes_SingleDocument(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: xmlMap})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(xmlBytes, got) {
		t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
	}
}
func TestXMLParser_ToBytes_SingleDocument_Colourise(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: xmlMap}, storage.ColouriseOption(true))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	_ = quick.Highlight(buf, string(xmlBytes), "xml", storage.ColouriseFormatter, storage.ColouriseStyle)
	exp := buf.Bytes()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestXMLParser_ToBytes_MultiDocument(t *testing.T) {
	got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{xmlMap, xmlMap}})
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
	got, err := (&storage.XMLParser{}).ToBytes("asd")
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
	got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: "asd"})
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
	got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{"asd", "123"}})
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
		doc = res.(storage.SingleDocument).Document()
		got := doc.(map[string]interface{})["systemList"].(map[string]interface{})["system"].(map[string]interface{})["command"]
		exp := "sudo /home/fozz/RetroPie-Setup/retropie_packages.sh retropiemenu launch %ROM% &lt;/dev/tty &gt;/dev/tty"
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})

	t.Run("ToBytes", func(t *testing.T) {
		gotBytes, err := p.ToBytes(doc)
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
