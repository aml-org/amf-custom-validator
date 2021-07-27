package test13


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

# collection functions

# collect next set of nodes
collect[r] {
  nodes = data.nodes
  n = nodes[_]
  rs = object.get(n,data.property,[])
  rss = nodes_array with data.nodes as rs
  rsss = nested with data.nodes as rss
  r = rsss[_]
}

# collect terminal values
collect_values[r] {
  nodes = data.nodes
  n = nodes[_]
  rs = object.get(n,data.property,[])
  rss = nodes_array with data.nodes as rs
  r = rss[_]
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

as_string(x) = x["@id"] {
  is_object(x)
  x["@id"]
}

as_string(x) = json.marshal(x) {
  is_object(x)
  not x["@id"]
}

# Traces one evaluation of a constraint
trace(constraint, path, node, value) = t {
  id := node["@id"]
  raw_lexical := input["@lexical"][id]
  lexical_parts := regex.find_n("\\d+", raw_lexical, 4)
  t := {
    "component": constraint,
    "path": path,
    "focusNode": id,
    "value": value,
	"lexical": {
      "start": {
        "line": lexical_parts[0],
        "column": lexical_parts[1]
      },
      "end": {
        "line": lexical_parts[2],
        "column": lexical_parts[3]
      }
    }
  }
}

trace(constraint, path, node, value) = t {
  id := node["@id"]
  not input["@lexical"][id]
  t := {
    "component": constraint,
    "path": path,
    "focusNode": id,
    "value": value
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
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_2[nodes] {
  init_p_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_p_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_15[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_4[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_5[nodes] {
  init_q_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_q_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_15[nodes] {
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
  init_y_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_15[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_10[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_11[nodes] {
  init_z_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__pattern_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_15[nodes] {
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  ps = gen_path_rule_1 with data.sourceNode as x
  ps_errors = [ ps_error|
    p = ps[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_2_node_3_array = gen_path_rule_2 with data.sourceNode as p
    gen_gen_path_rule_2_node_3 = gen_gen_path_rule_2_node_3_array[_]
    not regex.match("^4[0-9]{2}$",gen_gen_path_rule_2_node_3)
    _result_0 := trace("pattern","apiContract.statusCode",p,{"negated":false,"argument": gen_gen_path_rule_2_node_3})
    ps_error := error("nested",p,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(ps) - count(ps_errors) >= 1
  _result_0 := trace("nested","apiContract.returns",x,{"negated":false, "expected":0, "actual":count(ps_errors)})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","apiContract.method",x,{"negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  qs = gen_path_rule_4 with data.sourceNode as x
  qs_errors = [ qs_error|
    q = qs[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_5_node_6_array = gen_path_rule_5 with data.sourceNode as q
    gen_gen_path_rule_5_node_6 = gen_gen_path_rule_5_node_6_array[_]
    not regex.match("^5[0-9]{2}$",gen_gen_path_rule_5_node_6)
    _result_0 := trace("pattern","apiContract.statusCode",q,{"negated":false,"argument": gen_gen_path_rule_5_node_6})
    qs_error := error("nested",q,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(qs) - count(qs_errors) >= 1
  _result_0 := trace("nested","apiContract.returns",x,{"negated":false, "expected":0, "actual":count(qs_errors)})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","apiContract.method",x,{"negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  ys = gen_path_rule_7 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_8_node_9_array = gen_path_rule_8 with data.sourceNode as y
    gen_gen_path_rule_8_node_9 = gen_gen_path_rule_8_node_9_array[_]
    not regex.match("^201$",gen_gen_path_rule_8_node_9)
    _result_0 := trace("pattern","apiContract.statusCode",y,{"negated":false,"argument": gen_gen_path_rule_8_node_9})
    ys_error := error("nested",y,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  count(ys) - count(ys_errors) >= 1
  _result_0 := trace("nested","apiContract.returns",x,{"negated":true, "expected":0, "actual":count(ys_errors)})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","apiContract.method",x,{"negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  zs = gen_path_rule_10 with data.sourceNode as x
  zs_errors = [ zs_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_11_node_12_array = gen_path_rule_11 with data.sourceNode as z
    gen_gen_path_rule_11_node_12 = gen_gen_path_rule_11_node_12_array[_]
    not regex.match("^2[0-9]{2}$",gen_gen_path_rule_11_node_12)
    _result_0 := trace("pattern","apiContract.statusCode",z,{"negated":false,"argument": gen_gen_path_rule_11_node_12})
    zs_error := error("nested",z,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(zs) - count(zs_errors) >= 1
  _result_0 := trace("nested","apiContract.returns",x,{"negated":false, "expected":0, "actual":count(zs_errors)})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","apiContract.method",x,{"negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}