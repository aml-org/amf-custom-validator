package profile_json_api


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
  t := {
	"@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue,
	"location": {
	  "@type": ["lexicalSchema:LocationNode", "lexical:Location"],
      "uri": uri,
      "range": range
	}
  }
}

trace(constraint, resultPath, focusNode, traceValue) = t {
  id := focusNode["@id"]
  not input["@lexical"][id]
  t := {
	"@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue
  }
}

# Builds an error message that can be returned to the calling client software
error(sourceShapeName, focusNode, resultMessage, traceLog) = e {
  id := focusNode["@id"]
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

gen_path_rule_91[nodes] {
  init_x_0__in_ = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as init_x_0__in_["apiContract:payload"]
  x_0__in_ = tmp_x_0__in_[_][_]
  nodes_tmp = object.get(x_0__in_,"core:mediaType",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2__in_ = nodes_tmp2[_]
  nodes = x_2__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Request"
  #  querying path: apiContract.payload / core.mediaType
  gen_x_check_90_array = gen_path_rule_91 with data.sourceNode as x
  gen_x_check_90_scalar = gen_x_check_90_array[_]
  gen_x_check_90 = as_string(gen_x_check_90_scalar)
  gen_inValues_89 = { "application/vnd.api+json"}
  not gen_inValues_89[gen_x_check_90]
  _result_0 := trace("in","apiContract.payload / core.mediaType",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_90,"expected": "[\"application/vnd.api+json\"]"})
  matches := error("json-api-media-type-request",x,"Clients MUST send all JSON:API data in request documents with the header Content-Type: application/vnd.api+json\nwithout any media type parameters.\nClients that include the JSON:API media type in their Accept header MUST specify the media type there at least once\nwithout any media type parameters.\nClients MUST ignore any parameters for the application/vnd.api+json media type received in the Content-Type header\nof response documents.\n",[_result_0])
}
# Path rules

gen_path_rule_94[nodes] {
  init_x_0__in_ = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as init_x_0__in_["apiContract:payload"]
  x_0__in_ = tmp_x_0__in_[_][_]
  nodes_tmp = object.get(x_0__in_,"core:mediaType",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2__in_ = nodes_tmp2[_]
  nodes = x_2__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Response"
  #  querying path: apiContract.payload / core.mediaType
  gen_x_check_93_array = gen_path_rule_94 with data.sourceNode as x
  gen_x_check_93_scalar = gen_x_check_93_array[_]
  gen_x_check_93 = as_string(gen_x_check_93_scalar)
  gen_inValues_92 = { "application/vnd.api+json"}
  not gen_inValues_92[gen_x_check_93]
  _result_0 := trace("in","apiContract.payload / core.mediaType",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_93,"expected": "[\"application/vnd.api+json\"]"})
  matches := error("json-api-media-type-response",x,"Servers MUST send all JSON:API data in response documents with the header Content-Type: application/vnd.api+json\nwithout any media type parameters.\n",[_result_0])
}
# Path rules

gen_path_rule_95[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_98[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0__in_ = nodes_tmp2[_]
  nodes = y_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  ys = gen_path_rule_95 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_y_check_97_array = gen_path_rule_98 with data.sourceNode as y
    gen_y_check_97_scalar = gen_y_check_97_array[_]
    gen_y_check_97 = as_string(gen_y_check_97_scalar)
    gen_inValues_96 = { "406"}
    not gen_inValues_96[gen_y_check_97]
    _result_0 := trace("in","apiContract.statusCode",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_97,"expected": "[\"406\"]"})
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
  _result_0 := trace("atLeast","apiContract.returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("406-mandatory-response",x,"Servers MUST respond with a 406 Not Acceptable status code if a request’s Accept header contains the JSON:API media\ntype and all instances of that media type are modified with media type parameters.\n",[_result_0])
}
# Path rules

gen_path_rule_99[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["apiContract:returns"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_102[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"apiContract:statusCode",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0__in_ = nodes_tmp2[_]
  nodes = y_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.returns
  ys = gen_path_rule_99 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_y_check_101_array = gen_path_rule_102 with data.sourceNode as y
    gen_y_check_101_scalar = gen_y_check_101_array[_]
    gen_y_check_101 = as_string(gen_y_check_101_scalar)
    gen_inValues_100 = { "415"}
    not gen_inValues_100[gen_y_check_101]
    _result_0 := trace("in","apiContract.statusCode",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_101,"expected": "[\"415\"]"})
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
  _result_0 := trace("atLeast","apiContract.returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("415-mandatory-response",x,"Servers MUST respond with a 415 Unsupported Media Type status code if a request specifies the header Content-Type:\napplication/vnd.api+json with any media type parameters.\n",[_result_0])
}
# Path rules

gen_path_rule_105[nodes] {
  init_x_0__in_ = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as init_x_0__in_["shapes:schema"]
  x_0__in_ = tmp_x_0__in_[_][_]
  nodes_tmp = object.get(x_0__in_,"@type",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2__in_ = nodes_tmp2[_]
  nodes = x_2__in_
}

gen_path_rule_107[nodes] {
  init_x_0__minCount_ = data.sourceNode
  tmp_x_0__minCount_ = nested_nodes with data.nodes as init_x_0__minCount_["shapes:schema"]
  x_0__minCount_ = tmp_x_0__minCount_[_][_]
  nodes_tmp = object.get(x_0__minCount_,"@type",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2__minCount_ = nodes_tmp2[_]
  nodes = x_2__minCount_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema / @type
  gen_x_check_103_array = gen_path_rule_105 with data.sourceNode as x
  count(gen_x_check_103_array) != 0 # validation applies if property was defined
  gen_x_check_103_string_set = { mapped |
    original := gen_x_check_103_array[_]
    mapped := as_string(original)
}

  gen_containsAll_104 = { "shacl:NodeShape"}
  count(gen_containsAll_104 - gen_x_check_103_string_set) != 0
  gen_x_check_103_quoted = [concat("", ["\"", res, "\""]) |  res := gen_x_check_103_string_set[_]]
  gen_x_check_103_string = concat("", ["[", concat(", ",gen_x_check_103_quoted), "]"])
  _result_0 := trace("hasValue","shapes.schema / @type",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_103_string,"expected": "[\"shacl:NodeShape\"]"})
  matches := error("json-object-top-level-request-response",x,"JSON object MUST be at the root of every JSON:API request and response containing data. This object defines a document’s “top level”.\n",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema / @type
  gen_propValues_106 = gen_path_rule_107 with data.sourceNode as x
  not count(gen_propValues_106) >= 1
  _result_0 := trace("minCount","shapes.schema / @type",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_106),"expected": 1})
  matches := error("json-object-top-level-request-response",x,"JSON object MUST be at the root of every JSON:API request and response containing data. This object defines a document’s “top level”.\n",[_result_0])
}
# Path rules

gen_path_rule_108[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_109[nodes] {
  init_y_0__nested_ = data.sourceNode
  tmp_y_0__nested_ = nested_nodes with data.nodes as init_y_0__nested_["shacl:property"]
  y_0__nested_ = tmp_y_0__nested_[_][_]
  nodes = y_0__nested_
}

gen_path_rule_112[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0__in_ = nodes_tmp2[_]
  nodes = z_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema
  ys = gen_path_rule_108 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.property
    zs = gen_path_rule_109 with data.sourceNode as y
    z_errorAcc0 = []
    zs_br_0 = [ zs_br_0_error|
      z = zs[_]
      #  querying path: shacl.name
      gen_z_check_111_array = gen_path_rule_112 with data.sourceNode as z
      gen_z_check_111_scalar = gen_z_check_111_array[_]
      gen_z_check_111 = as_string(gen_z_check_111_scalar)
      gen_inValues_110 = { "data","errors","meta"}
      not gen_inValues_110[gen_z_check_111]
      _result_0 := trace("in","shacl.name",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_111,"expected": "[\"data\",\"errors\",\"meta\"]"})
      zs_br_0_inner_error := error("nested",z,"error in nested nodes under shacl.property",[_result_0])
      zs_br_0_error = [z["@id"],zs_br_0_inner_error]
    ]
    zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
    zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
    z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
    z_errorAcc = z_errorAcc1
    # let's accumulate results
    zs_error_node_variables_agg = zs_br_0_errors
    not count(zs) - count(zs_error_node_variables_agg) >= 1
    _result_0 := trace("atLeast","shacl.property",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under shapes.schema",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  not count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","shapes.schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("json-object-required-fields",x,"A document MUST contain at least one of the following top-level members:\n\n    data: the document’s “primary data”\n    errors: an array of error objects\n    meta: a meta object that contains non-standard meta-information.\n",[_result_0])
}
# Path rules

gen_path_rule_113[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_114[nodes] {
  init_p_0__nested_ = data.sourceNode
  tmp_p_0__nested_ = nested_nodes with data.nodes as init_p_0__nested_["shacl:property"]
  p_0__nested_ = tmp_p_0__nested_[_][_]
  nodes = p_0__nested_
}

gen_path_rule_117[nodes] {
  init_q_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_q_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  q_0__in_ = nodes_tmp2[_]
  nodes = q_0__in_
}

gen_path_rule_118[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_119[nodes] {
  init_y_0__nested_ = data.sourceNode
  tmp_y_0__nested_ = nested_nodes with data.nodes as init_y_0__nested_["shacl:property"]
  y_0__nested_ = tmp_y_0__nested_[_][_]
  nodes = y_0__nested_
}

gen_path_rule_122[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0__in_ = nodes_tmp2[_]
  nodes = z_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema
  ps = gen_path_rule_113 with data.sourceNode as x
  p_errorAcc0 = []
  ps_br_0 = [ ps_br_0_error|
    p = ps[_]
    #  querying path: shacl.property
    qs = gen_path_rule_114 with data.sourceNode as p
    q_errorAcc0 = []
    qs_br_0 = [ qs_br_0_error|
      q = qs[_]
      #  querying path: shacl.name
      gen_q_check_116_array = gen_path_rule_117 with data.sourceNode as q
      gen_q_check_116_scalar = gen_q_check_116_array[_]
      gen_q_check_116 = as_string(gen_q_check_116_scalar)
      gen_inValues_115 = { "errors"}
      not gen_inValues_115[gen_q_check_116]
      _result_0 := trace("in","shacl.name",q,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_q_check_116,"expected": "[\"errors\"]"})
      qs_br_0_inner_error := error("nested",q,"error in nested nodes under shacl.property",[_result_0])
      qs_br_0_error = [q["@id"],qs_br_0_inner_error]
    ]
    qs_br_0_errors = { nodeId | n = qs_br_0[_]; nodeId = n[0] }
    qs_br_0_errors_errors = [ node | n = qs_br_0[_]; node = n[1] ]
    q_errorAcc1 = array.concat(q_errorAcc0,qs_br_0_errors_errors)
    q_errorAcc = q_errorAcc1
    # let's accumulate results
    qs_error_node_variables_agg = qs_br_0_errors
    not count(qs) - count(qs_error_node_variables_agg) >= 1
    _result_0 := trace("atLeast","shacl.property",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(qs_error_node_variables_agg), "successfulNodes":(count(qs)-count(qs_error_node_variables_agg)), "cardinality":1, "subResult": q_errorAcc})
    ps_br_0_inner_error := error("nested",p,"error in nested nodes under shapes.schema",[_result_0])
    ps_br_0_error = [p["@id"],ps_br_0_inner_error]
  ]
  ps_br_0_errors = { nodeId | n = ps_br_0[_]; nodeId = n[0] }
  ps_br_0_errors_errors = [ node | n = ps_br_0[_]; node = n[1] ]
  p_errorAcc1 = array.concat(p_errorAcc0,ps_br_0_errors_errors)
  p_errorAcc = p_errorAcc1
  # let's accumulate results
  ps_error_node_variables_agg = ps_br_0_errors
  count(ps) - count(ps_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","shapes.schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ps_error_node_variables_agg), "successfulNodes":(count(ps)-count(ps_error_node_variables_agg)), "cardinality":1, "subResult": p_errorAcc})
  #  querying path: shapes.schema
  ys = gen_path_rule_118 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.property
    zs = gen_path_rule_119 with data.sourceNode as y
    z_errorAcc0 = []
    zs_br_0 = [ zs_br_0_error|
      z = zs[_]
      #  querying path: shacl.name
      gen_z_check_121_array = gen_path_rule_122 with data.sourceNode as z
      gen_z_check_121_scalar = gen_z_check_121_array[_]
      gen_z_check_121 = as_string(gen_z_check_121_scalar)
      gen_inValues_120 = { "data"}
      not gen_inValues_120[gen_z_check_121]
      _result_0 := trace("in","shacl.name",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_121,"expected": "[\"data\"]"})
      zs_br_0_inner_error := error("nested",z,"error in nested nodes under shacl.property",[_result_0])
      zs_br_0_error = [z["@id"],zs_br_0_inner_error]
    ]
    zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
    zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
    z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
    z_errorAcc = z_errorAcc1
    # let's accumulate results
    zs_error_node_variables_agg = zs_br_0_errors
    not count(zs) - count(zs_error_node_variables_agg) >= 1
    _result_0 := trace("atLeast","shacl.property",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under shapes.schema",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","shapes.schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("json-object-no-error-and-data",x,"Validation error",[_result_0,_result_1])
}
# Path rules

gen_path_rule_125[nodes] {
  init_x_0__in_ = data.sourceNode
  tmp_x_0__in_ = nested_nodes with data.nodes as init_x_0__in_["shapes:schema"]
  x_0__in_ = tmp_x_0__in_[_][_]
  tmp_x_2__in_ = nested_nodes with data.nodes as x_0__in_["shacl:property"]
  x_2__in_ = tmp_x_2__in_[_][_]
  nodes_tmp = object.get(x_2__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3__in_ = nodes_tmp2[_]
  nodes = x_3__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema / shacl.property / shacl.name
  gen_x_check_124_array = gen_path_rule_125 with data.sourceNode as x
  gen_x_check_124_scalar = gen_x_check_124_array[_]
  gen_x_check_124 = as_string(gen_x_check_124_scalar)
  gen_inValues_123 = { "data","errors","meta","jsonapi","links","included"}
  not gen_inValues_123[gen_x_check_124]
  _result_0 := trace("in","shapes.schema / shacl.property / shacl.name",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_124,"expected": "[\"data\",\"errors\",\"meta\",\"jsonapi\",\"links\",\"included\"]"})
  matches := error("json-object-top-level-fields",x,"A document MUST contain at least one of the following top-level members:\n\n    data: the document’s “primary data”\n    errors: an array of error objects\n    meta: a meta object that contains non-standard meta-information.\n\nA document MAY contain any of these top-level members:\n\n    jsonapi: an object describing the server’s implementation\n    links: a links object related to the primary data.\n    included: an array of resource objects that are related to the primary data and/or each other (“included resources”).\n",[_result_0])
}
# Path rules

gen_path_rule_126[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_127[nodes] {
  init_p_0__nested_ = data.sourceNode
  tmp_p_0__nested_ = nested_nodes with data.nodes as init_p_0__nested_["shacl:property"]
  p_0__nested_ = tmp_p_0__nested_[_][_]
  nodes = p_0__nested_
}

gen_path_rule_130[nodes] {
  init_q_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_q_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  q_0__in_ = nodes_tmp2[_]
  nodes = q_0__in_
}

gen_path_rule_131[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  nodes = x_0__nested_
}

gen_path_rule_132[nodes] {
  init_y_0__nested_ = data.sourceNode
  tmp_y_0__nested_ = nested_nodes with data.nodes as init_y_0__nested_["shacl:property"]
  y_0__nested_ = tmp_y_0__nested_[_][_]
  nodes = y_0__nested_
}

gen_path_rule_135[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0__in_ = nodes_tmp2[_]
  nodes = z_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema
  ps = gen_path_rule_126 with data.sourceNode as x
  p_errorAcc0 = []
  ps_br_0 = [ ps_br_0_error|
    p = ps[_]
    #  querying path: shacl.property
    qs = gen_path_rule_127 with data.sourceNode as p
    q_errorAcc0 = []
    qs_br_0 = [ qs_br_0_error|
      q = qs[_]
      #  querying path: shacl.name
      gen_q_check_129_array = gen_path_rule_130 with data.sourceNode as q
      gen_q_check_129_scalar = gen_q_check_129_array[_]
      gen_q_check_129 = as_string(gen_q_check_129_scalar)
      gen_inValues_128 = { "included"}
      not gen_inValues_128[gen_q_check_129]
      _result_0 := trace("in","shacl.name",q,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_q_check_129,"expected": "[\"included\"]"})
      qs_br_0_inner_error := error("nested",q,"error in nested nodes under shacl.property",[_result_0])
      qs_br_0_error = [q["@id"],qs_br_0_inner_error]
    ]
    qs_br_0_errors = { nodeId | n = qs_br_0[_]; nodeId = n[0] }
    qs_br_0_errors_errors = [ node | n = qs_br_0[_]; node = n[1] ]
    q_errorAcc1 = array.concat(q_errorAcc0,qs_br_0_errors_errors)
    q_errorAcc = q_errorAcc1
    # let's accumulate results
    qs_error_node_variables_agg = qs_br_0_errors
    not count(qs) - count(qs_error_node_variables_agg) >= 1
    _result_0 := trace("atLeast","shacl.property",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(qs_error_node_variables_agg), "successfulNodes":(count(qs)-count(qs_error_node_variables_agg)), "cardinality":1, "subResult": q_errorAcc})
    ps_br_0_inner_error := error("nested",p,"error in nested nodes under shapes.schema",[_result_0])
    ps_br_0_error = [p["@id"],ps_br_0_inner_error]
  ]
  ps_br_0_errors = { nodeId | n = ps_br_0[_]; nodeId = n[0] }
  ps_br_0_errors_errors = [ node | n = ps_br_0[_]; node = n[1] ]
  p_errorAcc1 = array.concat(p_errorAcc0,ps_br_0_errors_errors)
  p_errorAcc = p_errorAcc1
  # let's accumulate results
  ps_error_node_variables_agg = ps_br_0_errors
  count(ps) - count(ps_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","shapes.schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ps_error_node_variables_agg), "successfulNodes":(count(ps)-count(ps_error_node_variables_agg)), "cardinality":1, "subResult": p_errorAcc})
  #  querying path: shapes.schema
  ys = gen_path_rule_131 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.property
    zs = gen_path_rule_132 with data.sourceNode as y
    z_errorAcc0 = []
    zs_br_0 = [ zs_br_0_error|
      z = zs[_]
      #  querying path: shacl.name
      gen_z_check_134_array = gen_path_rule_135 with data.sourceNode as z
      gen_z_check_134_scalar = gen_z_check_134_array[_]
      gen_z_check_134 = as_string(gen_z_check_134_scalar)
      gen_inValues_133 = { "data"}
      not gen_inValues_133[gen_z_check_134]
      _result_0 := trace("in","shacl.name",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_134,"expected": "[\"data\"]"})
      zs_br_0_inner_error := error("nested",z,"error in nested nodes under shacl.property",[_result_0])
      zs_br_0_error = [z["@id"],zs_br_0_inner_error]
    ]
    zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
    zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
    z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
    z_errorAcc = z_errorAcc1
    # let's accumulate results
    zs_error_node_variables_agg = zs_br_0_errors
    not count(zs) - count(zs_error_node_variables_agg) <= 0
    _result_0 := trace("atMost","shacl.property",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":0, "subResult": z_errorAcc})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under shapes.schema",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","shapes.schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("json-object-no-included-without-data",x,"If a document does not contain a top-level data key, the included member MUST NOT be present either.",[_result_0,_result_1])
}
# Path rules

