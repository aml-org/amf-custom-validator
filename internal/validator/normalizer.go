package validator

import (
	"github.com/aml-org/amfopa/internal/parser/profile"
	"github.com/piprate/json-gold/ld"
)

func Normalize(json interface{}, prefixes profile.ProfileContext) interface{} {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	context := map[string]interface{}{
		"data":        "http://a.ml/vocabularies/data#",
		"shacl":       "http://www.w3.org/ns/shacl#",
		"shapes":      "http://a.ml/vocabularies/shapes#",
		"doc":         "http://a.ml/vocabularies/document#",
		"apiContract": "http://a.ml/vocabularies/apiContract#",
		"core":        "http://a.ml/vocabularies/core#",
		"xsd":         "http://www.w3.org/2001/XMLSchema#",
	}
	for n, p := range prefixes {
		context[n] = p
	}
	flattened, err := proc.Flatten(json, context, options)
	if err != nil {
		panic(err)
	}

	return flattened
}
