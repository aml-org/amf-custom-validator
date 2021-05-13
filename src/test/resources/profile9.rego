package test9
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
# Path rules




#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  version = object.get(x, "apiContract:version", null)
  rego_result_cd883184d5bb12e20f9175da0c491b1c = (version != null)
  
  rego_result_cd883184d5bb12e20f9175da0c491b1c == true
  _result_0 := trace("rego", "", x, "Violation in native Rego constraint")
  _result_1 := trace("or", "", x, "Failed or constraint")
  matches := error("simple-rego", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1])
}
# Path rules



#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  version = object.get(x, "apiContract:version", null)
  rego_result_69aeea3b9e2466f13cb18ad8df1f9e03 = (version != null)
  
  rego_result_69aeea3b9e2466f13cb18ad8df1f9e03 != true
  _result_0 := trace("rego", "", x, "api without version")
  matches := error("simple-rego2", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0])
}
# Path rules

gen_path_rule_3[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:version"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  gen_path_rule_3_node_array = gen_path_rule_3 with data.sourceNode as x
  gen_path_rule_3_node = gen_path_rule_3_node_array[_]
  rego_result_a51d50017b7cfc11eb9a7fd6607c63d0 = (gen_path_rule_3_node != null)
  rego_result_a51d50017b7cfc11eb9a7fd6607c63d0 != true
  _result_0 := trace("rego", "apiContract:version", gen_path_rule_3_node, "Violation in native Rego constraint")
  matches := error("simple-rego3", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0])
}