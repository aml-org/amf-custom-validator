package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/validator/config"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func processResult(result *rego.ResultSet, eventChan *chan e.Event, validationConfig config.ValidationConfiguration, reportConfig config.ReportConfiguration) (string, error) {
	dispatchEvent(e.NewEvent(e.BuildReportStart), eventChan)
	report, err := BuildReport(result, validationConfig, reportConfig)
	dispatchEvent(e.NewEvent(e.BuildReportDone), eventChan)
	return report, err
}
