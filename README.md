# dasel

[![Gitbook](https://badges.aleen42.com/src/gitbook_1.svg)](https://tomwright.gitbook.io/dasel/)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub All Releases Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)](https://github.com/TomWright/dasel/releases/latest)
[![Homebrew tag (latest by date)](https://img.shields.io/homebrew/v/dasel)](https://formulae.brew.sh/formula/dasel)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

Comparable to [jq](https://github.com/stedolan/jq) / [yq](https://github.com/kislyuk/yq), but supports JSON, YAML, TOML, XML and CSV with zero runtime dependencies.

## One tool to rule them all

Say good bye to learning new tools just to work with a different data format.

Dasel uses a standard selector syntax no matter the data format. This means that once you learn how to use dasel you immediately have the ability to query/modify any of the supported data types without any additional tools or effort. 

![Update Kubernetes Manifest](update_kubernetes.gif)

### Issue vs Discussion

I have enabled [discussions](https://github.com/TomWright/dasel/discussions) on this repository.

I am aware there may be some confusion when deciding where you should communicate when reporting issues, asking questions or raising feature requests so this section aims to help us align on that.

Please [raise an issue](https://github.com/TomWright/dasel/issues) if:
- You find a bug.
- You have a feature request and can clearly describe your request.

Please [open a discussion](https://github.com/TomWright/dasel/discussions) if:
- You have a question.
- You're not sure how to achieve something with dasel.
- You have an idea but don't quite know how you would like it to work.
- You have achieved something cool with dasel and want to show it off.
- Anything else!

## Features
- [Query/select data from structured data files](#select).
- [Update data in structured data files](#put).
- [Create data files](#creating-properties).
- [Supports multiple data formats/types](#supported-file-types).
- [Convert between data formats/types](#converting-between-formats).
- Uses a [standard query/selector syntax](#selectors) across all data formats.
- Zero runtime dependencies.
- [Available on Linux, Mac and Windows](#binary-on-release).
- Available to [import and use in your own projects](#using-dasel-as-a-package).
- [Run via Docker](#docker).
- [Faster than jq/yq](#benchmarks).
- [Self update](#self-update).

## Table of contents
* [Dasel](#dasel)
* [One tool to rule them all](#one-tool-to-rule-them-all)
* [Issue vs discussion](#issue-vs-discussion)
* [Features](#features)
* [Table of contents](#table-of-contents)
* [Documentation](#documentation)
* [Playground](#playground)
* [Usage](#usage)
  * [Select](#select)
  * [Put](#put)
  * [Put Object](#put-object)
  * [Put Document](#put-document)
* [Examples](#examples)
  * [General](#general)
    * [Filter JSON API results](#filter-json-api-results)
  * [jq to dasel](#jq-to-dasel)
  * [yq to dasel](#yq-to-dasel)
  * [Kubernetes](#kubernetes)
  * [XML](#xml-examples)
* [Benchmarks](#benchmarks)

## Documentation

The official dasel docs can be found at [tomwright.gitbook.io/dasel/](https://tomwright.gitbook.io/dasel/).

## Playground

You can test out dasel commands using the [playground](https://dasel.tomwright.me).

Source code for the playground can be found at  [github.com/TomWright/daselplayground](https://github.com/TomWright/daselplayground).

## Examples

### General

#### Filter JSON API results

The following line will return the download URL for the latest macOS dasel release:

```bash
curl https://api.github.com/repos/tomwright/dasel/releases/latest | dasel -p json --plain '.assets.(name=dasel_darwin_amd64).browser_download_url'
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

echo '{"users": [{"name": "Tom"}]}' | dasel put object -p json -s '.users.[]' -t string name=Frank
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
- name: Tom' | dasel put object -p yaml -s '.users.[]' -t string name=Frank
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

## Benchmarks

In my tests dasel has been up to 3x faster than jq and 15x faster than yq.

See the [benchmark directory](./benchmark/README.md).
