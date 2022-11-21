# Benchmarks

These benchmarks are auto generated using `./benchmark/run.sh`.

```
brew install dasel
go build -o "$GOPATH/bin/daselv2" cmd/dasel/main.go
brew install hyperfine
pip install matplotlib
./benchmark/run.sh
```

I have put together what I believe to be equivalent commands in daselv2/dasel/jq/yq.

If you have any feedback or wish to add new benchmarks please submit a PR.
## Benchmarks

### Root Object

<img src="diagrams/root_object.jpg" alt="Root Object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json` | 6.8 ± 0.8 | 5.7 | 9.9 | 1.00 |
| `dasel -f benchmark/data.json` | 8.9 ± 1.1 | 7.5 | 13.2 | 1.32 ± 0.22 |
| `jq '.' benchmark/data.json` | 27.8 ± 0.7 | 26.8 | 30.7 | 4.10 ± 0.49 |
| `yq --yaml-output '.' benchmark/data.yaml` | 129.0 ± 3.9 | 124.9 | 160.4 | 18.97 ± 2.30 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'id'` | 6.8 ± 0.6 | 6.1 | 9.2 | 1.00 |
| `dasel -f benchmark/data.json '.id'` | 9.1 ± 1.1 | 8.0 | 12.8 | 1.34 ± 0.20 |
| `jq '.id' benchmark/data.json` | 28.4 ± 1.1 | 27.2 | 32.8 | 4.19 ± 0.38 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 128.9 ± 2.9 | 124.4 | 140.9 | 19.04 ± 1.65 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'user.name.first'` | 6.7 ± 0.6 | 6.0 | 9.2 | 1.00 |
| `dasel -f benchmark/data.json '.user.name.first'` | 8.8 ± 0.8 | 7.8 | 11.5 | 1.32 ± 0.17 |
| `jq '.user.name.first' benchmark/data.json` | 28.7 ± 2.3 | 27.2 | 43.3 | 4.29 ± 0.51 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 128.2 ± 3.4 | 124.0 | 143.7 | 19.17 ± 1.79 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 -f benchmark/data.json 'favouriteNumbers.[1]'` | 6.9 ± 0.6 | 6.0 | 9.2 | 1.00 |
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 9.3 ± 1.3 | 7.9 | 13.5 | 1.36 ± 0.23 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 36.4 ± 13.1 | 27.1 | 106.0 | 5.30 ± 1.96 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 140.2 ± 22.3 | 124.2 | 202.2 | 20.39 ± 3.75 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[]'` | 8.3 ± 2.9 | 4.7 | 29.8 | 1.00 |
| `dasel put string -f benchmark/data.json -o - 'favouriteColours.[]' 'blue'` | 12.3 ± 2.8 | 6.4 | 23.6 | 1.48 ± 0.61 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 30.4 ± 5.9 | 25.7 | 46.5 | 3.66 ± 1.46 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 127.4 ± 3.6 | 123.1 | 155.5 | 15.33 ± 5.35 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -t string -v 'blue' -o - 'favouriteColours.[0]'` | 6.7 ± 0.7 | 6.0 | 9.5 | 1.00 |
| `dasel put string -f benchmark/data.json -o - 'favouriteColours.[0]' 'blue'` | 8.7 ± 0.8 | 7.7 | 12.1 | 1.31 ± 0.18 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 28.6 ± 1.5 | 27.3 | 36.9 | 4.27 ± 0.51 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 151.4 ± 28.6 | 125.9 | 252.9 | 22.61 ± 4.89 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 put -f benchmark/data.json -o - -t json -v '{"first":"Frank","last":"Jones"}' 'user.name'` | 6.8 ± 0.8 | 6.0 | 11.3 | 1.00 |
| `dasel put document -f benchmark/data.json -o - -d json '.user.name' '{"first":"Frank","last":"Jones"}'` | 9.9 ± 2.5 | 7.7 | 25.2 | 1.47 ± 0.40 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 34.0 ± 6.4 | 27.6 | 49.1 | 5.04 ± 1.13 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 144.6 ± 27.5 | 124.6 | 240.1 | 21.42 ± 4.81 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>


### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `daselv2 delete -f benchmark/data.json -o - 'id'` | 6.9 ± 1.4 | 6.0 | 12.6 | 1.00 |
| `dasel delete -f benchmark/data.json -o - '.id'` | 8.2 ± 0.6 | 7.7 | 11.5 | 1.20 ± 0.25 |
| `jq 'del(.id)' benchmark/data.json` | 28.1 ± 0.8 | 27.1 | 31.7 | 4.09 ± 0.82 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 140.5 ± 28.5 | 124.3 | 250.7 | 20.49 ± 5.79 |
