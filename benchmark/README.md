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

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 7.7 ± 0.3 | 7.2 | 8.4 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.8 ± 0.5 | 24.9 | 27.9 | 3.34 ± 0.15 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 120.1 ± 3.0 | 112.9 | 127.9 | 15.55 ± 0.70 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.4 ± 0.3 | 6.8 | 8.7 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.8 ± 0.9 | 24.7 | 29.4 | 3.48 ± 0.19 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 121.1 ± 3.6 | 114.3 | 130.8 | 16.37 ± 0.86 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.4 ± 0.2 | 7.0 | 8.3 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.7 ± 0.7 | 24.8 | 29.4 | 3.49 ± 0.15 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 118.6 ± 3.9 | 109.8 | 141.5 | 16.09 ± 0.75 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.5 ± 0.3 | 7.0 | 9.0 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.1 ± 1.3 | 25.1 | 34.5 | 3.48 ± 0.22 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 119.7 ± 3.3 | 112.6 | 128.7 | 15.99 ± 0.72 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.5 ± 0.2 | 7.1 | 8.5 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.2 ± 1.7 | 24.9 | 36.1 | 3.51 ± 0.25 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 119.6 ± 3.2 | 111.4 | 127.8 | 16.04 ± 0.68 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.6 ± 0.2 | 7.1 | 8.6 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 25.9 ± 0.7 | 25.0 | 29.4 | 3.41 ± 0.14 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 121.5 ± 3.7 | 116.0 | 142.1 | 15.98 ± 0.68 |
