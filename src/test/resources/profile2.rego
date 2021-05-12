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

gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}
gen_path_rule_5[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_3[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}
gen_path_rule_5[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}
gen_path_rule_7[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_3[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}
gen_path_rule_7[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}


#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"subscribe"}
  not gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'subscribe'}")
  x_check_array = gen_path_rule_5 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_6 = {"get"}
  not gen_invalues_6[x_check]
  _result_1 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_2 := trace("or", "", x_check, "Failed or constraint")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1,_result_2])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  gen_propValues_4 = gen_path_rule_3 with data.sourceNode as x
  not count(gen_propValues_4) >= 1
  _result_0 := trace("minCount", "apiContract:method", count(gen_propValues_4), "Value not matching minCount 1")
  x_check_array = gen_path_rule_5 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_6 = {"get"}
  not gen_invalues_6[x_check]
  _result_1 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_2 := trace("or", "", x_check, "Failed or constraint")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1,_result_2])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"subscribe"}
  not gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'subscribe'}")
  gen_propValues_8 = gen_path_rule_7 with data.sourceNode as x
  not count(gen_propValues_8) >= 1
  _result_1 := trace("minCount", "apiContract:method", count(gen_propValues_8), "Value not matching minCount 1")
  _result_2 := trace("or", "", gen_propValues_8, "Failed or constraint")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1,_result_2])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  gen_propValues_4 = gen_path_rule_3 with data.sourceNode as x
  not count(gen_propValues_4) >= 1
  _result_0 := trace("minCount", "apiContract:method", count(gen_propValues_4), "Value not matching minCount 1")
  gen_propValues_8 = gen_path_rule_7 with data.sourceNode as x
  not count(gen_propValues_8) >= 1
  _result_1 := trace("minCount", "apiContract:method", count(gen_propValues_8), "Value not matching minCount 1")
  _result_2 := trace("or", "", gen_propValues_8, "Failed or constraint")
  matches := error("validation1", x, "This is the message", [_result_0,_result_1,_result_2])
}