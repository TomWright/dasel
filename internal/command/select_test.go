package command

import (
	"testing"
)

func standardJsonSelectTestData() []byte {
	return []byte(`{
  "users": [
    {
      "name": {
        "first": "Tom",
        "last": "Wright"
      },
      "flags": {
        "isBanned": false
      }
    },
    {
      "name": {
        "first": "Jim",
        "last": "Wright"
      },
      "flags": {
        "isBanned": true
      }
    },
    {
      "name": {
        "first": "Joe",
        "last": "Blogs"
      },
      "flags": {
        "isBanned": false
      }
    }
  ]
}`)
}

func TestSelectCommand(t *testing.T) {

	t.Run("TotalUsersLen", runTest(
		[]string{"-r", "json", "--pretty=false", "users.len()"},
		standardJsonSelectTestData(),
		newline([]byte(`3`)),
		nil,
		nil,
	))

	t.Run("TotalUsersCount", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().count()"},
		standardJsonSelectTestData(),
		newline([]byte(`3`)),
		nil,
		nil,
	))

	t.Run("TotalBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,true)).count()"},
		standardJsonSelectTestData(),
		newline([]byte(`1`)),
		nil,
		nil,
	))

	t.Run("TotalNotBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,false)).count()"},
		standardJsonSelectTestData(),
		newline([]byte(`2`)),
		nil,
		nil,
	))

	t.Run("NotBannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,false)).name.first"},
		standardJsonSelectTestData(),
		newline([]byte(`"Tom"
"Joe"`)),
		nil,
		nil,
	))

	t.Run("BannedUsers", runTest(
		[]string{"-r", "json", "--pretty=false", "users.all().filter(equal(flags.isBanned,true)).name.first"},
		standardJsonSelectTestData(),
		newline([]byte(`"Jim"`)),
		nil,
		nil,
	))

	t.Run("VerifyCorrectIndentionForJSON", runTest(
		[]string{"-r", "json", "--indent", "6", "--pretty=true", "users.all().filter(equal(flags.isBanned,true)).name"},
		standardJsonSelectTestData(),
		newline([]byte("{\n      \"first\": \"Jim\",\n      \"last\": \"Wright\"\n}")),
		nil,
		nil,
	))

	t.Run("VerifyCorrectIndentionForYAML", runTest(
		[]string{"-r", "json", "-w", "yaml", "--indent", "6", "--pretty=true", "users.all().filter(equal(flags.isBanned,true))"},
		standardJsonSelectTestData(),
		newline([]byte("name:\n      first: Jim\n      last: Wright\nflags:\n      isBanned: true")),
		nil,
		nil,
	))

	t.Run("VerifyCorrectIndentionForTOML", runTest(
		[]string{"-r", "json", "-w", "toml", "--indent", "6", "--pretty=true", "users.all().filter(equal(flags.isBanned,true))"},
		standardJsonSelectTestData(),
		newline([]byte("[flags]\n      isBanned = true\n\n[name]\n      first = 'Jim'\n      last = 'Wright'")),
		nil,
		nil,
	))

	t.Run("Issue258", runTest(
		[]string{"-r", "json", "--pretty=false", "-w", "csv", "phones.all().mapOf(make,make,model,model,first,parent().parent().user.name.first,last,parent().parent().user.name.last).merge()"},
		[]byte(`{
	  "id": "1234",
	  "user": {
	    "name": {
	      "first": "Tom",
	      "last": "Wright"
	    }
	  },
	  "favouriteNumbers": [
	    1, 2, 3, 4
	  ],
	  "favouriteColours": [
	    "red", "green"
	  ],
	  "phones": [
	    {
	      "make": "OnePlus",
	      "model": "8 Pro"
	    },
	    {
	      "make": "Apple",
	      "model": "iPhone 12"
	    }
	  ]
	}`),
		newline([]byte(`first,last,make,model
Tom,Wright,OnePlus,8 Pro
Tom,Wright,Apple,iPhone 12`)),
		nil,
		nil,
	))

	t.Run("Issue181", runTest(
		[]string{"-r", "json", "--pretty=false", "all().filter(equal(this(),README.md))"},
		[]byte(`[
  "README.md",
  "tbump.toml"
]`),
		newline([]byte(`"README.md"`)),
		nil,
		nil,
	))

	// Flaky test due to ordering
	// t.Run("Discussion242", runTest(
	// 	[]string{"-r", "json", "--pretty=false", "-w", "plain", "all().filter(equal(type(),array)).key()"},
	// 	[]byte(`{
	// "array1": [
	//   {
	//     "a": "aaa",
	//     "b": "bbb",
	//     "c": "ccc"
	//   }
	// ],
	// "array2": [
	//   {
	//     "a": "aaa",
	//     "b": "bbb",
	//     "c": "ccc"
	//   }
	// ]
	// }`),
	// 		newline([]byte(`array1
	// array2`)),
	// 		nil,
	// 		nil,
	// 	))

	t.Run("YamlMultiDoc/Issue314", runTest(
		[]string{"-r", "yaml", ""},
		[]byte(`a: x
b: foo
---
a: y
c: bar
`),
		newline([]byte(`a: x
b: foo
---
a: y
c: bar`)),
		nil,
		nil,
	))

	t.Run("Issue316", runTest(
		[]string{"-r", "json"},
		[]byte(`{
  "a": "alice",
  "b": null,
  "c": [
    {
      "d": 9,
      "e": null
    },
    null
  ]
}`),
		newline([]byte(`{
  "a": "alice",
  "b": null,
  "c": [
    {
      "d": 9,
      "e": null
    },
    null
  ]
}`)),
		nil,
		nil,
	))

	// Hex, binary and octal values in YAML
	t.Run("Issue326", runTest(
		[]string{"-r", "yaml"},
		[]byte(`hex: 0x1234
binary: 0b1001
octal: 0o10
`),
		newline([]byte(`hex: 4660
binary: 9
octal: 8`)),
		nil,
		nil,
	))

	t.Run("Issue331 - YAML to JSON", runTest(
		[]string{"-r", "yaml", "-w", "json"},
		[]byte(`createdAt: 2023-06-13T20:19:48.531620935Z`),
		newline([]byte(`{
  "createdAt": "2023-06-13T20:19:48.531620935Z"
}`)),
		nil,
		nil,
	))

	t.Run("OrDefaultString", runTest(
		[]string{"-r", "json", "all().orDefault(locale,string(nope))"},
		[]byte(`{
  "-LCr5pXw_fN32IqNDr4E": {
    "bookCategory": "poetry",
    "locale": "en-us",
    "mediaType": "book",
    "publisher": "Pomelo Books",
    "title": "Sound Waves",
    "type": "poetry"
  },
  "-LDDHjkdY0306fZdvhEQ": {
    "ISBN13": "978-1534402966",
    "bookCategory": "fiction",
    "title": "What Can You Do with a Toolbox?",
    "type": "picturebook"
  }
}`),
		newline([]byte(`"en-us"
"nope"`)),
		nil,
		nil,
	))

	t.Run("OrDefaultLookup", runTest(
		[]string{"-r", "json", "all().orDefault(locale,bookCategory)"},
		[]byte(`{
  "-LCr5pXw_fN32IqNDr4E": {
    "bookCategory": "poetry",
    "locale": "en-us",
    "mediaType": "book",
    "publisher": "Pomelo Books",
    "title": "Sound Waves",
    "type": "poetry"
  },
  "-LDDHjkdY0306fZdvhEQ": {
    "ISBN13": "978-1534402966",
    "bookCategory": "fiction",
    "title": "What Can You Do with a Toolbox?",
    "type": "picturebook"
  }
}`),
		newline([]byte(`"en-us"
"fiction"`)),
		nil,
		nil,
	))

	t.Run("Issue364 - CSV root element part 1", runTest(
		[]string{"-r", "csv", "-w", "csv", "all().merge()"},
		[]byte(`A,B,C
a,b,c
d,e,f`),
		newline([]byte(`A,B,C
a,b,c
d,e,f`)),
		nil,
		nil,
	))

	t.Run("Issue364 - CSV root element part 2", runTest(
		[]string{"-r", "csv", "-w", "csv"},
		[]byte(`A,B,C
a,b,c
d,e,f`),
		newline([]byte(`A,B,C
a,b,c
d,e,f`)),
		nil,
		nil,
	))

	t.Run("CSV custom separator", runTest(
		[]string{"-r", "csv", "-w", "csv", "--csv-comma", "."},
		[]byte(`A.B.C
a.b.c
d.e.f`),
		newline([]byte(`A.B.C
a.b.c
d.e.f`)),
		nil,
		nil,
	))

	t.Run("CSV change separator", runTest(
		[]string{"-r", "csv", "-w", "csv", "--csv-comma", ".", "--csv-write-comma", ","},
		[]byte(`A.B.C
a.b.c
d.e.f`),
		newline([]byte(`A,B,C
a,b,c
d,e,f`)),
		nil,
		nil,
	))

	t.Run("CSV change from default separator", runTest(
		[]string{"-r", "csv", "-w", "csv", "--csv-write-comma", "."},
		[]byte(`A,B,C
a,b,c
d,e,f`),
		newline([]byte(`A.B.C
a.b.c
d.e.f`)),
		nil,
		nil,
	))

	t.Run("Issue351 incorrectly escaped html, default false", runTest(
		[]string{"-r", "json"},
		[]byte(`{
  "field1": "A",
  "field2": "A > B && B > C"
}`),
		newline([]byte(`{
  "field1": "A",
  "field2": "A > B && B > C"
}`)),
		nil,
		nil,
	))

	t.Run("Issue351 incorrectly escaped html, specific false", runTest(
		[]string{"-r", "json", "--escape-html=false"},
		[]byte(`{
  "field1": "A",
  "field2": "A > B && B > C"
}`),
		newline([]byte(`{
  "field1": "A",
  "field2": "A > B && B > C"
}`)),
		nil,
		nil,
	))

	t.Run("Issue351 correctly escaped html", runTest(
		[]string{"-r", "json", "--escape-html=true"},
		[]byte(`{
  "field1": "A",
  "field2": "A > B && B > C"
}`),
		newline([]byte(`{
  "field1": "A",
  "field2": "A \u003e B \u0026\u0026 B \u003e C"
}`)),
		nil,
		nil,
	))

	t.Run("Issue 374 empty input", func(t *testing.T) {
		tests := []struct {
			format string
			exp    []byte
		}{
			{
				format: "json",
				exp:    []byte("{}\n"),
			},
			{
				format: "toml",
				exp:    []byte("\n"),
			},
			{
				format: "yaml",
				exp:    []byte("{}\n"),
			},
			{
				format: "xml",
				exp:    []byte("<doc/>\n"),
			},
			{
				format: "csv",
				exp:    []byte(""),
			},
		}

		for _, test := range tests {
			tc := test
			t.Run(tc.format, runTest(
				[]string{"-r", tc.format},
				[]byte(``),
				tc.exp,
				nil,
				nil,
			))
		}
	})

	t.Run("Issue 392 panic", runTest(
		[]string{"-r", "csv", "--csv-comma", ";", "-w", "json", "equal([], )"},
		[]byte(`Hello;There;
1;2;`),
		[]byte("false\n"),
		nil,
		nil,
	))
  
	t.Run("Issue346", func(t *testing.T) {
		t.Run("Select null or default string", runTest(
			[]string{"-r", "json", "orDefault(foo,string(nope))"},
			[]byte(`{
  "foo": null
}`),
			newline([]byte(`"nope"`)),
			nil,
			nil,
		))

		t.Run("Select null or default null", runTest(
			[]string{"-r", "json", "orDefault(foo,null())"},
			[]byte(`{
  "foo": null
}`),
			newline([]byte(`null`)),
			nil,
			nil,
		))

		t.Run("Select null value", runTest(
			[]string{"-r", "json", "foo"},
			[]byte(`{
  "foo": null
}`),
			newline([]byte(`null`)),
			nil,
			nil,
		))
	})

}
