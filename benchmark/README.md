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
| `daselv2 -f benchmark/data.json` | 7.5 ± 2.5 | 6.1 | 25.1 | 1.00 |
| `dasel -f benchmark/data.json` | 9.5 ± 1.2 | 8.4 | 13.4 | 1.27 ± 0.46 |
| `jq '.' benchmark/data.json` | 29.2 ± 3.3 | 27.2 | 52.2 | 3.91 ± 1.39 |
| `yq --yaml-output '.' benchmark/data.yaml` | 134.5 ± 23.5 | 124.8 | 257.7 | 18.00 ± 6.84 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'id'` | 6.6 ± 1.2 | 5.9 | 14.6 | 1.00 |
| `dasel -f benchmark/data.json '.id'` | 9.4 ± 0.9 | 8.3 | 12.6 | 1.41 ± 0.29 |
| `jq '.id' benchmark/data.json` | 28.0 ± 0.8 | 27.0 | 32.0 | 4.21 ± 0.75 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 126.9 ± 2.1 | 123.8 | 138.1 | 19.09 ± 3.38 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'user.name.first'` | 6.8 ± 0.9 | 6.1 | 10.3 | 1.00 |
| `dasel -f benchmark/data.json '.user.name.first'` | 8.8 ± 0.8 | 8.0 | 11.9 | 1.30 ± 0.20 |
| `jq '.user.name.first' benchmark/data.json` | 28.7 ± 0.7 | 27.6 | 31.7 | 4.23 ± 0.55 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 127.4 ± 2.5 | 124.4 | 142.3 | 18.76 ± 2.43 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'favouriteNumbers.[1]'` | 6.4 ± 0.3 | 6.0 | 7.7 | 1.00 |
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 8.3 ± 0.5 | 7.8 | 11.5 | 1.29 ± 0.10 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 28.4 ± 1.0 | 27.3 | 33.2 | 4.41 ± 0.27 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 136.5 ± 29.0 | 124.5 | 325.9 | 21.19 ± 4.62 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[]'` | 3.1 ± 0.6 | 2.3 | 5.3 | 1.00 |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 4.9 ± 0.7 | 4.1 | 7.4 | 1.58 ± 0.36 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 24.4 ± 0.9 | 23.6 | 28.3 | 7.93 ± 1.46 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 123.5 ± 1.9 | 120.6 | 135.5 | 40.12 ± 7.25 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 6.5 ± 0.3 | 6.1 | 7.7 | 1.00 |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 8.5 ± 0.3 | 8.0 | 9.9 | 1.29 ± 0.07 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 28.0 ± 0.7 | 27.2 | 32.3 | 4.28 ± 0.22 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 128.2 ± 4.2 | 124.8 | 152.1 | 19.59 ± 1.07 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' 'user.name'` | 6.5 ± 0.7 | 6.0 | 12.5 | 1.00 |
| `dasel put document -f benchmark/data.json -o - -d json '.user.name' '{"first":"Frank","last":"Jones"}'` | 8.3 ± 0.3 | 7.9 | 9.4 | 1.28 ± 0.14 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 28.2 ± 0.8 | 27.3 | 33.0 | 4.34 ± 0.46 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 132.9 ± 12.9 | 124.6 | 194.2 | 20.44 ± 2.88 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'all().key()'` | 6.5 ± 0.3 | 6.1 | 7.4 | 1.00 |
| `dasel -f benchmark/data.json -m '.-'` | 8.5 ± 0.4 | 7.9 | 10.8 | 1.29 ± 0.08 |
| `jq 'keys[]' benchmark/data.json` | 28.2 ± 0.7 | 27.3 | 31.2 | 4.31 ± 0.20 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 136.9 ± 20.1 | 124.9 | 230.9 | 20.95 ± 3.18 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 delete -f benchmark/data.json -o - 'id'` | 7.7 ± 0.9 | 6.2 | 9.4 | 1.00 |
| `dasel delete -f benchmark/data.json -o - '.id'` | 9.8 ± 1.3 | 8.0 | 12.3 | 1.28 ± 0.22 |
| `jq 'del(.id)' benchmark/data.json` | 28.1 ± 0.7 | 27.4 | 31.0 | 3.66 ± 0.43 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 127.0 ± 1.4 | 125.0 | 132.4 | 16.55 ± 1.91 |
