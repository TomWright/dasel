# Benchmarks

These benchmarks are auto generated using `./benchmark/run.sh`.

I have put together what I believe to be equivalent commands in dasel/jq/yq.

If you have any feedback or wish to add new benchmarks please submit a PR.
## dasel vs jq

### Top level property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.id'` | 6.8 ± 0.8 | 5.7 | 8.5 | 1.00 |
| `jq '.id' benchmark/data.json` | 26.3 ± 0.6 | 25.1 | 28.5 | 3.89 ± 0.46 |
### Nested property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.user.name.first'` | 6.8 ± 0.8 | 5.7 | 8.2 | 1.00 |
| `jq '.user.name.first' benchmark/data.json` | 25.8 ± 0.5 | 24.9 | 27.7 | 3.80 ± 0.45 |
### Array index

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json '.favouriteNumbers.[1]'` | 6.3 ± 0.4 | 5.8 | 7.4 | 1.00 |
| `jq '.favouriteNumbers[1]' benchmark/data.json` | 26.2 ± 1.4 | 24.8 | 35.0 | 4.16 ± 0.36 |
### Append to array of strings

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[]' blue` | 6.4 ± 0.5 | 5.8 | 8.0 | 1.00 |
| `jq '.favouriteColours += ["blue"]' benchmark/data.json` | 26.0 ± 0.5 | 25.0 | 27.5 | 4.05 ± 0.35 |
### Update a string value

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.json -o - '.favouriteColours.[0]' blue` | 6.5 ± 0.6 | 5.8 | 10.0 | 1.00 |
| `jq '.favouriteColours[0] = "blue"' benchmark/data.json` | 26.1 ± 1.0 | 24.9 | 30.9 | 4.04 ± 0.42 |
### Overwrite an object

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.json -o - -t string -t string '.user.name' first=Frank last=Jones` | 6.4 ± 0.5 | 5.8 | 8.0 | 1.00 |
| `jq '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.json` | 26.0 ± 0.5 | 25.1 | 27.9 | 4.08 ± 0.30 |
### List keys of an array

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.json -m '.-'` | 6.6 ± 0.7 | 5.9 | 8.9 | 1.00 |
| `jq 'keys[]' benchmark/data.json` | 25.9 ± 0.7 | 25.0 | 28.7 | 3.94 ± 0.43 |
## dasel vs yq

### Top level property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.id'` | 6.5 ± 0.6 | 5.9 | 8.1 | 1.00 |
| `yq --yaml-output '.id' benchmark/data.yaml` | 109.0 ± 3.3 | 106.5 | 139.1 | 16.77 ± 1.57 |
### Nested property

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.user.name.first'` | 6.3 ± 0.4 | 5.7 | 8.2 | 1.00 |
| `yq --yaml-output '.user.name.first' benchmark/data.yaml` | 108.5 ± 1.4 | 105.8 | 112.8 | 17.15 ± 1.24 |
### Array index

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml '.favouriteNumbers.[1]'` | 6.4 ± 0.5 | 5.7 | 9.3 | 1.00 |
| `yq --yaml-output '.favouriteNumbers[1]' benchmark/data.yaml` | 109.2 ± 2.7 | 106.5 | 131.8 | 17.15 ± 1.53 |
### Append to array of strings

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.yaml -o - '.favouriteColours.[]' blue` | 6.4 ± 0.5 | 5.8 | 7.5 | 1.00 |
| `yq --yaml-output '.favouriteColours += ["blue"]' benchmark/data.yaml` | 117.8 ± 11.7 | 107.8 | 151.3 | 18.46 ± 2.34 |
### Update a string value

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put string -f benchmark/data.yaml -o - '.favouriteColours.[0]' blue` | 6.9 ± 0.7 | 5.9 | 8.9 | 1.00 |
| `yq --yaml-output '.favouriteColours[0] = "blue"' benchmark/data.yaml` | 110.0 ± 1.6 | 107.8 | 117.0 | 15.94 ± 1.69 |
### Overwrite an object

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel put object -f benchmark/data.yaml -o - -t string -t string '.user.name' first=Frank last=Jones` | 7.2 ± 1.1 | 5.9 | 10.7 | 1.00 |
| `yq --yaml-output '.user.name = {"first":"Frank","last":"Jones"}' benchmark/data.yaml` | 110.8 ± 3.4 | 107.5 | 131.4 | 15.33 ± 2.39 |
### List keys of an array

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `dasel -f benchmark/data.yaml -m '.-'` | 6.4 ± 0.5 | 5.8 | 8.1 | 1.00 |
| `yq --yaml-output 'keys[]' benchmark/data.yaml` | 109.3 ± 2.0 | 106.6 | 122.6 | 17.08 ± 1.32 |
