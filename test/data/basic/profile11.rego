package test11


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

gen_path_rule_1[nodes] {
  init_x_0__ = data.sourceNode
  nodes_tmp = object.get(init_x_0__,"apiContract:scheme",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  #  querying path: apiContract.scheme
  gen_datatype_check_2_elem = gen_path_rule_1 with data.sourceNode as x
  gen_datatype_check_2 = gen_datatype_check_2_elem[_]
  not check_datatype(gen_datatype_check_2,"xsd:string")
  _result_0 := trace("datatype","apiContract.scheme",gen_datatype_check_2,"Error with datatype gen_datatype_check_2")
  matches := error("allowed-protocols",x,"Validation error",[_result_0])
}