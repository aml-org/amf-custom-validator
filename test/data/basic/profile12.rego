package profile_test13


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
  id := focusNode["@id"]
  location := input["@lexical"][id]
  raw_range := location["range"]
  uri := location["uri"]	
  range_parts := regex.find_n("\\d+", raw_range, 4)
  range := {
	"@type": ["lexical:Range"],
    "start": {
	  "@type": ["lexical:Position"],
  	  "line": to_number(range_parts[0]),
  	  "column": to_number(range_parts[1])
    },
    "end": {
	  "@type": ["lexical:Position"],
  	  "line": to_number(range_parts[2]),
  	  "column": to_number(range_parts[3])
    }
  }
  t := {
	"@type": ["validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue,
	"location": {
	  "@type": ["lexical:Location"],
      "uri": uri,
      "range": range
	}
  }
}

trace(constraint, resultPath, focusNode, traceValue) = t {
  id := focusNode["@id"]
  not input["@lexical"][id]
  t := {
	"@type": ["validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue
  }
}

# Builds an error message that can be returned to the calling client software
error(sourceShapeName, focusNode, resultMessage, traceLog) = e {
  id := focusNode["@id"]
  e := {
	"@type": ["shacl:ValidationResult"],
    "sourceShapeName": sourceShapeName,
    "focusNode": {
		"@id": id,
	},
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
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_4[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_5[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_8[nodes] {
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
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_y_check_3_array = gen_path_rule_4 with data.sourceNode as y
    gen_y_check_3_scalar = gen_y_check_3_array[_]
    gen_y_check_3 = as_string(gen_y_check_3_scalar)
    gen_inValues_2 = { "200"}
    not gen_inValues_2[gen_y_check_3]
    _result_0 := trace("in","apiContract.statusCode",y,{"negated":false,"actual": gen_y_check_3,"expected": "[\"200\"]"})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under apiContract.returns",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  not count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","apiContract.returns",x,{"negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  #  querying path: apiContract.returns
  zs = gen_path_rule_5 with data.sourceNode as x
  z_errorAcc0 = []
  zs_br_0 = [ zs_br_0_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_z_check_7_array = gen_path_rule_8 with data.sourceNode as z
    gen_z_check_7_scalar = gen_z_check_7_array[_]
    gen_z_check_7 = as_string(gen_z_check_7_scalar)
    gen_inValues_6 = { "429"}
    not gen_inValues_6[gen_z_check_7]
    _result_0 := trace("in","apiContract.statusCode",z,{"negated":false,"actual": gen_z_check_7,"expected": "[\"429\"]"})
    zs_br_0_inner_error := error("nested",z,"error in nested nodes under apiContract.returns",[_result_0])
    zs_br_0_error = [z["@id"],zs_br_0_inner_error]
  ]
  zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
  zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
  z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
  z_errorAcc = z_errorAcc1
  # let's accumulate results
  zs_error_node_variables_agg = zs_br_0_errors
  not count(zs) - count(zs_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","apiContract.returns",x,{"negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
  matches := error("lack-of-resources-and-rate-limiting-too-many-requests",x,"Notify the client when the limit is exceeded by providing the limit number and the time at which the limit will\nbe reset.\n",[_result_0,_result_1])
}