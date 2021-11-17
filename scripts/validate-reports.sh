#!/usr/bin/env bash
input_reports=$(grep -rw './' -e '"shacl:ValidationReport"' | cut -d ":" -f1 | grep .jsonld | sed 's/.\/\///' | sed 's/$.\///')
declare -i exitCode
exitCode=0
for input_report in $input_reports
do
  output_report=$(java -jar amf.jar validate "$input_report")
  if grep -q 'Conforms: true' <<< "$output_report";
  then
    echo -e "\033[32m $input_report conforms \033[0m"
  else
    echo -e "\033[31m $input_report does not conform \033[0m"
    echo "$output_report"
    exitCode=1
  fi
done
exit $exitCode

