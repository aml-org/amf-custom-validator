#!/usr/bin/env bash

for production_dir in ./test/data/production/*
do
  rm $production_dir/*.jsonld
  for entry in $production_dir/*
  do
    re="positive|negative"
    if [[ $entry =~ $re ]]; then
      re="oas"
      if [[ $entry =~ $re ]]; then
        java -jar $AMF parse -in "OAS 3.0" -mime-in "application/yaml" $entry | jq > $entry.jsonld
      fi
      re="raml"
      if [[ $entry =~ $re ]]; then
        java -jar $AMF parse -in "RAML 1.0" -mime-in "application/yaml" $entry | jq > $entry.jsonld
      fi
    fi
  done
done