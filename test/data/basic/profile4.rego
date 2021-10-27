package profile_test_4


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

gen_path_rule_2[nodes] {
  init_x_0__maxCount_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__maxCount_,"raml-shapes:schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_4[nodes] {
  init_x_0__minCount_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__minCount_,"raml-shapes:schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_5[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x = nested_nodes with data.nodes as init_x_0__nested_["raml-shapes:schema"]
  x = tmp_x[_][_]
  nodes = x
}

gen_path_rule_7[nodes] {
  init_y_0__minCount_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__minCount_,"shacl:minLength",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_1 = gen_path_rule_2 with data.sourceNode as x
  not count(gen_propValues_1) <= 3
  _result_0 := trace("maxCount","raml-shapes.schema",x,{"negated":false,"condition":"<=","actual": count(gen_propValues_1),"expected": 3})
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_3 = gen_path_rule_4 with data.sourceNode as x
  not count(gen_propValues_3) >= 1
  _result_0 := trace("minCount","raml-shapes.schema",x,{"negated":false,"condition":">=","actual": count(gen_propValues_3),"expected": 1})
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Parameter"
  #  querying path: raml-shapes.schema
  ys = gen_path_rule_5 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.minLength
    gen_propValues_6 = gen_path_rule_7 with data.sourceNode as y
    not count(gen_propValues_6) >= 1
    _result_0 := trace("minCount","shacl.minLength",y,{"negated":false,"condition":">=","actual": count(gen_propValues_6),"expected": 1})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under raml-shapes.schema",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys_error_node_variables_agg) > 0
  _result_0 := trace("nested","raml-shapes.schema",x,{"negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)),"subResult": y_errorAcc})
  matches := error("validation1",x,"Scalars in parameters must have minLength defined",[_result_0])
}