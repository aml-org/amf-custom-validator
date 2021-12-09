#!/usr/bin/env bash
parse () {
  local file=$1
  echo "Processing $file"
  rm "$file.jsonld"
  java -jar amf.jar parse "$file" --with-lexical > "$file.jsonld"
}


for subdir in ./test/data/production/*
do
  oas_30_files=$(grep -rw "$subdir" -e 'openapi: .*$' | cut -d ":" -f1 | grep .yaml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $oas_30_files
  do
    parse "$file"
  done

  raml_10_files=$(grep -rw "$subdir" -e '#%RAML 1.0 *$' | cut -d ":" -f1 | grep .raml | sed 's/.\/\///' | sed 's/$.\///')
  for file in $raml_10_files
  do
    parse "$file"
  done
done
