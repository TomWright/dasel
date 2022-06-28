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

| Command                                    |   Mean [ms] | Min [ms] | Max [ms] |    Relative |
| :----------------------------------------- | ----------: | -------: | -------: | ----------: |
| `dasel -f benchmark/data.json`             |  14.8 ± 3.7 |     10.4 |     31.7 |        1.00 |
| `jq '.' benchmark/data.json`               |  31.8 ± 4.4 |     25.7 |     40.0 | 2.15 ± 0.62 |
| `yq --yaml-output '.' benchmark/data.yaml` | 133.5 ± 6.3 |    123.1 |    157.0 | 9.03 ± 2.33 |

### Top level property

<img src="diagrams/top_level_property.jpg" alt="Top level property" width="500"/>

| Command                                      |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :------------------------------------------- | ----------: | -------: | -------: | -----------: |
| `dasel -f benchmark/data.json '.id'`         |  11.6 ± 0.5 |     10.9 |     13.8 |         1.00 |
| `jq '.id' benchmark/data.json`               |  27.0 ± 2.5 |     25.7 |     44.1 |  2.34 ± 0.23 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 132.6 ± 5.3 |    122.0 |    152.1 | 11.47 ± 0.67 |

### Nested property

<img src="diagrams/nested_property.jpg" alt="Nested property" width="500"/>

| Command                                                   |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :-------------------------------------------------------- | ----------: | -------: | -------: | -----------: |
| `dasel -f benchmark/data.json '.user.name.first'`         |  11.5 ± 0.5 |     11.0 |     13.9 |         1.00 |
| `jq '.user.name.first' benchmark/data.json`               |  26.5 ± 0.8 |     25.3 |     32.2 |  2.31 ± 0.12 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 133.0 ± 4.1 |    125.8 |    145.1 | 11.58 ± 0.59 |

### Array index

<img src="diagrams/array_index.jpg" alt="Array index" width="500"/>

| Command                                                       |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :------------------------------------------------------------ | ----------: | -------: | -------: | -----------: |
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'`        |  11.3 ± 0.4 |     10.8 |     13.8 |         1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json`               |  26.8 ± 1.8 |     25.5 |     35.5 |  2.37 ± 0.18 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 133.8 ± 5.2 |    125.5 |    156.0 | 11.85 ± 0.66 |

### Append to array of strings

<img src="diagrams/append_array_of_strings.jpg" alt="Append to array of strings" width="500"/>

| Command                                                                    |    Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :------------------------------------------------------------------------- | -----------: | -------: | -------: | -----------: |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` |   11.5 ± 0.3 |     10.6 |     12.7 |         1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json`                   |   26.9 ± 1.6 |     25.7 |     40.1 |  2.33 ± 0.16 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml`     | 137.8 ± 11.0 |    122.0 |    184.2 | 11.98 ± 1.02 |

### Update a string value

<img src="diagrams/update_string.jpg" alt="Update a string value" width="500"/>

| Command                                                                     |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :-------------------------------------------------------------------------- | ----------: | -------: | -------: | -----------: |
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` |  11.9 ± 0.8 |     10.9 |     16.4 |         1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json`                    |  27.4 ± 2.2 |     25.8 |     37.0 |  2.31 ± 0.24 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml`      | 133.9 ± 4.2 |    126.2 |    148.2 | 11.30 ± 0.82 |

### Overwrite an object

<img src="diagrams/overwrite_object.jpg" alt="Overwrite an object" width="500"/>

| Command                                                                                                |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :----------------------------------------------------------------------------------------------------- | ----------: | -------: | -------: | -----------: |
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` |  11.5 ± 0.7 |     10.5 |     13.9 |         1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json`                               |  27.4 ± 3.4 |     25.2 |     42.0 |  2.39 ± 0.33 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml`                 | 133.2 ± 3.9 |    122.7 |    144.8 | 11.60 ± 0.76 |

### List keys of an array

<img src="diagrams/list_array_keys.jpg" alt="List keys of an array" width="500"/>

| Command                                         |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :---------------------------------------------- | ----------: | -------: | -------: | -----------: |
| `dasel -f benchmark/data.json -m '.-'`          |  11.7 ± 0.7 |     10.9 |     16.2 |         1.00 |
| `jq 'keys[]' benchmark/data.json`               |  26.8 ± 1.2 |     25.4 |     32.8 |  2.29 ± 0.17 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 133.6 ± 4.9 |    124.6 |    155.8 | 11.45 ± 0.81 |

### Delete property

<img src="diagrams/delete_property.jpg" alt="Delete property" width="500"/>

| Command                                           |   Mean [ms] | Min [ms] | Max [ms] |     Relative |
| :------------------------------------------------ | ----------: | -------: | -------: | -----------: |
| `dasel delete -f benchmark/data.json -o - '.id'`  |  11.9 ± 0.8 |     11.0 |     15.6 |         1.00 |
| `jq 'del(.id)' benchmark/data.json`               |  26.7 ± 1.1 |     25.5 |     34.0 |  2.24 ± 0.17 |
| `yq --yaml-output 'del(.id)' benchmark/data.yaml` | 134.6 ± 4.5 |    124.7 |    155.1 | 11.29 ± 0.84 |
