# AMF Rule-sets

This tutorial tells how to use the AMF Rule-sets format and the AMF Custom Validator tooling that enables you to 
define, group, document, reuse and execute common validation rules over any JSON-LD model. 
By making use of AMF parsing capabilities, you can define rules for RAML, OAS, 
AsyncAPI, GraphQL ans AsyncAPI specifications that AMF then can translate into a common JSON-LD model.

You can think of custom validations as a version of validation languages like XML Schema or JSON Schema but instead of
targeting the syntax tree of a particular document, targeting the semantic graph encoded in the JSON-LD document.

## Setting up the environment

To follow this tutorial we will be using the `ruleset-development-cli` tooling that can be used to easily author new rules
and test them.

After installing the [`ruleset-development-cli` module from NPM](https://www.npmjs.com/package/@aml-org/ruleset-development-cli). we can start creating a new project for the tutorial
rules:

```sh-session
$ npm install -g @aml-org/ruleset-development-cli
$ ruleset-development-cli new ruleset_tutorial
```

This command will create a new project with a `rules` directory and a project configuration file `configuration.json`.
It will also create a default rule under `rules/sample-rule`.

We can test that the validator is installed and working by running the `test` CLI command that will validate the sample
rule against some test examples:

```sh-session
$ ruleset-development-cli test
* Processing rule directory: rules/sample-rule
  ✓ rules/sample-rule/negative1.raml.yaml
  ✓ rules/sample-rule/positive1.oas.yaml
All examples validate
```

Now we can remove the sample rule directory by executing:

```sh-session
$ rm -rf rules/sample-rule
```
At this point we are ready to start writing some rules.


## 1. Writing and executing a custom validation

The AMF custom validation mechanism requires that you define a set of validation rules known as **"Validation Profiles"**.

A validation profile defines the name of the rule-set, documentation, a set of rules and the severity
levels associated with each rule.

Profiles are encoded into YAML documents.

Let's use the `ruleset-development-cli` to create a new rule:

```sh-session
$ ruleset-development-cli new rule example1
Generating rule directory skeleton for rule example1 with severity violation
  - rule directory: rules/example1
  - rule profile: rules/example1/profile.yaml
```

The CLI has created a new directory with a full profile for the rule named `example1`.

Let's write a simple profile for the rule:


File: *.examples/rules/example1/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example1
violation:
  - example1
validations:
  example1:
    targetClass: apiContract.WebAPI
    propertyConstraints:
      core.version:
        pattern: ^[0-9]+\.[0-9]+\.[0-9]+$
```

This profile defines a new profile, `ruleset_tutorial/example1` with a single validation rule called `example1`.
The validation checks that the version of the API model being parsed matches a regular expression `^[0-9]+\.[0-9]+\.[0-9]+$`.

The definition of the validation involves the following parts:

- `targetClass`: Defines a class of nodes in the parsed graph and that class is the target of the validation. This means that all 
nodes in the graph with that class, `apiContract.WebAPI` in this case, will be checked for all of the validation 
rules defined in the validation.

- `propertyConstraints`: Defines validation constraints for properties in the target node. In this case we are targeting
the `core.version` property and setting a `pattern` constraint.

Additionally the profile is setting a `violation` severity for this validation rule using the `violation` entry in the
document.

Let's try to validate an example API. Let's ask the `ruleset-development-cli` to generate a new example for the `example1` rule.

```sh-session
$ ruleset-development-cli new example example1 1  
  - positive example file: rules/example1/positive1.oas.yaml
  - negative example file: rules/example1/negative1.raml.yaml
```

Let's edit now the `rules/example/positive1.oas.yaml` file and provide a positive example for the rule we are testing:

File: *.examples/rules/example1/positive1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths: {}
``` 

In the same way, we can now define a negative example in the `rules/example/negative1.oas.yaml` example file:

File: *.examples/rules/example1/negative1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "v1.0"
paths: {}
``` 

We can now test these example with the `ruleset-development-cli` and check the result:

```sh-session
% ruleset-development-cli test -f example1
* Processing rule directory: rules/example1
  ✓ rules/example1/negative1.raml.yaml
  ✓ rules/example1/positive1.oas.yaml
All examples validate
```

As we can see all the examples validate correctly.

The CLI has also generated a validation report for each teste example, in the case of the positive example, the report
is located at `rules/example1/positive1.oas.yaml.report.jsonld`. The `conforms` property marks if the input data was 
compliant with the tested rule-set.

File: *.examples/rules/example1/positive1.oas.yaml.report.jsonld*
```json
[
  {
    "@context": {
      "conforms": {
        "@id": "http://www.w3.org/ns/shacl#conforms"
      },
      "doc": "http://a.ml/vocabularies/document#",
      "meta": "http://a.ml/vocabularies/meta#",
      "reportSchema": "file:///dialects/validation-report.yaml#/declarations/",
      "shacl": "http://www.w3.org/ns/shacl#"
    },
    "@id": "dialect-instance",
    "@type": [
      "meta:DialectInstance",
      "doc:Document",
      "doc:Fragment",
      "doc:Module",
      "doc:Unit"
    ],
    "doc:encodes": [
      {
        "@id": "validation-report",
        "@type": [
          "reportSchema:ReportNode",
          "shacl:ValidationReport"
        ],
        "conforms": true,
        "profileName": "ruleset_tutorial/example1"
      }
    ],
    "doc:processingData": [
      {
        "@id": "processing-data",
        "@type": [
          "doc:DialectInstanceProcessingData"
        ],
        "doc:sourceSpec": "Validation Report 1.0"
      }
    ]
  }
]
```

Great! We have validated our first API specification with a custom rule.

Notice that the same validation profile can also be applied to RAML or AsyncAPI specifications.
For example, let's generate an additional negative example in RAML:

```sh-session
$ ruleset-development-cli new example example1 2 -f raml --only negative
  - example file: rules/example1/negative2.raml.yaml
```

File: *.examples/rules/example1/negative2.raml.yaml*
```yaml
#%RAML 1.0
title: example API
version: v1.0
``` 

If we run the tests, we can find that this example also validates:

```sh-session
% ruleset-development-cli test -f example1                          
* Processing rule directory: rules/example1
  ✓ rules/example1/negative2.raml.yaml
  ✓ rules/example1/negative1.raml.yaml
  ✓ rules/example1/positive1.oas.yaml
All examples validate
```

The reason both examples validate with the same rule is that the validation logic is targeting the underlying common 
model generated by the AMF parser.
We can ask the `ruleset-development-cli` to dump the model being validated using the `model dump` command:

```sh-session
% ruleset-development-cli model dump -f example1
* Processing rule directory: rules/example1
    - JSON-LD model: rules/example1/positive1.oas.yaml.jsonld
    - JSON-LD model: rules/example1/negative2.raml.yaml.jsonld
    - JSON-LD model: rules/example1/negative1.raml.yaml.jsonld
```

In this case three model files have been created, we can take a look at the OAS example for the negative test case:

File: *.examples/rules/example1/negative1.oas.yaml.jsonld*
```json
{
  "@context": {
    "@base": "amf://id#",
    "data": "http://a.ml/vocabularies/data#",
    "shacl": "http://www.w3.org/ns/shacl#",
    "shapes": "http://a.ml/vocabularies/shapes#",
    "doc": "http://a.ml/vocabularies/document#",
    "meta": "http://a.ml/vocabularies/meta#",
    "apiContract": "http://a.ml/vocabularies/apiContract#",
    "core": "http://a.ml/vocabularies/core#",
    "xsd": "http://www.w3.org/2001/XMLSchema#",
    "rdfs": "http://www.w3.org/2000/01/rdf-schema",
    "rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
    "security": "http://a.ml/vocabularies/security#",
    "sourcemaps": "http://a.ml/vocabularies/document-source-maps#"
  },
  "@graph": [
    {
      "@id": "./",
      "@type": [
        "doc:Document",
        "doc:Fragment",
        "doc:Module",
        "doc:Unit"
      ],
      "doc:encodes": {
        "@id": "#2",
        "@type": [
          "apiContract:WebAPI",
          "apiContract:API",
          "doc:RootDomainElement",
          "doc:DomainElement"
        ],
        "apiContract:endpoint": [],
        "core:name": "example API",
        "core:version": "v1.0"
      },
      "doc:processingData": {
        "@id": "#1",
        "@type": "doc:APIContractProcessingData",
        "apiContract:modelVersion": "3.6.0",
        "doc:sourceSpec": "OAS 3.0",
        "doc:transformed": true
      },
      "doc:root": true
    }
  ]
}
```


As we can see, here we can find the node with `@id` `#2` that ahs a `@type` `apiContract:WebAPI` and includes a property
`core:version` with value `v1.0`.
These were the class of node and property we were targeting with rule `example1` in `rules/example1/profile.yaml`.

If we open the model generated from RAML in `rules/example1/negative2.raml.yaml.jsonld`, we will find the exact same model
for the API.

As we will see, is not always possible to normalize the input API specification into exactly the same common model
(especially in GraphQL and gRPC cases), so it is important to always check the rules with different examples in 
different formats if we are trying to write rule-sets valid for different types of APIs.

Now that we understand the basics of how to validate API specifications and the role that the common JSON-LD model
plays in it, we are ready to learn more about the kind of rules that can be expressed using the AML Rule-set language.

## 2. Basic scalar validations

An easier way of writing validation rules is using simple property constraints through the `propertyConstraints` facet 
in the validation profile document.

There are a number of property constraints over scalar properties that can be defined:

- `pattern`: Validates the value of a property in a target node against the provided regular expression
- `maxCount`: Validates the maximum number of values that the target node can have for a property
- `minCount`: Validates the minimum number of values that the target node can have for a property
- `exactCount`: Validates the exact number of values that the target node can have for a property
- `maxLength`: Validates the maximum length of the string value that a property of the target node can have
- `minLength`: Validates the minimum length of the string value that a property of the target node can have
- `exactLength`: Validates the exact length of the string value that a property of the target node can have
- `minExclusive`: Validates the minimum value that a value in a property of the target node can have
- `maxExclusive`: Validates the maximum value that a value in a property of the target node can have
- `minInclusive`:Validates the minimum or equal value that a value in a property of the target node can have
- `maxInclusive`: Validates the maximum or equal value that a value in a property of the target node can have
- `datatype`: Validates the type of scalar value (integer, string, float, etc.) a value for a property of the target node must have
- `in`: Validates that the set of values for a property in a target node is a subset of the values provided as an array in the validation rule
- `containsAll`: Validates that the set of matched input values is equal or a superset of the values provided as arguments in the constraint
- `containsSome`: Validates that the intersection of the set of matched input values and the values provided as arguments in the constraint is not empty

All these validations must be associated to a particular property under the `propertyConstraints` property in a validation
profile rule. The key of the `propertyConstraints` node must be a namespaced version of the property URI.

You can find the name of the properties that can potentially be constrained in the JSON-LD output generated by the parser as a URI or CURIE. You can also find these properties in 
the standard description of the [API model](https://github.com/aml-org/amf/blob/develop/documentation/model.md) that is  
generated by AMF as a YAML file.

Let's unpack each of these validations.

### 2.1 Pattern

We already explained the way pattern works. It allows you to define a regular expression that will constrain any property
in any node holding a string value. If the property has multiple values, all of them will be validated.

For example, the following profile will constrain the possible values for the protocols associated with the API using a 
regular expression over the `apiContract.scheme` property of the `apiContract.WebAPI` node class:

```sh-session
% ruleset-development-cli new rule example2                             
Generating rule directory skeleton for rule example2 with severity violation
- rule directory: rules/example2
- rule profile: rules/example2/profile.yaml
```

File: *./examples/example2/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example2
violation:
  - example2
validations:
  example2:
    targetClass: apiContract.WebAPI
    propertyConstraints:
      apiContract.scheme:
        pattern: ^http|https$
```

Let's take a look at the following OAS 2.0 API spec that defines `ws` as valid for the protocols:

```sh-session
$ ruleset-development-cli new example example2 1 -f oas --only negative 
  - example file: rules/example2/negative1.raml.yaml
```

File: *.examples/rules/example2/negative.oas.yaml*
```yaml
swagger: "2.0"
info:
  title: Basic servers
  version: "1.0"
schemes:
  - https
  - ws
paths: {}
```

If we look at the JSON-LD graph generated by the parser, we can find two values, `http` and `ws`, for the `apiContract.schemes`
property:

```sh-session
$ ruleset-development-cli model dump -f example2/negative1.raml.yaml
* Processing rule directory: rules/example2
    - JSON-LD model: rules/example2/negative1.raml.yaml.jsonld
```

Only the relevant node:

File: *.examples/rules/example2/negative1.oas.yaml.jsonld*
```json
{
  "@id": "#2",
  "@type": [
    "apiContract:WebAPI",
    "apiContract:API",
    "doc:RootDomainElement",
    "doc:DomainElement"
  ],
  "apiContract:endpoint": [],
  "apiContract:scheme": [
    "https",
    "ws"
  ],
  "core:name": "Basic servers",
  "core:version": "1.0"
}
```

If we now try to validate using the profile we have just defined, we will obtain a validation constraint that points to the
`ws` protocol.

```sh-session
% ruleset-development-cli test -f example2
* Processing rule directory: rules/example2
  ✓ rules/example2/negative1.raml.yaml
All examples validate
```

The detail from the validation report:

File: *.examples/rules/example2/negative1.oas.yaml.report.jsonld*
```json
{
        "@id": "validation-report",
        "@type": [
          "reportSchema:ReportNode",
          "shacl:ValidationReport"
        ],
        "conforms": false,
        "profileName": "ruleset_tutorial/example2",
        "result": [
          {
            "@id": "violation_0",
            "@type": [
              "reportSchema:ValidationResultNode",
              "shacl:ValidationResult"
            ],
            "focusNode": "amf://id#2",
            "resultMessage": "Validation error",
            "resultSeverity": "http://www.w3.org/ns/shacl#Violation",
            "sourceShapeName": "example2",
            "trace": [
              {
                "@id": "violation_0_0",
                "@type": [
                  "reportSchema:TraceMessageNode",
                  "validation:TraceMessage"
                ],
                "component": "pattern",
                "resultPath": "http://a.ml/vocabularies/apiContract#scheme",
                "traceValue": {
                  "@id": "violation_0_0_traceValue",
                  "@type": [
                    "reportSchema:TraceValueNode",
                    "validation:TraceValue"
                  ],
                  "argument": "ws",
                  "negated": false
                }
              }
            ]
          }
        ]
      }
```
Notice how in this case the `pattern` validation constraint has been applied to both values of the `apiContract.scheme`
property but it has only failed for the `ws` value.

**Escape character usage**: When using escape characters in your regular expression, make sure to avoid double quotes so that the special character conserves its raw value when parsed in yaml.

```pattern: ".*\?.*"``` will fail with a parse error when processing the yaml document.

```pattern: '.*\?.*'``` or ```pattern: .*\?.*``` will work as expected.

### 2.2 minCount, maxCount and exactCount

Validation constraints `minCount`, `maxCount`, and `exactCount` can be used to limit how many values a property in any target node can have.

`minCount` is especially interesting since it can be used to make part of the spec optional, if set with a `minCount` value
of zero or mandatory if `minCount` has a value major than zero.

For example, the following profile makes it mandatory to provide the name of an operation, parsed as a node of class 
`apiContract.Operation` through the `core.name` property:

File: *.examples/rules/example3/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example3

violation:
  - example3

validations:

  mandatory-operation-name:
    targetClass: apiContract.Operation
    propertyConstraints:
      core.name:
        minCount: 1
```

In OAS 3.0, names for operations are provided through the `operationId` property, so the following API should trigger a
validation error:

File: *.examples/rules/example3/negative1.oas.yaml*
```yaml
openapi: 3.0.0

info:
  title: Example API

paths:
  /test:
    get:
      summary: test path
```

If we provide the `operationId` value for the `get` operation in the `/test` endpoint, the validation will disappear.

`maxCount` can be used to limit the maximum number of values a property can have in the parsed graph.

For example, we could modify the profile discussed in section 2.1 to limit not only the value of the protocol schemes but
also the number of protocols that can defined on a single Web API:

File: *.examples/rules/example3b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example3b

violation:
  - example3b

validations:

  allowed-protocols:
    targetClass: apiContract.WebAPI
    propertyConstraints:
      apiContract.scheme:
        pattern: ^http|https$
        maxCount: 1
```

After these changes, validating the previous OAS spec defined in section 2.1 should produce two errors under the `trace` 
property, one about the value of the schemes and another one about the maximum number of schemes defined.

```sh-session
% ruleset-development-cli test -f example3b
* Processing rule directory: rules/example3b
  ✓ rules/example3b/negative1.raml.yaml
All examples validate
```

File: *.examples/rules/example3b/negative1.oas.yaml.report.jsonld*
```json
{
  "@id": "validation-report",
  "@type": [
    "reportSchema:ReportNode",
    "shacl:ValidationReport"
  ],
  "conforms": false,
  "profileName": "ruleset_tutorial/example3b",
  "result": [
    {
      "@id": "violation_0",
      "@type": [
        "reportSchema:ValidationResultNode",
        "shacl:ValidationResult"
      ],
      "focusNode": "amf://id#2",
      "resultMessage": "Validation error",
      "resultSeverity": "http://www.w3.org/ns/shacl#Violation",
      "sourceShapeName": "example3b",
      "trace": [
        {
          "@id": "violation_0_0",
          "@type": [
            "reportSchema:TraceMessageNode",
            "validation:TraceMessage"
          ],
          "component": "maxCount",
          "resultPath": "http://a.ml/vocabularies/apiContract#scheme",
          "traceValue": {
            "@id": "violation_0_0_traceValue",
            "@type": [
              "reportSchema:TraceValueNode",
              "validation:TraceValue"
            ],
            "actual": 2,
            "condition": "<=",
            "expected": 1,
            "negated": false
          }
        }
      ]
    },
    {
      "@id": "violation_1",
      "@type": [
        "reportSchema:ValidationResultNode",
        "shacl:ValidationResult"
      ],
      "focusNode": "amf://id#2",
      "resultMessage": "Validation error",
      "resultSeverity": "http://www.w3.org/ns/shacl#Violation",
      "sourceShapeName": "example3b",
      "trace": [
        {
          "@id": "violation_1_0",
          "@type": [
            "reportSchema:TraceMessageNode",
            "validation:TraceMessage"
          ],
          "component": "pattern",
          "resultPath": "http://a.ml/vocabularies/apiContract#scheme",
          "traceValue": {
            "@id": "violation_1_0_traceValue",
            "@type": [
              "reportSchema:TraceValueNode",
              "validation:TraceValue"
            ],
            "argument": "ws",
            "negated": false
          }
        }
      ]
    }
  ]
}
```
### 2.3 minLength, maxLength and exactLength

This set of constraints makes it possible to control the length of string values in the parsed graph.

For example, the following profile uses these constraints to limit the length of the string the user can provide for a description:

```yaml
#%Validation Profile 1.0

profile: string length example

violation:
 - description-length-validation

validations:

 description-length-validation:
  targetClass: apiContract.WebAPI
  propertyConstraints:
   apiContract.description:
    minLength: 40
    maxLength: 100
 ```

If we validate the following API, we should obtain a validation error due to the length of the string provided in the description of our API.

```yaml
#%RAML 1.0

title: Example API
version: 1.0
description: short description
```

### 2.4 minExclusive, maxExclusive, minInclusive, maxInclusive

This set of constraints makes it possible to control the ranges of numeric values in the parsed graph.

For example, the following profile uses these constraints to limit the range of values that are valid for the elements
defined for a RAML Array type:

File: *.examples/rules/example4/profile.yaml*
 ```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example4

violation:
  - example4

validations:

  example4:
    targetClass: raml-shapes.ArrayShape
    propertyConstraints:
      shacl.minCount:
        minInclusive: 25
        maxExclusive: 50
 ```

In this `array-limits` validation rule, we want to target all RAML arrays, parsed as `raml-shapes.ArrayShape` nodes. We have set up
two constraints, `minInclusive` and `maxExclusive`, to constrain the possible values for the elements in the array, parsed
by AMF as `shacl.minCount` properties in the output graph.
The values are constrained betweeen the values 25 (inclusive) and 50.

Provided a simple API defining some array types:

File: *.examples/rules/example4/negative1.raml.yaml*
```yaml
#%RAML 1.0

title: Example API
version: 1.0

types:
  Emails:
    type: any[]
    minItems: 100
```

If we try to parse it using the previous validation profile, we should obtain a validation error, since the `minItems` value for
the defined array (100) does not fall between the [25,50) range.

### 2.5 datatype

`datatype` is a constraint that limits the valid scalar value for a property in the parsed graph. This constraint is not
particularly useful in custom validations, since RAML, OAS and AsyncAPI have well defined types for all the properties. However, 
it can still used to modify standard type definitions, such as making it mandatory for a version to be an integer
instead of a string.

### 2.6 in

`in` makes it possible to specify an enumeration of values that constrain the possible values for a certain property in
a node. Values can be booleans, numeric values, or strings.

The following example rewrites the profile used as an example in section 2.1. It uses an `in` constraint instead of `pattern`
to constrain the possible values for the default `schemes` in an API spec:

File: *.examples/rules/example5/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example5

violation:
  - example5

validations:

  example5:
    targetClass: apiContract.WebAPI
    propertyConstraints:
      apiContract.scheme:
        in: [ http, https ]
``` 

If we try to validate the OAS 2 API defined in section 2.1 with this profile we will obtain an equivalent violation.

### 2.7 containsAll

`in` makes sure that the set of matching values in the input data is a subset or equal to the set of values provided
in the constraint.
Sometimes we want to require that the matching values must be a superset or equal to a different set of values.
`containsAll` can be used in these situations to indicate the values that must be extended by the values in the input data.

In the following example, we require that the operations for any `apiContract:EndPoint` in an API must include at least the GET and POST
operations using `containsAll`.

File: *.examples/rules/example6/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example6

violation:
  - example6

validations:

  example6:
    targetClass: apiContract.EndPoint
    propertyConstraints:
      apiContract.supportedOperation / apiContract.method:
        containsAll: [ get, post ]
``` 

The following example will validate since the `/op1` path has three operations including the mandatory `get` and `post`:

File: *.examples/rules/example6/positive1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /op1:
    get:
      responses:
        "200":
          description:
    post:
      responses:
        "200":
          description:
    delete:
      responses:
        "200":
          description:
```

On the other hand, the following negative example will fail because even if it has defined the `get` operation and two
more (`put` and `delete`) it does not include the mandatory `post` method:


File: *.examples/rules/example6/negative1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /op1:
    get:
      responses:
        "200":
          description:
    put:
      responses:
        "200":
          description:
    delete:
      responses:
        "200":
          description:
```

### 2.8 containsSome

`containsSome` is an alternative to `containsAll` where instead of checking that the input selected values are equal or 
a superset of the provided values, we are checking that the interesection between the values is not empty.

In other words we are checking that at least one of the provided values is in the set of matching input
values.

The previous example could be rewritten using `containsSome` and in this case, instead of checking that both of the 
operations GET and POST are included, we would be checking that at least the GET or the POST operation are defined for 
each `apiContract:EndPoint`:

File: *.examples/rules/example6b/profileyaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example6b

violation:
  - example6b

validations:

  example6:
    targetClass: apiContract.EndPoint
    propertyConstraints:
      apiContract.supportedOperation / apiContract.method:
        containsSome: [ get, post ]
```

And now for this rule, both previous examples are valid positive example, since both contains at least the `get` and
`post` methods.

Only this negative example without `get` or `post` operation will throw a validation error:

File: *.examples/rules/example6b/negative1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /op1:    
    put:
      responses:
        "200":
          description:
    delete:
      responses:
        "200":
          description:
```

## 3. Property pairs validations

Validation constraints discussed in section 2 are all validations over a single scalar property. In this section we will
review validations constraining pairs of scalar properties:

- *lessThanProperty*: Establishes a less than constraint over the values of two scalar properties in a node
- *lessThanOrEqualsToProperty*: Establishes a more than constraint over the values of two scalar properties in a noe
- *equalsToProperty*: Establishes an equality constraint over tha values of two scalar properties in a node
- *disjointWithProperty*: Establishes an inequality constraint over the values of two scalar properties in a node

The rest of this section will review and provide examples of these constraints.

### 3.1 lessThanProperty, lessThanOrEqualsToProperty

`lessThanProperty` and `lessThanOrEqualsToProperty` make it possible to define that the values in one property of a node
must be less than or less than or equal to the values in another property of the same node.

The following OAS 3.0.0 API spec defines maximum and minimum length for a string schema called `name`:

File: *.examples/rules/example7/negative1.oas.yaml*
```yaml
openapi: 3.0.0

info:
  title: Example API

components:
  schemas:
    name:
      type: string
      minLength: 500
      maxLength: 100

paths: {}
```

Here we can see how there is an error over those limits, making `minLength` greater than the `maxLength`.

If we parse the specification, we can see how both JSON Schema constraints are stored in the parsed graph using the 
`shacl:minLength` and `shacl:maxLength` properties of a `shapes:ScalarShape` node:

```bash
$ ruleset-development-cli model dump -f example7
* Processing rule directory: rules/example7
    - JSON-LD model: rules/example7/negative1.oas.yaml.jsonld
```

```json
{
  "@id": "#1",
  "@type": [
    "shapes:ScalarShape",
    "shapes:AnyShape",
    "shacl:Shape",
    "shapes:Shape",
    "doc:DomainElement"
  ],
  "shacl:datatype": {
    "@id": "xsd:string"
  },
  "shacl:maxLength": 100,
  "shacl:minLength": 500,
  "shacl:name": "name"
}
```

We could write a validation profile to capture this kind of errors:

```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example7

violation:
  - example7

validations:
  example7:
    targetClass: shapes.ScalarShape
    message: Min length must be less than max length
    propertyConstraints:
      shacl.minLength:
        lessThanProperty: shacl.maxLength
```

Notice how the value for the `lessThanProperty` is another property that is the target of the comparison, in this case `shacl.maxLength`.


### 3.2 equalsToProperty, disjointWithProperty

`equalsToProperty` and `disjointWithProperty` makes it possible to state that the values in two properties of the same
node must have the same or different values.

The following profiles will define two validation rules, one stating that `minLength` and `maxLength` for a string, as shown
in example 3.1, must be equal or different:

File: *.examples/rules/example8/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example8
violation:
  - example8
validations:
  example8:
    targetClass: shapes.ScalarShape
    message: Min and max length must match in scalar
    propertyConstraints:
      shacl.maxLength:
        equalsToProperty: shacl.minLength
```

File: *.examples/rules/example8b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example8b
violation:
  - example8b
validations:
  example8b:
    targetClass: shapes.ScalarShape
    message: Min and max length must not match in scalar
    propertyConstraints:
      shacl.maxLength:
        disjointWithProperty: shacl.minLength

```

The following example will be used as positive and negative example for both rules, since they implement a complementary
validation. In the first one we check that both values are equal and in the second one that both values are different:

File: *.examples/rules/example8/positive1.oas.yaml*
File: *.examples/rules/example8b/negative1.oas.yaml*
```yaml
openapi: 3.0.0

info:
  title: Example API

components:
  schemas:
    name:
      type: string
      minLength: 100
      maxLength: 100

paths: {}
```

```sh-session
$ ruleset-development-cli test -f example8  
* Processing rule directory: rules/example8b
  ✓ rules/example8b/negative1.oas.yaml
* Processing rule directory: rules/example8
  ✓ rules/example8/positive1.oas.yaml
All examples validate
```

## 4. Nested node validations and property paths

So far we have discussed examples where constraints were set for properties for a single node.

The custom validation mechanism also supports defining validations for multiple nodes in the same validation rules,
connecting them with nesting constraints:

- *nested*: Constrains all nested nodes connected to the target node through some property

As an alternative, the properties for constraints can be defined as property paths that make it possible to express
nesting traversing the output graph in a simple and efficient way.

The following sections review each of these constraints.

### 4.1 nested

`nested` can be used to add additional constraints to nodes nested under the target node and connected through a
specific property.

Consider the following simple RAML API:

File: *./examples/rules/example9/positive1.raml.yaml*
```yaml
#%RAML 1.0

title: Test API

/endpoint1:
  get:
    queryParameters:
      a:
        required: true
        type: string
        maxLength: 20
```

In this API spec the required query parameter `a` has an associated schema defining a shape validation for an scalar of
type `string` and with a `maxLength` value.

```sh-session
$ ruleset-development-cli model dump -f example9
* Processing rule directory: rules/example9
    - JSON-LD model: rules/example9/positive1.raml.yaml.jsonld
```

File: *./examples/rules/example9/positive1.raml.yaml.jsonld*
```json
{
  "@id": "#5",
  "@type": [
    "apiContract:Request",
    "core:Request",
    "apiContract:Message",
    "doc:DomainElement"
  ],
  "apiContract:parameter": {
    "@id": "#6",
    "@type": [
      "apiContract:Parameter",
      "core:Parameter",
      "doc:DomainElement"
    ],
    "apiContract:binding": "query",
    "apiContract:paramName": "a",
    "apiContract:required": true,
    "core:name": "a",
    "shapes:schema": {
      "@id": "#7",
      "@type": [
        "shapes:ScalarShape",
        "shapes:AnyShape",
        "shacl:Shape",
        "shapes:Shape",
        "doc:DomainElement"
      ],
      "shacl:datatype": {
        "@id": "xsd:string"
      },
      "shacl:maxLength": 20,
      "shacl:name": "schema"
    }
  }
}
```

As we can see, the node for the schema with `@id` `#7` is nested under the parqameter node with `@id` `#6`.

If we wanted to describe a rule checking that all query parameters must have some `maxLength`, we will need to correlate conditions
in two nested nodes. The `nested` constraint can be used exactly for these scenarios where a constraint must correlate
multiple nested nodes:

File: *./examples/rules/example9/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example9
violation:
  - example9
validations:
  example9:
    targetClass: apiContract.Parameter
    message: Scalars in parameters must have minLength defined
    propertyConstraints:
      shapes.schema:
        nested:
          propertyConstraints:
            shacl.maxLength:
              minCount: 1
```

Nested can be applied multiple times in a validation rule definition. For example the previous rule could be rewritten
targeting the `apiContract:Request` node (`id: #5`) that is the parent of the `apiContract:Parameter` node:

File: *./examples/rules/example9b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example9b

violation:
  - example9b

validations:
  example9b:
    targetClass: apiContract.Request
    message: Scalars in parameters must have minLength defined
    propertyConstraints:
      apiContract.parameter:
        nested:
          propertyConstraints:
            shapes.schema:
              nested:
                propertyConstraints:
                  shacl.maxLength:
                    minCount: 1
```

Both rules are equivalent, we are just targeting a different source node and as a consequence requiring different levels
of nesting in the definition of the rule.

### 4.2 Property paths

At the end of section 4.1 we saw how multiple nested clauses can be used to describe a path of validations through
the parsed output graph.

This way can be useful if we want to define additional constraints at different nodes that are being traversed, but as a
mechanism to reach a target nested node, this way is too verbose and error-prone.

Property paths are a simpler way to traverse and reach the target node of the graph that is being validated.

Property paths are built using a simple subset of SPARQL property path syntax:

- *Predicates*: Any property identifier, like `core.name` or `apiContract.expects` are valid property paths
- *Sequence paths*: Sequences in the form `a / b / c / ...` where `a`, `a` and `c` are valid property paths
- *Alternate paths*: Alternative paths in the form `a | b | c | ...` where `a`, `a` and `c` are valid property paths
- *Inverse paths*: Expressed as `a^` where `a` is a valid predicate like `core.name` or `apiContract.expects`

All these types of paths can be combined in complex expressions reaching any part of the output graph from a target node.
Notice that `sequence paths` have a greater priority than `alternate paths` in path expressions. Parenthesis can be used
to change the associativity in a path expression.

Let's review each of these types of path expressions. 

#### 4.2.1 Sequence paths

Sequence paths are a concatenation of properties from one node to a set of target nodes.

The last example of section 4.1 could be rewritten to simply use a path `apiContract.parameter / raml-shapes.schema` 
instead of two nested clauses:

File: *./examples/rules/example9c/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example9c
violation:
  - example9c
validations:
  example9c:
    targetClass: apiContract.Request
    message: Scalars in parameters must have minLength defined
    propertyConstraints:
      apiContract.parameter / shapes.schema:
        nested:
          propertyConstraints:
            shacl.maxLength:
              minCount: 1
```

#### 4.2.2 Alternate paths

Alternate paths make it possible to reach target nodes in different parts of the output graph that must be validated in
the same way.

For example, consider the following variation of the API spec used as example in section 4.1 where scalar RAML types are 
used in parameters and also as properties in the payload. 

File: *./examples/rules/example10/positive1.raml.yaml*
```yaml
#%RAML 1.0

title: Test API

/endpoint1:
  get:
    queryParameters:
      a:
        required: true
        type: string
        maxLength: 20
    body: 
      application/json:
        properties:
          b:
            type: string
            maxLength: 20
```

The following simplified JSON-LD document shows the node structure of the model for this data example, starting in the GET
operation:

File: *./examples/rules/example10/positive1.raml.yaml.jsonld*
```json
{
  "@id": "#4",
  "@type": [
    "apiContract:Operation"
  ],
  "apiContract:expects": {
    "@id": "#5",
    "@type": [
      "apiContract:Request"
    ],
    "apiContract:parameter": {
      "@id": "#6",
      "@type": [
        "apiContract:Parameter"
      ],
      "shapes:schema": {
        "@id": "#7",
        "@type": [
          "shapes:ScalarShape"
        ],
        "shacl:maxLength": 20
      }
    },
    "apiContract:payload": {
      "@id": "#8",
      "@type": [
        "apiContract:Payload"
      ],
      "shapes:schema": {
        "@id": "#9",
        "@type": [
          "shacl:NodeShape"
        ],
        "shacl:property": {
          "@id": "#10",
          "@type": [
            "shacl:PropertyShape"
          ],
          "shapes:range": {
            "@id": "#11",
            "@type": [
              "shapes:ScalarShape"
            ],
            "shacl:maxLength": 20	      
          }
        }
      }
    }
  }   
}
```

If we want to constrain both types of scalar types, and only the ones at those positions, 
with the same constraint for `maxLength`, we could use an `alterante path` to apply it in both cases:

File: *./examples/rules/example10/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example10
violation:
  - example10
validations:
  example10:
    targetClass: apiContract.Operation
    message: Scalars in parameters must have minLength defined
    propertyConstraints:
      apiContract.expects / ( apiContract.parameter / shapes.schema ) | ( apiContract.payload / shapes.schema / shacl.property / shapes.range ):
        nested:
          propertyConstraints:
            shacl.maxLength:
              minCount: 1
```

In this profile we first select scalar type nodes in the parameters and request body using the path 
expression `apiContract.expects / ( apiContract.parameter / shapes.schema ) | ( apiContract.payload / shapes.schema / shacl.property / shapes.range )`.

This complex path can be split in two chains thanks to the `|` alternative connector in the middle of the path:

The first one connecting the operation with the schema of the paramters:

`apiContract.expects / apiContract.parameter / shapes.schema`

and the second one connecting the operation with the schema of the properties in the schema associated to the request payload:

`apiContract.expects / apiContract.payload / shapes.schema / shacl.property / shapes.range`

If we apply manually these two paths using the JSON-LD input for the previous example we can see how the set of nodes 
with `@id` `#7` and `#11` are selected by the path.

Once the target nodes have been selected we validate the `shacl.maxLength` property is present using the `minCount` constraint

#### 4.2.3 Inverse paths

Inverse paths traverse the graph in the opposite direction from target node to parent node instead of from target node
to nested node.

Consider the following example OAS specification:

File: *./examples/rules/example11/positive1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /test:
    get:
      parameters:
        - name: a
          in: query
          schema:
            type: string
      responses:
        "200":
          description: a response
```

Imagine we would like to check that only get operations has parameters. The following profile shows a way of achieving
this validation by targeting the `apiContract.Parameter` node and then navigating back in the graph of nodes to the
`apiContract.Operation` that holds the parameter and finally selecting the `apiContract.method` property of the operation
to check that is a `get` value:

File: *./examples/rules/example11/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example11
violation:
  - example11
validations:
  example11:
    targetClass: apiContract.Parameter
    propertyConstraints:
      apiContract.parameter^ / apiContract.expects^ / apiContract.method:
        pattern: get
```

As we can see, the property path `apiContract.parameter^ / apiContract.expects^ / apiContract.method` is using tow inverse
properties `apiContract.parameter^ / apiContract.expects^` to do the inverse navigation through the graph.

One interesting reason to use inverse paths in the definition of rules is that the target node reported in the validation report
is the node targeted through the `targetClasss` property of the rule. Using inverse paths we can use a `targetClass` node 
that might be more useful in the reporting of the error.

### 4.3 Custom properties

Different API specifications provide mechanisms to extend the kind of information that can be expressed in the specification.
Annotations in RAML, vendor extensions in OAS, directives in GraphQL and custom options in gRPC protobuffers are examples
of these extensions mechanisms.

Consider the following examples:

File: *./examples/rules/example16/positive1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths: {}
x-wadus: "value"
```

File: *./examples/rules/example16/positive2.raml.yaml*
```raml
#%RAML 1.0
title: example api
version: 1.0.0

annotationTypes:
  ext:
    type: string

(wadus): "value"
```

File: *./examples/rules/example16/positive3.graphql*
```graphql
directive @wadus(message: String) on SCHEMA

schema @wadus(message: "value") {
  query: Query
}

type Query {
  field: String
}
```

File: *./examples/rules/example16/positive4.grpc.proto*
```protobuf
syntax = "proto3";

package example;


import "options.proto";

option (wadus) = "value";

service GRPCMinimal {
  rpc one_to_one (Request) returns (Reply) {}
  rpc one_to_many (Request) returns (stream Reply) {}
  rpc many_to_one (stream Request) returns (Reply) {}
  rpc many_to_many (stream Request) returns (stream Reply) {}
}

message Request {
  string message = 1;
}

message Reply {
  string message = 1;
}
```

with `options.proto` being:

File: *./examples/rules/example16/options.grpc.proto*
```protobuf
syntax = "proto2";

package example;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FileOptions {
    optional string wadus = 50000;
}
```

All these examples define (if possible) and apply a custom property called `wadus` that extend the kind of information
that can be defined at the top-level object of the API.

Note that the syntax for the actual extension can be different, for example in OAS, the prefix `x-` must be used while in
GraphQL the extensions require using a directive prefixed by `@` and parameter name for the value.

When AMF parses all these specs it generates the same element in the JSON-LD model. a `data:Object` connected by custom
`doc:DomainProperty` with a name matching the name of the extension.

In order to refer to the values connected via custom properties in the API model, we can use a special prefix `apiExt`
to mix standard and custom properties in the `propertyConstraints` paths of a rule.

For example if we would like to specify a rule that checks that the `wadus` custom property is applied at the top-level
element of the API model we could write the following profile:

```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example16
violation:
  - example16
validations:
  example16:
    message: wadus is a mandatory extension
    targetClass: apiContract.WebAPI
    propertyConstraints:
      apiExt.wadus:
        minCount: 1
```

As we can see here, the property `apiExt.wadus` instructs the validator to look for custom domain property with name
`wadus` from the target node, in this case the `apiContract.WebAPI.

Custom domain properties, specified via the prefix `apiExt` can be mixed with regular properties in any property path,
including inverse navigation statements.

## 5. Qualified constraints

Qualified constraints make it possible to express validation rules that match only a minimum or maximum number of the target
nodes selected by a particular constraint.

These are the qualified constraints supported:

- *atLeast*: Makes it possible to check that a particular validation rule matches a minimum number of the target nodes
- *atMost*: Makes it possible to check that a particular validation rule matches a maximum number of the target nodes

The following API describes a OAS API with different endpoints supporting each different set of HTTP operations:

File: *./examples/rules/example12/negative1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /endpoint1:
    get:
      responses:
        "200":
          description:
  /endpoint2:
    get:
      responses:
        "200":
          description:
    post:
      responses:
        "201":
          description:
```

This API generates the following JSON-LD graph (only showing the relevant fragment) when parsed by AMF:

```bash
$ ruleset-development-cli model dump -f example12                        
* Processing rule directory: rules/example12
    - JSON-LD model: rules/example12/positive1.oas.yaml.jsonld
    - JSON-LD model: rules/example12/negative1.oas.yaml.jsonld
```

File: *./examples/rules/example12/negative1.oas.yaml.jsonld*
```json
{
  "@id": "#2",
  "@type": [
    "apiContract:WebAPI"
  ],
  "apiContract:endpoint": [
    {
      "@id": "#3",
      "@type": [
        "apiContract:EndPoint"
      ],
      "apiContract:path": "/endpoint1",
      "apiContract:supportedOperation": {
        "@id": "#4",
        "@type": [
          "apiContract:Operation"
        ],
        "apiContract:method": "get",
        "apiContract:returns": { ... }
      }
    },
    {
      "@id": "#6",
      "@type": [
        "apiContract:EndPoint"
      ],
      "apiContract:path": "/endpoint2",
      "apiContract:supportedOperation": [
        {
          "@id": "#7",
          "@type": [
            "apiContract:Operation"
          ],
          "apiContract:method": "get",
          "apiContract:returns": { ... }
        },
        {
          "@id": "#9",
          "@type": [
            "apiContract:Operation"
          ],
          "apiContract:method": "post",
          "apiContract:returns": { ... }
        }
      ]
    }
  ]
}
```

Notice how each parsed `apiContract:EndPoint` has an associated `apiContrat.Operation`, each with a number of defined 
`apiContract:method` (`get` or `post`).

We could define a profile to check that each endpoint has *at least* one POST operation:

File: *./examples/rules/example12/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example12
violation:
  - example12
validations:
  example12:
    message: Endpoints must have a POST method
    targetClass: apiContract.EndPoint
    propertyConstraints:
      apiContract.supportedOperation:
        atLeast:
          count: 1
          validation:
            propertyConstraints:
              apiContract.method:
                in: [ post ]
```

If we try to validate the API against this profile, we will get an error about the endpoints without the post operation.

In the same way as the preceding example, we could generate a profile validating that an API is read-only by validating that no endpoint has a 
put, patch, post or delete methods using a `atMost` qualified constraint with value `0`:

File: *./examples/rules/example12b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example12b
violation:
  - example12b
validations:
  example12b:
    message: Endpoints must be read-only
    targetClass: apiContract.EndPoint
    propertyConstraints:
      apiContract.supportedOperation:
        atMost:
          count: 0
          validation:
            propertyConstraints:
              apiContract.method:
                in: [ post, put, patch, delete ]
```

Trying to parse with this new profile will result in a validation error for the endpoint with a `post` method.

## 6. Logical constraints

Logical constraints make it possible to combine set of constraints using basic boolean logic operators like: `and`, `or` and
`not`.

- *and*: Combines a set of validation rules using a logical and
- *or*: Combines a set of validation rules using a logical or
- *not*: Negates a validation rule

Logical constraints are introduced at the top level definition of a validation rule and can be combined and nested to achieve
complex validation logic.

In the following sections we will review each of these constraints with some examples.

### 6.1 and

`And` combines multiple rules using a logical and to compute the final validation result.

The following RAML API shows a simple OAS API where a GET operation defines multiple status codes for the operation
responses: 

File: *./examples/rules/example13/negative1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /endpoint1:
    get:
      responses:
        "200":
          description:
        "201":
          description:
        "300":
          description:
        "400":
          description:
        "401":
          description:  
```

The JSON-LD model generated by the AMF parser a single operation with multiple response nodes for each status code that is defined.

We could write a complex rule to check that every get operation defines status codes in the ranges 2XX, 4XX and 5XX using the `and`
boolean connector:
 
 File: *./examples/rules/example13/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example13
violation:
  - example13
validations:
  example13:
    message: Operations must have 2xx, 4xx and 5xx status codes
    targetClass: apiContract.Operation
    and:
      - propertyConstraints:
          apiContract.returns:
            atLeast:
              count: 1
              validation:
                propertyConstraints:
                  apiContract.statusCode:
                    pattern: ^2[0-9]{2}$
      - propertyConstraints:
          apiContract.returns:
            atLeast:
              count: 1
              validation:
                propertyConstraints:
                  apiContract.statusCode:
                    pattern: ^4[0-9]{2}$
      - propertyConstraints:
          apiContract.returns:
            atLeast:
              count: 1
              validation:
                propertyConstraints:
                  apiContract.statusCode:
                    pattern: ^5[0-9]{2}$
```

If we validate the API spec with this profile, an error will be reported for the missing error in the 5XX range.

### 6.2 or

In section 6.1 we have defined a validation profile with a rule to validate that all operations have status codes defined
in the ranges 2XX, 3XX and 5XX.

However sometimes we would like to express that one must valid against at least one of multiple conditions. `or` logical 
constraints can be used to achieve this behavior. 
For example, let's refine our rule to validate status codes in the ranges 2XX, 3XX
and 5XX only for GET operations.

We can achieve this expressing in the validation rule that the operation must validate the status code condition or be
a PUT, POST, DELETE or PATCH operation:

File: *./examples/rules/example13b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example13b
violation:
  - example13b
validations:
  example13b:
    message: GET operations must have 2xx, 4xx and 5xx status codes
    targetClass: apiContract.Operation
    or:
      - propertyConstraints:
          apiContract.method:
            in: [ post, put, patch, delete ]
      - and:
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^2[0-9]{2}$
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^4[0-9]{2}$
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^5[0-9]{2}$
```

If we use this profile the API discussed in section 6.1 will still fail.

However, if we define a new API spec with the same operation and list of status codes but for a `post` method, it will validate correctly:

File: *./examples/rules/example13b/positive1.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /endpoint1:
    post:
      responses:
        "200":
          description:
        "201":
          description:
        "300":
          description:
        "400":
          description:
        "401":
          description: 
```

### 6.3 not

`not`makes it possible to negate a validation rule and combine it logically with other rules.

For example, let's continue refining the example discussed in section 6.2. We can use `not` to simplify the constraint
about the method being either `put`, `post`, `delete` or `patch` by simply asserting that the method must not be `get`.

Additionally, we can add another rule to the top-level conjunction to avoid `get` operations returning `201` (created) status
codes:

File: *./examples/rules/example13c/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example13c
violation:
  - example13c
validations:
  example13c:
    message: GET operations must have 2xx, 4xx and 5xx status codes but no 201
    targetClass: apiContract.Operation
    or:
      - not:
          propertyConstraints:
            apiContract.method:
              in: [ get ]
      - and:
          - not:
              propertyConstraints:
                apiContract.returns:
                  atLeast:
                    count: 1
                    validation:
                      propertyConstraints:
                        apiContract.statusCode:
                          pattern: "^201$"
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^2[0-9]{2}$
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^4[0-9]{2}$
          - propertyConstraints:
              apiContract.returns:
                atLeast:
                  count: 1
                  validation:
                    propertyConstraints:
                      apiContract.statusCode:
                        pattern: ^5[0-9]{2}$
```

The same examples tested for the rule in the previous section will still validate this logically equivalent version of the rule.

### 6.5 conditionals

Sometimes constraints must be expressed as conditional statements. For example, in the following profile we are verifying that
if a field of type scalar (`shapes.ScalarShape`) is named `modified_at`, it has format `dateTime`:

File: *./examples/rules/example14/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example14
violation:
  - example14
validations:
  example14:
    message: Modified-at fields must be date-times
    targetClass: shapes.ScalarShape
    or:
      - not:
          propertyConstraints:
            shacl.name:
              in:
                - modified_at
      - propertyConstraints:
          shapes.format:
            minCount: 1
            in:
              - date-time
```

Now the following example will validate correctly:

File: *./examples/rules/example14/positive.oas.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
paths:
  /test:
    get:
      responses:
        "200":
          description: an operation
          content:
            application/json:
              schema:
                type: object
                properties:
                  modified_at:
                    type: string
                    format: date-time
```

In the example we are using [material implication](https://en.wikipedia.org/wiki/Material_implication_(rule_of_inference)) to express the condition
using an `or` constraint.

The same constraint can be rewritten to use the `if` / `then` conditional constraint directly:

File: *./examples/rules/example14b/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example14b
violation:
  - example14b
validations:
  example14b:
    message: Modified-at fields must be date-times
    targetClass: shapes.ScalarShape
    if:
      propertyConstraints:
        shacl.name:
          in:
            - modified_at
    then:
      propertyConstraints:
        shapes.format:
          minCount: 1
          in:
            - date-time
```

Optionally, the `else` clause can also be used to define a constraint for the case where the `if` condition does not match.

## 7. Rego rules

Sometimes you might need to express some complex validation logic not supported by any combination of the constraints described
in this tutorial.

In this situation one workaround is to resort to the use of the underlying technology used to execute the validation logic: 
 [Rego policy language](https://www.openpolicyagent.org/docs/latest/policy-language/) provided by the [Open Policy Agent](https://www.openpolicyagent.org/)
initiative.

Rego rules are defined as a template that works over a provided node in the input JSON-LD document, performs some kind of validation, and returns a result.

The input node is passed to the Rego template code through the `$node` variable and the positive result of the check must be stored 
in the `$result` variable.

The code must contain a validation that holds true for all nodes in the spec being selected by the validation containing the Rego rule.

The following profile writes a simple rule to check that the tags defined in an OAS API are being used in all the
operations in the API:

File: *./examples/rules/example15/profile.yaml*
```yaml
#%Validation Profile 1.0

profile: ruleset_tutorial/example15
violation:
  - example15
validations:
  example15:
    message: Operation tags should be defined in global tags.
    targetClass: apiContract.WebAPI
    rego: |
      o1 = collect with data.nodes as [$node] with data.property as "http://a.ml/vocabularies/apiContract#tag"
      top_level_tags = collect_values with data.nodes as o1 with data.property as "http://a.ml/vocabularies/core#name"

      p1 = collect with data.nodes as [$node] with data.property as "http://a.ml/vocabularies/apiContract#endpoint"
      p2 = collect with data.nodes as p1 with data.property as  "http://a.ml/vocabularies/apiContract#supportedOperation"
      p3 = collect with data.nodes as p2 with data.property as "http://a.ml/vocabularies/apiContract#tag"
      operation_tags = collect_values with data.nodes as p3 with data.property as "http://a.ml/vocabularies/core#name"

      common_tags = operation_tags & top_level_tags
      $result = (count(common_tags) == count(top_level_tags))
```

In the preceding example rule, the `targetClass: apiContract.WebAPI` validation target selects the top level API node, which is 
referenced as `$node` in the Repo snippet. The snippet verifies that the provided `$node` (top level API node) declares all 
the tags that are being used in each operation.

The final result is being stored in the `$result` variable. 


With this rule, an API spec will not validate if after collecting all the tags in all the operations, there is some
declared tag at the top level not being used:

File: *./examples/rules/example15/profile.yaml*
```yaml
openapi: "3.0.0"
info:
  title: example API
  version: "1.0.0"
tags:
  - name: a
  - name: b
paths:
  /endpoint1:
    get:
      tags:
        - a
      responses:
        "200":
          description:
  /endpoint2:
    post:
      tags:
        - c
      responses:
        "200":
          description: 
```

This negative example will fail because the collected tags (`a`, `c`) does not match the declared tags (`a`, `c`).

Through Rego almost any check over the API metadata can be accomplished at the cost of dealing with a lower level of
abstraction and understanding the details of the input representation.

The `ruleset-development-cli` provide the input data for Rego and the translation to Rego of any profile by passing the
`--debug` flag to the test command:

```sh-session
% ruleset-development-cli test -f example15 --debug
* Processing rule directory: rules/example15
    - (debug) generating OPA Rego profile code at rules/example15/profile.rego
  ✓ rules/example15/negative1.oas.yaml
    - (debug) OPA input data: rules/example15/negative1.oas.yaml.input
  ✓ rules/example15/positive1.oas.yaml
    - (debug) OPA input data: rules/example15/positive1.oas.yaml.input
All examples validate
```

This output can be directly tested in a Rego editing tool or in the online [Rego Playground](https://play.openpolicyagent.org/)

To write the Rego rule, a set of auxiliary Rego rules to navigate the JSON-LD graph can be used, such as the `collect` and `collect_values` rules used in the example.

Rego constraints can be used anywhere a regular constraint can be used. They can be nested in other complex constraints.

In general, Rego rules should be considered a last-resort option if the check cannot be expressed using the regular
declarative syntax. They are harder to write, maintain and they are a blackbox for reporting.
As the Ruleset syntax expands and grows, you can expect and expanded catalog of declarative rules that will make less
necessary the use of Rego.
