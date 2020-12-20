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

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 7.7 ± 0.3 | 7.1 | 8.7 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.9 ± 0.8 | 24.9 | 31.5 | 3.37 ± 0.17 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 118.3 ± 3.0 | 110.3 | 127.1 | 15.40 ± 0.72 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.5 ± 0.3 | 6.9 | 9.0 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.9 ± 0.8 | 24.7 | 29.5 | 3.47 ± 0.19 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 118.3 ± 2.5 | 113.3 | 125.2 | 15.86 ± 0.81 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.4 ± 0.3 | 6.9 | 8.4 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.8 ± 0.8 | 24.8 | 29.5 | 3.51 ± 0.17 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 120.2 ± 3.6 | 110.0 | 136.9 | 16.35 ± 0.80 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.5 ± 0.2 | 7.1 | 8.1 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 25.9 ± 0.8 | 24.9 | 29.7 | 3.44 ± 0.14 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 118.9 ± 2.9 | 112.1 | 131.0 | 15.83 ± 0.61 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.5 ± 0.3 | 7.0 | 9.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.1 ± 1.6 | 24.8 | 35.9 | 3.47 ± 0.27 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 119.1 ± 2.7 | 113.7 | 128.7 | 15.82 ± 0.82 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="200"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.5 ± 0.4 | 7.0 | 9.7 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 25.7 ± 0.6 | 24.8 | 28.7 | 3.45 ± 0.20 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 122.5 ± 3.8 | 113.6 | 135.4 | 16.43 ± 1.01 |
