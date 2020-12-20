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

![Top level property](./benchmark/diagrams/top_level_property.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 7.8 ± 0.4 | 7.2 | 9.0 | 1.00 |
| `jq '.id' benchmark/data.json` | 26.6 ± 1.4 | 25.1 | 32.3 | 3.42 ± 0.24 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 124.2 ± 7.2 | 113.2 | 144.0 | 15.94 ± 1.19 |

### Nested property

![Nested property](./benchmark/diagrams/nested_property.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.4 ± 0.2 | 7.0 | 8.3 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.9 ± 0.9 | 24.9 | 31.2 | 3.48 ± 0.17 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 119.4 ± 4.0 | 113.3 | 140.8 | 16.06 ± 0.76 |

### Array index

![Array index](./benchmark/diagrams/array_index.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.4 ± 0.2 | 6.9 | 8.0 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 26.0 ± 1.2 | 24.7 | 32.2 | 3.52 ± 0.20 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 118.1 ± 2.6 | 113.1 | 126.8 | 16.01 ± 0.61 |

### Append to array of strings

![Append to array of strings](./benchmark/diagrams/append_array_of_strings.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.5 ± 0.5 | 5.9 | 8.8 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.2 ± 0.7 | 25.1 | 29.3 | 3.50 ± 0.26 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 119.8 ± 3.7 | 110.8 | 142.2 | 16.03 ± 1.24 |

### Update a string value

![Update a string value](./benchmark/diagrams/update_string.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.5 ± 0.3 | 7.0 | 8.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 25.9 ± 0.7 | 24.8 | 28.5 | 3.44 ± 0.16 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 129.7 ± 8.5 | 113.1 | 164.3 | 17.27 ± 1.31 |

### Overwrite an object

![Overwrite an object](./benchmark/diagrams/overwrite_object.jpg)

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.6 ± 0.7 | 7.1 | 11.4 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.2 ± 1.7 | 25.0 | 37.0 | 3.44 ± 0.40 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 122.3 ± 6.2 | 114.3 | 140.9 | 16.02 ± 1.76 |
