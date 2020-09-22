# dasel

[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![Documentation](https://godoc.org/github.com/TomWright/dasel?status.svg)](https://godoc.org/github.com/TomWright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

## Usage 
You can import dasel as a package and use it in your applications, or you can use a pre-built binary to modify files from the command line.

### Import
As with any other go package, just use `go get`.
```
go get github.com/tomwright/dasel/cmd/dasel
```

### Command line
You can `go get` the `main` package and go should automatically build and install dasel for you.
```
go get github.com/tomwright/dasel/cmd/dasel
```

Alternatively you can download a compiled executable from the [latest release](https://github.com/TomWright/dasel/releases/latest):
This one liner should work for you - be sure to change the targeted release executable if needed. It currently targets `dasel_linux_amd64`.
```
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep linux_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_linux_amd64 dasel && chmod +x dasel
```

## Support data types
Dasel attempts to find the correct parser for the given file type, but if that fails you can choose which parser to use with the `-p` or `--parser` flag. 

- YAML - `-p yaml`
- JSON - `-p json`

## Selectors

Selectors are used to define a path through a set of data. This path is usually defined as a chain of nodes.

A selector is made up of different parts separated by a dot `.`, each part being used to identify the next node in the chain.

The following YAML data structure will be used as a reference in the following examples.
```
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
```
$ dasel select -f ./tests/assets/example.yaml -s "name"
Tom
```
- `name` == `Tom`

### Child Elements
Just separate the child element from the parent element using a `.`:
```
$ dasel select -f ./tests/assets/example.yaml -s "preferences.favouriteColour"
red
```
- `preferences.favouriteColour` == `red`

#### Index
When you have a list, you can use square brackets to access a specific item in the list by its index.
```
$ dasel select -f ./tests/assets/example.yaml -s "colours.[1]"
green
```
- `colours.[0]` == `red`
- `colours.[1]` == `green`
- `colours.[2]` == `blue`

#### Next Available Index - WIP
Next available index selector is used when adding to a list of items. It allows you to append to a list.
- `colours.[]`

#### Dynamic
Dynamic selectors are used with lists when you don't know the index of the item, but instead know the value of a property of an object within the list. 

Look ups are defined in brackets. You can use multiple dynamic selectors within the same part to perform multiple checks.
```
$ dasel select -f ./tests/assets/example.yaml -s ".colourCodes.(name=red).rgb"
ff0000

$ dasel select -f ./tests/assets/example.yaml -s ".colourCodes.(name=blue)(rgb=0000ff)"
map[name:blue rgb:0000ff]
```
- `.colourCodes.(name=red).rgb` == `ff0000`
- `.colourCodes.(name=green).rgb` == `00ff00`
- `.colourCodes.(name=blue).rgb` == `0000ff`
- `.colourCodes.(name=blue)(rgb=0000ff).rgb` == `0000ff`