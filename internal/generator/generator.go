package generator

import (
	"fmt"
	"github.com/aml-org/amfopa/internal/parser/profile"
	"strings"
)

type RegoUnit struct {
	Name       string
	Entrypoint string
	Code       string
	Prefixes   profile.ProfileContext
}

// Main entry point generating a valid Rego unit from a parsed profile.
// It uses the encoded Rego preamble code to provide all the library
// code to execute the profile.
func Generate(profile profile.Profile) RegoUnit {
	acc := []string{pkg(profile), preamble(profile)}
	for _, r := range ruleSet(profile) {
		acc = append(acc, GenerateTopLevelExpression(r))
	}
	return RegoUnit{
		Name:       packageName(profile),
		Entrypoint: entrypoint(profile),
		Code:       strings.Join(acc, "\n"),
		Prefixes:   profile.Prefixes,
	}
}

func ruleSet(prof profile.Profile) []profile.Rule {
	acc := make([]profile.Rule, 0)
	for _, r := range prof.Violation {
		acc = append(acc, r)
	}
	for _, r := range prof.Warning {
		acc = append(acc, r)
	}
	for _, r := range prof.Info {
		acc = append(acc, r)
	}

	return acc
}

func pkg(profile profile.Profile) string {
	return fmt.Sprintf("package %s\n", packageName(profile))
}

func packageName(profile profile.Profile) string {
	return strings.ReplaceAll(strings.ToLower(profile.Name), " ", "")
}

func entrypoint(profile profile.Profile) string {
	return "report"
}

func preamble(profile profile.Profile) string {
	acc := make([]string, 0)
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

	return strings.Join(acc, "\n\n")
}

const preambleRaw = `
# Finds a node in the graph, following a link in the flatten JSON-LD node
find = node {
  id := data.link["@id"]
  node := input["@ids"][id]
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

# collection functions

# collect next set of nodes
collect[r] {
  nodes = data.nodes
  n = nodes[_]
  rs = object.get(n,data.property,[])
  rss = nodes_array with data.nodes as rs
  rsss = nested with data.nodes as rss
  r = rsss[_]
}

# collect terminal values
collect_values[r] {
  nodes = data.nodes
  n = nodes[_]
  rs = object.get(n,data.property,[])
  rss = nodes_array with data.nodes as rs
  r = rss[_]
}

# helper to check datatype constraints

check_datatype(x,dt) = true {
  dt == "xsd:string"
  is_string(x)
}

check_datatype(x,dt) = true {
  dt == "xsd:integer"
  is_number(x)
}

check_datatype(x,dt) = true {
  dt == "xsd:float"
  is_number(x)
}

check_datatype(x,dt) = true {
  dt == "xsd:boolean"
  is_boolean(x)
}

check_datatype(x,dt) = true {
  is_object(x)
  t = object.get(x,"@type","")
  t == dt
}

check_datatype(x,dt) = false {
  not is_object(x)
  dt != "xsd:string"
  dt != "xsd:integer"
  dt != "xsd:float"
  dt != "xsd:boolean"
}

# Fetches all the nodes for a given RDF class
target_class[node] {
  class = data.class
  id = input["@types"][class][_]
  node = input["@ids"][id] 
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
