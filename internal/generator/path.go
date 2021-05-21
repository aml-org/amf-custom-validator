package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/profile"
	"strings"
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
	id            string
	variable      string
	hint          string
	counter       *int
	rego          []string
	pathVariables []string
	paths         []string
}

func newTraversal(source string, variable string, hint string) traversal {
	c := 0
	return traversal{
		id:            valueHash(source),
		variable:      variable,
		hint:          hint,
		counter:       &c,
		rego:          make([]string, 0),
		pathVariables: make([]string, 0),
		paths:         make([]string, 0),
	}
}

func valueHash(source string) string {
	s := strings.ToLower(source)
	return internal.HashString(s)
}

func internalResultToTraversal(p traversal, r regoPathResultInternal) traversal {
	return traversal{
		id:            p.id,
		variable:      p.variable,
		hint:          p.hint,
		counter:       r.counter,
		rego:          r.rego,
		pathVariables: r.pathVariables,
		paths:         r.paths,
	}
}

// Traversed the path path, starting at the provided variable and returns an array of reached values
func GeneratePropertyArray(path path.PropertyPath, variable string, hint string) RegoPathResult {
	t := newTraversal(path.Source(), variable, hint)
	var acc []regoPathResultInternal
	for _, tr := range traverse(path, t) {
		effectiveRego := tr.rego[0 : len(tr.rego)-2] // we remove the last element so we can get the array of values instead

		previousBinding := tr.pathVariables[len(tr.pathVariables)-2]

		nextPath := tr.paths[len(tr.paths)-1]
		effectiveRego = append(effectiveRego, fmt.Sprintf("nodes_tmp = object.get(%s,\"%s\",[])", previousBinding, nextPath))
		effectiveRego = append(effectiveRego, "nodes_tmp2 = nodes_array with data.nodes as nodes_tmp") // this returns and array
		effectiveRego = append(effectiveRego, "nodes = nodes_tmp2[_]")                                 // I need to iterate to each element in the array so it can be wrapped in the rule result
		tr.rego = effectiveRego
		acc = append(acc, tr)
	}

	regoResult := accumulatePaths(acc)
	regoResult.source = path.Source()
	regoResult.variable = variable

	return regoResult
}

// Traverses the path but just returns a generator of nodes instead of a set of values
func GenerateNodeArray(path path.PropertyPath, variable string, hint string) RegoPathResult {
	t := newTraversal(path.Source(), variable, hint)
	var acc []regoPathResultInternal
	for _, tr := range traverse(path, t) {
		effectiveRego := tr.rego[0 : len(tr.rego)-2] // we remove the last element so we can get the array of values instead
		previousBinding := variable
		previousBinding = tr.pathVariables[len(tr.pathVariables)-2]
		nextPath := tr.paths[len(tr.paths)-1]
		// we just get the nested nodes and return a comprehension over them
		effectiveRego = append(effectiveRego, fmt.Sprintf("tmp_%s = nested_nodes with data.nodes as %s[\"%s\"]", variable, previousBinding, nextPath))
		effectiveRego = append(effectiveRego, fmt.Sprintf("%s = tmp_%s[_][_]", variable, variable))
		effectiveRego = append(effectiveRego, fmt.Sprintf("nodes = %s", variable))

		tr.rego = effectiveRego
		acc = append(acc, tr)
	}

	regoResult := accumulatePaths(acc)
	regoResult.source = path.Source()
	regoResult.variable = variable

	return regoResult
}

// Since there are many alternative paths to reach the nodes, we need to provide a single
// collection of nodes of the rest of the checks.
func accumulatePaths(paths []regoPathResultInternal) RegoPathResult {
	// Let's generate a rule that will return the flat list of nodes in the path
	// If there are more than one path (because of ORs) a rule with multiple clauses
	// will be generated and the final list of nodes will be the UNION of all the clauses
	rego := make([]string, 0)
	ruleName := profile.Genvar("path_rule")
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

// Different traversals based on the type of path element
func traverse(propPath path.PropertyPath, traversed traversal) []regoPathResultInternal {
	switch p := propPath.(type) {
	case path.Property:
		return traverseProperty(p, traversed)
	case path.AndPath:
		return traverseAnd(p, traversed)
	case path.OrPath:
		return traverseOr(p, traversed)
	default:
		return make([]regoPathResultInternal, 0)
	}
}

// Traverses in parallel each of the elements in the OR path, creating new branches for each one
func traverseOr(or path.OrPath, t traversal) []regoPathResultInternal {
	acc := make([]regoPathResultInternal, 0)
	for _, p := range or.Or {
		traversed := traverse(p, t)
		for _, tr := range traversed {
			acc = append(acc, tr)
		}
	}
	return acc
}

// Traverses an AND branch in the property path.
// it sequentially traverses the path until all the elements of the AND
// has been traversed and the variables accumulated.
func traverseAnd(and path.AndPath, t traversal) []regoPathResultInternal {
	first := and.And[0]
	firstTraversed := traverse(first, t)
	if len(and.And) > 1 {
		remaining := and.And[1:len(and.And)]
		acc := make([]regoPathResultInternal, 0)
		for _, tr := range firstTraversed {
			next := path.AndPath{
				And: remaining,
			}
			for _, ntr := range traverse(next, internalResultToTraversal(t, tr)) {
				acc = append(acc, ntr)
			}
		}
		return acc
	} else {
		return firstTraversed
	}
}

// Traverses the leaf components of the path expression, always a property.
// TODO: We don't take into account inverse or transitive paths yet.
func traverseProperty(property path.Property, t traversal) []regoPathResultInternal {

	// We use IDX go generate a unique property for the Rego computation
	idx := fmt.Sprintf("%d", len(t.pathVariables))
	if *t.counter > 0 {
		idx = fmt.Sprintf("%s_%d", idx, t.counter)
	}
	binding := fmt.Sprintf("%s_%s_%s_%s", t.variable, idx, t.id, t.hint)

	if len(t.pathVariables) == 0 {
		// If this is the first element in the path, we start computing the path from the previous variable passed
		// to the path generator, usually a classTarget.
		t.rego = append(t.rego, fmt.Sprintf("init_%s = data.sourceNode", binding)) // this is the connection to the variable past to the generator
		t.pathVariables = append(t.pathVariables, fmt.Sprintf("init_%s", binding)) // initial value
		t.rego = append(t.rego, fmt.Sprintf("tmp_%s = nested_nodes with data.nodes as init_%s[\"%s\"]", binding, binding, property.Iri))
		t.rego = append(t.rego, fmt.Sprintf("%s = tmp_%s[_][_]", binding, binding))
	} else {
		// Internal path element, we generate the next set of values from the previous bound property
		previousBinding := t.pathVariables[len(t.pathVariables)-1]
		t.rego = append(t.rego, fmt.Sprintf("tmp_%s = nested_nodes with data.nodes as %s[\"%s\"]", binding, previousBinding, property.Iri))
		t.rego = append(t.rego, fmt.Sprintf("%s = tmp_%s[_][_]", binding, binding))
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
