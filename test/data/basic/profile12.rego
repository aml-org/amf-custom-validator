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

as_string(x) = json.marshal(x) {
  not is_string(x)
}


# Traces one evaluation of a constraint
trace(constraint, path, node, value) = t {
  t := {
    "component": constraint,
    "path": path,
    "focusNode": node["@id"],
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

gen_path_rule_3[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_4[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_6[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  ys = gen_path_rule_1 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    y_check_array = gen_path_rule_3 with data.sourceNode as y
    y_check_scalar = y_check_array[_]
    y_check = as_string(y_check_scalar)
    gen_inValues_2 = { "200"}
    not gen_inValues_2[y_check]
    _result_0 := trace("in","apiContract.statusCode",y,{"negated":false,"actual": gen_inValues_2,"expected": "y_check"})
    ys_error := error("nested",y,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(ys) - count(ys_errors) >= 1
  _result_0 := trace("nested","apiContract.returns",x,{"negated":false, "expected":0, "actual":count(ys_errors)})
  #  querying path: apiContract.returns
  zs = gen_path_rule_4 with data.sourceNode as x
  zs_errors = [ zs_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    z_check_array = gen_path_rule_6 with data.sourceNode as z
    z_check_scalar = z_check_array[_]
    z_check = as_string(z_check_scalar)
    gen_inValues_5 = { "429"}
    not gen_inValues_5[z_check]
    _result_0 := trace("in","apiContract.statusCode",z,{"negated":false,"actual": gen_inValues_5,"expected": "z_check"})
    zs_error := error("nested",z,"error in nested nodes under apiContract.returns",[_result_0])
  ]
  not count(zs) - count(zs_errors) >= 1
  _result_1 := trace("nested","apiContract.returns",x,{"negated":false, "expected":0, "actual":count(zs_errors)})
  matches := error("lack-of-resources-and-rate-limiting-too-many-requests",x,"Notify the client when the limit is exceeded by providing the limit number and the time at which the limit will\nbe reset.\n",[_result_0,_result_1])
}