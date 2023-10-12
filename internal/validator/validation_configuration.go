package validator

import "time"

type ValidationConfiguration interface {
	CurrentTime() time.Time
}

type DefaultConfiguration struct{}

func (d DefaultConfiguration) CurrentTime() time.Time {
	return time.Now()
}

type TestConfiguration struct{}

func (d TestConfiguration) CurrentTime() time.Time {
	return time.Date(2000, time.November, 28, 0, 0, 0, 0, time.UTC)
}
