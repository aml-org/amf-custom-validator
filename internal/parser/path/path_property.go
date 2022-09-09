package path

import (
	"errors"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	"strings"
)

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

func (p Property) IsCustom(iriExpander *misc.IriExpander) bool {
	expanded, err := iriExpander.Expand(p.Iri)
	if err != nil {
		return false
	}
	return strings.Index(expanded, contexts.ApiExtensionUri) == 0
}

func (p Property) CustomName(iriExpander *misc.IriExpander) (string, error) {
	expanded, err := iriExpander.Expand(p.Iri)
	if err != nil {
		return "", errors.New("Invalid custom property")
	}

	return strings.Split(expanded, "#")[1], nil
}