gen_path_rule_136[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shacl:property"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  nodes = x_2__nested_
}

gen_path_rule_139[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0__in_ = nodes_tmp2[_]
  nodes = y_0__in_
}

gen_path_rule_140[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shacl:property"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  nodes = x_2__nested_
}

gen_path_rule_143[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0__in_ = nodes_tmp2[_]
  nodes = z_0__in_
}

gen_path_rule_144[nodes] {
  init_z_0__nested_ = data.sourceNode
  tmp_z_0__nested_ = nested_nodes with data.nodes as init_z_0__nested_["shapes:range"]
  z_0__nested_ = tmp_z_0__nested_[_][_]
  nodes = z_0__nested_
}

gen_path_rule_147[nodes] {
  init_p_0__in_ = data.sourceNode
  tmp_p_0__in_ = nested_nodes with data.nodes as init_p_0__in_["shacl:property"]
  p_0__in_ = tmp_p_0__in_[_][_]
  nodes_tmp = object.get(p_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  p_2__in_ = nodes_tmp2[_]
  nodes = p_2__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema / shacl.property
  ys = gen_path_rule_136 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.name
    gen_y_check_138_array = gen_path_rule_139 with data.sourceNode as y
    gen_y_check_138_scalar = gen_y_check_138_array[_]
    gen_y_check_138 = as_string(gen_y_check_138_scalar)
    gen_inValues_137 = { "links"}
    not gen_inValues_137[gen_y_check_138]
    _result_0 := trace("in","shacl.name",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_138,"expected": "[\"links\"]"})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  #  querying path: shapes.schema / shacl.property
  zs = gen_path_rule_140 with data.sourceNode as x
  z_errorAcc0 = []
  zs_br_0 = [ zs_br_0_error|
    z = zs[_]
    #  querying path: shacl.name
    gen_z_check_142_array = gen_path_rule_143 with data.sourceNode as z
    gen_z_check_142_scalar = gen_z_check_142_array[_]
    gen_z_check_142 = as_string(gen_z_check_142_scalar)
    gen_inValues_141 = { "links"}
    not gen_inValues_141[gen_z_check_142]
    _result_0 := trace("in","shacl.name",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_142,"expected": "[\"links\"]"})
    zs_br_0_inner_error := error("nested",z,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    zs_br_0_error = [z["@id"],zs_br_0_inner_error]
  ]
  zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
  zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
  z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
  zs_br_1 = [ zs_br_1_error|
    z = zs[_]
    #  querying path: shapes.range
    ps = gen_path_rule_144 with data.sourceNode as z
    p_errorAcc0 = []
    ps_br_0 = [ ps_br_0_error|
      p = ps[_]
      #  querying path: shacl.property / shacl.name
      gen_p_check_146_array = gen_path_rule_147 with data.sourceNode as p
      gen_p_check_146_scalar = gen_p_check_146_array[_]
      gen_p_check_146 = as_string(gen_p_check_146_scalar)
      gen_inValues_145 = { "self","related"}
      not gen_inValues_145[gen_p_check_146]
      _result_0 := trace("in","shacl.property / shacl.name",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_p_check_146,"expected": "[\"self\",\"related\"]"})
      ps_br_0_inner_error := error("nested",p,"error in nested nodes under shapes.range",[_result_0])
      ps_br_0_error = [p["@id"],ps_br_0_inner_error]
    ]
    ps_br_0_errors = { nodeId | n = ps_br_0[_]; nodeId = n[0] }
    ps_br_0_errors_errors = [ node | n = ps_br_0[_]; node = n[1] ]
    p_errorAcc1 = array.concat(p_errorAcc0,ps_br_0_errors_errors)
    p_errorAcc = p_errorAcc1
    # let's accumulate results
    ps_error_node_variables_agg = ps_br_0_errors
    count(ps_error_node_variables_agg) > 0
    _result_0 := trace("nested","shapes.range",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ps_error_node_variables_agg), "successfulNodes":(count(ps)-count(ps_error_node_variables_agg)),"subResult": p_errorAcc})
    zs_br_1_inner_error := error("nested",z,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    zs_br_1_error = [z["@id"],zs_br_1_inner_error]
  ]
  zs_br_1_errors = { nodeId | n = zs_br_1[_]; nodeId = n[0] }
  zs_br_1_errors_errors = [ node | n = zs_br_1[_]; node = n[1] ]
  z_errorAcc2 = array.concat(z_errorAcc1,zs_br_1_errors_errors)
  z_errorAcc = z_errorAcc2
  # let's accumulate results
  zs_error_node_variables_agg = zs_br_0_errors | zs_br_1_errors
  not count(zs) - count(zs_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
  matches := error("json-object-links-field",x,"The top-level links object MAY contain the following members:\n\n    self: the link that generated the current response document.\n    related: a related resource link when the primary data represents a resource relationship.\n    pagination links for the primary data.\n",[_result_0,_result_1])
}
# Path rules

