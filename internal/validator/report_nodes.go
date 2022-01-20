package validator

import "github.com/aml-org/amf-custom-validator/internal/types"

func DialectInstance(report *types.ObjectMap, context *types.ObjectMap) []types.ObjectMap {
	dialectInstance := types.ObjectMap{
		"@context": *context,
		"@id": "dialect-instance",
		"@type": []string{"meta:DialectInstance", "doc:Document", "doc:Fragment", "doc:Module", "doc:Unit"},
		"doc:encodes": []types.ObjectMap{*report},
		"doc:processingData": []types.ObjectMap{processingDataNode},
	}
	return []types.ObjectMap{dialectInstance}
}

func ValidationReportNode(profileName string, results []interface{}) types.ObjectMap {
	reportTypes := []string{"reportSchema:ReportNode", "shacl:ValidationReport"}
	commonReport := types.ObjectMap{
		"@id": "validation-report",
		"@type":    reportTypes,
		"profileName": profileName,
	}
	if len(results) == 0 { // TODO: conforms does not take into account severities!
		result := types.ObjectMap{
			"conforms": true,
		}
		types.MergeObjectMap(&result, &commonReport)
		return result
	} else {
		result := types.ObjectMap{
			"conforms": false,
			"result":   results,
		}
		types.MergeObjectMap(&result, &commonReport)
		return result
	}
}

var processingDataNode = types.ObjectMap{
	"@id": "processing-data",
	"@type": []string{"doc:DialectInstanceProcessingData"},
	"doc:sourceSpec": "Validation Report 1.0",
}