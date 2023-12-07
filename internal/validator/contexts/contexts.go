package contexts

import "github.com/aml-org/amf-custom-validator/internal/types"

var ApiExtensionUri = "http://a.ml/vocabularies/api-extension#"
var DefaultAMFContext = types.ObjectMap{
	"data":        "http://a.ml/vocabularies/data#",
	"shacl":       "http://www.w3.org/ns/shacl#",
	"shapes":      "http://a.ml/vocabularies/shapes#",
	"raml-shapes": "http://a.ml/vocabularies/shapes#",
	"doc":         "http://a.ml/vocabularies/document#",
	"meta":        "http://a.ml/vocabularies/meta#",
	"apiContract": "http://a.ml/vocabularies/apiContract#",
	"core":        "http://a.ml/vocabularies/core#",
	"xsd":         "http://www.w3.org/2001/XMLSchema#",
	"rdfs":        "http://www.w3.org/2000/01/rdf-schema",
	"rdf":         "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
	"security":    "http://a.ml/vocabularies/security#",
	"sourcemaps":  "http://a.ml/vocabularies/document-source-maps#",
	"apiExt":      ApiExtensionUri,
	"gcl":         "http://anypoint.com/vocabs/gcl#",
	"management":  "http://anypoint.com/vocabs/management#",
	"api":         "http://anypoint.com/vocabs/api#",
	"catalog":     "http://anypoint.com/vocabs/digital-repository#",
}

var DefaultValidationContext = types.ObjectMap{
	"actual": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#actual",
	},
	"condition": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#condition",
	},
	"expected": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#expected",
	},
	"negated": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#negated",
	},
	"argument": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#argument",
	},
	"focusNode": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#focusNode",
	},
	"trace": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#trace",
	},
	"component": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#component",
	},
	"resultPath": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#resultPath",
	},
	"traceValue": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#traceValue",
	},
	"location": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#location",
	},
	"uri": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#uri",
	},
	"start": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#start",
	},
	"end": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#end",
	},
	"range": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#range",
	},
	"line": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#line",
	},
	"column": types.StringMap{
		"@id": "http://a.ml/vocabularies/lexical#column",
	},
	"sourceShapeName": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#sourceShapeName",
	},
	"conforms": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#conforms",
	},
	"dateCreated": types.StringMap{
		"@id": "http://a.ml/vocabularies/core#dateCreated",
	},
	"profileName": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#profileName",
	},
	"result": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#result",
	},
	"subResult": types.StringMap{
		"@id": "http://a.ml/vocabularies/validation#subResult",
	},
	"resultSeverity": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#resultSeverity",
	},
	"resultMessage": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#resultMessage",
	},
	"shacl":         "http://www.w3.org/ns/shacl#",
	"doc":           "http://a.ml/vocabularies/document#",
	"meta":          "http://a.ml/vocabularies/meta#",
	"validation":    "http://a.ml/vocabularies/validation#",
	"lexical":       "http://a.ml/vocabularies/lexical#",
	"reportSchema":  reportPath,
	"lexicalSchema": lexicalPath,
}

var reportPath = "file:///dialects/validation-report.yaml#/declarations/"
var lexicalPath = "file:///dialects/lexical.yaml#/declarations/"

var ConformsContext = types.ObjectMap{
	"conforms": types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#conforms",
	},
	"shacl":        "http://www.w3.org/ns/shacl#",
	"doc":          "http://a.ml/vocabularies/document#",
	"reportSchema": reportPath,
	"meta":         "http://a.ml/vocabularies/meta#",
	"dateCreated": types.StringMap{
		"@id": "http://a.ml/vocabularies/core#dateCreated",
	},
}
