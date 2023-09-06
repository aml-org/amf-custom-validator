package validator

import "github.com/aml-org/amf-custom-validator/internal/types"

func DialectInstance(report *types.ObjectMap, context *types.ObjectMap) []types.ObjectMap {
	dialectInstance := types.ObjectMap{
		"@context":           *context,
		"@id":                "dialect-instance",
		"@type":              []string{"meta:DialectInstance", "doc:Document", "doc:Fragment", "doc:Module", "doc:Unit"},
		"doc:encodes":        []types.ObjectMap{*report},
		"doc:processingData": []types.ObjectMap{processingDataNode},
	}
	return []types.ObjectMap{dialectInstance}
}

func ValidationReportNode(profileName string, results []interface{}, conforms bool) types.ObjectMap {
	reportTypes := []string{"reportSchema:ReportNode", "shacl:ValidationReport"}
	report := types.ObjectMap{
		"@id":         "validation-report",
		"@type":       reportTypes,
		"profileName": profileName,
		"conforms":    conforms,
	}
	if len(results) != 0 {
		report["result"] = results
	}
	return report
}

var processingDataNode = types.ObjectMap{
	"@id":            "processing-data",
	"@type":          []string{"doc:DialectInstanceProcessingData"},
	"doc:sourceSpec": "Validation Report 1.0",
}
