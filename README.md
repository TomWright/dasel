# dasel

[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/dasel)](https://goreportcard.com/report/github.com/TomWright/dasel)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Dasel (short for data-selector) allows you to query and modify data structures using selector strings.

### Installation
You can import dasel as a package and use it in your applications, or you can use a pre-built binary to modify files from the command line.

#### Import
As with any other go package, just use `go get`.
```
go get github.com/tomwright/dasel
```

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
```

##### Mac OS amd64
You may have to `brew install wget` in order for this to work.
```
curl -s https://api.github.com/repos/tomwright/dasel/releases/latest | grep browser_download_url | grep macos_amd64 | cut -d '"' -f 4 | wget -qi - && mv dasel_macos_amd64 dasel && chmod +x dasel
```

## Usage 

### Select
```
$ dasel select -h
```

The following should select the image within a kubernetes deployment manifest.
```
$ dasel select -f deployment.yaml -s "spec.template.spec.containers.(name=auth).image"
tomwright/auth:v1.0.0
```

You can also pipe data into dasel, although you then need to tell dasel which parser to use with the `-p` flag.
```
$ cat deployment.yaml | dasel select -p yaml -s "spec.template.spec.containers.(name=auth).image"
tomwright/auth:v1.0.0
```

### Put
```
$ dasel put -h
```

Basic usage is:
```
dasel put <string|int|bool|object> -f <file> -s "<selector>" <value>
```

#### Piping Data
You can pipe data both in and out of dasel.

```
$ echo "name: Tom" | ./dasel put string -p yaml -s "name" Jim
name: Jim
```

It's important to remember than if you are piping data you must provide a parser using the `-p` flag.

##### Input
Input is taken from `stdin` if you do not pass a file using the `-f` flag.

##### Output
The `select` commands will always output to `stdout`.

The default functionality for `put` commands is to edit the file in place, unless input is from stdin in which case it will be written to stdout.

You can choose a new output file by passing `-o <filepath>`. Alternatively passing `-o stdout` will result in the results being written to stdout.

#### Putting Objects
If putting an object, you can pass multiple arguments in the format of `KEY=VALUE`, each of which needs a related `-t <string|int|bool>` flag passed in the same order as the arguments.
This tells dasel which data types to parse the values as.

```
$ dasel put object -f preferences.yaml -s "my.favourites" -t string -t int colour=red number=3
```

Results in the following:

```
my:
  favourites:
    colour: red
    number: 3
```

#### Kubernetes
The following should work on a kubernetes deployment manifest. While kubernetes isn't for everyone, it does give me some good example use-cases. 

##### Change the image for a container named `auth`
```
$ dasel put string -f deployment.yaml -s "spec.template.spec.containers.(name=auth).image" "tomwright/x:v2.0.0"
```

##### Update replicas to 3
```
$ dasel put int -f deployment.yaml -s "spec.replicas" 3
```

##### Add a new env var
```
$ dasel put object -f deployment.yaml -s "spec.template.spec.containers.(name=auth).env.[]" -t string -t string name=MY_NEW_ENV_VAR value=MY_NEW_VALUE
```

##### Update an existing env var
```
$ dasel put string -f deployment.yaml -s "spec.template.spec.containers.(name=auth).env.(name=MY_NEW_ENV_VAR).value" NEW_VALUE
```

## Supported data types
Dasel attempts to find the correct parser for the given file type, but if that fails you can choose which parser to use with the `-p` or `--parser` flag. 

- JSON - `-p json`
- TOML - `-p toml`
- YAML - `-p yaml`

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

#### Next Available Index
Next available index selector is used when adding to a list of items. It allows you to append to a list.
- `colours.[]`

#### Dynamic
Dynamic selectors are used with lists when you don't know the index of the item, but instead know the value of a property of an object within the list. 

Look ups are defined in brackets. You can use multiple dynamic selectors within the same part to perform multiple checks.
```
$ dasel select -f ./tests/assets/example.yaml -s "colourCodes.(name=red).rgb"
ff0000

$ dasel select -f ./tests/assets/example.yaml -s "colourCodes.(name=blue)(rgb=0000ff)"
map[name:blue rgb:0000ff]
```
- `colourCodes.(name=red).rgb` == `ff0000`
- `colourCodes.(name=green).rgb` == `00ff00`
- `colourCodes.(name=blue).rgb` == `0000ff`
- `colourCodes.(name=blue)(rgb=0000ff).rgb` == `0000ff`

If you want to dynamically target a value in a list when it isn't a list of objects, just define the dynamic selector with `(value=<some_value>)` instead.
