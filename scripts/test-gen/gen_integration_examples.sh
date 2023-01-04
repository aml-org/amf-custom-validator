#!/usr/bin/env bash
parse () {
  local file=$1
  echo "Processing $file"
  target=$(echo $file | sed 's/.raml//' | sed 's/.yaml//')
  java -jar amf.jar parse "$file" > "$target.jsonld"
  java -jar amf.jar parse "$file" --with-lexical > "$target.lexical.jsonld"
}


for subdir in ./test/data/integration/*
do
  oas_30_files=$(grep -rw "$subdir" -e 'openapi: .*$' | grep -v '.ignore' | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $oas_30_files
  do
    parse "$file"
  done

  raml_10_files=$(grep -rw "$subdir" -e '#%RAML 1.0 *$' | grep -v '.ignore' | cut -d ":" -f1 | grep .raml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $raml_10_files
  do
    parse "$file"
  done
done
