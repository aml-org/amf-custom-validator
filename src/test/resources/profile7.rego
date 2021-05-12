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
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_4[nodes] {
  x = data.sourceNode
  tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199 = nested_nodes with data.nodes as x["apiContract:returns"]
  x_0_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199 = tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199[_][_]
  nodes = x_0_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199
}

gen_path_rule_3[nodes] {
  y = data.sourceNode
  nodes_tmp = y["apiContract:statusCode"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}


gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_6[nodes] {
  x = data.sourceNode
  tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9 = nested_nodes with data.nodes as x["apiContract:returns"]
  x_0_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9 = tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9[_][_]
  nodes = x_0_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9
}

gen_path_rule_5[nodes] {
  z = data.sourceNode
  nodes_tmp = z["apiContract:statusCode"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_8[nodes] {
  x = data.sourceNode
  tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73 = nested_nodes with data.nodes as x["apiContract:returns"]
  x_0_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73 = tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73[_][_]
  nodes = x_0_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73
}

gen_path_rule_7[nodes] {
  p = data.sourceNode
  nodes_tmp = p["apiContract:statusCode"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = x["apiContract:method"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_10[nodes] {
  x = data.sourceNode
  tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2 = nested_nodes with data.nodes as x["apiContract:returns"]
  x_0_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2 = tmp_x_0_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2[_][_]
  nodes = x_0_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2
}

gen_path_rule_9[nodes] {
  q = data.sourceNode
  nodes_tmp = q["apiContract:statusCode"]
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}


#Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"get"}
  gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_check, "Failed or constraint")
  ys = gen_path_rule_4 with data.sourceNode as x
  ys_errors = [ ys_error |
    y = ys[_]
    gen_path_rule_3_node_array = gen_path_rule_3 with data.sourceNode as y
    gen_path_rule_3_node = gen_path_rule_3_node_array[_]
    not regex.match("^201$",gen_path_rule_3_node)
    _result_0 := trace("pattern", "apiContract:statusCode", gen_path_rule_3_node, "Value does not match regular expression {'^201$'}")
    ys_error := error("null", y, "null", [_result_0])
  ]
  count(ys) - count(ys_errors) >= 1
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(ys_errors), "success":(count(ys) - count(ys_errors))}, [e | e := ys_errors[_].trace])
  _result_3 := trace("or", "", ys_errors, "Failed or constraint")
  _result_4 := trace("or", "", ys_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3,_result_4])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"get"}
  gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_check, "Failed or constraint")
  zs = gen_path_rule_6 with data.sourceNode as x
  zs_errors = [ zs_error |
    z = zs[_]
    gen_path_rule_5_node_array = gen_path_rule_5 with data.sourceNode as z
    gen_path_rule_5_node = gen_path_rule_5_node_array[_]
    not regex.match("^2[0-9]{2}$",gen_path_rule_5_node)
    _result_0 := trace("pattern", "apiContract:statusCode", gen_path_rule_5_node, "Value does not match regular expression {'^2[0-9]{2}$'}")
    zs_error := error("null", z, "null", [_result_0])
  ]
  not(count(zs) - count(zs_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(zs_errors), "success":(count(zs) - count(zs_errors))}, [e | e := zs_errors[_].trace])
  _result_3 := trace("or", "", zs_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"get"}
  gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_check, "Failed or constraint")
  ps = gen_path_rule_8 with data.sourceNode as x
  ps_errors = [ ps_error |
    p = ps[_]
    gen_path_rule_7_node_array = gen_path_rule_7 with data.sourceNode as p
    gen_path_rule_7_node = gen_path_rule_7_node_array[_]
    not regex.match("^4[0-9]{2}$",gen_path_rule_7_node)
    _result_0 := trace("pattern", "apiContract:statusCode", gen_path_rule_7_node, "Value does not match regular expression {'^4[0-9]{2}$'}")
    ps_error := error("null", p, "null", [_result_0])
  ]
  not(count(ps) - count(ps_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(ps_errors), "success":(count(ps) - count(ps_errors))}, [e | e := ps_errors[_].trace])
  _result_3 := trace("or", "", ps_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}
violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  x_check_array = gen_path_rule_1 with data.sourceNode as x
  x_check = x_check_array[_]
  gen_invalues_2 = {"get"}
  gen_invalues_2[x_check]
  _result_0 := trace("in", "apiContract:method", x_check, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_check, "Failed or constraint")
  qs = gen_path_rule_10 with data.sourceNode as x
  qs_errors = [ qs_error |
    q = qs[_]
    gen_path_rule_9_node_array = gen_path_rule_9 with data.sourceNode as q
    gen_path_rule_9_node = gen_path_rule_9_node_array[_]
    not regex.match("^5[0-9]{2}$",gen_path_rule_9_node)
    _result_0 := trace("pattern", "apiContract:statusCode", gen_path_rule_9_node, "Value does not match regular expression {'^5[0-9]{2}$'}")
    qs_error := error("null", q, "null", [_result_0])
  ]
  not(count(qs) - count(qs_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(qs_errors), "success":(count(qs) - count(qs_errors))}, [e | e := qs_errors[_].trace])
  _result_3 := trace("or", "", qs_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}