package validator

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/pkg/milestones"
	"strings"
)

const debug = false

func conforms(report string) bool {
	return strings.Index(report, "\"http://www.w3.org/ns/shacl#conforms\": true") > -1
}

func printMilestones(profile interface{}, data string, milestonesChan *chan milestones.Milestone) {
	for m := range *milestonesChan {
		fmt.Printf("%s,%s,%s,%d\n", profile, data, m.Operation, m.Duration.Microseconds())
	}
}
