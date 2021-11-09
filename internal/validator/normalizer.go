package validator

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/internal/types"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	"github.com/piprate/json-gold/ld"
)

func Normalize(json interface{}, prefixes profile.ProfileContext) interface{} {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	context := make(types.ObjectMap)
	types.MergeObjectMap(&context, &contexts.DefaultAMFContext)
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
	nodeIndex := make(types.ObjectMap)

	g := json.(types.ObjectMap)["@graph"]
	nodes := g.([]interface{})

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
	lexicalIndex := make(types.ObjectMap)
	for _, sourceMapId := range classIndex["sourcemaps:SourceMap"] {
		sourceMap := nodeIndex[sourceMapId].(types.ObjectMap)
		lexicalContainer := sourceMap["sourcemaps:lexical"] // can be map or array of maps
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
	id := lexicalEntry["sourcemaps:element"].(string)
	value := lexicalEntry["sourcemaps:value"]

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
	sourceInformation := (*classIndex)["doc:BaseUnitSourceInformation"]
	if len(sourceInformation) > 0 {
		sourceInformationNode := (*nodeIndex)[sourceInformation[0]].(types.ObjectMap)
		defaultLocation := sourceInformationNode["doc:rootLocation"].(string)
		additionalLocations := sourceInformationNode["doc:additionalLocations"]
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
	locationValue := locationNode["doc:location"].(string)
	elementIds := locationNode["doc:elements"]
	handleSingleOrMultipleNodes(&elementIds, func(node *types.ObjectMap) {
		(*idToLocation)[(*node)["@id"].(string)] = locationValue
	})
}

func handleSingleOrMultipleNodes(node *interface{}, operation func(*types.ObjectMap)) {
	switch v := (*node).(type) {
	case types.ObjectMap: // single node
		operation(&v)
	case []interface{}: // array with multiple nodes
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
