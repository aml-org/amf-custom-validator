package profile_test_4

report["profile"] = "Test 4"

# Import future keywords
import future.keywords.in
import future.keywords.every
import future.keywords.if
import future.keywords.contains

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

search_custom_property_subjects[valid_subject] {
  extension = data.property_extension
  object = data.object
  
  object["http://a.ml/vocabularies/core#extensionName"] = extension
  object_id = object["@id"]
  
  node = input["@ids"][_]
  
  node[prop] = {"@id": object_id}
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

# Finds a nested custom domain property for a given node an custom domain property name
gen_path_extension[nodes] {
  sourceNode := data.custom_property_data[0]
  extensionName := data.custom_property_data[1]
  customPropertiesTarget := object.get(sourceNode, "http://a.ml/vocabularies/document#customDomainProperties", [])
  customProperties = nodes_array with data.nodes as customPropertiesTarget
  extensionFound = [found |
    annotationLink = customProperties[_]
    annotationLinkId = annotationLink["@id"]
    annotation = object.get(sourceNode, annotationLinkId, {})
    annotationNode = find with data.link as annotation
  
    name = annotationNode["http://a.ml/vocabularies/core#extensionName"]
    name = extensionName
  
    found = annotationNode
  ]
  
  nodes = extensionFound
}

#creates an array from a comma separated values str
split_values[values] {
   values := split(data.nodes, ",")
}

#check if a given target is present at the list serialized as comma separated values(chain)
values_contains(chain, target) = matches {
  split_values[values] with data.nodes as chain
    values[_] = target
    matches := true
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

gen_path_set_rule_2[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/shapes#schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_set_rule_4[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/shapes#schema",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_set_rule_5[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/shapes#schema"]
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_set_rule_7[nodes] {
  init_y_0 = data.sourceNode
  nodes_tmp = object.get(init_y_0,"http://www.w3.org/ns/shacl#minLength",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  y_0 = nodes_tmp2[_]
  nodes = y_0
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_1 = gen_path_set_rule_2 with data.sourceNode as x
  not count(gen_propValues_1) <= 3
  _result_0 := trace("maxCount","http://a.ml/vocabularies/shapes#schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":"<=","actual": count(gen_propValues_1),"expected": 3})
  message := "Scalars in parameters must have minLength defined"
  matches := error("validation1",x, message ,[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Parameter"
  #  querying path: raml-shapes.schema
  gen_propValues_3 = gen_path_set_rule_4 with data.sourceNode as x
  not count(gen_propValues_3) >= 1
  _result_0 := trace("minCount","http://a.ml/vocabularies/shapes#schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_3),"expected": 1})
  message := "Scalars in parameters must have minLength defined"
  matches := error("validation1",x, message ,[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Parameter"
  #  querying path: raml-shapes.schema
  ys = gen_path_set_rule_5 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: shacl.minLength
    gen_propValues_6 = gen_path_set_rule_7 with data.sourceNode as y
    not count(gen_propValues_6) >= 1
    _result_0 := trace("minCount","http://www.w3.org/ns/shacl#minLength",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_6),"expected": 1})
    message := "error in nested nodes under http://a.ml/vocabularies/shapes#schema"
    ys_br_0_inner_error := error("nested",y, message ,[_result_0])
    ys_br_0_error = [y["@id"],ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  count(ys_error_node_variables_agg) > 0
  _result_0 := trace("nested","http://a.ml/vocabularies/shapes#schema",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)),"subResult": y_errorAcc})
  message := "Scalars in parameters must have minLength defined"
  matches := error("validation1",x, message ,[_result_0])
}