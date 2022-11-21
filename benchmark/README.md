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
| `daselv2 -f benchmark/data.json` | 7.6 ± 1.9 | 6.2 | 18.6 | 1.00 |
| `dasel -f benchmark/data.json` | 9.1 ± 1.2 | 7.9 | 14.2 | 1.19 ± 0.33 |
| `jq '.' benchmark/data.json` | 28.3 ± 1.1 | 27.1 | 32.9 | 3.73 ± 0.93 |
| `yq --yaml-output '.' benchmark/data.yaml` | 128.9 ± 2.9 | 125.6 | 144.1 | 16.95 ± 4.20 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'id'` | 8.8 ± 4.5 | 6.1 | 41.9 | 1.03 ± 0.54 |
| `dasel -f benchmark/data.json '.id'` | 8.5 ± 0.7 | 7.7 | 11.1 | 1.00 |
| `jq '.id' benchmark/data.json` | 28.1 ± 0.8 | 27.1 | 31.4 | 3.30 ± 0.28 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 127.6 ± 2.4 | 124.4 | 136.5 | 15.00 ± 1.24 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'user.name.first'` | 7.0 ± 0.7 | 6.2 | 9.5 | 1.00 |
| `dasel -f benchmark/data.json '.user.name.first'` | 8.8 ± 0.7 | 7.9 | 10.8 | 1.25 ± 0.16 |
| `jq '.user.name.first' benchmark/data.json` | 28.4 ± 1.1 | 27.4 | 32.4 | 4.04 ± 0.42 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 127.5 ± 2.5 | 124.6 | 136.0 | 18.13 ± 1.80 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'favouriteNumbers.[1]'` | 6.9 ± 0.5 | 6.2 | 8.7 | 1.00 |
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 8.4 ± 0.5 | 7.5 | 10.1 | 1.23 ± 0.12 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 28.2 ± 0.8 | 27.2 | 31.9 | 4.11 ± 0.33 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 130.0 ± 11.6 | 124.4 | 231.5 | 18.97 ± 2.20 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[]'` | 6.9 ± 0.6 | 6.2 | 8.7 | 1.00 |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 8.6 ± 0.7 | 7.4 | 11.0 | 1.25 ± 0.15 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 29.4 ± 6.8 | 27.2 | 92.8 | 4.29 ± 1.06 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 130.4 ± 8.2 | 124.9 | 190.8 | 19.01 ± 1.95 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 7.3 ± 1.3 | 6.1 | 12.1 | 1.00 |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 8.5 ± 0.5 | 7.7 | 10.0 | 1.15 ± 0.21 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 28.3 ± 1.0 | 26.9 | 32.1 | 3.86 ± 0.69 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 130.2 ± 10.7 | 124.8 | 206.9 | 17.72 ± 3.43 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' 'user.name'` | 7.1 ± 0.9 | 6.2 | 10.5 | 1.00 |
| `dasel put document -f benchmark/data.json -o - -d json '.user.name' '{"first":"Frank","last":"Jones"}'` | 8.5 ± 0.7 | 7.8 | 11.6 | 1.20 ± 0.18 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 28.2 ± 1.0 | 27.2 | 32.7 | 3.98 ± 0.52 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 129.0 ± 4.0 | 125.5 | 161.9 | 18.18 ± 2.37 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'all().key()'` | 6.9 ± 0.6 | 6.1 | 9.5 | 1.00 |
| `dasel -f benchmark/data.json -m '.-'` | 8.4 ± 0.5 | 7.6 | 10.2 | 1.21 ± 0.13 |
| `jq 'keys[]' benchmark/data.json` | 28.1 ± 0.9 | 27.0 | 31.8 | 4.06 ± 0.38 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 131.5 ± 12.3 | 124.5 | 184.9 | 19.01 ± 2.46 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 delete -f benchmark/data.json -o - 'id'` | 6.8 ± 0.5 | 6.0 | 9.3 | 1.00 |
| `dasel delete -f benchmark/data.json -o - '.id'` | 8.5 ± 0.6 | 7.5 | 10.7 | 1.25 ± 0.13 |
| `jq 'del(.id)' benchmark/data.json` | 28.3 ± 0.8 | 27.3 | 31.1 | 4.19 ± 0.36 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 128.4 ± 2.0 | 124.8 | 136.8 | 18.98 ± 1.57 |
