package test_1
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

default warning = []

default info = []
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = x["apiContract:method"]
  gen_invalues_1 = {"publish","subscribe"}
  not gen_invalues_1[x_0_b87947f49ae3eed9ba2e63e2c81fd029]
  _result_0 := trace("in", "apiContract:method", x_0_b87947f49ae3eed9ba2e63e2c81fd029, "Value no in set {'publish','subscribe'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = object.get(x,"apiContract:method",[])
  gen_propValues_2 = nodes_array with data.nodes as x_0_b87947f49ae3eed9ba2e63e2c81fd029
  not count(gen_propValues_2) >= 1
  _result_0 := trace("minCount", "apiContract:method", count(gen_propValues_2), "Value not matching minCount 1")
  matches := error("validation1", x, "This is the message", [_result_0])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_a82db48390e82e6cd3d806595c67bd32 = x["shacl:name"]
  not regex.match(x_0_a82db48390e82e6cd3d806595c67bd32, "^put|post$")
  _result_0 := trace("match", "shacl:name", x_0_a82db48390e82e6cd3d806595c67bd32, "Value does not match regular expression {'^put|post$'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}