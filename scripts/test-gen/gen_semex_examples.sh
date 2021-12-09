#!/usr/bin/env bash
files=(
  api.async.yaml
  api.oas20.yaml
  api.oas30.yaml
  api.raml
)
for dir in ./test/data/semex/*
do
  for file in "${files[@]}"
  do
    if [ -f "$dir/$file" ]; then
      echo "Processing $dir/$file"
      java -jar amf.jar parse "$dir/$file" --with-lexical --extensions "$dir"/dialect.yaml > "$dir/$file.jsonld"
    else
        echo "Skipped $dir/$file"
    fi
  done
done
