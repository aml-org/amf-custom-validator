package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/notifications"
	c "github.com/aml-org/amf-custom-validator/pkg/config"
	e "github.com/aml-org/amf-custom-validator/pkg/events"
	"github.com/open-policy-agent/opa/rego"
)

func processResult(result *rego.ResultSet, eventChan *chan e.Event, validationConfig c.ValidationConfiguration, reportConfig c.ReportConfiguration) (string, error) {
	notifications.DispatchEvent(e.NewEvent(e.BuildReportStart), eventChan)
	report, err := BuildReport(result, validationConfig, reportConfig)
	notifications.DispatchEvent(e.NewEvent(e.BuildReportDone), eventChan)
	return report, err
}
