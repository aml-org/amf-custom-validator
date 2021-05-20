package test4


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

gen_path_rule_9[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"raml-shapes:schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_11[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"raml-shapes:schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_12[nodes] {
  x = data.sourceNode
  tmp_x = nested_nodes with data.nodes as x["raml-shapes:schema"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_14[nodes] {
  y = data.sourceNode
  nodes_tmp = object.get(y,"shacl:minLength",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_8 = gen_path_rule_9 with data.sourceNode as x
  not count(gen_propValues_8) >= 1
  _result_0 := trace("minCount","raml-shapes.schema",count(gen_propValues_8),"value not matching rule 1")
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_10 = gen_path_rule_11 with data.sourceNode as x
  not count(gen_propValues_10) <= 3
  _result_0 := trace("maxCount","raml-shapes.schema",count(gen_propValues_10),"value not matching rule 3")
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  ys = gen_path_rule_12 with data.sourceNode as x
  ys_errors = [ ys_error|
    y = ys[_]
    #  querying path: shacl.minLength
    gen_propValues_13 = gen_path_rule_14 with data.sourceNode as y
    not count(gen_propValues_13) >= 1
    _result_0 := trace("minCount","shacl.minLength",count(gen_propValues_13),"value not matching rule 1")
    ys_error := error("nested",y,"error in nested nodes under raml-shapes.schema",[_result_0])
  ]
  not count(ys_errors) > 0
  _result_0 := trace("nested","raml-shapes.schema",{"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))},"")
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}