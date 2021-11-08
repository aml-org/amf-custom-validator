#!/usr/bin/env bash
rm -rf dialects vocabularies
mkdir dialects vocabularies
curl https://raw.githubusercontent.com/aml-org/models/master/src/main/dialects/validation-profile.yaml >> dialects/validation-profile.yaml
curl https://raw.githubusercontent.com/aml-org/models/master/src/main/dialects/validation-report.yaml >> ./dialects/validation-report.yaml
curl https://raw.githubusercontent.com/aml-org/models/master/src/main/dialects/lexical.yaml >> ./dialects/lexical.yaml
curl https://raw.githubusercontent.com/aml-org/models/master/src/main/vocabularies/lexical.yaml >> ./vocabularies/lexical.yaml