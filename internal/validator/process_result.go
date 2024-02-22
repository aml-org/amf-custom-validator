package validator

import (
	config2 "github.com/aml-org/amf-custom-validator/pkg/config"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func processResult(result *rego.ResultSet, eventChan *chan e.Event, validationConfig config2.ValidationConfiguration, reportConfig config2.ReportConfiguration) (string, error) {
	dispatchEvent(e.NewEvent(e.BuildReportStart), eventChan)
	report, err := BuildReport(result, validationConfig, reportConfig)
	dispatchEvent(e.NewEvent(e.BuildReportDone), eventChan)
	return report, err
}
