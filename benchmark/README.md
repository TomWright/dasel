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
| `dasel -f benchmark/data.json` | 6.4 ± 0.3 | 6.0 | 7.8 | 1.00 |
| `jq '.' benchmark/data.json` | 28.0 ± 0.9 | 27.2 | 31.8 | 4.38 ± 0.22 |
| `yq --yaml-output '.' benchmark/data.yaml` | 127.1 ± 2.6 | 124.0 | 141.2 | 19.87 ± 0.90 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'id'` | 6.4 ± 0.4 | 5.9 | 8.3 | 1.00 |
| `jq '.id' benchmark/data.json` | 28.0 ± 0.8 | 27.1 | 31.3 | 4.35 ± 0.29 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 126.2 ± 2.2 | 123.3 | 136.6 | 19.60 ± 1.21 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'user.name.first'` | 6.4 ± 0.3 | 6.0 | 7.9 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 28.4 ± 1.2 | 27.1 | 36.0 | 4.41 ± 0.28 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 126.3 ± 1.8 | 123.4 | 131.2 | 19.63 ± 0.94 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'favouriteNumbers.[1]'` | 6.4 ± 0.3 | 6.0 | 7.2 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 28.2 ± 1.0 | 27.2 | 32.1 | 4.44 ± 0.24 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 126.6 ± 2.7 | 123.6 | 141.1 | 19.91 ± 0.92 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[]'` | 6.4 ± 0.2 | 6.0 | 7.3 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 28.4 ± 1.0 | 27.4 | 32.0 | 4.46 ± 0.21 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 127.3 ± 2.7 | 124.3 | 142.0 | 20.00 ± 0.80 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 6.6 ± 0.9 | 6.1 | 15.3 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 28.3 ± 1.0 | 27.2 | 32.4 | 4.31 ± 0.63 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 127.3 ± 1.9 | 124.6 | 131.9 | 19.41 ± 2.77 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' '.user.name'` | 6.4 ± 0.3 | 5.9 | 7.7 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 28.2 ± 0.9 | 27.3 | 31.6 | 4.39 ± 0.22 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 130.9 ± 13.4 | 124.4 | 229.2 | 20.38 ± 2.24 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'all().key()'` | 6.5 ± 0.4 | 6.0 | 8.4 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 28.1 ± 0.8 | 27.1 | 31.0 | 4.34 ± 0.28 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 127.2 ± 5.4 | 123.2 | 171.9 | 19.63 ± 1.42 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel delete -f benchmark/data.json -o - 'id'` | 6.3 ± 0.2 | 5.9 | 7.2 | 1.00 |
| `jq 'del(.id)' benchmark/data.json` | 28.2 ± 1.0 | 27.1 | 31.4 | 4.45 ± 0.22 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 127.0 ± 2.1 | 124.4 | 134.9 | 20.01 ± 0.73 |
