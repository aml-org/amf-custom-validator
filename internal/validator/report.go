package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/open-policy-agent/opa/rego"
	"strconv"
	"strings"
)

func BuildReport(result rego.ResultSet) (string, error) {
	if len(result) == 0 {
		return "", errors.New("empty result from evaluation")
	}

	raw := result[0]
	m := raw.Expressions[0].Value.(map[string]interface{})

	violations := m["violation"].([]interface{})
	warnings := m["warning"].([]interface{})
	infos := m["info"].([]interface{})

	if (len(violations) + len(warnings) + len(infos)) == 0 {
		res := map[string]interface{}{
			"@type":                               "http://www.w3.org/ns/shacl#ValidationReport",
			"http://www.w3.org/ns/shacl#conforms": true,
		}

		return Encode(res), nil
	} else {
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

		res := map[string]interface{}{
			"@type":                               "http://www.w3.org/ns/shacl#ValidationReport",
			"http://www.w3.org/ns/shacl#conforms": false,
			"http://www.w3.org/ns/shacl#result":   results,
		}
		return Encode(res), nil
	}
}

func Encode(data interface{}) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	return b.String()
}

func buildViolation(level string, raw interface{}) map[string]interface{} {
	violation := raw.(map[string]interface{})
	msg := violation["message"].(string)
	shapeId := violation["shapeId"].(string)
	target := violation["target"].(string)
	traces := violation["trace"].([]interface{})

	acc := make([]interface{}, 0)
	for _, t := range traces {
		acc = append(acc, buildTrace(t))
	}

	res := map[string]interface{}{
		"@type": []string{"http://www.w3.org/ns/shacl#ValidationResult"},
		"http://www.w3.org/ns/shacl#resultSeverity": map[string]string{
			"@id": "http://www.w3.org/ns/shacl#" + strings.Title(level),
		},
		"http://www.w3.org/ns/shacl#focusNode": map[string]string{
			"@id": target,
		},
		"http://a.ml/vocabularies/validation#trace": acc,
		"http://www.w3.org/ns/shacl#resultMessage":  msg,
		"http://www.w3.org/ns/shacl#sourceShape": map[string]string{
			"@id": shapeId,
		},
	}

	return res

}

func buildTrace(raw interface{}) interface{} {
	trace := raw.(map[string]interface{})
	component := trace["component"]
	focusNode := trace["focusNode"]
	path := trace["path"]
	value := trace["value"]

	res := map[string]interface{}{
		"@type": []string{"http://a.ml/vocabularies/validation#TraceMessage"},
		"http://a.ml/vocabularies/validation#component": component,
		"http://www.w3.org/ns/shacl#focusNode":          focusNode,
		"http://www.w3.org/ns/shacl#resultPath":         path,
		"http://www.w3.org/ns/shacl#traceValue":         value,
	}

	switch trace["lexical"].(type) {
	case map[string]interface{}:
		res["http://a.ml/vocabularies/amf/parser#lexicalPosition"] = buildLexicalPosition(trace)
	}

	return res
}

func buildLexicalPosition(trace map[string]interface{}) map[string]interface{} {
	lexical := trace["lexical"].(map[string]interface{})
	start := lexical["start"].(map[string]interface{})
	end := lexical["end"].(map[string]interface{})

	res := map[string]interface{}{
		"@type": []string{"http://a.ml/vocabularies/amf/parser#Position"},
		"http://a.ml/vocabularies/amf/parser#start": map[string]interface{}{
			"@type": "http://a.ml/vocabularies/amf/parser#Location",
			"http://a.ml/vocabularies/amf/parser#line":   intFrom(start["line"]),
			"http://a.ml/vocabularies/amf/parser#column": intFrom(start["column"]),
		},
		"http://a.ml/vocabularies/amf/parser#end": map[string]interface{}{
			"@type": "http://a.ml/vocabularies/amf/parser#Location",
			"http://a.ml/vocabularies/amf/parser#line":   intFrom(end["line"]),
			"http://a.ml/vocabularies/amf/parser#column": intFrom(end["column"]),
		},
	}
	return res
}

func intFrom(any interface{}) int {
	res, _ := strconv.Atoi(any.(string))
	return res
}
