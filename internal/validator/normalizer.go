package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/types"
	"github.com/piprate/json-gold/ld"
)

func Normalize(json any) any {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	context := make(types.ObjectMap)
	flattened, err := proc.Flatten(json, context, options)
	if err != nil {
		panic(err)
	}

	return flattened
}

func Index(json any) any {
	classIndex := make(map[string][]string)
	nodeIndex := make(types.ObjectMap)

	g := json.(types.ObjectMap)["@graph"]
	nodes := g.([]any)

	for _, nn := range nodes {
		n := nn.(types.ObjectMap)
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
		case []any:
			for _, cc := range classes.([]any) {
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

	var locationIndex = createLocationIndex(&nodeIndex, &classIndex)

	// Build lexical index
	lexicalIndex := make(types.ObjectMap)
	for _, sourceMapId := range classIndex["http://a.ml/vocabularies/document-source-maps#SourceMap"] {
		sourceMap := nodeIndex[sourceMapId].(types.ObjectMap)
		lexicalContainer := sourceMap["http://a.ml/vocabularies/document-source-maps#lexical"] // can be map or array of maps
		handleSingleOrMultipleNodes(&lexicalContainer, func(node *types.ObjectMap) {
			addLexicalEntryFrom(node, &nodeIndex, &lexicalIndex, locationIndex)
		})
	}

	return types.ObjectMap{
		"@ids":     nodeIndex,
		"@types":   classIndex,
		"@lexical": lexicalIndex,
	}
}

func addLexicalEntryFrom(node, nodeIndex, lexicalIndex *types.ObjectMap, locIndex *LocationIndex) {
	lexicalEntry := (*nodeIndex)[(*node)["@id"].(string)].(types.ObjectMap)
	id := lexicalEntry["http://a.ml/vocabularies/document-source-maps#element"].(string)
	value := lexicalEntry["http://a.ml/vocabularies/document-source-maps#value"]

	/**
	Index:
		nodeId -> (range , uri)

	Cannot index property lexical info (property URI -> lexical) because property URIs are not unique and will
	get overwritten by each node
	*/
	if _, ok := (*nodeIndex)[id]; ok {
		(*lexicalIndex)[id] = types.ObjectMap{
			"range": value,
			"uri":   locIndex.Location(id),
		}
	}
}

func createLocationIndex(nodeIndex *types.ObjectMap, classIndex *map[string][]string) *LocationIndex {
	sourceInformation := (*classIndex)["http://a.ml/vocabularies/document#BaseUnitSourceInformation"]
	if len(sourceInformation) > 0 {
		sourceInformationNode := (*nodeIndex)[sourceInformation[0]].(types.ObjectMap)
		defaultLocation := sourceInformationNode["http://a.ml/vocabularies/document#rootLocation"].(string)
		additionalLocations := sourceInformationNode["http://a.ml/vocabularies/document#additionalLocations"]
		idToLocation := make(types.StringMap)
		handleSingleOrMultipleNodes(&additionalLocations, func(node *types.ObjectMap) {
			addElementsOfLoc(node, nodeIndex, &idToLocation)
		})
		return &LocationIndex{DefaultLocation: defaultLocation, IdToLocation: idToLocation}

	} else {
		return &LocationIndex{IdToLocation: make(types.StringMap), DefaultLocation: ""}
	}
}

func addElementsOfLoc(node *types.ObjectMap, nodeIndex *types.ObjectMap, idToLocation *types.StringMap) {
	locationNode := (*nodeIndex)[(*node)["@id"].(string)].(types.ObjectMap)
	locationValue := locationNode["http://a.ml/vocabularies/document#location"].(string)
	elementIds := locationNode["http://a.ml/vocabularies/document#elements"]
	handleSingleOrMultipleNodes(&elementIds, func(node *types.ObjectMap) {
		(*idToLocation)[(*node)["@id"].(string)] = locationValue
	})
}

func handleSingleOrMultipleNodes(node *any, operation func(*types.ObjectMap)) {
	switch v := (*node).(type) {
	case types.ObjectMap: // single node
		operation(&v)
	case []any: // array with multiple nodes
		for _, e := range v {
			switch vv := e.(type) {
			case types.ObjectMap:
				operation(&vv)
			}
		}
	default:
	}
}

type LocationIndex struct {
	DefaultLocation string
	IdToLocation    types.StringMap
}

func (locIndex *LocationIndex) Location(id string) string {
	value, exists := locIndex.IdToLocation[id]
	if exists {
		return value
	} else {
		return locIndex.DefaultLocation
	}
}
