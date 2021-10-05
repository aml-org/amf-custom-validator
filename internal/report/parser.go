package report

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func parseResults(results []any, severity *Severity) []Result {
	var parsed []Result
	for _, result := range results {
		parsed = append(parsed, parseResult(result.(object), severity))
	}
	return parsed
}

func parseResult(result object, severity *Severity) Result {
	return Result{
		Typed: Typed{
			Type: []string{"shacl:ValidationResult"},
		},
		SourceShapeName: result["sourceShapeName"].(string),
		FocusNode: Ref{
			Id: result["focusNode"].(object)["@id"].(string),
		},
		ResultSeverity: nilSafeCastRefPtr(ifNotNil(severity, func(severity any) any {
			s := severity.(*Severity)
			return &Ref{
				Id: *s,
			}
		})),
		ResultMessage: result["resultMessage"].(string),
		Trace:         parseTraces(result["trace"].([]any)),
	}
}

func parseTraces(traces []any) []Trace {
	var parsed []Trace
	for _, trace := range traces {
		parsed = append(parsed, parseTrace(trace.(object)))
	}
	return parsed
}

func parseTrace(trace object) Trace {
	return Trace{
		Typed: Typed{
			Type: []string{"validation:TraceMessage"},
		},
		Component:  trace["component"].(string),
		ResultPath: trace["resultPath"].(string),
		TraceValue: parseTraceValue(trace["traceValue"].(object)),
		Location: nilSafeCastLocation(ifNotNil(trace["location"], func(location any) any {
			return parseLocation(location.(object))
		})),
	}
}

func parseTraceValue(traceValue object) TraceValue {
	return TraceValue{
		Typed: Typed{
			Type: []string{"validation:TraceValue"},
		},
		Actual:    traceValue["actual"],
		Condition: nilSafeCastString(traceValue["condition"]),
		Expected:  traceValue["expected"],
		Negated:   traceValue["negated"].(bool),
		Argument:  nilSafeCastString(traceValue["argument"]),
		SubResult: nilSafeCastResults(ifNotNil(traceValue["subResult"], func(subResult any) any {
			return parseResults(subResult.([]any), nil)
		})),
	}
}

func parseLocation(location object) Location {
	return Location{
		Typed: Typed{
			Type: []string{"lexical:Location"},
		},
		Uri:   nil, // TODO implement
		Range: parseRange(location["range"].(object)),
	}
}

func parseRange(range_ object) Range {
	return Range{
		Typed: Typed{
			Type: []string{"lexical:Range"},
		},
		Start: parsePosition(range_["start"].(object)),
		End:   parsePosition(range_["end"].(object)),
	}

}

func parsePosition(position object) Position {
	return Position{
		Typed: Typed{
			Type: []string{"lexical:Position"},
		},
		Line:   parseJsonNumber(position["line"]),
		Column: parseJsonNumber(position["column"]),
	}
}

func nilSafeCastString(value any) *string {
	if value != nil {
		ret := value.(string)
		return &ret
	} else {
		return nil
	}
}

func nilSafeCastLocation(value any) *Location {
	if value != nil {
		ret := value.(Location)
		return &ret
	} else {
		return nil
	}
}

func nilSafeCastResults(value any) *[]Result {
	if value != nil {
		ret := value.([]Result)
		return &ret
	} else {
		return nil
	}
}

func nilSafeCastRefPtr(value any) *Ref {
	if value != nil {
		return value.(*Ref)
	} else {
		return nil
	}
}

func parseJsonNumber(raw any) int {
	switch v := raw.(type) {
	case json.Number:
		value, _ := v.Int64()
		return int(value)
	default:
		panic(fmt.Sprintf("Cannot parse type %s to int", v))
	}
}

func ifNotNil(value any, fn func(any) any) any {
	if !isNilFixed(value) {
		return fn(value)
	} else {
		return nil
	}
}

func isNilFixed(value any) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		//use of IsNil method
		return reflect.ValueOf(value).IsNil()
	}
	return false
}
