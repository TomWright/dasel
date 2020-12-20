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

![Top level property](benchmark/diagrams/top_level_property.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 7.9 ± 0.3 | 7.3 | 9.6 | 1.00 |
| `jq '.id' benchmark/data.json` | 26.3 ± 1.1 | 24.8 | 30.6 | 3.34 ± 0.19 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 125.9 ± 9.0 | 112.4 | 177.0 | 16.01 ± 1.32 |

### Nested property

![Nested property](benchmark/diagrams/nested_property.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 7.9 ± 0.9 | 6.8 | 11.8 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 28.6 ± 3.2 | 24.6 | 38.0 | 3.63 ± 0.57 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 127.3 ± 8.6 | 114.0 | 175.0 | 16.14 ± 2.09 |

### Array index

![Array index](benchmark/diagrams/array_index.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 7.3 ± 0.2 | 6.8 | 8.3 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.8 ± 0.8 | 24.8 | 30.1 | 3.55 ± 0.16 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 123.4 ± 8.1 | 112.1 | 142.3 | 17.01 ± 1.25 |

### Append to array of strings

![Append to array of strings](benchmark/diagrams/append_array_of_strings.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 7.6 ± 0.3 | 7.0 | 8.7 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.4 ± 1.8 | 24.8 | 36.6 | 3.49 ± 0.28 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 128.6 ± 9.8 | 115.6 | 186.9 | 16.97 ± 1.52 |

### Update a string value

![Update a string value](benchmark/diagrams/update_string.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 7.6 ± 0.3 | 6.9 | 8.9 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.8 ± 1.8 | 24.9 | 32.0 | 3.55 ± 0.27 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 131.5 ± 8.8 | 114.0 | 159.2 | 17.40 ± 1.32 |

### Overwrite an object

![Overwrite an object](benchmark/diagrams/overwrite_object.jpg)
| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.5 ± 0.3 | 7.0 | 8.6 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.2 ± 0.7 | 25.2 | 29.3 | 3.48 ± 0.17 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 132.3 ± 12.2 | 110.7 | 158.3 | 17.59 ± 1.76 |
