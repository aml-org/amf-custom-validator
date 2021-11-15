package generator

import "fmt"

func BuildTraceValueNode(mapContent string) string {
	return fmt.Sprintf("{\"@type\": [\"reportSchema:TraceValueNode\", \"validation:TraceValue\"], %s}", mapContent)
}