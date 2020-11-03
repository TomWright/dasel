# dasel

[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub All Releases Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

Comparable to [jq](https://github.com/stedolan/jq) / [yq](https://github.com/kislyuk/yq), but supports JSON, YAML, TOML and XML with zero runtime dependencies.

![Update Kubernetes Manifest](update_kubernetes.gif)

## Features
- [Query/select data from structured data files](#select).
- [Update data in structured data files](#put).
- [Supports multiple data formats/types](#supported-file-types).
- Uses a [standard query/selector syntax](#selectors) across all data formats.
- Zero runtime dependencies.
- Available on Linux, Mac and Windows.
- Available to import and use in your own projects.

## Table of contents
* [Dasel](#dasel)
* [Features](#features)
* [Installation](#installation)
* [Notes](#notes)
  * [Preserved formatting and ordering](#preserved-formatting-and-ordering)
  * [Memory Usage](#memory-usage)
* [Usage](#usage)
  * [Select](#select)
  * [Put](#put)
  * [Put Object](#put-object)
* [Supported file types](#supported-file-types)
* [Selectors](#selectors)
  * [Property](#property)
  * [Child](#child-elements)
  * [Index](#index)
  * [Next available index](#next-available-index)
  * [Dynamic](#dynamic)
    * [Using queries in dynamic selectors](#using-queries-in-dynamic-selectors)
* [Examples](#examples)
  * [General](#general)
    * [Filter JSON API results](#filter-json-api-results)
  * [jq to dasel](#jq-to-dasel)
  * [yq to dasel](#yq-to-dasel)
  * [Kubernetes](#kubernetes)
  * [XML](#xml-examples)

## Installation
You can import dasel as a package and use it in your applications, or you can use a pre-built binary to modify files from the command line.

### Command line

#### Go
You can `go get` the `main` package and go should automatically build and install dasel for you.
```bash
go get github.com/tomwright/dasel/cmd/dasel
```

#### Binary on release
You can download a compiled executable from the [latest release](https://github.com/TomWright/dasel/releases/latest).

##### Linux amd64
This one liner should work for you - be sure to change the targeted release executable if needed. It currently targets `dasel_linux_amd64`.
```bash
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep linux_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_linux_amd64 dasel && chmod +x dasel
mv ./dasel /usr/local/bin/dasel
```

##### Mac OS amd64
You may have to `brew install wget` in order for this to work.
```bash
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep macos_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_macos_amd64 dasel && chmod +x dasel
mv ./dasel /usr/local/bin/dasel
```

#### Docker
You also have the option of using the docker image to run dasel for you.
```bash
echo '{"name": "Tom"}' | docker run -i --rm ghcr.io/tomwright/dasel:latest -p json '.name'
"Tom"
```

If you want to use a specific version of dasel simply change `latest` to the desired version.

- `latest` - The latest released version.
- `dev` - The latest build from `master` branch.
- `v*.*.*` - The build from the given release.

### Go get
As with any other go package, just use `go get`.
```bash
go get github.com/tomwright/dasel
```

The documentation for this is still a WIP. Please raise an issue if you have a specific need for this and I'll do my best to help out.

## Notes

### Preserved formatting and ordering

The formatting of files can be changed while being processed. Dasel itself doesn't make these changes, rather the act of marshaling the results.

In short, the output files may have properties in a different order but the actual contents will be as expected.

### Memory usage

Dasel's method of querying data requires that the entire input document is stored in memory.

You should keep this in mind as the maximum filesize it can process will be limited by your system's available resources (specifically RAM).

## Usage 

```bash
dasel -h
```

An important note is that if no sub-command is given, dasel will default to `select`.

### Select
```bash
dasel select -f <file> -p <parser> <selector>
```

#### Arguments

##### `-f`, `--file`

Specify the file to query. This is required unless you are piping in data.

If piping in data you can optionally pass `-f stdin`/`-f -`.

##### `-p`, `--parser`

Specify the parser to use when reading the file.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

See [supported file types](#supported-file-types).

##### `-s`, `--selector`, `<selector>`

Specify the selector to use. See [Selectors](#selectors) for more information.

If no selector flag is given, dasel assumes the first argument given is the selector.

This is required.

##### `--plain`

By default, dasel formats the output using the specified parser.

If this flag is used no formatting occurs and the results output as a string.

#### Example

##### Select the image within a kubernetes deployment manifest file:
```bash
dasel select -f deployment.yaml "spec.template.spec.containers.(name=auth).image"
"tomwright/auth:v1.0.0"
```

##### Piping data into the select:
```bash
cat deployment.yaml | dasel select -p yaml "spec.template.spec.containers.(name=auth).image"
"tomwright/auth:v1.0.0"
```

### Put
```bash
dasel put <type> -f <file> -o <out> -p <parser> <selector> <value>
```

```bash
echo "name: Tom" | ./dasel put string -p yaml "name" Jim
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

If piping in data you can optionally pass `-f stdin`/`-f -`.

##### `-o`, `--out`

Specify the output file. If present, results will be written to the given file. If not present, results will be written to the input file (or stdout if none given).

To force output to be written to stdout, pass `-o stdout`/`-o -`.

##### `-p`, `--parser`

Specify the parser to use when reading/writing the input/output files.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

See [supported file types](#supported-file-types).

##### `-s`, `--selector`, `<selector>`

Specify the selector to use. See [Selectors](#selectors) for more information.

If no selector flag is given, dasel assumes the first argument given is the selector.

This is required.

##### `value`

The value to write.

Dasel will parse this value as a string, int, or bool from this value depending on the given `type`.

This is required.

### Put Object

Putting objects works slightly differently to a standard put, but the same principles apply.

```bash
dasel put object -f <file> -o <out> -p <parser> -t <type> <selector> <values>
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

If piping in data you can optionally pass `-f stdin`/`-f -`.

##### `-o`, `--out`

Specify the output file. If present, results will be written to the given file. If not present, results will be written to the input file (or stdout if none given).

To force output to be written to stdout, pass `-o stdout`/`-o -`.

##### `-p`, `--parser`

Specify the parser to use when reading/writing the input/output files.

This is required if you are piping in data, otherwise dasel will use the given file extension to guess which parser to use.

See [supported file types](#supported-file-types).

##### `-s`, `--selector`, `<selector>`

Specify the selector to use. See [Selectors](#selectors) for more information.

If no selector flag is given, dasel assumes the first argument given is the selector.

This is required.

##### `values`

A space-separated list of `key=value` pairs.

Dasel will parse each value as a string, int, or bool depending on the related `type`.

This is required.

#### Example

```bash
echo "" | dasel put object -p yaml -t string -t int "my.favourites" colour=red number=3
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

### JSON
```bash
-p json
```
Using [golang.org/pkg/encoding/json](https://golang.org/pkg/encoding/json/).

### TOML
```bash
-p toml
```
Using [github.com/pelletier/go-toml](https://github.com/pelletier/go-toml).

### YAML
```bash
-p yaml
```
Using [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2).

#### Multi-document files

Multi-document files are not supported. This is due to [limitations in the yaml parser used](https://github.com/go-yaml/yaml/tree/v2.3.0#compatibility).

### XML
```bash
-p xml
```
Using [github.com/clbanning/mxj](https://github.com/clbanning/mxj).

#### Arrays/Lists

Due to the way that XML is decoded, dasel can only detect something as a list if there are at least 2 items.

If you try to use list selectors (dynamic, index, append) when there are less than 2 items in the list you will get an error.

There are no plans to introduce a workaround for this but if there is enough demand it may be worked on in the future.

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
The next available index selector is used when adding to a list of items. It allows you to append to a list.
- `colours.[]`

#### Dynamic
Dynamic selectors are used with lists when you don't know the index of the item, but instead want to find the index based on some other criteria.
 
Dasel currently supports `key/query=value` checks but I aim to support more check types in the future.

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

If you want to refer to the value of a non-object value in a list, just define the key as `value` or `.`, meaning the current value. This may look something like `(value=2)`.

##### Using queries in dynamic selectors
When performing a check dasel creates a new root node at the current position and then selects data using the given key as the query.

This allows you to perform complex queries such as...

```bash
echo `{
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
}` | dasel -p json '.users.(.addresses.(.primary=true).number=123).name.first'
"Tom"
```

The above query in plain English may read as...

> Give me the first name of the user
> who's primary address is at number 123

The resolution of that query looks something like this:
```
.users.(.addresses.(.primary=true).number=123).name.first
.users.(.addresses.[0].number=123).name.first
.users.[0].name.first
```

## Examples

### General

#### Filter JSON API results

The following line will return the download URL for the latest macOS dasel release:

```bash
curl https://api.github.com/repos/tomwright/dasel/releases/latest | dasel -p json --plain '.assets.(name=dasel_macos_amd64).browser_download_url'
```

### jq to dasel

The follow examples show a set of [jq](https://github.com/stedolan/jq) commands and the equivalent in dasel.

#### Select a single value

```bash
echo '{"name": "Tom"}' | jq '.name'
"Tom"

echo '{"name": "Tom"}' | dasel -p json '.name'
"Tom"
```

#### Select a nested value

```bash
echo '{"user": {"name": "Tom", "age": 27}}' | jq '.user.age'
27

echo '{"user": {"name": "Tom", "age": 27}}' | dasel -p json '.user.age'
27
```

#### Select an array index

```bash
echo '[1, 2, 3]' | jq '.[1]'
2

echo '[1, 2, 3]' | dasel -p json '.[1]'
2
```

#### Append to an array of strings

```bash
echo '["a", "b", "c"]' | jq '. += ["d"]'
[
  "a",
  "b",
  "c",
  "d"
]

echo '["a", "b", "c"]' | dasel put string -p json -s '.[]' d
[
  "a",
  "b",
  "c",
  "d"
]
```

#### Update a string value

```bash
echo '["a", "b", "c"]' | jq '.[1] = "d"'
[
  "a",
  "d",
  "c"
]

echo '["a", "b", "c"]' | dasel put string -p json -s '.[1]' d
[
  "a",
  "d",
  "c"
]
```

#### Update an int value

```bash
echo '[1, 2, 3]' | jq '.[1] = 5'
[
  1,
  5,
  3
]

echo '[1, 2, 3]' | dasel put int -p json -s '.[1]' 5
[
  1,
  5,
  3
]
```

#### Overwrite an object

```bash
echo '{"user": {"name": "Tom", "age": 27}}' | jq '.user = {"name": "Frank", "age": 25}'
{
  "user": {
    "name": "Frank",
    "age": 25
  }
}

echo '{"user": {"name": "Tom", "age": 27}}' | dasel put object -p json -s '.user' -t string -t int name=Frank age=25
{
  "user": {
    "age": 25,
    "name": "Frank"
  }
}
```

#### Append to an array of objects

```bash
echo '{"users": [{"name": "Tom"}]}' | jq '.users += [{"name": "Frank"}]'
{
  "users": [
    {
      "name": "Tom"
    },
    {
      "name": "Frank"
    }
  ]
}

echo '{"users": [{"name": "Tom"}]}' | dasel put object -p json -s '.users[]' -t string name=Frank
{
  "users": [
    {
      "name": "Tom"
    },
    {
      "name": "Frank"
    }
  ]
}
```

### yq to dasel

The follow examples show a set of [yq](https://github.com/kislyuk/yq) commands and the equivalent in dasel.

#### Select a single value

```bash
echo 'name: Tom' | yq '.name'
"Tom"

echo 'name: Tom' | dasel -p yaml '.name'
Tom
```

#### Select a nested value

```bash
echo 'user:
  name: Tom
  age: 27' | yq '.user.age'
27

echo 'user:
       name: Tom
       age: 27' | dasel -p yaml '.user.age'
27
```

#### Select an array index

```bash
echo '- 1
- 2
- 3' | yq '.[1]'
2

echo '- 1
- 2
- 3' | dasel -p yaml '.[1]'
2
```

#### Append to an array of strings

```bash
echo '- a
- b
- c' | yq --yaml-output '. += ["d"]'
- a
- b
- c
- d

echo '- a
- b
- c' | dasel put string -p yaml -s '.[]' d
- a
- b
- c
- d

```

#### Update a string value

```bash
echo '- a
- b
- c' | yq --yaml-output '.[1] = "d"'
- a
- d
- c

echo '- a
- b
- c' | dasel put string -p yaml -s '.[1]' d
- a
- d
- c
```

#### Update an int value

```bash
echo '- 1
- 2
- 3' | yq --yaml-output '.[1] = 5'
- 1
- 5
- 3

echo '- 1
- 2
- 3' | dasel put int -p yaml -s '.[1]' 5
- 1
- 5
- 3
```

#### Overwrite an object

```bash
echo 'user:
  name: Tom
  age: 27' | yq --yaml-output '.user = {"name": "Frank", "age": 25}'
user:
  name: Frank
  age: 25


echo 'user:
  name: Tom
  age: 27' | dasel put object -p yaml -s '.user' -t string -t int name=Frank age=25
user:
  age: 25
  name: Frank
```

#### Append to an array of objects

```bash
echo 'users:
- name: Tom' | yq --yaml-output '.users += [{"name": "Frank"}]'
users:
  - name: Tom
  - name: Frank


echo 'users:
- name: Tom' | dasel put object -p yaml -s '.users[]' -t string name=Frank
users:
- name: Tom
- name: Frank
```

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

### XML Examples

XML has some slight differences (such as attributes) that should be documented.

#### Query attributes

Decoded attributes are set as properties on the related object with a prefix of `-`.

```bash
echo '<data>
    <users primary="true">
        <name>Tom</name>
    </users>
    <users primary="false">
        <name>Frank</name>
    </users>
</data>' | dasel -p xml '.data.users[0].-primary'
true
```

#### Filtering on attributes

We can also filter on attributes since they are defined against the related object.

```bash
echo '<data>
    <users primary="true">
        <name>Tom</name>
    </users>
    <users primary="false">
        <name>Frank</name>
    </users>
</data>' | dasel -p xml '.data.users.(-primary=true).name'
Tom
``` 
