package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/config"
	p "github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/pkg/milestones"
	"io/ioutil"
	"strings"
	"testing"
)

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
	report, err := Validate(profileText, dataText, config.Debug, nil)
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

/**
Directory file naming convention
Profile -> profile.yaml
Data -> data.jsonld
Report -> report.jsonld
*/
func validateAndCompareDirectory(directory relativePath, t *testing.T) {
	resolvedDirectory := fmt.Sprintf("%s", directory)
	profile := fmt.Sprintf("%s/profile.yaml", resolvedDirectory)
	data := fmt.Sprintf("%s/data.jsonld", resolvedDirectory)
	actualText := validate(profile, data)
	expected := fmt.Sprintf("%s/report.jsonld", resolvedDirectory)
	if config.Override {
		write(actualText, expected)
	} else {
		if !compare(actualText, expected) {
			t.Errorf("Failed %s. Actual did not match expexted", directory)
		}
	}
}

func validateAndCompare(profile, data, expected string, t *testing.T) {
	actualText := validate(profile, data)
	if config.Override {
		write(actualText, expected)
	} else {
		if !compare(actualText, expected) {
			println("Expected")
			println("==================")
			println(expected)
			println("\n\n\nActual")
			println("==================")
			println(actualText)
			t.Errorf("Failed %s. Actual did not match expexted", data)
		}
	}
}

func ignoreDirectory(directory relativePath, t *testing.T) {
	t.Skipf("Ignored %s", directory)
}

type relativePath = string
