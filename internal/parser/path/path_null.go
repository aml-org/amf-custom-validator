package path

import "github.com/aml-org/amf-custom-validator/internal/misc"

type NullPath struct {
	BasePath
}

func (p NullPath) Trace(iriExpander *misc.IriExpander) (string, error) {
	return p.source, nil
}

func (p NullPath) Source() string {
	return p.source
}

func (p NullPath) Expanded(iriExpander *misc.IriExpander) (string, error) {
	return p.source, nil // should always be ""?
}
