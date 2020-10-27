# dasel

[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

Comparable to [jq](https://github.com/stedolan/jq) / [yq](https://github.com/kislyuk/yq), but supports JSON, YAML and TOML with zero dependencies.

# Table of contents
* [Dasel](#dasel)
* [Installation](#installation)
* [Usage](#usage)
  * [Select](#select)
  * [Put](#put)
  * [Put Object](#put-object)
* [Supported file types](#supported-file-types)
* [Selectors](#selectors)
* [Examples](#examples)
  * [jq to dasel](#jq-to-dasel)
  * [Kubernetes](#kubernetes)

### Installation
You can import dasel as a package and use it in your applications, or you can use a pre-built binary to modify files from the command line.

#### Command line
You can `go get` the `main` package and go should automatically build and install dasel for you.
```
go get github.com/tomwright/dasel/cmd/dasel
```

Alternatively you can download a compiled executable from the [latest release](https://github.com/TomWright/dasel/releases/latest).
##### Linux amd64
This one liner should work for you - be sure to change the targeted release executable if needed. It currently targets `dasel_linux_amd64`.
```
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep linux_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_linux_amd64 dasel && chmod +x dasel
mv ./dasel /usr/local/bin/dasel
```

##### Mac OS amd64
You may have to `brew install wget` in order for this to work.
```
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep macos_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_macos_amd64 dasel && chmod +x dasel
mv ./dasel /usr/local/bin/dasel
```

#### Import
As with any other go package, just use `go get`.
```
go get github.com/tomwright/dasel
```

## Notes

The formatting of files can be changed while being processed. Dasel itself doesn't make these changes, rather the act of marshaling the results.

In short, the output files may have properties in a different order but the actual contents will be as expected.

## Usage 

```bash
dasel -h
```

### Select
```bash
dasel select -f <file> -p <json|yaml|toml> -s <selector>
```

#### Arguments

##### `-f`, `--file`

Specify the file to query. This is required unless you are piping in data.

If piping in data you can optionally pass `-f stdin`.

##### `-p`, `--parser`

Specify the parser to use when reading the file.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

##### `-s`, `--selector`

Specify the selector to use. See [Selectors](#selectors) for more information.

This is required.

#### Example

Select the image within a kubernetes deployment manifest file:
```bash
dasel select -f deployment.yaml -s "spec.template.spec.containers.(name=auth).image"
tomwright/auth:v1.0.0
```

Piping data into the select:
```bash
cat deployment.yaml | dasel select -p yaml -s "spec.template.spec.containers.(name=auth).image"
tomwright/auth:v1.0.0
```

### Put
```bash
dasel put <type> -f <file> -o <out> -p <parser> -s <selector> <value>
```

```bash
echo "name: Tom" | ./dasel put string -p yaml -s "name" Jim
name: Jim
```

#### Arguments

##### `type`

The type of value you want to put.

Available arguments:
- `string`
- `int`
- `bool`
- `object` - see [Put Object](#put-object)

##### `-f`, `--file`

Specify the file to query. This is required unless you are piping in data.

If piping in data you can optionally pass `-f stdin`.

##### `-o`, `--out`

Specify the output file. If present, results will be written to the given file. If not present, results will be written to the input file (or stdout if none given).

To force output to be written to stdout, pass `-o stdout`.

##### `-p`, `--parser`

Specify the parser to use when reading/writing the input/output files.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

##### `-s`, `--selector`

Specify the selector to use. See [Selectors](#selectors) for more information.

##### `value`

The value to write.

Dasel will parse this value as a string, int, or bool from this value depending on the given `type`.

This is required.

### Put Object

Putting objects works slightly differently to a standard put, but the same principles apply.

```bash
dasel put object -f <file> -o <out> -p <parser> -s <selector> -t <type> <values>
```

#### Arguments

##### `-t`, `--type`

The type of value you want to put.

You must repeat this argument for each value provided.

Available arguments:
- `string`
- `int`
- `bool`

##### `-f`, `--file`

Specify the file to query. This is required unless you are piping in data.

If piping in data you can optionally pass `-f stdin`.

##### `-o`, `--out`

Specify the output file. If present, results will be written to the given file. If not present, results will be written to the input file (or stdout if none given).

To force output to be written to stdout, pass `-o stdout`.

##### `-p`, `--parser`

Specify the parser to use when reading/writing the input/output files.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

##### `-s`, `--selector`

Specify the selector to use. See [Selectors](#selectors) for more information.

##### `values`

A space separated list of `key=value` pairs.

Dasel will parse each value as a string, int, or bool depending on the related `type`.

This is required.

#### Example

```bash
echo "" | dasel put object -s "my.favourites" -t string -t int colour=red number=3
```
Results in the following:
```yaml
my:
  favourites:
    colour: red
    number: 3
```

## Supported file types
Dasel attempts to find the correct parser for the given file type, but if that fails you can choose which parser to use with the `-p` or `--parser` flag. 

- JSON - `-p json`
- TOML - `-p toml`
- YAML - `-p yaml`

## Selectors

Selectors define a path through a set of data.

Selectors are made up of different parts separated by a dot `.`, each part being used to identify the next node in the chain.

The following YAML data structure will be used as a reference in the following examples.
```yaml
name: Tom
preferences:
  favouriteColour: red
colours:
- red
- green
- blue
colourCodes:
- name: red
  rgb: ff0000
- name: green
  rgb: 00ff00
- name: blue
  rgb: 0000ff
```

### Property
Property selectors are used to reference a single property of an object.

Just use the property name as a string.
```bash
dasel select -f ./tests/assets/example.yaml -s "name"
Tom
```
- `name` == `Tom`

### Child Elements
Just separate the child element from the parent element using a `.`:
```bash
dasel select -f ./tests/assets/example.yaml -s "preferences.favouriteColour"
red
```
- `preferences.favouriteColour` == `red`

#### Index
When you have a list, you can use square brackets to access a specific item in the list by its index.
```bash
dasel select -f ./tests/assets/example.yaml -s "colours.[1]"
green
```
- `colours.[0]` == `red`
- `colours.[1]` == `green`
- `colours.[2]` == `blue`

#### Next Available Index
Next available index selector is used when adding to a list of items. It allows you to append to a list.
- `colours.[]`

#### Dynamic
Dynamic selectors are used with lists when you don't know the index of the item, but instead know the value of a property of an object within the list. 

Look ups are defined in brackets. You can use multiple dynamic selectors within the same part to perform multiple checks.
```bash
dasel select -f ./tests/assets/example.yaml -s "colourCodes.(name=red).rgb"
ff0000

dasel select -f ./tests/assets/example.yaml -s "colourCodes.(name=blue)(rgb=0000ff)"
map[name:blue rgb:0000ff]
```
- `colourCodes.(name=red).rgb` == `ff0000`
- `colourCodes.(name=green).rgb` == `00ff00`
- `colourCodes.(name=blue).rgb` == `0000ff`
- `colourCodes.(name=blue)(rgb=0000ff).rgb` == `0000ff`

If you want to dynamically target a value in a list when it isn't a list of objects, just define the dynamic selector with `(value=<some_value>)` instead.

## Examples

### jq to dasel

The follow examples show a set of commands and the equivalent in dasel.

<table>
    <thead>
        <tr>
            <th>Tool</th>
            <th>Input</th>
            <th>Output</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <th style="text-align: center" colspan="3">Select a single value</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '{"name": "Tom"}' | jq '.name'</pre></td>
            <td><pre>"Tom"</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '{"name": "Tom"}' | dasel select -p json -s '.name'</pre></td>
            <td><pre>Tom</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Select a nested value</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '{"user": {"name": "Tom", "age": 27}}' | jq '.user.age'</pre></td>
            <td><pre>27</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '{"user": {"name": "Tom", "age": 27}}' | dasel select -p json -s '.user.age'</pre></td>
            <td><pre>27</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Select an array index</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '[1, 2, 3]' | jq '.[1]'</pre></td>
            <td><pre>2</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '[1, 2, 3]' | dasel select -p json -s '.[1]'</pre></td>
            <td><pre>2</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Update a string value</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '["a", "b", "c"]' | jq '.[1] = "d"'</pre></td>
            <td><pre>["a", "d", "c"]</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '["a", "b", "c"]' | dasel put string -p json -s '.[1]' d</pre></td>
            <td><pre>["a", "d", "c"]</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Update an int value</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '[1, 2, 3]' | jq '.[1] = 5'</pre></td>
            <td><pre>[1, 5, 3]</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '[1, 2, 3]' | dasel put int -p json -s '.[1]' 5</pre></td>
            <td><pre>[1, 5, 3]</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Overwrite an object</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '{"user": {"name": "Tom", "age": 27}}' | jq '.user = {"name": "Frank", "age": 25}'</pre></td>
            <td><pre>{"user": {"name": "Frank", "age": 25}}</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '{"user": {"name": "Tom", "age": 27}}' | dasel put object -p json -s '.user' -t string -t int name=Frank age=25</pre></td>
            <td><pre>{"user": {"name": "Frank", "age": 25}}</pre></td>
        </tr>
        <tr>
            <th style="text-align: center" colspan="3">Append to an array of objects</th>
        </tr>
        <tr>
            <td>jq</td>
            <td><pre>echo '{"users": [{"name": "Tom"}]}' | jq '.users += [{"name": "Frank"}]'</pre></td>
            <td><pre>{"users": [{"name": "Tom"}, {"name": "Frank"}]}</pre></td>
        </tr>
        <tr>
            <td>dasel</td>
            <td><pre>echo '{"users": [{"name": "Tom"}]}' | dasel put object -p json -s '.users[]' -t string name=Frank</pre></td>
            <td><pre>{"users": [{"name": "Tom"}, {"name": "Frank"}]}</pre></td>
        </tr>
    </tbody>
</table>

### Kubernetes
The following should work on a kubernetes deployment manifest. While kubernetes isn't for everyone, it does give some good example use-cases. 

#### Select the image for a container named `auth`
```bash
dasel select -f deployment.yaml -s "spec.template.spec.containers.(name=auth).image"
tomwright/x:v2.0.0
```

#### Change the image for a container named `auth`
```bash
dasel put string -f deployment.yaml -s "spec.template.spec.containers.(name=auth).image" "tomwright/x:v2.0.0"
```

#### Update replicas to 3
```bash
dasel put int -f deployment.yaml -s "spec.replicas" 3
```

#### Add a new env var
```bash
dasel put object -f deployment.yaml -s "spec.template.spec.containers.(name=auth).env.[]" -t string -t string name=MY_NEW_ENV_VAR value=MY_NEW_VALUE
```

#### Update an existing env var
```bash
dasel put string -f deployment.yaml -s "spec.template.spec.containers.(name=auth).env.(name=MY_NEW_ENV_VAR).value" NEW_VALUE
```