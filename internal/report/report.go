package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/open-policy-agent/opa/rego"
)

func BuildReport(result rego.ResultSet) (string, error) {
	if len(result) == 0 {
		return "", errors.New("empty result from evaluation")
	}

	raw := result[0]
	m := raw.Expressions[0].Value.(object)

	violation := Violation
	warning := Warning
	info := Info

	violations := parseResults(m["violation"].([]any), &violation)
	warnings := parseResults(m["warning"].([]any), &warning)
	infos := parseResults(m["info"].([]any), &info)

	if (len(violations) + len(warnings) + len(infos)) == 0 {
		res := object{
			"@type":    "shacl:ValidationReport",
			"@context": conformsContext(),
			"conforms": true,
		}

		return Encode(res), nil
	} else {
		var results []Result
		results = append(results, violations...)
		results = append(results, warnings...)
		results = append(results, infos...)

		res := object{
			"@type":    "shacl:ValidationReport",
			"@context": fullContext(),
			"conforms": false,
			"result":   results,
		}
		return Encode(res), nil
	}
}

func conformsContext() object {
	return object{
		"conforms": object{
			"@id": "http://www.w3.org/ns/shacl#conforms",
		},
		"shacl": "http://www.w3.org/ns/shacl#",
	}

}

func fullContext() object {
	return object{
		"actual": object{
			"@id": "http://a.ml/vocabularies/validation#actual",
		},
		"condition": object{
			"@id": "http://a.ml/vocabularies/validation#condition",
		},
		"expected": object{
			"@id": "http://a.ml/vocabularies/validation#expected",
		},
		"negated": object{
			"@id": "http://a.ml/vocabularies/validation#negated",
		},
		"argument": object{
			"@id": "http://a.ml/vocabularies/validation#argument",
		},
		"focusNode": object{
			"@id": "http://www.w3.org/ns/shacl#focusNode",
		},
		"trace": object{
			"@id": "http://a.ml/vocabularies/validation#trace",
		},
		"component": object{
			"@id": "http://a.ml/vocabularies/validation#component",
		},
		"resultPath": object{
			"@id": "http://www.w3.org/ns/shacl#resultPath",
		},
		"traceValue": object{
			"@id": "http://www.w3.org/ns/shacl#traceValue",
		},
		"location": object{
			"@id": "http://a.ml/vocabularies/validation#location",
		},
		"uri": object{
			"@id": "http://a.ml/vocabularies/lexical#uri",
		},
		"start": object{
			"@id": "http://a.ml/vocabularies/lexical#start",
		},
		"end": object{
			"@id": "http://a.ml/vocabularies/lexical#end",
		},
		"range": object{
			"@id": "http://a.ml/vocabularies/lexical#range",
		},
		"line": object{
			"@id": "http://a.ml/vocabularies/lexical#line",
		},
		"column": object{
			"@id": "http://a.ml/vocabularies/lexical#column",
		},
		"sourceShapeName": object{
			"@id": "http://a.ml/vocabularies/validation#sourceShapeName",
		},
		"conforms": object{
			"@id": "http://www.w3.org/ns/shacl#conforms",
		},
		"result": object{
			"@id": "http://www.w3.org/ns/shacl#result",
		},
		"subResult": object{
			"@id": "http://a.ml/vocabularies/validation#subResult",
		},
		"resultSeverity": object{
			"@id": "http://www.w3.org/ns/shacl#resultSeverity",
		},
		"resultMessage": object{
			"@id": "http://www.w3.org/ns/shacl#resultMessage",
		},
		"shacl":      "http://www.w3.org/ns/shacl#",
		"validation": "http://a.ml/vocabularies/validation#",
		"lexical":    "http://a.ml/vocabularies/lexical#",
	}
}

func Encode(data any) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	return b.String()
}
