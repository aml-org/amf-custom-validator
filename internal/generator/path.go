package generator

import (
	"errors"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
)

type RegoPathResult struct {
	source   string
	rule     string
	rego     []string
	variable string
}

type regoPathResultInternal struct {
	rego          []string
	pathVariables []string
	paths         []string
	variable      string
	counter       *int
}

type traversal struct {
	variable      string
	counter       *int
	rego          []string
	pathVariables []string
	paths         []string
}

func newTraversal(variable string) traversal {
	c := 0
	return traversal{
		variable:      variable,
		counter:       &c,
		rego:          make([]string, 0),
		pathVariables: make([]string, 0),
		paths:         make([]string, 0),
	}
}

func internalResultToTraversal(p traversal, r regoPathResultInternal) traversal {
	return traversal{
		variable:      p.variable,
		counter:       r.counter,
		rego:          r.rego,
		pathVariables: r.pathVariables,
		paths:         r.paths,
	}
}

// GeneratePropertySet Traversed the path, starting at the provided variable and returns a set of reached values
func GeneratePropertySet(path path.PropertyPath, variable string, iriExpander *misc.IriExpander) RegoPathResult {
	return generateResult(path, variable, false, iriExpander, aggregateResultsIntoSet)
}

// GenerateNodeSet Traverses the path but just returns a generator of nodes instead of a set of values
func GenerateNodeSet(path path.PropertyPath, variable string, iriExpander *misc.IriExpander) RegoPathResult {
	return generateResult(path, variable, true, iriExpander, aggregateResultsIntoSet)
}

// GeneratePropertyArray Traversed the path, starting at the provided variable and returns an array of reached values (allows duplicate values)
func GeneratePropertyArray(path path.PropertyPath, variable string, iriExpander *misc.IriExpander) RegoPathResult {
	return generateResult(path, variable, false, iriExpander, aggregateResultsIntoArray)
}

// GenerateNodeArray Traverses the path but just returns a generator of nodes instead of an array of values (allows duplicate values)
func GenerateNodeArray(path path.PropertyPath, variable string, iriExpander *misc.IriExpander) RegoPathResult {
	return generateResult(path, variable, false, iriExpander, aggregateResultsIntoArray)
}

// fetchNodes defines if the path result will fetch resulting nodes (fetching each @id reference) or simple return the values.
func generateResult(path path.PropertyPath, variable string, fetchNodes bool, iriExpander *misc.IriExpander, resultAggregator func([]regoPathResultInternal) RegoPathResult) RegoPathResult {
	source, err := path.Expanded(iriExpander)
	if err != nil {
		panic(err)
	}

	acc := traversePath(path, variable, fetchNodes, iriExpander)
	regoResult := resultAggregator(acc)

	regoResult.source = source
	regoResult.variable = variable

	return regoResult
}

func traversePath(path path.PropertyPath, variable string, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	t := newTraversal(variable)
	var acc []regoPathResultInternal
	for _, tr := range traverse(path, t, fetchNodes, iriExpander) {
		effectiveRego := tr.rego
		effectiveRego = append(effectiveRego, fmt.Sprintf("nodes = %s", tr.variable))
		tr.rego = effectiveRego
		acc = append(acc, tr)
	}
	return acc
}

// Since there are many alternative paths to reach the nodes, we need to provide a single
// collection of nodes of the rest of the checks.
func aggregateResultsIntoSet(paths []regoPathResultInternal) RegoPathResult {
	// Let's generate a rule that will return the flat list of nodes in the path
	// If there are more than one path (because of ORs) a rule with multiple clauses
	// will be generated and the final list of nodes will be the UNION of all the clauses
	rego := make([]string, 0)
	ruleName := profile.Genvar("path_set_rule")
	for i, p := range paths {
		if i == 0 {
			rego = append(rego, fmt.Sprintf("%s[nodes] {", ruleName)) // header of the rule
		} else {
			rego = append(rego, "} {") // add another clause to the rule
		}
		for _, r := range p.rego {
			rego = append(rego, "  "+r) // add the rego code to the final rule
		}
	}
	if len(rego) > 0 {
		rego = append(rego, "}")
	}

	return RegoPathResult{
		rego: rego,
		rule: ruleName,
	}
}

func aggregateResultsIntoArray(paths []regoPathResultInternal) RegoPathResult {
	rego := make([]string, 0)
	ruleName := profile.Genvar("path_array_rule")
	for i, p := range paths {
		if i == 0 {
			rego = append(rego, fmt.Sprintf("%s = [ nodes | ", ruleName)) // header of the rule
		} else {
			rego = append(rego, "} {") // add another clause to the rule // TODO ?
		}
		for _, r := range p.rego {
			rego = append(rego, "  "+r) // add the rego code to the final rule
		}
	}
	if len(rego) > 0 {
		rego = append(rego, "]")
	}

	return RegoPathResult{
		rego: rego,
		rule: ruleName,
	}
}

// Different traversals based on the type of path element
func traverse(propPath path.PropertyPath, traversed traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	switch p := propPath.(type) {
	case path.Property:
		return traverseProperty(p, traversed, fetchNodes, iriExpander)
	case path.AndPath:
		return traverseAnd(p, traversed, fetchNodes, iriExpander)
	case path.OrPath:
		return traverseOr(p, traversed, fetchNodes, iriExpander)
	default:
		return make([]regoPathResultInternal, 0)
	}
}

// Traverses in parallel each of the elements in the OR path, creating new branches for each one
func traverseOr(or path.OrPath, t traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	acc := make([]regoPathResultInternal, 0)
	for _, p := range or.Or {
		traversed := traverse(p, t, fetchNodes, iriExpander)
		for _, tr := range traversed {
			acc = append(acc, tr)
		}
	}
	return acc
}

