package profile_test9


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



# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  # custom1
  version = object.get(x, "apiContract:version", null)
  gen_rego_result_2 = (version != null)
  
  gen_rego_result_2 == true
  _result_0 := trace("rego","",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":true})
  matches := error("simple-rego",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0])
}
# Path rules



# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  # custom2
  version = object.get(x, "apiContract:version", null)
  gen_rego_result_4 = (version != null)
  
  gen_rego_result_4 != true
  _result_0 := trace("rego","",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false})
  matches := error("simple-rego2",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0])
}
# Path rules

gen_path_rule_5[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"apiContract:version",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:WebAPI"
  #  querying path: apiContract.version
  gen_gen_path_rule_5_node_6_array = gen_path_rule_5 with data.sourceNode as x
  gen_gen_path_rule_5_node_6 = gen_gen_path_rule_5_node_6_array
  gen_rego_result_7 = (gen_gen_path_rule_5_node_6 != null) # custom 3
  gen_rego_result_7 != true
  _result_0 := trace("rego","apiContract.version",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false})
  matches := error("simple-rego3",x,"GET operations must have 2xx, 4xx and 5xx status codes but no 201",[_result_0])
}