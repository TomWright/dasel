outputFile="benchmark/README.md"
mdOutputFile="benchmark/tmp_results.md"

function run_file() {
  counter=0
  echo "## ${1}" >> "${outputFile}"

  name=""
  key=""
  daselCmd=""
  jqCmd=""
  yqCmd=""

  while IFS= read -r line
  do
    counter=$(($counter+1))
    if [ -z "$line" ]
    then
      counter=0
      jsonFile="benchmark/data/${key}.json"
      imagePath="benchmark/diagrams/${key}.jpg"
      readmeImagePath="diagrams/${key}.jpg"

      hyperfine --warmup 10 --runs 100 --export-json="${jsonFile}" --export-markdown="${mdOutputFile}" "${daselCmd}" "${jqCmd}" "${yqCmd}"
      python benchmark/plot_barchart.py "${jsonFile}" --title "${name}" --out "${imagePath}"

      echo "\n### ${name}\n" >> "${outputFile}"
      echo "<img src=\"${readmeImagePath}\" alt=\"${name}\" width=\"500\"/>\n" >> "${outputFile}"
      cat "${mdOutputFile}" >> "${outputFile}"

      rm "${mdOutputFile}"

    else
      case $counter in
        1)  name=$line
            ;;
        2)  key=$line
            ;;
        3)  daselCmd=$line
            ;;
        4) jqCmd=$line
           ;;
        5) yqCmd=$line
           ;;
      esac
    fi
  done < $2
}

rm -rf benchmark/data
rm -rf benchmark/diagrams

mkdir -p benchmark/data
mkdir -p benchmark/diagrams

cat benchmark/partials/top.md > "${outputFile}"

run_file "Benchmarks" "benchmark/tests.txt"

cat benchmark/partials/bottom.md >> "${outputFile}"
