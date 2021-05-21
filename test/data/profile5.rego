package testprofile5


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

gen_path_rule_15[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["apiContract:supportedOperation"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_17[nodes] {
  y = data.sourceNode
  nodes_tmp = object.get(y,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:EndPoint"
  #  querying path: apiContract.supportedOperation
  ys = gen_path_rule_15 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: apiContract.method
    y_check_array = gen_path_rule_17 with data.sourceNode as y
    y_check_scalar = y_check_array[_]
    y_check = as_string(y_check_scalar)
    gen_inValues_16 = { "post"}
    not gen_inValues_16[y_check]
    _result_0 := trace("in","apiContract.method",y_check,"Error with value gen_inValues_16 and enumeration ['post']")
    ys_error := error("nested",y,"error in nested nodes under apiContract.supportedOperation",[_result_0])
  ]
  not count(ys) - count(ys_errors) >= 1
  _result_0 := trace("nested","apiContract.supportedOperation",{"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))},"")
  matches := error("validation1",x,"Endpoints must have a POST method",[_result_0])
}