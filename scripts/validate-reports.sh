#!/usr/bin/env bash
input_reports=$(grep -rw './' -e '"@type": "shacl:ValidationReport"' | cut -d ":" -f1 | grep .jsonld | sed 's/.\/\///' | sed 's/$.\///')
declare -i exitCode
exitCode=0
for input_report in $input_reports
do
  output_report=$(java -jar amf.jar validate -in "AML 1.0" -mime-in "application/yaml" -ds file://dialects/validation-report.yaml "file://$input_report")
  if grep -q '"http://www.w3.org/ns/shacl#conforms": true' <<< "$output_report";
  then
    echo -e "\033[32m $input_report conforms \033[0m"
  else
    echo -e "\033[31m $input_report does not conform \033[0m"
    echo "$output_report"
    exitCode=1
  fi
done
exit $exitCode

