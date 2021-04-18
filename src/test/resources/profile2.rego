package test_2
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
trace(component, value, traceMessage) = t {
  t := {
    "component": component,
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
  gen_invalues_1 = {"subscribe"}
  not gen_invalues_1[x_0_b87947f49ae3eed9ba2e63e2c81fd029]
  _result_0 := trace("in", x_0_b87947f49ae3eed9ba2e63e2c81fd029, "Value no in set {'subscribe'}")
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = x["apiContract:method"]
  gen_invalues_2 = {"get"}
  not gen_invalues_2[x_0_b87947f49ae3eed9ba2e63e2c81fd029]
  _result_1 := trace("in", x_0_b87947f49ae3eed9ba2e63e2c81fd029, "Value no in set {'get'}")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = object.get(x,"apiContract:method",[])
  gen_propValues_3 = nodes_array with data.nodes as x_0_b87947f49ae3eed9ba2e63e2c81fd029
  not count(gen_propValues_3) >= 1
  _result_0 := trace("minCount", count(gen_propValues_3), "Value not matching minCount 1")
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = x["apiContract:method"]
  gen_invalues_4 = {"get"}
  not gen_invalues_4[x_0_b87947f49ae3eed9ba2e63e2c81fd029]
  _result_1 := trace("in", x_0_b87947f49ae3eed9ba2e63e2c81fd029, "Value no in set {'get'}")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = x["apiContract:method"]
  gen_invalues_5 = {"subscribe"}
  not gen_invalues_5[x_0_b87947f49ae3eed9ba2e63e2c81fd029]
  _result_0 := trace("in", x_0_b87947f49ae3eed9ba2e63e2c81fd029, "Value no in set {'subscribe'}")
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = object.get(x,"apiContract:method",[])
  gen_propValues_6 = nodes_array with data.nodes as x_0_b87947f49ae3eed9ba2e63e2c81fd029
  not count(gen_propValues_6) >= 1
  _result_1 := trace("minCount", count(gen_propValues_6), "Value not matching minCount 1")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = object.get(x,"apiContract:method",[])
  gen_propValues_7 = nodes_array with data.nodes as x_0_b87947f49ae3eed9ba2e63e2c81fd029
  not count(gen_propValues_7) >= 1
  _result_0 := trace("minCount", count(gen_propValues_7), "Value not matching minCount 1")
  x_0_b87947f49ae3eed9ba2e63e2c81fd029 = object.get(x,"apiContract:method",[])
  gen_propValues_8 = nodes_array with data.nodes as x_0_b87947f49ae3eed9ba2e63e2c81fd029
  not count(gen_propValues_8) >= 1
  _result_1 := trace("minCount", count(gen_propValues_8), "Value not matching minCount 1")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1])
}