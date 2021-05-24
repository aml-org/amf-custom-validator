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
		"rdfs":        "http://www.w3.org/2000/01/rdf-schema",
		"rdf":         "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
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

func Index(json interface{}) interface{} {
	classIndex := make(map[string][]string)
	nodeIndex := make(map[string]interface{})

	g := json.(map[string]interface{})["@graph"]
	nodes := g.([]interface{})

	for _, nn := range nodes {
		n := nn.(map[string]interface{})
		id := n["@id"].(string)
		classes := n["@type"]
		nodeIndex[id] = n
		switch cc := classes.(type) {
		case string:
			acc, ok := classIndex[cc]
			if !ok {
				acc = make([]string, 0)
			}
			acc = append(acc, id)
			classIndex[cc] = acc
		case []interface{}:
			for _, cc := range classes.([]interface{}) {
				c := cc.(string)
				acc, ok := classIndex[c]
				if !ok {
					acc = make([]string, 0)
				}
				acc = append(acc, id)
				classIndex[c] = acc
			}
		}
	}

	return map[string]interface{}{
		"@ids":   nodeIndex,
		"@types": classIndex,
	}
}
