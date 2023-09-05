# AMF Custom Validator

AMF Custom Validator is a tool for writing and validating semantic graphs.

It takes two inputs, data and constraints, and produces a validation report as output. The data input is a semantic
graph written in JSON-LD syntax. The constraints input is called a _validation profile_ (or _ruleset_) which contains a
series of constraints written in YAML syntax with a custom vocabulary.

AMF Custom Validator is implemented in Go and can be consumed as:

* A Go library
* A JS library, using a JS runtime with WASM support
* A standalone CLI

# Getting started

Requirements:

* Go 1.21
* Make

## Using the CLI

To install the AMF Custom Validator (`acv`) CLI, run:

```shell
make install
```

This will compile and add `acv` to your compiled Go binaries.

After installing run:

```shell
acv help
acv validate PROFILE DATA
```

## Using the Go library

To add the AMF Custom Validator library to your Go project, run the following command in a shell:

```shell
go get github.com/aml-org/amf-custom-validator
```

To use the validator in your project:

```go
package validation

import (
	acv "github.com/aml-org/amf-custom-validator/pkg"
)

func Validate() (string, error) {
	// Read validation profile YAML file
	var profile string = ""

	// Read data JSON-LD file
	var data string = ""

	// Call validation
	report, err := acv.Validate(profile, data, false, nil)
	
	return report, err
}
```

## Using the JS library

To add the library to your JS project, run the following command in a shell:

```shell
npm i @aml-org/amf-custom-validator
```

To use the validator in your project:

```js
var acv = require('@aml-org/amf-custom-validator')

// Validator must be explicitly initialized
acv.initialize(() => {

    // Validation must be called after initialization in a callback
    acv.validate(profile, data, false, (report, err) => {

        // Validator must be exited on the last validator call
        acv.exit();
    });
}); 
```

# How does it work

The AMF Custom Validator is based on OPA (https://www.openpolicyagent.org/). OPA is the core validation engine that
executes Rego code, a DSL for writing policies. The AMF Custom Validator works by:

* Generating Rego code from validation profiles that OPA can execute.
* Normalizing input data to facilitate how the Rego code manages data.

After executing the validation, the AMF Custom Validator produces a validation report in JSON-LD syntax.

---

To obtain output from the data normalization and rego code generation processes, run the following `acv` commands in a shell:

Data normalization

```shell
acv normalize DATA
```

Rego code generation

```shell
acv generate PROFILE
```

Compile profile

```shell
acv compile PROFILE
```

# Relation with AMF

AMF is a framework capable of producing semantic graphs in JSON-LD syntax from API specification documents or arbitrary
YAML/JSON documents through AML. The AMF output can be used as the data input for the AMF Custom Validator. You can also
check validation profile definitions with AMF.

For more information on AMF check:
* The [AMF GitHub repository](https://github.com/aml-org/amf)
* The [AMF Documentation](https://a.ml/docs/)

To integrate with AMF, download an AMF CLI JAR with the following script:

```shell
./scripts/download-amf-cli.sh
```

## Obtaining semantic data from API specs with AMF

To parse API specification files with the CLI using Java (Java 17 required):

```shell
# Simple parse
java -jar amf.jar parse API_SPEC_FILE

# Parse with lexical information to be used in validation report
java -jar amf.jar parse API_SPEC_FILE --with-lexical

# Parse with semantic extensions for parsing
java -jar amf.jar parse API_SPEC_FILE --extensions SEMANTIC_EXTENSION_FILE
```

## Validating validation profile definitions with AMF

For a complete guide on how to write validation profiles, read
the [validation tutorial](docs/validation_tutorial/validation.md).

To validate your validation profile definition, run the following command:

```shell
java -jar amf.jar validate PROFILE
```

# Contributing

For details on how to develop and contribute, please refer to [contributing guide](docs/contributing.md)
