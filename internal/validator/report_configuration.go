package validator

type ReportConfiguration struct {
	IncludeReportCreationTime bool
}

func DefaultReportConfiguration() ReportConfiguration {
	return ReportConfiguration{
		IncludeReportCreationTime: true,
	}
}
