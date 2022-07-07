package command_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"io"
	"strings"
	"testing"
)

const jsonDataSingle = `{"x": "asd"}`
const yamlDataSingle = `x: asd`
const tomlDataSingle = `x="asd"`
const xmlDataSingle = `<x>asd</x>`

const jsonData = `{
  "id": "1111",
  "details": {
    "name": "Tom",
  	"age": 27,
    "addresses": [
      {
        "street": "101 Some Street",
        "town": "Some Town",
        "county": "Some Country",
        "postcode": "XXX XXX",
        "primary": true
      },
      {
        "street": "34 Another Street",
        "town": "Another Town",
        "county": "Another County",
        "postcode": "YYY YYY"
      }
    ]
  }
}`

const yamlData = `
id: 1111
details:
  name: Tom
  age: 27
  addresses:
  - street: 101 Some Street
    town: Some Town
    county: Some County
    postcode: XXX XXX
    primary: true
  - street: 34 Another Street
    town: Another Town
    county: Another County
    postcode: YYY YYY
`

const tomlData = `id = "1111"
[details]
  name = "Tom"
  age = 27
  [[details.addresses]]
    street =  "101 Some Street"
    town = "Some Town"
    county = "Some County"
    postcode = "XXX XXX"
    primary = true
  [[details.addresses]]
    street = "34 Another Street"
    town = "Another Town"
    county = "Another County"
    postcode = "YYY YYY"
`

const xmlData = `<data>
	<id>1111</id>
	<details>
		<name>Tom</name>
		<age>27</age>
		<addresses primary="true">
			<street>101 Some Street</street>
			<town>Some Town</town>
			<county>Some County</county>
			<postcode>XXX XXX</postcode>
		</addresses>
		<addresses>
			<street>34 Another Street</street>
			<town>Another Town</town>
			<county>Another County</county>
			<postcode>YYY YYY</postcode>
		</addresses>
	</details>
</data>
`

const csvData = `id,name
1,Tom
2,Jim
`

func newline(x string) string {
	return x + "\n"
}

func TestRootCMD_Select(t *testing.T) {
	t.Run("InvalidFile", expectErr(
		[]string{"select", "-f", "bad.json", "-s", "x"},
		"could not open input file",
	))
	t.Run("MissingParser", expectErr(
		[]string{"select", "-s", "x"},
		"parser flag required when reading from stdin",
	))
	t.Run("Stdin", expectOutput(
		`{"name": "Tom"}`,
		[]string{"select", "-f", "stdin", "-p", "json", "-s", ".name"},
		`"Tom"
`,
	))
	t.Run("StdinAlias", expectOutput(
		`{"name": "Tom"}`,
		[]string{"select", "-f", "-", "-p", "json", "-s", ".name"},
		`"Tom"
`,
	))

	t.Run("InvalidSingleSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"select", "-p", "json", "-s", "[-]"},
		"invalid index: -",
	))
	t.Run("InvalidMultiSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"select", "-p", "json", "-m", "-s", "[-]"},
		"invalid index: -",
	))
}

func selectTest(in string, parser string, selector string, output string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return selectTestCheck(in, parser, selector, func(out string) error {
		if out != output {
			return fmt.Errorf("expected %v, got %v", output, out)
		}
		return nil
	}, expErr, additionalArgs...)
}

func selectTestContainsLines(in string, parser string, selector string, output []string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return selectTestCheck(in, parser, selector, func(out string) error {
		splitOut := strings.Split(out, "\n")
		for _, s := range output {
			found := false
			for _, got := range splitOut {
				if s == got {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("required value not found: %s", s)
			}
		}
		return nil
	}, expErr, additionalArgs...)
}

func selectTestCheck(in string, parser string, selector string, checkFn func(out string) error, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"select", "-p", parser,
		}
		if additionalArgs != nil {
			args = append(args, additionalArgs...)
		}
		args = append(args, selector)

		cmd.SetOut(outputBuffer)
		cmd.SetIn(strings.NewReader(in))
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := io.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if err := checkFn(string(output)); err != nil {
			t.Errorf("unexpected output: %s", err)
		}
	}
}

func selectTestFromFile(inputPath string, selector string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"select", "-f", inputPath, "-s", selector,
		}

		cmd.SetOut(outputBuffer)
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := io.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		if out != string(output) {
			t.Errorf("expected result %v, got %v", out, string(output))
		}
	}
}

