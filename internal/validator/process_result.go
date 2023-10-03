package validator

import (
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func processResult(result *rego.ResultSet, eventChan *chan e.Event, configuration ValidationConfiguration) (string, error) {
	dispatchEvent(e.NewEvent(e.BuildReportStart), eventChan)
	report, err := BuildReport(result, configuration)
	dispatchEvent(e.NewEvent(e.BuildReportDone), eventChan)
	return report, err
}
