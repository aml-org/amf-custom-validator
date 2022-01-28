package profile_name_with_3_non_alphanumeric_characters_like_

report["profile"] = "name with 3 non-alphanumeric characters like %, ^ & *"

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

gen_path_rule_3[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_5[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://www.w3.org/ns/shacl#name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_7[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

gen_path_rule_8[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://www.w3.org/ns/shacl#name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.method
  gen_x_check_2_array = gen_path_rule_3 with data.sourceNode as x
  gen_x_check_2_scalar = gen_x_check_2_array[_]
  gen_x_check_2 = as_string(gen_x_check_2_scalar)
  gen_inValues_1 = { "publish","subscribe","1","2"}
  not gen_inValues_1[gen_x_check_2]
  _result_0 := trace("in","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_2,"expected": "[\"publish\",\"subscribe\",\"1\",\"2\"]"})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: shacl.name
  gen_propValues_4 = gen_path_rule_5 with data.sourceNode as x
  not count(gen_propValues_4) <= 1
  _result_0 := trace("maxCount","http://www.w3.org/ns/shacl#name",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":"<=","actual": count(gen_propValues_4),"expected": 1})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: apiContract.method
  gen_propValues_6 = gen_path_rule_7 with data.sourceNode as x
  not count(gen_propValues_6) >= 1
  _result_0 := trace("minCount","http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_6),"expected": 1})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"
  #  querying path: shacl.name
  gen_gen_path_rule_8_node_9_array = gen_path_rule_8 with data.sourceNode as x
  gen_gen_path_rule_8_node_9 = gen_gen_path_rule_8_node_9_array[_]
  not regex.match("^put|post$",gen_gen_path_rule_8_node_9)
  _result_0 := trace("pattern","http://www.w3.org/ns/shacl#name",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"argument": gen_gen_path_rule_8_node_9})
  matches := error("validation1",x,"This is the message",[_result_0])
}