gen_path_rule_148[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shacl:property"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  nodes = x_2__nested_
}

gen_path_rule_151[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0__in_ = nodes_tmp2[_]
  nodes = y_0__in_
}

gen_path_rule_152[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["shapes:schema"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shacl:property"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  nodes = x_2__nested_
}

gen_path_rule_155[nodes] {
  init_z_0__in_ = data.sourceNode
  tmp_z_0__in_ = nested_nodes with data.nodes as init_z_0__in_["shapes:range"]
  z_0__in_ = tmp_z_0__in_[_][_]
  nodes_tmp = object.get(z_0__in_,"@type",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_2__in_ = nodes_tmp2[_]
  nodes = z_2__in_
}

gen_path_rule_158[nodes] {
  init_z_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_z_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  z_0__in_ = nodes_tmp2[_]
  nodes = z_0__in_
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Payload"
  #  querying path: shapes.schema / shacl.property
  ys = gen_path_rule_148 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.name
    gen_y_check_150_array = gen_path_rule_151 with data.sourceNode as y
    gen_y_check_150_scalar = gen_y_check_150_array[_]
    gen_y_check_150 = as_string(gen_y_check_150_scalar)
    gen_inValues_149 = { "data"}
    not gen_inValues_149[gen_y_check_150]
    _result_0 := trace("in","shacl.name",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_150,"expected": "[\"data\"]"})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  #  querying path: shapes.schema / shacl.property
  zs = gen_path_rule_152 with data.sourceNode as x
  z_errorAcc0 = []
  zs_br_0 = [ zs_br_0_error|
    z = zs[_]
    #  querying path: shapes.range / @type
    gen_z_check_153_array = gen_path_rule_155 with data.sourceNode as z
    count(gen_z_check_153_array) != 0 # validation applies if property was defined
    gen_z_check_153_string_set = { mapped |
    original := gen_z_check_153_array[_]
    mapped := as_string(original)
}

    gen_containsAll_154 = { "shacl:NodeShape","shapes:ArrayShape","shapes:NilShape"}
    count(gen_containsAll_154 - gen_z_check_153_string_set) != 0
    gen_z_check_153_quoted = [concat("", ["\"", res, "\""]) |  res := gen_z_check_153_string_set[_]]
    gen_z_check_153_string = concat("", ["[", concat(", ",gen_z_check_153_quoted), "]"])
    _result_0 := trace("containsAll","shapes.range / @type",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_153_string,"expected": "[\"shacl:NodeShape\",\"shapes:ArrayShape\",\"shapes:NilShape\"]"})
    zs_br_0_inner_error := error("nested",z,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    zs_br_0_error = [z["@id"],zs_br_0_inner_error]
  ]
  zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
  zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
  z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
  zs_br_1 = [ zs_br_1_error|
    z = zs[_]
    #  querying path: shacl.name
    gen_z_check_157_array = gen_path_rule_158 with data.sourceNode as z
    gen_z_check_157_scalar = gen_z_check_157_array[_]
    gen_z_check_157 = as_string(gen_z_check_157_scalar)
    gen_inValues_156 = { "data"}
    not gen_inValues_156[gen_z_check_157]
    _result_0 := trace("in","shacl.name",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_157,"expected": "[\"data\"]"})
    zs_br_1_inner_error := error("nested",z,"error in nested nodes under shapes.schema / shacl.property",[_result_0])
    zs_br_1_error = [z["@id"],zs_br_1_inner_error]
  ]
  zs_br_1_errors = { nodeId | n = zs_br_1[_]; nodeId = n[0] }
  zs_br_1_errors_errors = [ node | n = zs_br_1[_]; node = n[1] ]
  z_errorAcc2 = array.concat(z_errorAcc1,zs_br_1_errors_errors)
  z_errorAcc = z_errorAcc2
  # let's accumulate results
  zs_error_node_variables_agg = zs_br_0_errors | zs_br_1_errors
  not count(zs) - count(zs_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
  matches := error("json-object-primary-data",x,"The document’s “primary data” is a representation of the resource or collection of resources targeted by a request.\n\nPrimary data MUST be either:\n\na single resource object, a single resource identifier object, or null, for requests that target single resources\nan array of resource objects, an array of resource identifier objects, or an empty array ([]), for requests that target resource collections\n\nFor example, the following primary data is a single resource object:\n\n{\n   'data': {\n     'type': 'articles',\n     'id': '1',\n     'attributes': {\n        // ... this article's attributes\n     },\n     'relationships': {\n       // ... this article's relationships\n     }\n   }\n}\n\nThe following primary data is a single resource identifier object that references the same resource:\n\n{\n  'data': {\n    'type': 'articles',\n    'id': '1'\n  }\n}\n\nA logical collection of resources MUST be represented as an array, even if it only contains one item or is empty.\n",[_result_0,_result_1])
}
# Path rules

gen_path_rule_159[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["apiContract:payload"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shapes:schema"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  tmp_x_3__nested_ = nested_nodes with data.nodes as x_2__nested_["shacl:property"]
  x_3__nested_ = tmp_x_3__nested_[_][_]
  nodes = x_3__nested_
}

gen_path_rule_162[nodes] {
  init_p_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_p_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  p_0__in_ = nodes_tmp2[_]
  nodes = p_0__in_
}

gen_path_rule_163[nodes] {
  init_p_0__nested_ = data.sourceNode
  tmp_p_0__nested_ = nested_nodes with data.nodes as init_p_0__nested_["shapes:range"]
  p_0__nested_ = tmp_p_0__nested_[_][_]
  nodes = p_0__nested_
}

gen_path_rule_164[nodes] {
  init_q_0__nested_ = data.sourceNode
  tmp_q_0__nested_ = nested_nodes with data.nodes as init_q_0__nested_["shacl:property"]
  q_0__nested_ = tmp_q_0__nested_[_][_]
  nodes = q_0__nested_
}

gen_path_rule_167[nodes] {
  init_r_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_r_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  r_0__in_ = nodes_tmp2[_]
  nodes = r_0__in_
}

gen_path_rule_168[nodes] {
  init_x_0__nested_ = data.sourceNode
  tmp_x_0__nested_ = nested_nodes with data.nodes as init_x_0__nested_["apiContract:payload"]
  x_0__nested_ = tmp_x_0__nested_[_][_]
  tmp_x_2__nested_ = nested_nodes with data.nodes as x_0__nested_["shapes:schema"]
  x_2__nested_ = tmp_x_2__nested_[_][_]
  tmp_x_3__nested_ = nested_nodes with data.nodes as x_2__nested_["shacl:property"]
  x_3__nested_ = tmp_x_3__nested_[_][_]
  nodes = x_3__nested_
}

gen_path_rule_171[nodes] {
  init_y_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_y_0__in_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0__in_ = nodes_tmp2[_]
  nodes = y_0__in_
}

gen_path_rule_172[nodes] {
  init_y_0__nested_ = data.sourceNode
  tmp_y_0__nested_ = nested_nodes with data.nodes as init_y_0__nested_["shapes:range"]
  y_0__nested_ = tmp_y_0__nested_[_][_]
  nodes = y_0__nested_
}



# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Response"
  #  querying path: apiContract.payload / shapes.schema / shacl.property
  ps = gen_path_rule_159 with data.sourceNode as x
  p_errorAcc0 = []
  ps_br_0 = [ ps_br_0_error|
    p = ps[_]
    #  querying path: shacl.name
    gen_p_check_161_array = gen_path_rule_162 with data.sourceNode as p
    gen_p_check_161_scalar = gen_p_check_161_array[_]
    gen_p_check_161 = as_string(gen_p_check_161_scalar)
    gen_inValues_160 = { "data"}
    not gen_inValues_160[gen_p_check_161]
    _result_0 := trace("in","shacl.name",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_p_check_161,"expected": "[\"data\"]"})
    ps_br_0_inner_error := error("nested",p,"error in nested nodes under apiContract.payload / shapes.schema / shacl.property",[_result_0])
    ps_br_0_error = [p["@id"],ps_br_0_inner_error]
  ]
  ps_br_0_errors = { nodeId | n = ps_br_0[_]; nodeId = n[0] }
  ps_br_0_errors_errors = [ node | n = ps_br_0[_]; node = n[1] ]
  p_errorAcc1 = array.concat(p_errorAcc0,ps_br_0_errors_errors)
  ps_br_1 = [ ps_br_1_error|
    p = ps[_]
    #  querying path: shapes.range
    qs = gen_path_rule_163 with data.sourceNode as p
    q_errorAcc0 = []
    qs_br_0 = [ qs_br_0_error|
      q = qs[_]
      #  querying path: shacl.property
      rs = gen_path_rule_164 with data.sourceNode as q
      r_errorAcc0 = []
      rs_br_0 = [ rs_br_0_error|
        r = rs[_]
        #  querying path: shacl.name
        gen_r_check_166_array = gen_path_rule_167 with data.sourceNode as r
        gen_r_check_166_scalar = gen_r_check_166_array[_]
        gen_r_check_166 = as_string(gen_r_check_166_scalar)
        gen_inValues_165 = { "id","type"}
        not gen_inValues_165[gen_r_check_166]
        _result_0 := trace("in","shacl.name",r,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_r_check_166,"expected": "[\"id\",\"type\"]"})
        rs_br_0_inner_error := error("nested",r,"error in nested nodes under shacl.property",[_result_0])
        rs_br_0_error = [r["@id"],rs_br_0_inner_error]
      ]
      rs_br_0_errors = { nodeId | n = rs_br_0[_]; nodeId = n[0] }
      rs_br_0_errors_errors = [ node | n = rs_br_0[_]; node = n[1] ]
      r_errorAcc1 = array.concat(r_errorAcc0,rs_br_0_errors_errors)
      r_errorAcc = r_errorAcc1
      # let's accumulate results
      rs_error_node_variables_agg = rs_br_0_errors
      not count(rs) - count(rs_error_node_variables_agg) >= 2
      _result_0 := trace("atLeast","shacl.property",q,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(rs_error_node_variables_agg), "successfulNodes":(count(rs)-count(rs_error_node_variables_agg)), "cardinality":2, "subResult": r_errorAcc})
      qs_br_0_inner_error := error("nested",q,"error in nested nodes under shapes.range",[_result_0])
      qs_br_0_error = [q["@id"],qs_br_0_inner_error]
    ]
    qs_br_0_errors = { nodeId | n = qs_br_0[_]; nodeId = n[0] }
    qs_br_0_errors_errors = [ node | n = qs_br_0[_]; node = n[1] ]
    q_errorAcc1 = array.concat(q_errorAcc0,qs_br_0_errors_errors)
    q_errorAcc = q_errorAcc1
    # let's accumulate results
    qs_error_node_variables_agg = qs_br_0_errors
    count(qs_error_node_variables_agg) > 0
    _result_0 := trace("nested","shapes.range",p,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(qs_error_node_variables_agg), "successfulNodes":(count(qs)-count(qs_error_node_variables_agg)),"subResult": q_errorAcc})
    ps_br_1_inner_error := error("nested",p,"error in nested nodes under apiContract.payload / shapes.schema / shacl.property",[_result_0])
    ps_br_1_error = [p["@id"],ps_br_1_inner_error]
  ]
  ps_br_1_errors = { nodeId | n = ps_br_1[_]; nodeId = n[0] }
  ps_br_1_errors_errors = [ node | n = ps_br_1[_]; node = n[1] ]
  p_errorAcc2 = array.concat(p_errorAcc1,ps_br_1_errors_errors)
  p_errorAcc = p_errorAcc2
  # let's accumulate results
  ps_error_node_variables_agg = ps_br_0_errors | ps_br_1_errors
  not count(ps) - count(ps_error_node_variables_agg) >= 1
  _result_0 := trace("atLeast","apiContract.payload / shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ps_error_node_variables_agg), "successfulNodes":(count(ps)-count(ps_error_node_variables_agg)), "cardinality":1, "subResult": p_errorAcc})
  #  querying path: apiContract.payload / shapes.schema / shacl.property
  ys = gen_path_rule_168 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.name
    gen_y_check_170_array = gen_path_rule_171 with data.sourceNode as y
    gen_y_check_170_scalar = gen_y_check_170_array[_]
    gen_y_check_170 = as_string(gen_y_check_170_scalar)
    gen_inValues_169 = { "data"}
    not gen_inValues_169[gen_y_check_170]
    _result_0 := trace("in","shacl.name",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_170,"expected": "[\"data\"]"})
    ys_br_0_inner_error := error("nested",y,"error in nested nodes under apiContract.payload / shapes.schema / shacl.property",[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  ys_br_1 = [ ys_br_1_error|
    y = ys[_]
    #  querying path: shapes.range
    zs = gen_path_rule_172 with data.sourceNode as y
    z_errorAcc0 = []
    zs_br_0 = [ zs_br_0_error|
      z = zs[_]
      types = object.get(z, "@type", [])
      nodeshapes = [ r |
        e = types[_]
        e == "shacl:NodeShape"
        r = e
      ]
      gen_rego_result_174 = (count(nodeshapes) == 1)
      
      gen_rego_result_174 != true
      _result_0 := trace("rego","",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false})
      zs_br_0_inner_error := error("nested",z,"error in nested nodes under shapes.range",[_result_0])
      zs_br_0_error = [z["@id"],zs_br_0_inner_error]
    ]
    zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
    zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
    z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
    z_errorAcc = z_errorAcc1
    # let's accumulate results
    zs_error_node_variables_agg = zs_br_0_errors
    count(zs_error_node_variables_agg) > 0
    _result_0 := trace("nested","shapes.range",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)),"subResult": z_errorAcc})
    ys_br_1_inner_error := error("nested",y,"error in nested nodes under apiContract.payload / shapes.schema / shacl.property",[_result_0])
    ys_br_1_error = [y["@id"],ys_br_1_inner_error]
  ]
  ys_br_1_errors = { nodeId | n = ys_br_1[_]; nodeId = n[0] }
  ys_br_1_errors_errors = [ node | n = ys_br_1[_]; node = n[1] ]
  y_errorAcc2 = array.concat(y_errorAcc1,ys_br_1_errors_errors)
  y_errorAcc = y_errorAcc2
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors | ys_br_1_errors
  count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_1 := trace("atLeast","apiContract.payload / shapes.schema / shacl.property",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  matches := error("resource-object-required-fields",x,"A resource object MUST contain at least the following top-level members:\n\n    id\n    type\n\nException: The id member is not required when the resource object originates at the client and represents a new resource to be created on the server.\n",[_result_0,_result_1])
}