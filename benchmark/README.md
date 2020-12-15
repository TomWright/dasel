# Benchmarks

These benchmarks are auto generated using `./benchmark/run.sh`.

I have build what I believe are equivalent commands in dasel/jq/yq. If you have any feedback or wish to add new benchmarks please submit a PR.
## dasel vs jq

### Top level property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 6.4 ± 0.6 | 5.7 | 8.8 | 1.00 |
| `jq '.id' benchmark/data.json` | 25.9 ± 1.2 | 24.8 | 33.3 | 4.07 ± 0.41 |
### Nested property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 6.3 ± 0.6 | 5.7 | 7.8 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.7 ± 0.6 | 24.9 | 28.0 | 4.06 ± 0.38 |
### Array index

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 6.3 ± 0.5 | 5.7 | 7.4 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 25.8 ± 0.5 | 25.1 | 27.6 | 4.07 ± 0.31 |
### Append to array of strings

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 6.3 ± 0.4 | 5.7 | 7.6 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.0 ± 0.7 | 24.9 | 29.1 | 4.14 ± 0.29 |
### Update a string value

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 6.5 ± 0.5 | 5.7 | 7.6 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.0 ± 0.5 | 25.2 | 28.1 | 3.99 ± 0.34 |
### Overwrite an object

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.0 ± 0.9 | 5.8 | 9.4 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.0 ± 0.6 | 25.0 | 28.2 | 3.72 ± 0.48 |
### List keys of an array

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json -m '.-'` | 6.2 ± 0.4 | 5.6 | 7.5 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 26.0 ± 0.8 | 24.9 | 30.7 | 4.19 ± 0.31 |
## dasel vs yq

### Top level property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.id'` | 6.5 ± 0.5 | 5.8 | 8.7 | 1.00 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 120.3 ± 14.0 | 107.2 | 164.3 | 18.62 ± 2.62 |
### Nested property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.user.name.first'` | 6.4 ± 0.4 | 5.9 | 7.6 | 1.00 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 113.9 ± 7.8 | 106.8 | 134.5 | 17.87 ± 1.75 |
### Array index

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.favouriteNumbers.[1]'` | 8.5 ± 0.6 | 6.6 | 9.5 | 1.00 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 109.9 ± 3.7 | 106.7 | 130.5 | 12.97 ± 1.05 |
### Append to array of strings

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.yaml -o - '.favouriteColours.[]' blue` | 6.8 ± 0.7 | 5.9 | 8.1 | 1.00 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 109.9 ± 2.7 | 107.2 | 133.9 | 16.19 ± 1.71 |
### Update a string value

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.yaml -o - '.favouriteColours.[0]' blue` | 6.5 ± 0.5 | 5.8 | 8.1 | 1.00 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 111.4 ± 4.4 | 107.4 | 137.1 | 17.20 ± 1.43 |
### Overwrite an object

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.yaml -o - -t string -t string '.user.name' first=Frank last=Jones` | 8.0 ± 0.5 | 6.1 | 9.1 | 1.00 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 110.0 ± 2.2 | 107.6 | 125.6 | 13.79 ± 0.86 |
### List keys of an array

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml -m '.-'` | 6.4 ± 0.5 | 5.8 | 7.8 | 1.00 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 109.0 ± 1.7 | 106.4 | 114.8 | 17.10 ± 1.25 |
