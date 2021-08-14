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

### Root Object

<img src="diagrams/root_object.jpg" alt="Root Object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json` | 10.7 ± 1.8 | 9.5 | 27.2 | 1.00 |
| `jq '.' benchmark/data.json` | 26.3 ± 1.2 | 24.9 | 33.3 | 2.46 ± 0.43 |
| `yq --yaml-output '.' benchmark/data.yaml` | 127.6 ± 20.9 | 108.0 | 273.2 | 11.94 ± 2.80 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 10.0 ± 2.0 | 9.1 | 28.6 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.8 ± 0.8 | 24.2 | 30.2 | 2.58 ± 0.52 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 111.4 ± 2.7 | 107.1 | 119.6 | 11.12 ± 2.21 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 9.7 ± 0.3 | 9.1 | 11.0 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.8 ± 1.0 | 24.4 | 29.9 | 2.66 ± 0.14 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 112.5 ± 4.4 | 107.1 | 135.4 | 11.60 ± 0.62 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 9.7 ± 0.4 | 9.1 | 11.5 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 26.4 ± 3.6 | 24.4 | 50.4 | 2.73 ± 0.38 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 111.8 ± 2.2 | 107.2 | 116.5 | 11.55 ± 0.50 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 9.8 ± 0.4 | 9.3 | 11.4 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 25.9 ± 0.7 | 24.8 | 28.8 | 2.65 ± 0.12 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 125.3 ± 23.8 | 108.6 | 272.1 | 12.83 ± 2.48 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 9.4 ± 0.7 | 8.2 | 11.1 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 25.1 ± 0.8 | 23.6 | 28.8 | 2.66 ± 0.22 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 111.8 ± 2.3 | 107.7 | 121.9 | 11.84 ± 0.95 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 11.5 ± 2.5 | 9.5 | 20.4 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.4 ± 2.3 | 24.5 | 40.8 | 2.30 ± 0.53 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 114.3 ± 2.9 | 108.2 | 128.0 | 9.95 ± 2.14 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json -m '.-'` | 9.9 ± 0.4 | 9.3 | 11.2 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 25.8 ± 0.8 | 24.8 | 29.0 | 2.60 ± 0.13 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 113.8 ± 3.4 | 107.7 | 132.7 | 11.46 ± 0.57 |
