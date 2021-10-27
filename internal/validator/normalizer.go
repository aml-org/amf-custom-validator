package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/piprate/json-gold/ld"
)

type Object = map[string]interface{}

func Normalize(json interface{}, prefixes profile.ProfileContext) interface{} {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	context := Object{
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
	nodeIndex := make(Object)

	g := json.(Object)["@graph"]
	nodes := g.([]interface{})

	for _, nn := range nodes {
		n := nn.(Object)
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

	var locationIndex *LocationIndex = createLocationIndex(&nodeIndex, &classIndex)

	// Build lexical index
	lexicalIndex := make(Object)
	for _, sourceMapId := range classIndex["sourcemaps:SourceMap"] {
		sourceMap := nodeIndex[sourceMapId].(Object)
		lexicalContainer := sourceMap["sourcemaps:lexical"] // can be map or array of maps
		handleSingleOrMultipleNodes(&lexicalContainer, func(node *Object){
			addLexicalEntryFrom(node, &nodeIndex, &lexicalIndex, locationIndex)
		})
	}

	return Object{
		"@ids":     nodeIndex,
		"@types":   classIndex,
		"@lexical": lexicalIndex,
	}
}

func addLexicalEntryFrom(node, nodeIndex, lexicalIndex *Object, locIndex *LocationIndex) {
	lexicalEntry := (*nodeIndex)[(*node)["@id"].(string)].(Object)
	id := lexicalEntry["sourcemaps:element"].(string)
	value := lexicalEntry["sourcemaps:value"]

	/**
	Index:
		nodeId -> (range , uri)

	Cannot index property lexical info (property URI -> lexical) because property URIs are not unique and will
	get overwritten by each node
	*/
	if _, ok := (*nodeIndex)[id]; ok {
		(*lexicalIndex)[id] = Object{
			"range": value,
			"uri":   locIndex.Location(id),
		}
	}
}

func createLocationIndex(nodeIndex *Object, classIndex *map[string][]string) *LocationIndex {
	sourceInformation := (*classIndex)["doc:BaseUnitSourceInformation"]
	if len(sourceInformation) > 0 {
		sourceInformationNode := (*nodeIndex)[sourceInformation[0]].(Object)
		defaultLocation := sourceInformationNode["doc:rootLocation"].(string)
		additionalLocations := sourceInformationNode["doc:additionalLocations"]
		idToLocation := make(map[string]string)
		handleSingleOrMultipleNodes(&additionalLocations, func(node *Object){
			addElementsOfLoc(node, nodeIndex, &idToLocation)
		})
		return &LocationIndex{DefaultLocation: defaultLocation, IdToLocation: idToLocation}

	} else {
		return &LocationIndex{IdToLocation: make(map[string]string), DefaultLocation: ""}
	}
}

func addElementsOfLoc(node *Object, nodeIndex *Object, idToLocation *map[string]string) {
	locationNode := (*nodeIndex)[(*node)["@id"].(string)].(Object)
	locationValue := locationNode["doc:location"].(string)
	elementIds := locationNode["doc:elements"]
	handleSingleOrMultipleNodes(&elementIds, func(node *Object){
		(*idToLocation)[(*node)["@id"].(string)] = locationValue
	})
}

func handleSingleOrMultipleNodes(node *interface{}, operation func(*Object)) {
	switch v := (*node).(type) {
	case Object: // single node
		operation(&v)
	case []interface{}: // array with multiple nodes
		for _, e := range v {
			switch vv := e.(type) {
			case Object:
				operation(&vv)
			}
		}
	default:
	}
}


type LocationIndex struct {
	DefaultLocation string
	IdToLocation map[string]string
}

func (locIndex *LocationIndex) Location(id string) string {
	value, exists := locIndex.IdToLocation[id]
	if exists {
		return value
	} else {
		return locIndex.DefaultLocation
	}
}
