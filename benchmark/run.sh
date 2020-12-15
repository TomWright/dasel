outputFile="benchmark/README.md"
mdOutputFile="benchmark/tmp_results.md"

function run_file() {
  echo "## ${1}\n" >> "${outputFile}"
  while IFS=, read -r name daselCmd otherCmd
  do
    echo "### ${name}\n" >> "${outputFile}"
    hyperfine --warmup 10 --runs 100 --export-markdown="${mdOutputFile}" "${daselCmd}" "${otherCmd}"
    cat "${mdOutputFile}" >> "${outputFile}"
    rm "${mdOutputFile}"
  done < $2
}

cat benchmark/partials/top.md > "${outputFile}"

run_file "dasel vs jq" "benchmark/dasel_jq.csv"
run_file "dasel vs yq" "benchmark/dasel_yq.csv"

cat benchmark/partials/bottom.md >> "${outputFile}"
