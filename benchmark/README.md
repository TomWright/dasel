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
| `dasel -f benchmark/data.json` | 11.0 ± 1.9 | 9.4 | 22.9 | 1.00 |
| `jq '.' benchmark/data.json` | 26.0 ± 1.0 | 24.8 | 31.7 | 2.36 ± 0.42 |
| `yq --yaml-output '.' benchmark/data.yaml` | 114.9 ± 4.0 | 110.1 | 139.7 | 10.41 ± 1.86 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 10.3 ± 0.8 | 9.2 | 14.1 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.8 ± 1.1 | 24.5 | 31.6 | 2.51 ± 0.23 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 113.8 ± 3.2 | 107.0 | 123.4 | 11.09 ± 0.93 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 9.7 ± 0.7 | 9.0 | 13.2 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.7 ± 0.9 | 24.2 | 30.7 | 2.65 ± 0.21 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 113.8 ± 3.1 | 108.3 | 127.1 | 11.73 ± 0.88 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 10.2 ± 0.6 | 9.1 | 12.4 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.9 ± 0.9 | 24.1 | 31.3 | 2.54 ± 0.17 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 116.2 ± 8.1 | 108.4 | 153.9 | 11.40 ± 1.03 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 10.7 ± 2.6 | 9.2 | 28.1 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 29.7 ± 7.0 | 24.9 | 81.5 | 2.77 ± 0.93 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 113.0 ± 3.6 | 107.8 | 131.3 | 10.51 ± 2.53 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 10.3 ± 1.3 | 9.1 | 21.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.3 ± 1.9 | 24.6 | 36.7 | 2.55 ± 0.36 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 116.8 ± 9.1 | 109.4 | 153.2 | 11.32 ± 1.66 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 10.4 ± 2.0 | 9.2 | 28.8 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 25.9 ± 0.7 | 24.9 | 28.6 | 2.49 ± 0.48 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 114.8 ± 6.0 | 107.3 | 142.4 | 11.01 ± 2.19 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json -m '.-'` | 10.1 ± 1.9 | 9.2 | 28.7 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 25.7 ± 0.7 | 24.5 | 27.8 | 2.55 ± 0.49 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 121.6 ± 36.7 | 108.2 | 399.6 | 12.06 ± 4.31 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel delete -f benchmark/data.json -o - '.id'` | 10.6 ± 1.6 | 9.1 | 21.3 | 1.00 |
| `jq 'del(.id)' benchmark/data.json` | 25.8 ± 0.7 | 24.4 | 28.0 | 2.44 ± 0.38 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 114.1 ± 3.0 | 109.0 | 122.6 | 10.79 ± 1.69 |
