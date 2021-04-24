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
violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b = x["apiContract:method"]
  gen_invalues_1 = {"get"}
  gen_invalues_1[x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b]
  _result_0 := trace("in", "apiContract:method", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Failed or constraint")
  nested_nodes[x_1_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199] with data.nodes as x["apiContract:returns"]
  ys = x_1_375ddfc577b719d3d870b481f14cd592_nested_2e7759a5cf0749f5fca54d528134e199
  ys_errors = [ ys_error |
    y = ys[_]
    y_1_f9751c730e01a629e82e0605d4de4e0b_pattern_8adff84d951bec40dcf58873c2c955c4 = y["apiContract:statusCode"]
    not regex.match("^201$",y_1_f9751c730e01a629e82e0605d4de4e0b_pattern_8adff84d951bec40dcf58873c2c955c4)
    _result_0 := trace("pattern", "apiContract:statusCode", y_1_f9751c730e01a629e82e0605d4de4e0b_pattern_8adff84d951bec40dcf58873c2c955c4, "Value does not match regular expression {'^201$'}")
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
  x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b = x["apiContract:method"]
  gen_invalues_1 = {"get"}
  gen_invalues_1[x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b]
  _result_0 := trace("in", "apiContract:method", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Failed or constraint")
  nested_nodes[x_1_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9] with data.nodes as x["apiContract:returns"]
  zs = x_1_375ddfc577b719d3d870b481f14cd592_nested_6c2fe142e2567f5c7c9591d899452ac9
  zs_errors = [ zs_error |
    z = zs[_]
    z_1_f9751c730e01a629e82e0605d4de4e0b_pattern_c438fa6dc9c84aba5bfbb8c481550be0 = z["apiContract:statusCode"]
    not regex.match("^2[0-9]{2}$",z_1_f9751c730e01a629e82e0605d4de4e0b_pattern_c438fa6dc9c84aba5bfbb8c481550be0)
    _result_0 := trace("pattern", "apiContract:statusCode", z_1_f9751c730e01a629e82e0605d4de4e0b_pattern_c438fa6dc9c84aba5bfbb8c481550be0, "Value does not match regular expression {'^2[0-9]{2}$'}")
    zs_error := error("null", z, "null", [_result_0])
  ]
  not(count(zs) - count(zs_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(zs_errors), "success":(count(zs) - count(zs_errors))}, [e | e := zs_errors[_].trace])
  _result_3 := trace("or", "", zs_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}

violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b = x["apiContract:method"]
  gen_invalues_1 = {"get"}
  gen_invalues_1[x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b]
  _result_0 := trace("in", "apiContract:method", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Failed or constraint")
  nested_nodes[x_1_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73] with data.nodes as x["apiContract:returns"]
  ps = x_1_375ddfc577b719d3d870b481f14cd592_nested_20c3bf1571cdf1eaec6643f102d82a73
  ps_errors = [ ps_error |
    p = ps[_]
    p_1_f9751c730e01a629e82e0605d4de4e0b_pattern_00df27e427732b6d83df12f40c7678b9 = p["apiContract:statusCode"]
    not regex.match("^4[0-9]{2}$",p_1_f9751c730e01a629e82e0605d4de4e0b_pattern_00df27e427732b6d83df12f40c7678b9)
    _result_0 := trace("pattern", "apiContract:statusCode", p_1_f9751c730e01a629e82e0605d4de4e0b_pattern_00df27e427732b6d83df12f40c7678b9, "Value does not match regular expression {'^4[0-9]{2}$'}")
    ps_error := error("null", p, "null", [_result_0])
  ]
  not(count(ps) - count(ps_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(ps_errors), "success":(count(ps) - count(ps_errors))}, [e | e := ps_errors[_].trace])
  _result_3 := trace("or", "", ps_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}

violation[matches] {
 target_class[x] with data.class as "apiContract:Operation"
  x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b = x["apiContract:method"]
  gen_invalues_1 = {"get"}
  gen_invalues_1[x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b]
  _result_0 := trace("in", "apiContract:method", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Value no in set {'get'}")
  _result_1 := trace("or", "", x_1_b87947f49ae3eed9ba2e63e2c81fd029_in_e1cc170ad25ef04281c98e035046e65b, "Failed or constraint")
  nested_nodes[x_1_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2] with data.nodes as x["apiContract:returns"]
  qs = x_1_375ddfc577b719d3d870b481f14cd592_nested_a19257deeb2838ab36c23ee012ed61f2
  qs_errors = [ qs_error |
    q = qs[_]
    q_1_f9751c730e01a629e82e0605d4de4e0b_pattern_1d7bd8e195d8ed802a49fb7a9b72cfd9 = q["apiContract:statusCode"]
    not regex.match("^5[0-9]{2}$",q_1_f9751c730e01a629e82e0605d4de4e0b_pattern_1d7bd8e195d8ed802a49fb7a9b72cfd9)
    _result_0 := trace("pattern", "apiContract:statusCode", q_1_f9751c730e01a629e82e0605d4de4e0b_pattern_1d7bd8e195d8ed802a49fb7a9b72cfd9, "Value does not match regular expression {'^5[0-9]{2}$'}")
    qs_error := error("null", q, "null", [_result_0])
  ]
  not(count(qs) - count(qs_errors) >= 1)
  _result_2 := trace("nested", "apiContract:returns", {"failed": count(qs_errors), "success":(count(qs) - count(qs_errors))}, [e | e := qs_errors[_].trace])
  _result_3 := trace("or", "", qs_errors, "Failed or constraint")
  matches := error("and-or-not-rule", x, "GET operations must have 2xx, 4xx and 5xx status codes but no 201", [_result_0,_result_1,_result_2,_result_3])
}