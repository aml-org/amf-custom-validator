package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
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
		"security":    "http://a.ml/vocabularies/security#",
		"sourcemaps":  "http://a.ml/vocabularies/document-source-maps#",
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

	// Build lexical index
	lexicalIndex := make(map[string]interface{})
	for _, sourceMapId := range classIndex["sourcemaps:SourceMap"] {
		sourceMap := nodeIndex[sourceMapId].(map[string]interface{})
		lexicalContainer := sourceMap["sourcemaps:lexical"] // can be map or array of maps

		switch v := lexicalContainer.(type) {
		case map[string]interface{}:
			addLexicalEntryFrom(&v, &nodeIndex, &lexicalIndex)
		case []interface{}:
			for _, e := range v {
				switch vv := e.(type) {
				case map[string]interface{}:
					addLexicalEntryFrom(&vv, &nodeIndex, &lexicalIndex)
				}
			}
		}
	}

	return map[string]interface{}{
		"@ids":     nodeIndex,
		"@types":   classIndex,
		"@lexical": lexicalIndex,
	}
}

func addLexicalEntryFrom(node, nodeIndex, lexicalIndex *map[string]interface{}) {
	lexicalEntry := (*nodeIndex)[(*node)["@id"].(string)].(map[string]interface{})
	id := lexicalEntry["sourcemaps:element"].(string)
	value := lexicalEntry["sourcemaps:value"]

	/**
	Index:
		nodeId -> lexical

	Cannot index property lexical info (property URI -> lexical) because property URIs are not unique and will
	get overwritten by each node
	*/
	if _, ok := (*nodeIndex)[id]; ok {
		(*lexicalIndex)[id] = value
	}
}
