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
| `dasel -f benchmark/data.json` | 6.7 ± 0.3 | 6.2 | 7.7 | 1.00 |
| `jq '.' benchmark/data.json` | 28.2 ± 1.0 | 27.3 | 32.7 | 4.19 ± 0.22 |
| `yq --yaml-output '.' benchmark/data.yaml` | 127.3 ± 2.1 | 124.2 | 134.2 | 18.94 ± 0.81 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'id'` | 7.0 ± 1.1 | 6.3 | 13.2 | 1.00 |
| `jq '.id' benchmark/data.json` | 28.4 ± 1.1 | 27.2 | 31.8 | 4.05 ± 0.66 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 129.0 ± 5.8 | 123.6 | 163.9 | 18.39 ± 3.05 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'user.name.first'` | 7.0 ± 0.8 | 6.4 | 9.8 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 28.7 ± 0.9 | 27.2 | 30.9 | 4.07 ± 0.47 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 133.2 ± 14.8 | 124.9 | 227.0 | 18.92 ± 2.97 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'favouriteNumbers.[1]'` | 6.4 ± 0.5 | 5.8 | 8.9 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 30.5 ± 4.8 | 27.2 | 57.1 | 4.75 ± 0.85 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 131.0 ± 14.6 | 123.8 | 219.5 | 20.39 ± 2.84 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 7.0 ± 0.7 | 6.4 | 9.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 30.3 ± 3.2 | 27.5 | 54.9 | 4.33 ± 0.61 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 128.6 ± 5.7 | 125.2 | 171.2 | 18.38 ± 1.92 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' '.user.name'` | 6.5 ± 0.7 | 5.5 | 8.7 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 27.8 ± 1.3 | 26.2 | 31.9 | 4.24 ± 0.47 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 127.8 ± 4.7 | 123.7 | 149.5 | 19.54 ± 2.11 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'all().key()'` | 7.7 ± 1.1 | 6.5 | 10.9 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 28.9 ± 2.9 | 27.3 | 51.4 | 3.75 ± 0.64 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 129.6 ± 8.7 | 123.4 | 185.8 | 16.78 ± 2.56 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel delete -f benchmark/data.json -o - 'id'` | 5.3 ± 0.2 | 4.9 | 6.3 | 1.00 |
| `jq 'del(.id)' benchmark/data.json` | 26.7 ± 1.0 | 25.8 | 30.6 | 5.05 ± 0.30 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 141.8 ± 28.8 | 123.3 | 272.6 | 26.82 ± 5.59 |
