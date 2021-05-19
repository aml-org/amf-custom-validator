package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
)

type RegoUnit struct {
	Name string
	Entrypoint string
	Code string
}

// Main entry point generating a valid Rego unit from a parsed profile.
// It uses the encoded Rego preamble code to provide all the library
// code to execute the profile.
func Generate(profile statements.Profile) RegoUnit {
	acc := []string{ pkg(profile), preamble(profile) }
	for _,r := range ruleSet(profile) {
		acc = append(acc, GenerateTopLevelExpression(r))
	}
	return RegoUnit{
		Name: packageName(profile),
		Entrypoint: entrypoint(profile),
		Code: strings.Join(acc,"\n"),
	}
}

func ruleSet(profile statements.Profile) []statements.Rule {
	acc := make([]statements.Rule,0)
	for _,r := range profile.Violation {
		acc = append(acc, r)
	}
	for _,r := range profile.Warning {
		acc = append(acc, r)
	}
	for _,r := range profile.Info {
		acc = append(acc, r)
	}

	return acc
}

func pkg(profile statements.Profile) string {
	return fmt.Sprintf("package %s\n",packageName(profile))
}

func packageName(profile statements.Profile) string {
	return strings.ReplaceAll(strings.ToLower(profile.Name), " ", "")
}

func entrypoint(profile statements.Profile) string {
	return "report"
}


func preamble(profile statements.Profile) string {
	acc := make([]string,0)
	acc = append(acc, preambleRaw)
	// If empty, we provide a default implementation so we can always generate a value report.
	// Same for other levels
	if len(profile.Violation) == 0 {
		acc = append(acc, "default violation = []")
	}
	if len(profile.Warning) == 0 {
		acc = append(acc, "default warning = []")
	}
	if len(profile.Info) == 0 {
		acc = append(acc, "default info = []")
	}

	return strings.Join(acc,"\n\n")
}

const preambleRaw = `
# Finds a node in the graph, following a link in the flatten JSON-LD node
find = node {
  node := input["@graph"][_]
  node["@id"] = data.link["@id"]
}

# Makes sure that value is wrapped in an array even if it is a single object property
nodes_array = data.nodes {
  is_array(data.nodes)
} else = [data.nodes]

# Navigates an object property for a given node
nested[nested_node] {
  nested_links := nodes_array with data.nodes as data.nodes
  link = nested_links[_]
  nested_node := find with data.link as link
}

# Fetches all the nodes for an object link
nested_nodes[nested_nodes] {
  nested_links := nodes_array with data.nodes as data.nodes
  nested_nodes := [nested_node |
    link = nested_links[_]
    nested_node := find with data.link as link
  ]
}

# Fetches all the scalars for a scalar link
nested_values[nested_values] {
  nested_values := {value | n = data.nodes[_]; value := n[data.property]}
}


# Fetches all the nodes for a given RDF class
target_class[node] {
  node := input["@graph"][_]
  node["@type"][_] == data.class
}

# Fetches all the nodes without the given RDF class
target_class_negated[result] {
  node = input["@graph"][_]
  classes = [type | c := node["@type"][_]; c == data.class; type := c]
  count(classes) == 0
  result := node
}

# Transform scalars to string, useful for 'in' constraints
as_string(x) = x {
  is_string(x)
}

as_string(x) = json.marshal(x) {
  not is_string(x)
}


# Traces one evaluation of a constraint
trace(component, path, value, traceMessage) = t {
  t := {
    "component": component,
    "path": path,
    "value": value,
    "message": traceMessage
  }
}

# Builds an error message that can be returned to the calling client software
error(shapeId, target, message, traceLog) = e {
  id := target["@id"]
  e := {
    "shapeId": shapeId,
    "target": id,
    "message": message,
    "trace": traceLog
  }
}

# generate the report for violation level
# default value must be added dynamically
report[level] = matches {
  vs = violation
  level := "violation"
  matches := vs
}

# generate the report for the info level
# default value must be added dynamically
report[level] = matches {
  vs = info
  level := "info"
  matches := vs
}

# generate the report for the info level
# default value must be added dynamically
report[level] = matches {
  vs = warning
  level := "warning"
  matches := vs
}
`