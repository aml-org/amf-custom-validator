package profile_test13

report["profile"] = "Test13"

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

# Fetches all the subject nodes that have certain predicate and object
search_subjects[valid_subject] {
  predicate = data.predicate
  object = data.object
  
  node = input["@ids"][_]
  node_predicate_values = nodes_array with data.nodes as object.get(node,predicate,[])
  node_predicate_value = node_predicate_values[_]
  node_predicate_value["@id"] == object["@id"]
  valid_subject = node
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
  dt == "http://www.w3.org/2001/XMLSchema#string"
  is_string(x)
}

check_datatype(x,dt) = true {
  dt == "http://www.w3.org/2001/XMLSchema#integer"
  is_number(x)
}

check_datatype(x,dt) = true {
  dt == "http://www.w3.org/2001/XMLSchema#float"
  is_number(x)
}

check_datatype(x,dt) = true {
  dt == "http://www.w3.org/2001/XMLSchema#boolean"
  is_boolean(x)
}

check_datatype(x,dt) = true {
  is_object(x)
  t = object.get(x,"@type","")
  t == dt
}

check_datatype(x,dt) = false {
  not is_object(x)
  dt != "http://www.w3.org/2001/XMLSchema#string"
  dt != "http://www.w3.org/2001/XMLSchema#integer"
  dt != "http://www.w3.org/2001/XMLSchema#float"
  dt != "http://www.w3.org/2001/XMLSchema#boolean"
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

as_string(x) = "false" {
  is_boolean(x)
  x == false
}

as_string(x) = "true" {
  is_boolean(x)
  x == true
}

as_string(x) = format_int(x, 10) {
  is_number(x)
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
trace(constraint, resultPath, focusNode, traceValue) = t {
  l := location(focusNode)  
  t := {
	"@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue,
	"location": l
  }
}

# Generates trace when lexical info is not available
trace(constraint, resultPath, focusNode, traceValue) = t {
  not location(focusNode)
  t := {
	"@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue
  }
}

location(focusNode) = l {
  id := focusNode["@id"]
  location := input["@lexical"][id]
  raw_range := location["range"]
  uri := location["uri"]	
  range_parts := regex.find_n("\\d+", raw_range, 4)
  range := {
	"@type": ["lexicalSchema:RangeNode", "lexical:Range"],
    "start": {
	  "@type": ["lexicalSchema:PositionNode", "lexical:Position"],
  	  "line": to_number(range_parts[0]),
  	  "column": to_number(range_parts[1])
    },
    "end": {
	  "@type": ["lexicalSchema:PositionNode", "lexical:Position"],
  	  "line": to_number(range_parts[2]),
  	  "column": to_number(range_parts[3])
    }
  }
  l := {
  	"@type": ["lexicalSchema:LocationNode", "lexical:Location"],
  	"uri": uri,
  	"range": range
  }
}

# Builds an error message that can be returned to the calling client software
error(sourceShapeName, focusNode, resultMessage, traceLog) = e {
  id := focusNode["@id"]
  locationNode := location(focusNode)
  e := {
	"@type": ["reportSchema:ValidationResultNode", "shacl:ValidationResult"],
    "sourceShapeName": sourceShapeName,
    "focusNode": id, # can potentially be wrapped in @id obj if report dialect is adjusted
    "resultMessage": resultMessage,
	"location": locationNode,	
    "trace": traceLog
  }
}

# Builds error message when lexical info is not available
error(sourceShapeName, focusNode, resultMessage, traceLog) = e {
  id := focusNode["@id"]
  not location(focusNode)
  e := {
	"@type": ["reportSchema:ValidationResultNode", "shacl:ValidationResult"],
    "sourceShapeName": sourceShapeName,
    "focusNode": id, # can potentially be wrapped in @id obj if report dialect is adjusted
    "resultMessage": resultMessage,
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
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_rule_2[nodes] {
  init_p_0 = data.sourceNode
  nodes_tmp = object.get(init_p_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  p_0 = nodes_tmp2[_]
  nodes = p_0
}

gen_path_rule_15[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_4[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_rule_5[nodes] {
  init_q_0 = data.sourceNode
  nodes_tmp = object.get(init_q_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  q_0 = nodes_tmp2[_]
  nodes = q_0
}

gen_path_rule_15[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_7[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_rule_8[nodes] {
  init_y_0 = data.sourceNode
  nodes_tmp = object.get(init_y_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0 = nodes_tmp2[_]
  nodes = y_0
}

gen_path_rule_15[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_10[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_rule_11[nodes] {
  init_z_0 = data.sourceNode
  nodes_tmp = object.get(init_z_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0 = nodes_tmp2[_]
  nodes = z_0
}

gen_path_rule_15[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.returns
  ps = gen_path_rule_1 with data.sourceNode as x
  p_errorAcc0 = []
  ps_br_0 = [ ps_br_0_error|
    p = ps[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_2_node_3_array = gen_path_rule_2 with data.sourceNode as p
    gen_gen_path_rule_2_node_3 = gen_gen_path_rule_2_node_3_array[_]
    not regex.match(`^4[0-9]{2}$`,gen_gen_path_rule_2_node_3)
    _result_0 := trace("pattern","http://a.ml/vocabularies/apiContract#statusCode",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"argument": gen_gen_path_rule_2_node_3})
    ps_br_0_inner_error := error("nested",p,"error in nested nodes under http://a.ml/vocabularies/apiContract#returns",[_result_0])
    ps_br_0_error = [p["@id"],ps_br_0_inner_error]
  ]
  ps_br_0_errors = { nodeId | n = ps_br_0[_]; nodeId = n[0] }
  ps_br_0_errors_errors = [ node | n = ps_br_0[_]; node = n[1] ]
  p_errorAcc1 = array.concat(p_errorAcc0,ps_br_0_errors_errors)
  p_errorAcc = p_errorAcc1
  # let's accumulate results
  ps_error_node_variables_agg = ps_br_0_errors
  not count(ps) - count(ps_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ps_error_node_variables_agg), "successfulNodes":(count(ps)-count(ps_error_node_variables_agg)), "cardinality":1, "subResult": p_errorAcc})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.returns
  qs = gen_path_rule_4 with data.sourceNode as x
  q_errorAcc0 = []
  qs_br_0 = [ qs_br_0_error|
    q = qs[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_5_node_6_array = gen_path_rule_5 with data.sourceNode as q
    gen_gen_path_rule_5_node_6 = gen_gen_path_rule_5_node_6_array[_]
    not regex.match(`^5[0-9]{2}$`,gen_gen_path_rule_5_node_6)
    _result_0 := trace("pattern","http://a.ml/vocabularies/apiContract#statusCode",q,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"argument": gen_gen_path_rule_5_node_6})
    qs_br_0_inner_error := error("nested",q,"error in nested nodes under http://a.ml/vocabularies/apiContract#returns",[_result_0])
    qs_br_0_error = [q["@id"],qs_br_0_inner_error]
  ]
  qs_br_0_errors = { nodeId | n = qs_br_0[_]; nodeId = n[0] }
  qs_br_0_errors_errors = [ node | n = qs_br_0[_]; node = n[1] ]
  q_errorAcc1 = array.concat(q_errorAcc0,qs_br_0_errors_errors)
  q_errorAcc = q_errorAcc1
  # let's accumulate results
  qs_error_node_variables_agg = qs_br_0_errors
  not count(qs) - count(qs_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(qs_error_node_variables_agg), "successfulNodes":(count(qs)-count(qs_error_node_variables_agg)), "cardinality":1, "subResult": q_errorAcc})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.returns
  ys = gen_path_rule_7 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_8_node_9_array = gen_path_rule_8 with data.sourceNode as y
    gen_gen_path_rule_8_node_9 = gen_gen_path_rule_8_node_9_array[_]
    not regex.match(`^201$`,gen_gen_path_rule_8_node_9)
    _result_0 := trace("pattern","http://a.ml/vocabularies/apiContract#statusCode",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"argument": gen_gen_path_rule_8_node_9})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under http://a.ml/vocabularies/apiContract#returns",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.returns
  zs = gen_path_rule_10 with data.sourceNode as x
  z_errorAcc0 = []
  zs_br_0 = [ zs_br_0_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_gen_path_rule_11_node_12_array = gen_path_rule_11 with data.sourceNode as z
    gen_gen_path_rule_11_node_12 = gen_gen_path_rule_11_node_12_array[_]
    not regex.match(`^2[0-9]{2}$`,gen_gen_path_rule_11_node_12)
    _result_0 := trace("pattern","http://a.ml/vocabularies/apiContract#statusCode",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"argument": gen_gen_path_rule_11_node_12})
    zs_br_0_inner_error := error("nested",z,"error in nested nodes under http://a.ml/vocabularies/apiContract#returns",[_result_0])
    zs_br_0_error = [z["@id"],zs_br_0_inner_error]
  ]
  zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
  zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
  z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
  z_errorAcc = z_errorAcc1
  # let's accumulate results
  zs_error_node_variables_agg = zs_br_0_errors
  not count(zs) - count(zs_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
  #  querying path: apiContract.method
  gen_x_check_14_array = gen_path_rule_15 with data.sourceNode as x
  gen_x_check_14_scalar = gen_x_check_14_array[_]
  gen_x_check_14 = as_string(gen_x_check_14_scalar)
  gen_inValues_13 = { "get"}
  gen_inValues_13[gen_x_check_14]
  _result_1 := trace("in","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true,"actual": gen_x_check_14,"expected": "[\"get\"]"})
  matches := error("and-or-not-rule",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0,_result_1])
}