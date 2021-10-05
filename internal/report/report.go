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

	violations := mergeTracesFrom(parseResults(m["violation"].([]any), &violation))
	warnings := mergeTracesFrom(parseResults(m["warning"].([]any), &warning))
	infos := mergeTracesFrom(parseResults(m["info"].([]any), &info))

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

type ResultIndexKey struct {
	sourceShapeName string
	focusNode       string
}

type ResultIndexValue struct {
	result     Result
	traceIndex map[TraceIndexKey]Trace
}

type TraceIndexKey struct {
	component  string
	resultPath string
}

func mergeTracesFrom(results []Result) []Result {
	var resultIndex = make(map[ResultIndexKey]ResultIndexValue)
	for _, result := range results {
		indexKey := ResultIndexKey{
			sourceShapeName: result.SourceShapeName,
			focusNode:       result.FocusNode.Id,
		}

		if indexValue, ok := resultIndex[indexKey]; ok {
			buildTraceIndex(result, indexValue.traceIndex) // update trace index
		} else {
			traceIndex := make(map[TraceIndexKey]Trace)
			buildTraceIndex(result, traceIndex)
			resultIndex[indexKey] = ResultIndexValue{
				result:     result,
				traceIndex: traceIndex,
			}
		}
	}

	var res []Result
	for _, vResult := range resultIndex {
		var traces []Trace
		for _, vTrace := range vResult.traceIndex {
			traces = append(traces, vTrace)
		}
		mergedResult := Result{
			Typed:           vResult.result.Typed,
			SourceShapeName: vResult.result.SourceShapeName,
			FocusNode:       vResult.result.FocusNode,
			ResultSeverity:  vResult.result.ResultSeverity,
			ResultMessage:   vResult.result.ResultMessage,
			Trace:           traces,
		}
		res = append(res, mergedResult)
	}
	return res
}

func buildTraceIndex(result Result, traceIndex map[TraceIndexKey]Trace) {
	for _, trace := range result.Trace {
		kTrace := TraceIndexKey{
			component:  trace.Component,
			resultPath: trace.ResultPath,
		}
		traceIndex[kTrace] = trace // always overwrite
	}
}