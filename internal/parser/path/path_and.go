package path

import "github.com/aml-org/amf-custom-validator/internal/misc"

type AndPath struct {
	BasePath
	And []PropertyPath
}

func (p AndPath) Trace(iriExpander *misc.IriExpander) (string, error) {
	return expandAndJoinParts(iriExpander, p.And, " / ", true)
}

func (p AndPath) Expanded(iriExpander *misc.IriExpander) (string, error) {
	return expandAndJoinParts(iriExpander, p.And, " / ", false)
}

func (p AndPath) Source() string {
	return p.source
}