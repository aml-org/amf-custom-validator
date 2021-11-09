package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/internal/types"
	"github.com/open-policy-agent/opa/rego"
	"strings"
)

func BuildReport(result rego.ResultSet, profileContext profile.ProfileContext) (string, error) {
	if len(result) == 0 {
		return "", errors.New("empty result from evaluation")
	}

	raw := result[0]
	m := raw.Expressions[0].Value.(types.ObjectMap)

	violations := m["violation"].([]interface{})
	warnings := m["warning"].([]interface{})
	infos := m["info"].([]interface{})

	conforms := (len(violations) + len(warnings) + len(infos)) == 0
	context := buildContext(conforms, profileContext)

	if conforms {
		res := types.ObjectMap{
			"@type":    "shacl:ValidationReport",
			"@context": context,
			"conforms": true,
		}

		return Encode(res), nil
	} else {
		results := buildResults(violations, warnings, infos)

		res := types.ObjectMap{
			"@type":    "shacl:ValidationReport",
			"@context": context,
			"conforms": false,
			"result":   results,
		}
		return Encode(res), nil
	}
}

func buildResults(violations []interface{}, warnings []interface{}, infos []interface{}) []interface{} {
	var results []interface{}
	for _, r := range violations {
		results = append(results, buildViolation("violation", r))
	}
	for _, r := range warnings {
		results = append(results, buildViolation("warning", r))
	}
	for _, r := range infos {
		results = append(results, buildViolation("info", r))
	}
	return results
}
func buildViolation(level string, raw interface{}) types.ObjectMap {
	violation := raw.(types.ObjectMap)
	violation["resultSeverity"] = types.StringMap{
		"@id": "http://www.w3.org/ns/shacl#" + strings.Title(level),
	}
	return violation
}

func buildContext(conforms bool, profileContext profile.ProfileContext) types.ObjectMap {
	if conforms {
		return buildConformsContext()
	} else {
		return buildFullContext(profileContext)
	}
}
func buildConformsContext() types.ObjectMap {
	return types.ObjectMap{
		"conforms": types.StringMap{
			"@id": "http://www.w3.org/ns/shacl#conforms",
		},
		"shacl": "http://www.w3.org/ns/shacl#",
	}
}
func buildFullContext(profileContext profile.ProfileContext) types.ObjectMap {
	context := types.ObjectMap{
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
		"shacl":      "http://www.w3.org/ns/shacl#",
		"validation": "http://a.ml/vocabularies/validation#",
		"lexical":    "http://a.ml/vocabularies/lexical#",
	}
	for k, v := range profileContext {
		context[k] = v
	}
	return context
}

func Encode(data interface{}) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	return b.String()
}