// Traverses an AND branch in the property path.
// it sequentially traverses the path until all the elements of the AND
// has been traversed and the variables accumulated.
func traverseAnd(and path.AndPath, t traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	first := and.And[0]
	if len(and.And) > 1 {
		firstTraversed := traverse(first, t, true, iriExpander)
		remaining := and.And[1:len(and.And)]
		acc := make([]regoPathResultInternal, 0)
		for _, tr := range firstTraversed {
			next := path.AndPath{
				And: remaining,
			}
			for _, ntr := range traverse(next, internalResultToTraversal(t, tr), fetchNodes, iriExpander) {
				acc = append(acc, ntr)
			}
		}
		return acc
	} else {
		return traverse(first, t, fetchNodes, iriExpander)
	}
}

func traverseProperty(property path.Property, t traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	if property.IsCustom(iriExpander) {
		return traverseCustomProperty(property, t, fetchNodes, iriExpander)
	} else {
		return traverseRegularProperty(property, t, fetchNodes, iriExpander)
	}
}

// Traverses the leaf components of the path expression, always a property.
// TODO: We don't take into transitive paths yet.
func traverseRegularProperty(property path.Property, t traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {

	propertyIri, err := property.Expanded(iriExpander)

	// We use IDX go generate a unique property for the Rego computation
	idx := fmt.Sprintf("%d", len(t.pathVariables))
	if *t.counter > 0 {
		idx = fmt.Sprintf("%s_%d", idx, t.counter)
	}
	binding := fmt.Sprintf("%s_%s", t.variable, idx)

	if len(t.pathVariables) == 0 {
		// If this is the first element in the path, we start computing the path from the previous variable passed
		// to the path generator, usually a classTarget.
		t.rego = append(t.rego, fmt.Sprintf("init_%s = data.sourceNode", binding)) // this is the connection to the variable past to the generator
		t.pathVariables = append(t.pathVariables, fmt.Sprintf("init_%s", binding)) // initial value
	}
	source := t.pathVariables[len(t.pathVariables)-1]

	if err != nil {
		panic(err)
	}

	if property.Inverse {
		t.rego = append(t.rego, fmt.Sprintf("search_subjects[%s] with data.predicate as \"%s\" with data.object as %s", binding, propertyIri, source))
	} else {
		// fetch resulting nodes (fetching each @id reference) or simply return the array values
		if fetchNodes {
			t.rego = append(t.rego, fmt.Sprintf("tmp_%s = nested_nodes with data.nodes as %s[\"%s\"]", binding, source, propertyIri))
			t.rego = append(t.rego, fmt.Sprintf("%s = tmp_%s[_][_]", binding, binding))
		} else {
			t.rego = append(t.rego, fmt.Sprintf("nodes_tmp = object.get(%s,\"%s\",[])", source, propertyIri))
			t.rego = append(t.rego, "nodes_tmp2 = nodes_array with data.nodes as nodes_tmp") // this returns and array
			t.rego = append(t.rego, fmt.Sprintf("%s = nodes_tmp2[_]", binding))
		}
	}

	r := regoPathResultInternal{
		rego:          t.rego,
		pathVariables: append(t.pathVariables, binding),
		paths:         append(t.paths, property.Iri),
		counter:       t.counter,
		variable:      binding,
	}

	return []regoPathResultInternal{r}
}

func traverseCustomProperty(property path.Property, t traversal, fetchNodes bool, iriExpander *misc.IriExpander) []regoPathResultInternal {
	customPropertyName, err := property.CustomName(iriExpander)
	if err != nil {
		panic(err)
	}
	// We use IDX go generate a unique property for the Rego computation
	idx := fmt.Sprintf("%d", len(t.pathVariables))
	if *t.counter > 0 {
		idx = fmt.Sprintf("%s_%d", idx, t.counter)
	}
	binding := fmt.Sprintf("%s_%s", t.variable, idx)

	if len(t.pathVariables) == 0 {
		// If this is the first element in the path, we start computing the path from the previous variable passed
		// to the path generator, usually a classTarget.
		t.rego = append(t.rego, fmt.Sprintf("init_%s = data.sourceNode", binding)) // this is the connection to the variable past to the generator
		t.pathVariables = append(t.pathVariables, fmt.Sprintf("init_%s", binding)) // initial value
	}
	source := t.pathVariables[len(t.pathVariables)-1]

	if property.Inverse {
		panic(errors.New("Inverse custom properties not supported yet"))
	} else {
		t.rego = append(t.rego, fmt.Sprintf("tmp_%s = gen_path_extension with data.custom_property_data as [%s, \"%s\"]", binding, source, customPropertyName))
		// fetch resulting nodes (fetching each @id reference) or simply return the array values
		if fetchNodes {
			t.rego = append(t.rego, fmt.Sprintf("%s = tmp_%s[_][_]", binding, binding))
		} else {
			t.rego = append(t.rego, fmt.Sprintf("tmp2_%s = tmp_%s[_][_]", binding, binding))
			t.rego = append(t.rego, fmt.Sprintf("%s = object.get(tmp2_%s,\"@id\",\"\")", binding, binding))
		}
	}

	r := regoPathResultInternal{
		rego:          t.rego,
		pathVariables: append(t.pathVariables, binding),
		paths:         append(t.paths, property.Iri),
		counter:       t.counter,
		variable:      binding,
	}

	return []regoPathResultInternal{r}
}
