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
 target_class[x] with data.class as "apiContract:EndPoint"
  x_1_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2 = x["apiContract:supportedOperation"]
  x_2_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2 = x_1_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2["apiContract:method"]
  gen_invalues_1 = {"publish","subscribe"}
  not gen_invalues_1[x_2_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2]
  _result_0 := trace("in", "apiContract:supportedOperation / apiContract:method", x_2_ea71737c8613215db5bdb23e4ee21161_in_1cd7c4508d18dac74835e4fe9b3f92a2, "Value no in set {'publish','subscribe'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}

violation[matches] {
 target_class[x] with data.class as "apiContract:EndPoint"
  x_1_ea71737c8613215db5bdb23e4ee21161_minCount_c4ca4238a0b923820dcc509a6f75849b = x["apiContract:supportedOperation"]
  x_2_ea71737c8613215db5bdb23e4ee21161_minCount_c4ca4238a0b923820dcc509a6f75849b = object.get(x_1_ea71737c8613215db5bdb23e4ee21161_minCount_c4ca4238a0b923820dcc509a6f75849b,"apiContract:method",[])
  gen_propValues_2 = nodes_array with data.nodes as x_2_ea71737c8613215db5bdb23e4ee21161_minCount_c4ca4238a0b923820dcc509a6f75849b
  not count(gen_propValues_2) >= 1
  _result_0 := trace("minCount", "apiContract:supportedOperation / apiContract:method", count(gen_propValues_2), "Value not matching minCount 1")
  matches := error("validation1", x, "This is the message", [_result_0])
}

violation[matches] {
 target_class[x] with data.class as "apiContract:EndPoint"
  x_1_a82db48390e82e6cd3d806595c67bd32_pattern_201c067558f8b4f7ce705935d9f8c304 = x["shacl:name"]
  not regex.match("^put|post$",x_1_a82db48390e82e6cd3d806595c67bd32_pattern_201c067558f8b4f7ce705935d9f8c304)
  _result_0 := trace("pattern", "shacl:name", x_1_a82db48390e82e6cd3d806595c67bd32_pattern_201c067558f8b4f7ce705935d9f8c304, "Value does not match regular expression {'^put|post$'}")
  matches := error("validation1", x, "This is the message", [_result_0])
}
violation[matches] {
 target_class[x] with data.class as "apiContract:EndPoint"
  x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = x["apiContract:expects"]
  x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["apiContract:parameter"]
  x_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = object.get(x_2_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,"shapes:schema",[])
  x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = x["apiContract:expects"]
  x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = x_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["apiContract:payload"]
  x_3_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = x_2_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b["shapes:schema"]
  x_4_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = object.get(x_3_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,"shacl:name",[])
  x_final_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b_0 = array.concat(x_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,[])
  x_final_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b = array.concat(x_4_1_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b,x_final_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b_0)
  gen_propValues_3 = nodes_array with data.nodes as x_final_3_1d7082cb5325d77ade9a7b4bc4637f09_minCount_c4ca4238a0b923820dcc509a6f75849b
  not count(gen_propValues_3) >= 1
  _result_0 := trace("minCount", "apiContract:expects / (apiContract:parameter / shapes:schema) | (apiContract:payload / shapes:schema) / shacl:name", count(gen_propValues_3), "Value not matching minCount 1")
  matches := error("validation2", x, "orPath test", [_result_0])
}