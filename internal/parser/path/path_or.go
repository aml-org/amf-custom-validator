package path

import "github.com/aml-org/amf-custom-validator/internal/misc"

type OrPath struct {
	BasePath
	Or []PropertyPath
}

func (p OrPath) Trace(iriExpander *misc.IriExpander) (string, error) {
	return expandAndJoinParts(iriExpander, p.Or, " | ", true)
}

func (p OrPath) Expanded(iriExpander *misc.IriExpander) (string, error) {
	return expandAndJoinParts(iriExpander, p.Or, " | ", false)
}

func (p OrPath) Source() string {
	return p.source
}
