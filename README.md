#AMF-OPA-Validator

##Requirements

In order to run the CLI version of the library you need to have installed the OPA binary in 
the root directory of this project.

Get OPA from: "https://github.com/open-policy-agent/opa/releases"

You must install the binary in `./opa` and provide an opa_capabilities file in `./opa_capabilities.json`.
A sample capabilities for OPA version v0.27.1 is provided with the project as an example.

You also need to have `tar` available in the system.

The project has only been tested on Linux at the moment.

## Running the validator

You can run the validator using NPM:

```shell
npm run validate -- -in "ASYNC 2.0" -mime-in "application/yaml" -cp file://profile.yaml file://async.yaml

  {
    "result": [
      {
        "constraintId": "in",
        "target": "file://async.yaml#/web-api/end-points/%2Fexample%2Fother_topic/publish",
        "shapeId": "validation1",
        "traceMessage": "Operation not permitted",
        "message": "Value no in set {'get'}",
        "value": "publish"
      },
      {
        "constraintId": "in",
        "target": "file://async.yaml#/web-api/end-points/%2Fexample%2Ftopic/publish",
        "shapeId": "validation1",
        "traceMessage": "Operation not permitted",
        "message": "Value no in set {'get'}",
        "value": "publish"
      }
    ]
  }
]

```

The script expects an input file and a profile. Syntax and format for the input file must be provided using the
same arguments as AMF.


## Support

Only `minCount`, `pattern` and `in` constraints are supported at the moment.