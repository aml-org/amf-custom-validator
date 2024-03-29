package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/types"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	c "github.com/aml-org/amf-custom-validator/pkg/config"
	"github.com/open-policy-agent/opa/rego"
	"strconv"
	"strings"
)

func BuildReport(resultPtr *rego.ResultSet, validationConfig c.ValidationConfiguration, reportConfig c.ReportConfiguration) (string, error) {
	result := *resultPtr
	if len(result) == 0 {
		return "", errors.New("empty result from evaluation")
	}
	raw := result[0]
	m := raw.Expressions[0].Value.(types.ObjectMap)

	profileName := m["profile"].(string)
	violations := m["violation"].([]any)
	warnings := m["warning"].([]any)
	infos := m["info"].([]any)
	results := buildResults(violations, warnings, infos)
	conforms := len(violations) == 0

	context := buildContext(len(results) == 0, reportConfig)
	reportNode := ValidationReportNode(profileName, results, conforms, validationConfig, reportConfig)
	instance := DialectInstance(&reportNode, &context)
	return Encode(instance), nil
}

func buildResults(violations []any, warnings []any, infos []any) []any {
	var results []any
	for i, r := range violations {
		results = append(results, buildValidation("violation", "violation_"+strconv.Itoa(i), r))
	}
	for i, r := range warnings {
		results = append(results, buildValidation("warning", "warning_"+strconv.Itoa(i), r))
	}
	for i, r := range infos {
		results = append(results, buildValidation("info", "info_"+strconv.Itoa(i), r))
	}
	return results
}
func buildValidation(level string, id string, raw any) types.ObjectMap {
	validation := raw.(types.ObjectMap)
	validation["resultSeverity"] = "http://www.w3.org/ns/shacl#" + strings.Title(level)
	defineIdRecursively(&validation, id)
	return validation
}

func defineIdRecursively(node *types.ObjectMap, id string) {
	if _, isTypeNode := (*node)["@type"]; isTypeNode {
		(*node)["@id"] = id
		for k, v := range *node {
			switch v := (v).(type) {
			case types.ObjectMap:
				defineIdRecursively(&v, fmt.Sprintf("%s_%s", id, k))
			case []any:
				for index, e := range v {
					switch vv := e.(type) {
					case types.ObjectMap:
						defineIdRecursively(&vv, fmt.Sprintf("%s_%d", id, index))
					}
				}
			default:
			}
		}
	}
}

func buildContext(emptyReport bool, reportConfig c.ReportConfiguration) types.ObjectMap {
	if emptyReport {
		return contexts.ConformsContext(reportConfig)
	} else {
		return contexts.DefaultValidationContext(reportConfig)
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
