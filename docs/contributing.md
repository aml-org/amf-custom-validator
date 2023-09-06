# Project structure

## Development requirements

* Go 1.16
* Java 17 (to validate validation profiles & reports only)
* Node 16
* Make

## Running tests

To run tests simply run

```shell
# All tests
make test

# Go tests
make test-go

# JS clients tests
make build-js # with every code change
make test-js
```

When writing new validation profiles or modifying validation reports you will need to validate both are correct. This 
requires Java 17 and an AMF CLI that can be downloaded with `./scripts/download-amf-cli.sh`

To find and validate all validation profiles and reports in the project run:

```shell
# Tests validation profiles
make test-profiles

# Tests validation reports
make test-reports
```


## Directory structure

`cmd`

Executable commands, the CLI implementation

---

`internal`

Code that is not exported to external clients 

* `internal/parser`

Internal logic to parse validation profiles encoded in YAML

* `internal/generator`

Transformation of a parsed validation profile into a Rego module

* `internal/validator`

Evaluation of the generated Rego code against a JSON-LD payload. It includes the
code to normalize and index the JSON-LD before evaluation.

* `internal/misc` & `internal/types`

Auxiliary code

---

`pkg`

Code that exposes entry points for Go clients

---

`js`

Code specific to generate the WebAssembly binary for JS. It requires to have enabled the
JS build profile since it is using experimental `syscall/js` features

---

`wrappers`

Contains JS code that wraps the generated WASM and exposes it in a friendlier manner to JS clients

* `wrappers/js`

Wrapper for Node clients

* `wrappers/js-web`

Wrapper for browser clients

---

`scripts`

Auxiliary scripts

* `scripts/gen_js_package.sh`

Generates the Node WASM package

* `scripts/bundle_js_web_package.sh`

Bundles the generated WASM package in a browser compatible bundle. Should be run after `scripts/gen_js_package.sh`

* `scripts/download-amf-cli.sh`

Downloads an AMF CLI

* `scripts/test-gen/*`

Generates semantic graph JSON-LD files from API specs files using the AMF CLI. Divided by testing suites

* `scripts/gen_property_parser.sh` 

Generates a property parser from the PEG grammar for the parsing property paths in validation profiles

* `scripts/validate-profiles.sh`

Finds validation profile files in the project and validates their definition with the AMF CLI

* `scripts/validate-reports.sh`

Finds validation report files in the project and validates these with the AMF CLI

---

`third_party`

Dependencies with other tools including the PEG grammar for the property path parser

---

`test`

Test data for the project

* `test/basic`

Basic test cases

* `test/integration`

More complex test cases

* `test/production`

Production-like test cases

* `test/semex`

Test cases which use the AMF Semantic Extensions feature 

* `test/shacl`

Test cases extracted from the SHACL TCK

* `test/tck`

Proprietary TCK cases
