# Benchmarks

These benchmarks are auto generated using `./benchmark/run.sh`.

```
brew install hyperfine
pip install matplotlib
./benchmark/run.sh
```

I have put together what I believe to be equivalent commands in dasel/jq/yq.

If you have any feedback or wish to add new benchmarks please submit a PR.
## Benchmarks

### Top level property

![Top level property](diagrams/top_level_property.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 8.4 ± 0.5 | 7.5 | 9.8 | 1.00 |
| `jq '.id' benchmark/data.json` | 28.4 ± 3.5 | 25.4 | 40.6 | 3.40 ± 0.46 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 134.0 ± 11.5 | 113.3 | 175.9 | 16.04 ± 1.66 |

### Nested property

![Nested property](diagrams/nested_property.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.5 ± 0.4 | 6.7 | 8.9 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 30.8 ± 3.3 | 25.0 | 37.4 | 4.10 ± 0.50 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 119.5 ± 7.3 | 111.4 | 149.1 | 15.91 ± 1.36 |

### Array index

![Array index](diagrams/array_index.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.6 ± 0.4 | 6.9 | 8.8 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.6 ± 0.7 | 24.6 | 29.1 | 3.38 ± 0.22 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 121.4 ± 8.2 | 110.6 | 147.7 | 16.03 ± 1.42 |

### Append to array of strings

![Append to array of strings](diagrams/append_array_of_strings.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.5 ± 0.2 | 7.1 | 8.5 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.2 ± 1.5 | 24.9 | 35.6 | 3.50 ± 0.22 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 123.8 ± 8.1 | 111.0 | 143.1 | 16.53 ± 1.18 |

### Update a string value

![Update a string value](diagrams/update_string.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.7 ± 0.4 | 7.1 | 8.6 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.8 ± 1.6 | 25.3 | 33.4 | 3.49 ± 0.26 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 125.2 ± 6.5 | 112.7 | 141.4 | 16.28 ± 1.14 |

### Overwrite an object

![Overwrite an object](diagrams/overwrite_object.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 8.3 ± 0.6 | 7.3 | 9.9 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 27.1 ± 2.1 | 25.1 | 32.5 | 3.28 ± 0.35 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 125.5 ± 8.3 | 113.7 | 146.2 | 15.20 ± 1.50 |
