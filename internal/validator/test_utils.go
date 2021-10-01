package validator

import (
	"fmt"
	p "github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/pkg/milestones"
	"io/ioutil"
	"strings"
)

const debug = false

func conforms(report string) bool {
	return strings.Index(report, "\"conforms\": true") > -1
}

func printMilestones(profile interface{}, data string, milestonesChan *chan milestones.Milestone) {
	for m := range *milestonesChan {
		fmt.Printf("%s,%s,%s,%d\n", profile, data, m.Operation, m.Duration.Microseconds())
	}
}

func read(path relativePath) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func validate(profile relativePath, data relativePath) string {
	p.GenReset()
	profileText := read(profile)
	dataText := read(data)
	report, err := Validate(profileText, dataText, false, nil)
	if err != nil {
		panic(err)
	}
	return report
}

func compare(actualText string, expected relativePath) bool {
	expectedText := read(expected)
	return strings.TrimSpace(actualText) == strings.TrimSpace(expectedText)
}

func write(content string, path relativePath) {
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

type relativePath = string
