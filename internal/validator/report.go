package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/open-policy-agent/opa/rego"
	"strings"
)

func BuildReport(result rego.ResultSet, profileContext profile.ProfileContext) (string, error) {
	if len(result) == 0 {
		return "", errors.New("empty result from evaluation")
	}

	raw := result[0]
	m := raw.Expressions[0].Value.(map[string]interface{})

	violations := m["violation"].([]interface{})
	warnings := m["warning"].([]interface{})
	infos := m["info"].([]interface{})

	conforms := (len(violations) + len(warnings) + len(infos)) == 0
	context := buildContext(conforms, profileContext)

	if conforms {
		res := map[string]interface{}{
			"@type":    "shacl:ValidationReport",
			"@context": context,
			"conforms": true,
		}

		return Encode(res), nil
	} else {
		results := buildResults(violations, warnings, infos)

		res := map[string]interface{}{
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
func buildViolation(level string, raw interface{}) map[string]interface{} {
	violation := raw.(map[string]interface{})
	violation["resultSeverity"] = map[string]string{
		"@id": "http://www.w3.org/ns/shacl#" + strings.Title(level),
	}
	return violation
}

func buildContext(conforms bool, profileContext profile.ProfileContext) map[string]interface{} {
	if conforms {
		return buildConformsContext()
	} else {
		return buildFullContext(profileContext)
	}
}
func buildConformsContext() map[string]interface{} {
	return map[string]interface{}{
		"conforms": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#conforms",
		},
		"shacl": "http://www.w3.org/ns/shacl#",
	}
}
func buildFullContext(profileContext profile.ProfileContext) map[string]interface{} {
	context := map[string]interface{}{
		"actual": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#actual",
		},
		"condition": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#condition",
		},
		"expected": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#expected",
		},
		"negated": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#negated",
		},
		"argument": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#argument",
		},
		"focusNode": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#focusNode",
		},
		"trace": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#trace",
		},
		"component": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#component",
		},
		"resultPath": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#resultPath",
		},
		"traceValue": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#traceValue",
		},
		"location": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#location",
		},
		"uri": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#uri",
		},
		"start": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#start",
		},
		"end": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#end",
		},
		"range": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#range",
		},
		"line": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#line",
		},
		"column": map[string]string{
			"@id": "http://a.ml/vocabularies/lexical#column",
		},
		"sourceShapeName": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#sourceShapeName",
		},
		"conforms": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#conforms",
		},
		"result": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#result",
		},
		"subResult": map[string]string{
			"@id": "http://a.ml/vocabularies/validation#subResult",
		},
		"resultSeverity": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#resultSeverity",
		},
		"resultMessage": map[string]string{
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