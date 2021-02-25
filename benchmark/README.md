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
| `dasel -f benchmark/data.json '.id'` | 7.5 ± 0.6 | 6.4 | 8.9 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.8 ± 0.7 | 24.3 | 28.2 | 3.44 ± 0.28 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 110.3 ± 3.9 | 106.2 | 133.9 | 14.73 ± 1.24 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.2 ± 0.5 | 5.8 | 9.4 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 26.9 ± 1.4 | 24.8 | 32.1 | 3.74 ± 0.33 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 122.3 ± 24.7 | 106.2 | 261.1 | 17.04 ± 3.64 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.3 ± 0.6 | 6.6 | 10.3 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.9 ± 0.9 | 24.8 | 29.9 | 3.54 ± 0.32 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 109.2 ± 2.5 | 104.0 | 118.2 | 14.92 ± 1.32 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.3 ± 0.4 | 6.4 | 8.5 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.2 ± 0.6 | 24.7 | 28.1 | 3.60 ± 0.20 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 122.2 ± 17.4 | 107.4 | 177.5 | 16.83 ± 2.55 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.0 ± 0.6 | 5.4 | 9.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 27.3 ± 3.0 | 24.2 | 41.5 | 3.89 ± 0.54 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 115.4 ± 13.3 | 104.6 | 180.9 | 16.44 ± 2.37 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 6.6 ± 1.1 | 4.5 | 9.4 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 24.6 ± 0.9 | 22.8 | 27.9 | 3.72 ± 0.65 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 112.4 ± 6.3 | 104.5 | 138.7 | 16.98 ± 3.05 |
