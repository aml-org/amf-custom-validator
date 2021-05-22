package test1


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


default warning = []

default info = []
# Path rules

gen_path_rule_2[nodes] {
  init_x_0__minCount_ = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as init_x_0__minCount_["apiContract:supportedOperation"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  nodes_tmp = object.get(x_0__minCount_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_4[nodes] {
  init_x_0__in_ = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as init_x_0__in_["apiContract:supportedOperation"]
  x_0__in_ = tmp_x_0__in_[_][_]
  nodes_tmp = object.get(x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_5[nodes] {
  init_x_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__pattern_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.supportedOperation / apiContract.method
  gen_propValues_1 = gen_path_rule_2 with data.sourceNode as x
  not count(gen_propValues_1) >= 1
  _result_0 := trace("minCount","apiContract.supportedOperation / apiContract.method",count(gen_propValues_1),"value not matching rule 1")
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.supportedOperation / apiContract.method
  x_check_array = gen_path_rule_4 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_3 = { "publish","subscribe"}
  not gen_inValues_3[x_check]
  _result_0 := trace("in","apiContract.supportedOperation / apiContract.method",x_check,"Error with value gen_inValues_3 and enumeration ['publish','subscribe']")
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: shacl.name
  gen_path_rule_5_node_array = gen_path_rule_5 with data.sourceNode as x
  gen_path_rule_5_node = gen_path_rule_5_node_array[_]
  not regex.match("^put|post$",gen_path_rule_5_node)
  _result_0 := trace("pattern","shacl.name",gen_path_rule_5_node,"Error with value gen_path_rule_5_node and matching regular expression '^put|post$'")
  matches := error("validation1",x,"This is the message",[_result_0])
}
# Path rules

gen_path_rule_7[nodes] {
  init_x_0__minCount_ = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as init_x_0__minCount_["apiContract:expects"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  tmp_x_2__minCount_ = nested_nodes with data.nodes as x_0__minCount_["apiContract:parameter"]
  x_2__minCount_ = tmp_x_2__minCount_[_][_]
  tmp_x_3__minCount_ = nested_nodes with data.nodes as x_2__minCount_["shapes:schema"]
  x_3__minCount_ = tmp_x_3__minCount_[_][_]
  nodes_tmp = object.get(x_3__minCount_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
} {
  init_x_0__minCount_ = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as init_x_0__minCount_["apiContract:expects"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  tmp_x_2__minCount_ = nested_nodes with data.nodes as x_0__minCount_["apiContract:payload"]
  x_2__minCount_ = tmp_x_2__minCount_[_][_]
  tmp_x_3__minCount_ = nested_nodes with data.nodes as x_2__minCount_["shapes:schema"]
  x_3__minCount_ = tmp_x_3__minCount_[_][_]
  nodes_tmp = object.get(x_3__minCount_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name
  gen_propValues_6 = gen_path_rule_7 with data.sourceNode as x
  not count(gen_propValues_6) >= 1
  _result_0 := trace("minCount","apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name",count(gen_propValues_6),"value not matching rule 1")
  matches := error("validation2",x,"orPath test",[_result_0])
}