package test1


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


default warning = []

default info = []
# Path rules

gen_path_rule_43[nodes] {
  x = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as x["apiContract:supportedOperation"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  nodes_tmp = object.get(x_0__minCount_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_45[nodes] {
  x = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as x["apiContract:supportedOperation"]
  x_0__in_ = tmp_x_0__in_[_][_]
  nodes_tmp = object.get(x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_46[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.supportedOperation / apiContract.method
  gen_propValues_42 = gen_path_rule_43 with data.sourceNode as x
  not count(gen_propValues_42) >= 1
  _result_0 := trace("minCount","apiContract.supportedOperation / apiContract.method",count(gen_propValues_42),"value not matching rule 1")
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.supportedOperation / apiContract.method
  x_check_array = gen_path_rule_45 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_44 = { "publish","subscribe"}
  not gen_inValues_44[x_check]
  _result_0 := trace("in","apiContract.supportedOperation / apiContract.method",x_check,"Error with value gen_inValues_44 and enumeration ['publish','subscribe']")
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: shacl.name
  gen_path_rule_46_node_array = gen_path_rule_46 with data.sourceNode as x
  gen_path_rule_46_node = gen_path_rule_46_node_array[_]
  not regex.match("^put|post$",gen_path_rule_46_node)
  _result_0 := trace("pattern","shacl.name",gen_path_rule_46_node,"Error with value gen_path_rule_46_node and matching regular expression '^put|post$'")
  matches := error("validation1",x,"This is the message",[_result_0])
}
# Path rules

gen_path_rule_48[nodes] {
  x = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as x["apiContract:expects"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  tmp_x_1__minCount_ = nested_nodes with data.nodes as x_0__minCount_["apiContract:parameter"]
  x_1__minCount_ = tmp_x_1__minCount_[_][_]
  tmp_x_2__minCount_ = nested_nodes with data.nodes as x_1__minCount_["shapes:schema"]
  x_2__minCount_ = tmp_x_2__minCount_[_][_]
  nodes_tmp = object.get(x_2__minCount_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
} {
  x = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as x["apiContract:expects"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  tmp_x_1__minCount_ = nested_nodes with data.nodes as x_0__minCount_["apiContract:payload"]
  x_1__minCount_ = tmp_x_1__minCount_[_][_]
  tmp_x_2__minCount_ = nested_nodes with data.nodes as x_1__minCount_["shapes:schema"]
  x_2__minCount_ = tmp_x_2__minCount_[_][_]
  nodes_tmp = object.get(x_2__minCount_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name
  gen_propValues_47 = gen_path_rule_48 with data.sourceNode as x
  not count(gen_propValues_47) >= 1
  _result_0 := trace("minCount","apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name",count(gen_propValues_47),"value not matching rule 1")
  matches := error("validation2",x,"orPath test",[_result_0])
}