package config

type ReportConfiguration struct {
	IncludeReportCreationTime bool
}

func DefaultReportConfiguration() ReportConfiguration {
	return ReportConfiguration{
		IncludeReportCreationTime: true,
	}
}
