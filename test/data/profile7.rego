package test13


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

gen_path_rule_2[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_3[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_4[nodes] {
  init_y_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_2[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_5[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_6[nodes] {
  init_z_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_2[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_7[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_8[nodes] {
  init_p_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_p_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_2[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_9[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_10[nodes] {
  init_q_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_q_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_2 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_1 = { "get"}
  gen_inValues_1[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_1 and enumeration ['get']")
  #  querying path: apiContract.returns
  ys = gen_path_rule_3 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_4_node_array = gen_path_rule_4 with data.sourceNode as y
    gen_path_rule_4_node = gen_path_rule_4_node_array[_]
    not regex.match("^201$",gen_path_rule_4_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_4_node,"Error with value gen_path_rule_4_node and matching regular expression '^201$'")
    ys_error := error("nested",y,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  count(ys) - count(ys_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_2 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_1 = { "get"}
  gen_inValues_1[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_1 and enumeration ['get']")
  #  querying path: apiContract.returns
  zs = gen_path_rule_5 with data.sourceNode as x
  zs_errors = [ zs_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_6_node_array = gen_path_rule_6 with data.sourceNode as z
    gen_path_rule_6_node = gen_path_rule_6_node_array[_]
    not regex.match("^2[0-9]{2}$",gen_path_rule_6_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_6_node,"Error with value gen_path_rule_6_node and matching regular expression '^2[0-9]{2}$'")
    zs_error := error("nested",z,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(zs) - count(zs_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(zs_errors), "success":(count(zs) - count(zs_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_2 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_1 = { "get"}
  gen_inValues_1[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_1 and enumeration ['get']")
  #  querying path: apiContract.returns
  ps = gen_path_rule_7 with data.sourceNode as x
  ps_errors = [ ps_error|
    p = ps[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_8_node_array = gen_path_rule_8 with data.sourceNode as p
    gen_path_rule_8_node = gen_path_rule_8_node_array[_]
    not regex.match("^4[0-9]{2}$",gen_path_rule_8_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_8_node,"Error with value gen_path_rule_8_node and matching regular expression '^4[0-9]{2}$'")
    ps_error := error("nested",p,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(ps) - count(ps_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(ps_errors), "success":(count(ps) - count(ps_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_2 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_1 = { "get"}
  gen_inValues_1[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_1 and enumeration ['get']")
  #  querying path: apiContract.returns
  qs = gen_path_rule_9 with data.sourceNode as x
  qs_errors = [ qs_error|
    q = qs[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_10_node_array = gen_path_rule_10 with data.sourceNode as q
    gen_path_rule_10_node = gen_path_rule_10_node_array[_]
    not regex.match("^5[0-9]{2}$",gen_path_rule_10_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_10_node,"Error with value gen_path_rule_10_node and matching regular expression '^5[0-9]{2}$'")
    qs_error := error("nested",q,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(qs) - count(qs_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(qs_errors), "success":(count(qs) - count(qs_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}