#!/usr/bin/env bash

production_dir="test/data/integration"

# remove jsonld
for entry in $production_dir/*
do
  for file in $entry/*
  do
    re="data\.lexical\.jsonld"
    if [[ $file =~ $re ]]; then
      echo $file
      rm $file
    fi
  done
done
# generate new jsonld

for entry in $production_dir/*
do
  for file in $entry/*
  do
    re="negative"
    if [[ $file =~ $re ]]; then
      re="oas$"
      if [[ $file =~ $re ]]; then
        echo $file
        java -jar $AMF parse -in "OAS 3.0" -sm true -si true -mime-in "application/yaml" --validate false $file | jq > $entry/negative.data.lexical.jsonld
      fi
      re="raml$"
      if [[ $file =~ $re ]]; then
        echo $file
        java -jar $AMF parse -in "RAML 1.0" -sm true -si true -mime-in "application/yaml" --validate false $file | jq > $entry/negative.data.lexical.jsonld
      fi
    fi
  done
done

