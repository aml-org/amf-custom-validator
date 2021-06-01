# AMF OPA Validator

Implementation of a validator for the AMF validation profile syntax backed by a rego translation.

## Evironment

In order for the build and test process the following environment variales must be
defined:

- _GO_: path to the go binary going to be used to compile and run tests
- _AMF_: path to an AMF jar going to be used to parse examples and validate the validation profile dialect

## Test

`make test`


## Directory structure

- `cmd`

Executable commands, currently only the validator entrypoint

- `internal/parser`

Internal logic to parse validation profiles encoded in YAML.

- `internal/generator`

Transformation of a parsed validtion profile into a Rego module

- `internal/validator`

Evaluation of the generated Rego code against a JSON-LD payload. It includes the
code to normalize and index the JSON-LD before evaluation.

- `pkg/validator`

Code specific to generate the WebAssembly binary for JS. It requires to have enabled the
JS build profile since it is using experimental `syscall/js` features

- `scripts`

Auxiliary scripts:

`gen_js_package`: generates the node WASM package

`gen_production_examples`: regenerates the JSON-LD for the RAML/OAS examples in the production tests. 
The `$AMF` environment variable pointing to the AMF jar file must be defined

`gen_property_parser`: parsed the PEG grammar for the property path parser and generates the module in the `internal/parser` package.

- `third_party`

Dependencies with other tools including the AML dialect for the custom validation profile and the grammar for the property path parser

- `test`

Test data for the project including `basic`, `integration` and `production` fixtures.