package config

type ReportConfiguration struct {
	IncludeReportCreationTime bool
	ReportSchemaIri           string
	LexicalSchemaIri          string
	BaseIri                   string
}

func DefaultReportConfiguration() ReportConfiguration {
	return ReportConfiguration{
		IncludeReportCreationTime: true,
		ReportSchemaIri:           "file:///dialects/validation-report.yaml",
		LexicalSchemaIri:          "file:///dialects/lexical.yaml",
		BaseIri:                   "http://a.ml/vocabularies/validation/report#",
	}
}
