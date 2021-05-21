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

gen_path_rule_33[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_34[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_35[nodes] {
  y = data.sourceNode
  nodes_tmp = object.get(y,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_33[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_36[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_37[nodes] {
  z = data.sourceNode
  nodes_tmp = object.get(z,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_33[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_38[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_39[nodes] {
  p = data.sourceNode
  nodes_tmp = object.get(p,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_33[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_40[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_41[nodes] {
  q = data.sourceNode
  nodes_tmp = object.get(q,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_33 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_32 = { "get"}
  not gen_inValues_32[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_32 and enumeration ['get']")
  #  querying path: apiContract.returns
  ys = gen_path_rule_34 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_35_node_array = gen_path_rule_35 with data.sourceNode as y
    gen_path_rule_35_node = gen_path_rule_35_node_array[_]
    not regex.match("^201$",gen_path_rule_35_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_35_node,"Error with value gen_path_rule_35_node and matching regular expression '^201$'")
    ys_error := error("nested",y,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(ys) - count(ys_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_33 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_32 = { "get"}
  not gen_inValues_32[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_32 and enumeration ['get']")
  #  querying path: apiContract.returns
  zs = gen_path_rule_36 with data.sourceNode as x
  zs_errors = [ zs_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_37_node_array = gen_path_rule_37 with data.sourceNode as z
    gen_path_rule_37_node = gen_path_rule_37_node_array[_]
    not regex.match("^2[0-9]{2}$",gen_path_rule_37_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_37_node,"Error with value gen_path_rule_37_node and matching regular expression '^2[0-9]{2}$'")
    zs_error := error("nested",z,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(zs) - count(zs_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(zs_errors), "success":(count(zs) - count(zs_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_33 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_32 = { "get"}
  not gen_inValues_32[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_32 and enumeration ['get']")
  #  querying path: apiContract.returns
  ps = gen_path_rule_38 with data.sourceNode as x
  ps_errors = [ ps_error|
    p = ps[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_39_node_array = gen_path_rule_39 with data.sourceNode as p
    gen_path_rule_39_node = gen_path_rule_39_node_array[_]
    not regex.match("^4[0-9]{2}$",gen_path_rule_39_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_39_node,"Error with value gen_path_rule_39_node and matching regular expression '^4[0-9]{2}$'")
    ps_error := error("nested",p,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(ps) - count(ps_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(ps_errors), "success":(count(ps) - count(ps_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  x_check_array = gen_path_rule_33 with data.sourceNode as x
  x_check_scalar = x_check_array[_]
  x_check = as_string(x_check_scalar)
  gen_inValues_32 = { "get"}
  not gen_inValues_32[x_check]
  _result_0 := trace("in","apiContract.method",x_check,"Error with value gen_inValues_32 and enumeration ['get']")
  #  querying path: apiContract.returns
  qs = gen_path_rule_40 with data.sourceNode as x
  qs_errors = [ qs_error|
    q = qs[_]
    #  querying path: apiContract.statusCode
    gen_path_rule_41_node_array = gen_path_rule_41 with data.sourceNode as q
    gen_path_rule_41_node = gen_path_rule_41_node_array[_]
    not regex.match("^5[0-9]{2}$",gen_path_rule_41_node)
    _result_0 := trace("pattern","apiContract.statusCode",gen_path_rule_41_node,"Error with value gen_path_rule_41_node and matching regular expression '^5[0-9]{2}$'")
    qs_error := error("nested",q,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(qs) - count(qs_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",{"failed": count(qs_errors), "success":(count(qs) - count(qs_errors))},"")
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}