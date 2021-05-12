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
# Path rules

gen_path_rule_1[nodes] {
  x = data.sourceNode
  tmp_x_0_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2 = nested_nodes with data.nodes as x["apiContract:supportedOperation"]
  x_0_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2 = tmp_x_0_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2[_][_]
  nodes_tmp = x_0_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2["apiContract:method"]
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
  nodes_tmp = x["shacl:name"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"publish","subscribe"}
  not gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:supportedOperation / apiContract:method", x_check, "Value no in set {'publish','subscribe'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  gen_propValues_4 = gen_path_rule_3 with data.sourceNode as x
  not count(gen_propValues_4) >= 1
  _result_0 := trace("minCount", "apiContract:supportedOperation / apiContract:method", count(gen_propValues_4), "Value not matching minCount 1")
  matches := error("validation1", x, "This is the message", [_result_0])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  gen_path_rule_5_node_array = gen_path_rule_5 with data.sourceNode as x
  gen_path_rule_5_node = gen_path_rule_5_node_array[_]
  not regex.match("^put|post$",gen_path_rule_5_node)
  _result_0 := trace("pattern", "shacl:name", gen_path_rule_5_node, "Value does not match regular expression {'^put|post$'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}
# Path rules

gen_path_rule_6[nodes] {
  x = data.sourceNode
  tmp_x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x["apiContract:expects"]
  x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  tmp_x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["apiContract:parameter"]
  x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  tmp_x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["shapes:schema"]
  x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  nodes_tmp = object.get(x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
} {
  x = data.sourceNode
  tmp_x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x["apiContract:expects"]
  x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  tmp_x_1_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x_0_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["apiContract:payload"]
  x_1_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_1_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  tmp_x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = nested_nodes with data.nodes as x_1_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["shapes:schema"]
  x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = tmp_x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b[_][_]
  nodes_tmp = object.get(x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  gen_propValues_7 = gen_path_rule_6 with data.sourceNode as x
  not count(gen_propValues_7) >= 1
  _result_0 := trace("minCount", "apiContract:expects / (apiContract:parameter / shapes:schema) | (apiContract:payload / shapes:schema) / shacl:name", count(gen_propValues_7), "Value not matching minCount 1")
  matches := error("validation2", x, "orPath test", [_result_0])
}