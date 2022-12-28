package validator

import (
	"errors"
	"github.com/aml-org/amf-custom-validator/internal/generator"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	p "github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

/**
We remove nodes that will not be "visited" during validation. Algorithm:
 1. Parse the profile and extract all the traversals it encodes
 2. Perform those traversals on the actual input and compute the list of visited nodes
 3. Remove unvisited nodes from JSON-LD

Considerations:
 - This should be done before indexing
 - We need to keep lexical nodes despite not actually being traversed
 - We cannot perform optimization with custom Rego code (we cannot extract the traversals)
*/

type Optimizer struct {
	traversalsIndex *map[string][]*Traversal // indexed by class
}

type Traversal struct {
	class string
	path  *path.PropertyPath
}

func NewOptimizer(profile *p.Profile) Optimizer {
	// Obtain all classes
	traversals := extractTraversals(&profile.Violation)
	iriExpander := generator.IriExpanderFrom(*profile)
	traversalsIndex := indexTraversals(traversals, iriExpander)
	return Optimizer{traversalsIndex: traversalsIndex}
}

// TODO add error handling!!
func (o *Optimizer) Optimize(input *interface{}) {
	nodes, _ := get[[]interface{}](input, "@graph")
	toKeepNodes := make(map[string]interface{})
	for _, node := range nodes {
		id, _ := get[string](&node, "@id")
		classes, _ := get[[]interface{}](&node, "@type")
		if toKeepNodes[id] == nil {
			for _, class := range classes {
				classStr := class.(string)
				if (*(o.traversalsIndex))[classStr] != nil {
					toKeepNodes[id] = node
					break
				}
			}
		}
	}
	var newGraph []interface{}
	for _, node := range toKeepNodes {
		newGraph = append(newGraph, node)
	}

	set(input, "@graph", newGraph)
}

func set(obj *interface{}, prop string, value interface{}) {
	switch o := (*obj).(type) {
	case map[string]interface{}:
		o[prop] = value
	}
}

func get[K interface{}](obj *interface{}, prop string) (K, error) {
	var empty K
	switch o := (*obj).(type) {
	case map[string]interface{}:
		switch oo := o[prop].(type) {
		case K:
			return oo, nil
		default:
			return empty, errors.New("'obj[prop]' is not of type K")
		}
	default:
		return empty, errors.New("argument 'obj' is not of type map[string]interface{}")
	}
}

func extractTraversals(rules *[]p.Rule) *[]Traversal {
	var result []Traversal
	for _, rule := range *rules {
		switch topRule := rule.(type) {
		case p.TopLevelExpression:
			class := topRule.ClassGenerator
			result = extractTraversalsFromRule(class, &topRule.Expression.Value, result)
		default:
			break
		}
	}
	return &result
}

func extractTraversalsFromRule(class string, rule *p.Rule, acc []Traversal) []Traversal {
	result := acc
	switch r := (*rule).(type) {
	case p.CountRule:
		result = append(result, Traversal{class: class, path: &r.Path})
		break
	case p.AndRule:
		for _, innerRule := range r.Body {
			c := extractTraversalsFromRule(class, &innerRule, result)
			result = append(result, c...)
		}
		break
	default:
		break
	}
	return result
}

func indexTraversals(traversals *[]Traversal, expander *misc.IriExpander) *map[string][]*Traversal {
	result := make(map[string][]*Traversal)
	for _, traversal := range *traversals {
		key, _ := expander.Expand(traversal.class)
		result[key] = append(result[key], &traversal)
	}
	return &result
}
