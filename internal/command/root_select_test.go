package command_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"io/ioutil"
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

		output, err := ioutil.ReadAll(outputBuffer)
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

		output, err := ioutil.ReadAll(outputBuffer)
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

	t.Run("ObjectKeysSelector", selectTestContainsLines(jsonData, "json", ".-", []string{`"id"`,
		`"details"`}, nil, "-m"))

	t.Run("ArrayIndexesSelector", selectTest(jsonData, "json", ".details.addresses.-", newline(`"0"
"1"`), nil, "-m"))

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
}
