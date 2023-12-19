package config

type ReportConfiguration struct {
	IncludeReportCreationTime bool
	ReportSchemaIri           string
	LexicalSchemaIri          string
}

func DefaultReportConfiguration() ReportConfiguration {
	return ReportConfiguration{
		IncludeReportCreationTime: true,
		ReportSchemaIri:           "file:///dialects/validation-report.yaml",
		LexicalSchemaIri:          "file:///dialects/lexical.yaml",
	}
}
