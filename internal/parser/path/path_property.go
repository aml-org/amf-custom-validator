package path

import "github.com/aml-org/amf-custom-validator/internal/misc"

type Property struct {
	BasePath
	Iri        string
	Inverse    bool
	Transitive bool
}

func (p Property) Trace(iriExpander *misc.IriExpander) (string, error) {
	expanded, err := p.Expanded(iriExpander)
	if p.Inverse && err == nil {
		return expanded + "^", err
	} else {
		return expanded, err
	}
}

func (p Property) Expanded(iriExpander *misc.IriExpander) (string, error) {
	if iriExpander != nil {
		return iriExpander.Expand(p.Iri)
	} else {
		return p.source, nil
	}
}

func (p Property) Source() string {
	return p.source
}
