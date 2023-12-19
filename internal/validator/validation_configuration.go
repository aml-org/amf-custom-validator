package validator

import "time"

type ValidationConfiguration interface {
	ReportCreationTime() time.Time
}

type DefaultValidationConfiguration struct{}

func (d DefaultValidationConfiguration) ReportCreationTime() time.Time {
	return time.Now()
}

type TestValidationConfiguration struct{}

func (d TestValidationConfiguration) ReportCreationTime() time.Time {
	return time.Date(2000, time.November, 28, 0, 0, 0, 0, time.UTC)
}
