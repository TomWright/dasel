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
| `dasel -f benchmark/data.json` | 7.9 ± 3.2 | 5.3 | 19.8 | 1.00 |
| `jq '.' benchmark/data.json` | 31.2 ± 2.9 | 26.6 | 37.6 | 3.93 ± 1.64 |
| `yq --yaml-output '.' benchmark/data.yaml` | 126.9 ± 4.2 | 123.2 | 145.6 | 15.97 ± 6.51 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'id'` | 6.4 ± 0.2 | 6.1 | 7.0 | 1.00 |
| `jq '.id' benchmark/data.json` | 28.2 ± 2.3 | 27.1 | 48.8 | 4.40 ± 0.39 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 126.4 ± 1.9 | 123.5 | 132.6 | 19.71 ± 0.69 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'user.name.first'` | 6.5 ± 0.3 | 6.1 | 8.3 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 28.2 ± 0.8 | 27.0 | 31.1 | 4.36 ± 0.23 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 126.6 ± 3.2 | 123.0 | 149.9 | 19.60 ± 1.01 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'favouriteNumbers.[1]'` | 7.3 ± 3.7 | 6.0 | 37.7 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 28.6 ± 2.8 | 27.1 | 50.9 | 3.90 ± 1.99 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 126.7 ± 2.3 | 123.4 | 134.1 | 17.24 ± 8.66 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 6.4 ± 0.3 | 6.0 | 8.4 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 28.4 ± 1.6 | 27.2 | 38.0 | 4.42 ± 0.32 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 127.0 ± 2.1 | 123.7 | 140.1 | 19.72 ± 1.01 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' '.user.name'` | 6.5 ± 0.4 | 5.9 | 8.4 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 28.3 ± 1.1 | 27.1 | 31.5 | 4.33 ± 0.31 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 127.3 ± 2.5 | 123.9 | 142.2 | 19.49 ± 1.22 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json 'all().key()'` | 6.4 ± 0.3 | 6.0 | 8.1 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 28.1 ± 0.9 | 27.0 | 32.7 | 4.39 ± 0.26 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 126.4 ± 2.0 | 123.6 | 132.9 | 19.76 ± 1.02 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel delete -f benchmark/data.json -o - 'id'` | 6.6 ± 0.9 | 6.0 | 10.7 | 1.00 |
| `jq 'del(.id)' benchmark/data.json` | 28.2 ± 0.9 | 27.1 | 31.5 | 4.26 ± 0.58 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 128.4 ± 3.9 | 124.0 | 152.0 | 19.36 ± 2.63 |
