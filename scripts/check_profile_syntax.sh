#!/usr/bin/env bash

set -e

for file in ./test/data/**/**/*
do
  re="profile[0-9]*.yaml"
  if [[ $file =~ $re ]]; then
    echo $file
    java -jar amf.jar parse -in "AML 1.0" -mime-in "application/yaml" -ds file://third_party/dialect.yaml $file 1>/dev/null
  fi
done

for file in ./test/data/**/*
do
  re="profile[0-9]*.yaml"
  if [[ $file =~ $re ]]; then
    echo $file
    java -jar amf.jar parse -in "AML 1.0" -mime-in "application/yaml" -ds file://third_party/dialect.yaml $file 1>/dev/null
  fi
done
