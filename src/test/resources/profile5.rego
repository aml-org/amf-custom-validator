package test_profile_5
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
 target_class[x] with data.class as "apiContract:EndPoint"
  nested_nodes[x_0_4d6e6b7923f1e16651ec6583344220eb] with data.nodes as x["apiContract:supportedOperation"]
  ys = x_0_4d6e6b7923f1e16651ec6583344220eb
  ys_errors = [ ys_error |
    y = ys[_]
    y_0_b87947f49ae3eed9ba2e63e2c81fd029_in_f9333456ffbea053fef5a20ad9deb8c0 = y["apiContract:method"]
    gen_invalues_1 = {"post"}
    not gen_invalues_1[y_0_b87947f49ae3eed9ba2e63e2c81fd029_in_f9333456ffbea053fef5a20ad9deb8c0]
    _result_0 := trace("in", "apiContract:method", y_0_b87947f49ae3eed9ba2e63e2c81fd029_in_f9333456ffbea053fef5a20ad9deb8c0, "Value no in set {'post'}")
    ys_error := error("null", y, "null", [_result_0])
  ]
  not(count(ys) - count(ys_errors) >= 1)
  _result_0 := trace("nested", "apiContract:supportedOperation", {"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))}, [e | e := ys_errors[_].trace])
  matches := error("validation1", x, "Endpoints must have a POST method", [_result_0])
}