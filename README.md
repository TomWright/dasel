# dasel

Read and modify data structures using selectors.

## Selectors
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

### Root Element
Just use the root element name as a string.
- `name` == `Tom`

### Child Element
Just separate the parent element from the parent element using a `.`:
- `preferences.favouriteColour` == `red`

#### Index
When you have a list, you can use square brackets to access or modify a specific item.
- `colours.[0]` == `red`
- `colours.[1]` == `green`
- `colours.[2]` == `blue`

#### Next Available Index
Next available index selector is used when adding to a list of items.
- `colours.[]`

#### Look up
Look ups are defined in brackets and allow you to dynamically select an object to use.
- `.colourCodes.(name=red).rgb` == `ff0000`
- `.colourCodes.(name=green).rgb` == `00ff00`
- `.colourCodes.(name=blue).rgb` == `0000ff`
- `.colourCodes.(name=blue)(rgb=0000ff).rgb` == `0000ff`