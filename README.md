# AMF-OPA-Validator

## Running the validator

You can run the validator using building the project as a docker image from the 'amf-opa-validator' script:

```shell
./amf-opa-validator -in "ASYNC 2.0" -mime-in "application/yaml" -cp file://profile.yaml file://async.yaml

{
  "@type": "http://www.w3.org/ns/shacl#ValidationReport",
  "http://www.w3.org/ns/shacl#conforms": false,
  "http://www.w3.org/ns/shacl#result": [
    {
      "@type": "http://www.w3.org/ns/shacl#ValidationResult",
      "http://www.w3.org/ns/shacl#resultSeverity": {
        "@id": "http://www.w3.org/ns/shacl#Violation"
      },
      "http://www.w3.org/ns/shacl#focusNode": {
        "@id": "file://src/test/resources/integration/profile1.negative.yaml#/web-api/end-points/%2Fendpoint1/get/request/parameter/a"
      },
      "http://a.ml/vocabularies/validation#trace": [
        {
          "http://a.ml/vocabularies/validation#component": "nested",
          "http://www.w3.org/ns/shacl#resultMessage": "Not nested matching constraints for parent ∀x and child ∀y under shapes:schema",
          "http://www.w3.org/ns/shacl#focusNode": {
            "@id": "file://src/test/resources/integration/profile1.negative.yaml#/web-api/end-points/%2Fendpoint1/get/request/parameter/a/scalar/schema"
          }
        },
        {
          "http://a.ml/vocabularies/validation#component": "minCount",
          "http://www.w3.org/ns/shacl#resultMessage": "Value not matching minCount 1",
          "http://www.w3.org/ns/shacl#focusNode": {
            "@value": 0
          }
        }
      ],
      "http://www.w3.org/ns/shacl#resultMessage": "Scalars in parameters must have minLength defined",
      "http://www.w3.org/ns/shacl#sourceShape": {
        "@id": "scalar-parameters"
      }
    }
  ]
}
```

The script expects an input file and a profile. Syntax and format for the input file must be provided using the
same arguments as AMF.


## Support

Only `and`, `or`, `minCount`, `pattern`, `in` and `nested` rules and constraints are supported at the moment.