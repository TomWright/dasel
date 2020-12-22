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
| `dasel -f benchmark/data.json '.id'` | 7.7 ± 0.3 | 7.2 | 8.5 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.8 ± 0.8 | 24.9 | 29.2 | 3.36 ± 0.15 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 125.1 ± 3.6 | 119.8 | 141.9 | 16.24 ± 0.72 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.8 ± 0.3 | 7.1 | 9.0 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 26.0 ± 1.0 | 24.8 | 30.6 | 3.35 ± 0.19 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 124.5 ± 3.3 | 114.8 | 140.8 | 16.05 ± 0.81 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.7 ± 0.3 | 7.1 | 8.7 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.8 ± 0.7 | 25.0 | 29.5 | 3.35 ± 0.15 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 124.7 ± 2.7 | 118.1 | 131.2 | 16.20 ± 0.68 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.7 ± 0.2 | 7.3 | 8.1 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 25.9 ± 0.5 | 25.1 | 27.7 | 3.37 ± 0.11 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 125.9 ± 3.8 | 118.5 | 150.4 | 16.37 ± 0.64 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.8 ± 0.3 | 7.3 | 9.2 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.2 ± 1.1 | 25.0 | 30.7 | 3.36 ± 0.20 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 125.6 ± 4.1 | 111.6 | 145.6 | 16.11 ± 0.85 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 8.2 ± 0.5 | 7.4 | 9.9 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.1 ± 0.7 | 24.9 | 29.0 | 3.18 ± 0.22 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 126.5 ± 3.8 | 119.8 | 144.2 | 15.43 ± 1.08 |
