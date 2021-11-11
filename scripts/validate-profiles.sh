#!/usr/bin/env bash
profiles=$(grep -rw './' -e "#%Validation Profile" | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
declare -i exitCode
exitCode=0
for profile in $profiles
do
  report=$(java -jar amf.jar validate "$profile")
  if grep -q 'Conforms: true' <<< "$report";
  then
    echo -e "\033[32m$profile conforms \033[0m"
  else
    echo -e "\033[31m$profile does not conform \033[0m"
    echo "$report"
    exitCode=1
  fi
done
exit $exitCode


