#!/usr/bin/env bash
parse () {
  local file=$1
  echo "Processing $file"
  java -jar amf.jar parse "$file" --with-lexical > "$file.jsonld"
}


for subdir in ./test/data/production
do
  oas_30_files=$(grep -rw "$subdir" -e 'openapi: .*$' | grep -v '.ignore' | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $oas_30_files
  do
    parse "$file"
  done

  asyncapi_files=$(grep -rw "$subdir" -e 'asyncapi: .*$' | grep -v '.ignore' | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $asyncapi_files
  do
    parse "$file"
  done

  swagger_20_files=$(grep -rw "$subdir" -e 'swagger: .*$' | grep -v '.ignore' | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $swagger_20_files
  do
    parse "$file"
  done

  raml_10_files=$(grep -rw "$subdir" -e '#%RAML*' | grep -v '.ignore' | cut -d ":" -f1 | grep .raml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $raml_10_files
  do
    parse "$file"
  done
done