func TestRootCmd_Select_JSON(t *testing.T) {
	t.Run("RootElement", selectTest(jsonDataSingle, "json", ".", newline(`{
  "x": "asd"
}`), nil))
	t.Run("SingleProperty", selectTest(jsonData, "json", ".id", newline(`"1111"`), nil))
	t.Run("ObjectProperty", selectTest(jsonData, "json", ".details.name", newline(`"Tom"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[0].street", newline(`"101 Some Street"`), nil))
	t.Run("Index", selectTest(jsonData, "json", ".details.addresses.[1].street", newline(`"34 Another Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=XXX XXX).street", newline(`"101 Some Street"`), nil))
	t.Run("DynamicString", selectTest(jsonData, "json", ".details.addresses.(postcode=YYY YYY).street", newline(`"34 Another Street"`), nil))
	t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.json", ".preferences.favouriteColour", newline(`"red"`), nil))

	t.Run("MultiProperty", selectTest(jsonData, "json", ".details.addresses.[*].street", newline(`"101 Some Street"
"34 Another Street"`), nil, "-m"))

	t.Run("MultiRoot", selectTest(jsonDataSingle, "json", ".", newline(`{
  "x": "asd"
}`), nil, "-m"))

	t.Run("SubSelector", selectTest(`{
  "users": [
	{
	  "primary": true,
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  }
	},
	{
	  "primary": false,
	  "name": {
		"first": "Jim",
		"last": "Wright"
	  }
	}
  ]
}`, "json", ".users.(name.first=Tom).primary", newline(`true`), nil))

	t.Run("SubSubSelector", selectTest(`{
  "users": [
	{
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  },
      "addresses": [
        {
          "primary": true,
          "number": 123
        },
        {
          "primary": false,
          "number": 456
        }
      ]
	}
  ]
}`, "json", ".users.(.addresses.(.primary=true).number=123).name.first", newline(`"Tom"`), nil))

	t.Run("SubSubAndSelector", selectTest(`{
  "users": [
	{
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  },
      "addresses": [
        {
          "primary": true,
          "number": 123
        },
        {
          "primary": false,
          "number": 456
        }
      ]
	}
  ]
}`, "json", ".users.(.addresses.(.primary=true).number=123)(.name.last=Wright).name.first", newline(`"Tom"`), nil))

	t.Run("KeySearch", selectTestContainsLines(`{
  "users": [
    {
      "primary": true,
      "name": {
        "first": "Tom",
        "last": "Wright"
      }
    },
    {
      "primary": false,
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
      "name": {
        "first": "Jim",
        "last": "Wright"
      }
    }
  ]
}`, "json", ".(?:-=name).first", []string{`"Tom"`, `"Joe"`, `"Jim"`}, nil, "-m"))

	t.Run("NullNotFound", selectTest(`{}`, "json", ".asd", newline(`null`), nil, "-n"))
	t.Run("NullNotFoundMulti", selectTest(`{}`, "json", ".asd", newline(`null`), nil, "-m", "-n"))

	t.Run("ObjectKeysSelector", selectTestContainsLines(jsonData, "json", ".-", []string{`"id"`,
		`"details"`}, nil, "-m"))

	t.Run("ArrayIndexesSelector", selectTest(jsonData, "json", ".details.addresses.-", newline(`"0"
"1"`), nil, "-m"))

	t.Run("RootElementCompactShortFlag", selectTest(`{
  "x": "asd"
}`, "json", ".", newline(`{"x":"asd"}`), nil, "-c"))
	t.Run("RootElementCompactLongFlag", selectTest(`{
  "x": "asd"
}`, "json", ".", newline(`{"x":"asd"}`), nil, "--compact"))

	t.Run("LengthFlagList", selectTest(`{
  "x": [ "a", "b", "c" ]
}`, "json", ".x", newline(`3`), nil, "--length"))
	t.Run("LengthFlagMap", selectTest(`{
  "x": { "a": 1, "b": 2, "c": 3 }
}`, "json", ".x", newline(`3`), nil, "--length"))
	t.Run("LengthFlagString", selectTest(`{
  "x": "asd"
}`, "json", ".x", newline(`3`), nil, "--length"))
	t.Run("LengthFlagInt", selectTest(`{
  "x": 123
}`, "json", ".x", newline(`3`), nil, "--length"))
	t.Run("LengthFlagBool", selectTest(`{
  "x": true
}`, "json", ".x", newline(`4`), nil, "--length"))
	t.Run("LengthFlagMulti", selectTest(`[
  [ "a", "b", "c" ],
  { "a": 1 },
  "hello there",
  12345,
  123.45,
  true
]`, "json", ".[*]", newline(`3
1
11
5
6
4`), nil, "--length", "-m"))

	t.Run("NullInput", selectTest(`null`, "json", `.`, newline("{}"), nil))
	t.Run("EmptyDocument", selectTest(`{}`, "json", `.`, newline("{}"), nil))
	t.Run("EmptyArray", selectTest(`[]`, "json", `.`, newline("[]"), nil))
	t.Run("BlankInput", selectTest(``, "json", `.`, newline("{}"), nil))

	t.Run("LengthSelector", selectTest(jsonData, "json", `.details.addresses.[#]`, newline("2"), nil))
	t.Run("LengthSelectorMultiple", selectTest(jsonData, "json", `.details.addresses.[#]`, newline("2"), nil, "-m"))
	t.Run("LengthSelectorMultiple", selectTest(jsonData, "json", `.details.addresses.[*].[#]`, newline("5\n4"), nil, "-m"))
	t.Run("LengthDynamicSelector", selectTest(`{
  "a": {
    "id": 1,
    "uses": [1]
  },
  "b": {
    "id": 2,
    "uses": [1, 2]
  },
  "c": {
    "id": 3,
    "uses": [1, 2, 3]
  }
}`, "json", `.(.uses.[#]=2).id`, newline("2"), nil))
	t.Run("LengthDynamicSelectorMultiple", selectTest(`[
  {
    "id": 1,
    "uses": [1]
  },
  {
    "id": 2,
    "uses": [1, 2]
  },
  {
    "id": 3,
    "uses": [1, 2, 3]
  },
  {
    "id": 4,
    "uses": [3, 4]
  }
]`, "json", `.(.uses.[#]=2).id`, newline("2\n4"), nil, "-m"))
	t.Run("LengthDynamicSelectorMoreThan", selectTest(`[
  {
    "id": 1,
    "uses": [1]
  },
  {
    "id": 2,
    "uses": [1, 2]
  },
  {
    "id": 3,
    "uses": [1, 2, 3]
  },
  {
    "id": 4,
    "uses": [3, 4]
  }
]`, "json", `.(.uses.[#]>2).id`, newline("3"), nil, "-m"))

	t.Run("MergeInputDocuments", selectTest(`{
  "number": 1
}
{
  "number": 2
}
{
  "number": 3
}
`, "json", `.`, `[
  {
    "number": 1
  },
  {
    "number": 2
  },
  {
    "number": 3
  }
]
`, nil, "--merge-input-documents"))

	t.Run("EscapeHTMLOn", selectTest(`{
  "user": "Tom <contact@tomwright.me>"
}
`, "json", `.`, `{
  "user": "Tom \u003ccontact@tomwright.me\u003e"
}
`, nil, "--escape-html=true"))

	t.Run("EscapeHTMLOff", selectTest(`{
  "user": "Tom <contact@tomwright.me>"
}
`, "json", `.`, `{
  "user": "Tom <contact@tomwright.me>"
}
`, nil, "--escape-html=false"))

	t.Run("MixedDynamicSelectors", selectTest(`{
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/gitlab",
    [
      "@semantic-release/git",
      {
        "assets": [
          "tbump.toml",
          "**/pyproject.toml",
          "**/setup.py",
          "README.md"
        ],
        "message": "chore(release): ${nextRelease.version}\n\n${nextRelease.notes}"
      }
    ],
	[
      "@semantic-release/git",
      {
        "assets": [
          "y"
        ],
        "message": "chore(release): ${nextRelease.version}\n\n${nextRelease.notes}"
      }
    ]
  ]
}`, "json", `.plugins.([@]=array).([@]=map).assets`, `[
  "tbump.toml",
  "**/pyproject.toml",
  "**/setup.py",
  "README.md"
]
[
  "y"
]
`, nil, "-m"))

	t.Run("SearchOptional", selectTest(`{
  "users": [
    {
      "name": "Tom",
      "blocked": true
    },
    {
      "name": "Jim",
      "blocked": false
    },
    {
      "name": "Frank"
    }
  ]
}`, "json", `.users.(#:blocked=true).name`, `Tom
`, nil, "-m", "--plain"))

}

func TestRootCmd_Select_YAML(t *testing.T) {
	t.Run("RootElement", selectTest(yamlDataSingle, "yaml", ".", newline(`x: asd`), nil))
	t.Run("SingleProperty", selectTest(yamlData, "yaml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(yamlData, "yaml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(yamlData, "yaml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(yamlData, "yaml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
	t.Run("QueryFromFile", selectTestFromFile("./../../tests/assets/example.yaml", ".preferences.favouriteColour", newline(`red`), nil))

	// Following test implemented as a result of issue #35.
	t.Run("MultipleSeparateDynamic", selectTest(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: harbor-exporter
  labels:
    app: harbor-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: harbor-exporter
  strategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: harbor-exporter
    spec:
      serviceAccountName: default
      restartPolicy: Always
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
      containers:
        - name: harbor-exporter
          image: "c4po/harbor-exporter:debug"
          imagePullPolicy: Always
          env:
            - name: HARBOR_URI
#            name of the Service for harbor-core
              value: http://harbor-core.harbor # change prefix to the name of your Helm release
            - name: HARBOR_USERNAME
              value: "admin"
            - name: HARBOR_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: harbor-core # change prefix to the name of your Helm release
                  key: HARBOR_ADMIN_PASSWORD

          securityContext:
            capabilities:
              drop:
                - SETPCAP
                - MKNOD
                - AUDIT_WRITE
                - CHOWN
                - NET_RAW
                - DAC_OVERRIDE
                - FOWNER
                - FSETID
                - KILL
                - SETGID
                - SETUID
                - NET_BIND_SERVICE
                - SYS_CHROOT
                - SETFCAP
            readOnlyRootFilesystem: true
          resources:
            limits:
              cpu: 400m
              memory: 256Mi
            requests:
              cpu: 100m
              memory: 64Mi
          ports:
            - containerPort: 9107
              name: http
          livenessProbe:
            httpGet:
              path: /-/healthy
              port: http
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 5
          readinessProbe:
            httpGet:
              path: /-/ready
              port: http
            initialDelaySeconds: 1
            timeoutSeconds: 5
            periodSeconds: 5
`, "yaml", "spec.template.spec.containers.(name=harbor-exporter).env.(name=HARBOR_URI).value", newline(`http://harbor-core.harbor`), nil))

	// https://github.com/TomWright/dasel/issues/99
	// Worked in v1.13.3
	t.Run("NullInput", selectTest(`null`, "yaml", `.`, newline("{}"), nil))
	t.Run("EmptyDocument", selectTest(`---`, "yaml", `.`, newline("{}"), nil))
	t.Run("BlankInput", selectTest(``, "yaml", `.`, newline("{}"), nil))
}

func TestRootCmd_Select_TOML(t *testing.T) {
	t.Run("RootElement", selectTest(tomlDataSingle, "toml", ".", newline(`x = "asd"`), nil))
	t.Run("SingleProperty", selectTest(tomlData, "toml", ".id", newline(`1111`), nil))
	t.Run("ObjectProperty", selectTest(tomlData, "toml", ".details.name", newline(`Tom`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[0].street", newline(`101 Some Street`), nil))
	t.Run("Index", selectTest(tomlData, "toml", ".details.addresses.[1].street", newline(`34 Another Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=XXX XXX).street", newline(`101 Some Street`), nil))
	t.Run("DynamicString", selectTest(tomlData, "toml", ".details.addresses.(postcode=YYY YYY).street", newline(`34 Another Street`), nil))
}

func TestRootCMD_Select_XML(t *testing.T) {
	t.Run("RootElement", selectTest(xmlDataSingle, "xml", ".", newline(`<x>asd</x>`), nil))
	t.Run("SingleProperty", selectTest(xmlData, "xml", ".data.id", "1111\n", nil))
	t.Run("ObjectProperty", selectTest(xmlData, "xml", ".data.details.name", "Tom\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[0].street", "101 Some Street\n", nil))
	t.Run("Index", selectTest(xmlData, "xml", ".data.details.addresses.[1].street", "34 Another Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=XXX XXX).street", "101 Some Street\n", nil))
	t.Run("DynamicString", selectTest(xmlData, "xml", ".data.details.addresses.(postcode=YYY YYY).street", "34 Another Street\n", nil))
	t.Run("Attribute", selectTest(xmlData, "xml", ".data.details.addresses.(-primary=true).street", "101 Some Street\n", nil))

	t.Run("KeySearch", selectTestContainsLines(`
<food>
  <tart>
    <apple color="yellow"/>
  </tart>
  <pie>
    <crust quality="flaky"/>
    <filling>
      <apple color="red"/>
    </filling>
  </pie>
  <apple color="green"/>
</food>
`, "xml", ".food.(?:keyValue=apple).-color", []string{"yellow", "red", "green"}, nil, "-m"))
}

func TestRootCMD_Select_CSV(t *testing.T) {
	t.Run("RootElement", selectTest(csvData, "csv", ".", csvData, nil))
	t.Run("SingleProperty", selectTest(csvData, "csv", ".[0].id", "1\n", nil))
	t.Run("SingleProperty", selectTest(csvData, "csv", ".[1].id", "2\n", nil))

	// https://github.com/TomWright/dasel/issues/110
	t.Run("ObjectArrayJSONToCSV", selectTest(`
[
  {
    "id": "ABS",
    "name": "Australian Bureau of Statistics"
  },
  {
    "id": "ECB",
    "name": "European Central Bank"
  },
  {
    "id": "ESTAT",
    "name": "Eurostat"
  },
  {
    "id": "ILO",
    "name": "International Labor Organization"
  }
]
`, "json", ".", `id,name
ABS,Australian Bureau of Statistics
ECB,European Central Bank
ESTAT,Eurostat
ILO,International Labor Organization
`, nil, "-w", "csv"))

	// https://github.com/TomWright/dasel/issues/159
	t.Run("SelectFilterOnUnicode", selectTest(`
{
   "data": [
     {"name": "Fu Shun", "ship_type": "Destroyer"},
     {"name": "Sheffield", "ship_type": "Light Cruiser"},
     {"name": "Ägir", "ship_type": "Large Cruiser"}
   ]
 }
`, "json", ".data.(name=Ägir).ship_type", `"Large Cruiser"
`, nil))
}

func TestRootCmd_Select_JSON_Format(t *testing.T) {
	t.Run("RootElementFormattedToProperty", selectTest(jsonData, "json", ".", newline(`1111`), nil,
		"--format", `{{ query ".id" }}`))
	t.Run("SelectorFormatted", selectTest(jsonData, "json", ".id", newline(`1111`), nil,
		"--format", `{{ . }}`))
	t.Run("SelectorFormattedMultiple", selectTest(jsonData, "json", ".details.addresses.[*]",
		newline(`101 Some Street
34 Another Street`), nil,
		"-m", "--format", `{{ query ".street" }}`))
	t.Run("SelectorFormattedToMultiple", selectTest(jsonData, "json", ".",
		newline(`101 Some Street
34 Another Street`), nil,
		"-m", "--format", `{{ queryMultiple ".details.addresses.[*]" | format "{{ .street }}{{ if not isLast }}{{ newline }}{{end}}" }}`))

	// https://github.com/TomWright/dasel/discussions/146
	t.Run("Discussion146", selectTest(
		`[{"name": "click", "version": "7.1.2", "latest_version": "8.0.1", "latest_filetype": "wheel"}, {"name": "decorator", "version": "4.4.2", "latest_version": "5.0.9", "latest_filetype": "wheel"}, {"name": "ipython", "version": "7.20.0", "latest_version": "7.25.0", "latest_filetype": "wheel"}, {"name": "pandas", "version": "1.3.0", "latest_version": "1.3.1", "latest_filetype": "wheel"}, {"name": "parso", "version": "0.8.1", "latest_version": "0.8.2", "latest_filetype": "wheel"}, {"name": "pip", "version": "21.1.3", "latest_version": "21.2.1", "latest_filetype": "wheel"}, {"name": "prompt-toolkit", "version": "3.0.14", "latest_version": "3.0.19", "latest_filetype": "wheel"}, {"name": "Pygments", "version": "2.7.4", "latest_version": "2.9.0", "latest_filetype": "wheel"}, {"name": "setuptools", "version": "49.2.1", "latest_version": "57.4.0", "latest_filetype": "wheel"}, {"name": "tomli", "version": "1.0.4", "latest_version": "1.1.0", "latest_filetype": "wheel"}]`,
		"json", ".(name!=setuptools)(name!=six)(name!=pip)(name!=pip-tools)",
		newline(`click
7.1.2
8.0.1
decorator
4.4.2
5.0.9
ipython
7.20.0
7.25.0
pandas
1.3.0
1.3.1
parso
0.8.1
0.8.2
prompt-toolkit
3.0.14
3.0.19
Pygments
2.7.4
2.9.0
tomli
1.0.4
1.1.0`), nil,
		"-m", "--format", `{{ query ".name" }}{{ newline }}{{ query ".version" }}{{ newline }}{{ query ".latest_version" }}`))
}
