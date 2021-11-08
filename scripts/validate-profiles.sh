rm validation-profile.yaml
mkdir dialects vocabularies
curl https://raw.githubusercontent.com/aml-org/models/master/src/main/dialects/validation-profile.yaml >> dialects/validation-profile.yaml
profiles=$(grep -rw './' -e "#%Validation Profile" | cut -b 4- | cut -d ":" -f1 | grep .yaml)
for profile in $profiles
do
  report=$(java -jar amf.jar validate -in "AML 1.0" -mime-in "application/yaml" -ds file://dialects/validation-profile.yaml file://$profile)
  if grep -q '"http://www.w3.org/ns/shacl#conforms": true' <<< "$report";
  then
    echo "\033[32m $profile conforms \033[0m"
  else
    echo "\033[31m $profile does not conform \033[0m"
#    echo $report
  fi
done


