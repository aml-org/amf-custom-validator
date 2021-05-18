package validator

import (
	"github.com/piprate/json-gold/ld"
)


func Normalize(json interface{}) interface{} {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	context := map[string]interface{}{
		"shacl": "http://www.w3.org/ns/shacl#",
		"shapes": "http://a.ml/vocabularies/shapes#",
		"doc": "http://a.ml/vocabularies/document#",
		"apiContract": "http://a.ml/vocabularies/apiContract#",
		"core": "http://a.ml/vocabularies/core#",
	}
	flattened, err := proc.Flatten(json,context, options)
	if err != nil {
		panic(err)
	}

	return flattened
}
