#!/usr/bin/env bash

for production_dir in ./test/data/production/*
do
  echo $production_dir
  #re="owasp"
  #if  [[ $production_dir =~ $re ]]; then
  # remove jsonld
  for entry in $production_dir/*
  do
    re="negative3|negative4" # this cannot be generated, ignore
    if [[ $entry =~ $re ]]; then
      re="spectral"
      if [[ $entry =~ $re ]]; then
        continue
      fi
    fi
    re="jsonld"
    if [[ $entry =~ $re ]]; then
      rm $entry
    fi
  done
  # generate new jsonld
  for entry in $production_dir/*
  do
    re="negative3|negative4" # this cannot be generated, ignore
    if [[ $entry =~ $re ]]; then
      continue
    fi
    re="report" # this cannot be generated, ignore
    if [[ $entry =~ $re ]]; then
      continue
    fi
    re="positive|negative"
    if [[ $entry =~ $re ]]; then
      re="oas"
      if [[ $entry =~ $re ]]; then
        echo $entry
        java -jar $AMF parse -in "OAS 3.0" -mime-in "application/yaml" --validate false $entry | jq > $entry.jsonld
      fi
      re="raml"
      if [[ $entry =~ $re ]]; then
        echo $entry
        java -jar $AMF parse -in "RAML 1.0" -mime-in "application/yaml" --validate false $entry | jq > $entry.jsonld
      fi
    fi
  done
  #fi
done