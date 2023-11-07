package misc

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/types"
	"regexp"
	"strings"
)

//IriExpander expands compacted IRIs in the form prefix.suffix
type IriExpander struct {
	Context types.ObjectMap
}

func (i *IriExpander) Expand(iri string) (string, error) {
	isReservedKeyword := strings.HasPrefix(iri, "@")
	compactForm := regexp.MustCompile("^[a-zA-Z-0-9\\-]+\\.[\\.(\\\\/)a-zA-Z-0-9\\-]+$")
	isCompact := compactForm.MatchString(iri)

	if isCompact {
		return i.expandCompactIri(iri)
	} else if isReservedKeyword {
		return iri, nil
	} else {
		return iri, errors.New(fmt.Sprintf("IRI %s is not in compact form", iri))
	}
}

func (i *IriExpander) expandCompactIri(iri string) (string, error) {
	split := strings.SplitN(iri, ".", 2)
	prefix := split[0]
	suffix := strings.ReplaceAll(split[1], "\\/", "/")
	switch e := i.Context[prefix].(type) {
	case string:
		return e + suffix, nil
	default:
		return iri, errors.New(fmt.Sprintf("Term %s not present in context", prefix))
	}
}
