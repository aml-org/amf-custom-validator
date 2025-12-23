package profile_test_1

report["profile"] = "Test 1"

# Rego v1
import rego.v1

############################
# Core helpers (functions)
############################

# Finds a node in the graph, following a link in the flattened JSON-LD node.
# Expects link to be an object with an "@id" key.
find(link) := node if {
  id := link["@id"]
  node := input["@ids"][id]
}

# Wraps a value in an array if it isn't already one.
nodes_array(x) := x if {
  is_array(x)
} else := [x]

# Navigates an object property for a given node/link list: returns a SET of nested nodes.
nested(links) := out if {
  out := {n |
    link := nodes_array(links)[_]
    n := find(link)
  }
}

# Fetches all the nodes for an object link: returns an ARRAY of nested nodes.
nested_nodes(links) := result if {
  result := [n |
    link := nodes_array(links)[_]
    n := find(link)
  ]
}

# Fetches all the scalars for a scalar link: returns a SET of values for property from each node in nodes.
nested_values(nodes, property) := result if {
  result := {value |
    n := nodes_array(nodes)[_]
    value := n[property]
  }
}

# Extracts in a type-safe way the @id of a node object.
node_id(y) := id if {
  is_object(y)
  id := object.get(y, "@id", "")
} else := ""

############################
# Subject search helpers
############################

# Fetches all the subject nodes that have certain predicate and object.
# Returns a SET of node objects.
search_subjects(predicate, data_obj) := result if {
  result := {node |
    node := input["@ids"][_]
    vals := nodes_array(object.get(node, predicate, []))
    v := vals[_]
    v["@id"] == data_obj["@id"]
  }
}

# Fetches all the subject nodes where a given custom property extension matches.
search_custom_property_subjects(property_extension, data_obj, prop) := result if {
  result := {node |
    data_obj["http://a.ml/vocabularies/core#extensionName"] == property_extension
    object_id := data_obj["@id"]

    node := input["@ids"][_]
    node[prop] == {"@id": object_id}
  }
}

############################
# Collection helpers
############################

# Collect next set of nodes: from nodes, follow property links, then resolve to node objects.
# Returns a SET of node objects.
collect(nodes, property) := result if {
  result := {r |
    n := nodes_array(nodes)[_]
    rs := object.get(n, property, [])
    r := nested(rs)[_]
  }
}

# Collect terminal values: from nodes, fetch property values and flatten to an array.
# Returns a SET of values.
collect_values(nodes, property) := result if {
  result := {r |
    n := nodes_array(nodes)[_]
    rs := object.get(n, property, [])
    r := nodes_array(rs)[_]
  }
}

############################
# Datatype / string helpers
############################

check_datatype(x, dt) := true if {
  dt == "http://www.w3.org/2001/XMLSchema#string"
  is_string(x)
} else := true if {
  dt == "http://www.w3.org/2001/XMLSchema#integer"
  is_number(x)
} else := true if {
  dt == "http://www.w3.org/2001/XMLSchema#float"
  is_number(x)
} else := true if {
  dt == "http://www.w3.org/2001/XMLSchema#boolean"
  is_boolean(x)
} else := true if {
  is_object(x)
  t := object.get(x, "@type", "")
  t == dt
} else := false if {
  not is_object(x)
  dt != "http://www.w3.org/2001/XMLSchema#string"
  dt != "http://www.w3.org/2001/XMLSchema#integer"
  dt != "http://www.w3.org/2001/XMLSchema#float"
  dt != "http://www.w3.org/2001/XMLSchema#boolean"
}

# Fetches all the nodes for a given RDF class.
# Returns a SET of node objects.
target_class(class) := result if {
  result := {node |
    id := input["@types"][class][_]
    node := input["@ids"][id]
  }
}

# Transform scalars to string, useful for 'in' constraints.
as_string(x) := x if {
  is_string(x)
} else := "false" if {
  is_boolean(x)
  x == false
} else := "true" if {
  is_boolean(x)
  x == true
} else := format_int(x, 10) if {
  is_number(x)
} else := x["@id"] if {
  is_object(x)
  x["@id"]
} else := json.marshal(x) if {
  is_object(x)
  not x["@id"]
}

############################
# Custom-domain-property helpers
############################

# Finds a nested custom domain property for a given node and custom domain property name.
# Returns an ARRAY of matching annotation nodes.
gen_path_extension(sourceNode, extensionName) := nodes if {
  customPropertiesTarget := object.get(
    sourceNode,
    "http://a.ml/vocabularies/document#customDomainProperties",
    []
  )
  customProperties := nodes_array(customPropertiesTarget)

  nodes := [found |
    annotationLink := customProperties[_]
    annotationLinkId := annotationLink["@id"]
    annotation := object.get(sourceNode, annotationLinkId, {})
    annotationNode := find(annotation)

    annotationNode["http://a.ml/vocabularies/core#extensionName"] == extensionName
    found := annotationNode
  ]
}

############################
# CSV helpers
############################

split_values(s) := values if {
  values := split(s, ",")
}

values_contains(chain, target) := true if {
  values := split_values(chain)
  values[_] == target
} else := false

############################
# Tracing + error helpers
############################

traceRuleExecution(constraint, resultPath, focusNode, traceValue) := t if {
  l := location(focusNode)
  t := {
    "@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue,
    "location": l,
  }
} else := t if {
  not location(focusNode)
  t := {
    "@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
    "component": constraint,
    "resultPath": resultPath,
    "traceValue": traceValue,
  }
}

location(focusNode) := l if {
  id := focusNode["@id"]
  loc := input["@lexical"][id]
  raw_range := loc["range"]
  uri := loc["uri"]

  range_parts := regex.find_n("\\d+", raw_range, 4)
  range := {
    "@type": ["lexicalSchema:RangeNode", "lexical:Range"],
    "start": {
      "@type": ["lexicalSchema:PositionNode", "lexical:Position"],
      "line": to_number(range_parts[0]),
      "column": to_number(range_parts[1]),
    },
    "end": {
      "@type": ["lexicalSchema:PositionNode", "lexical:Position"],
      "line": to_number(range_parts[2]),
      "column": to_number(range_parts[3]),
    },
  }

  l := {
    "@type": ["lexicalSchema:LocationNode", "lexical:Location"],
    "uri": uri,
    "range": range,
  }
}

error(sourceShapeName, focusNode, resultMessage, traceLog) := e if {
  id := focusNode["@id"]
  locationNode := location(focusNode)
  e := {
    "@type": ["reportSchema:ValidationResultNode", "shacl:ValidationResult"],
    "sourceShapeName": sourceShapeName,
    "focusNode": id,
    "resultMessage": resultMessage,
    "location": locationNode,
    "trace": traceLog,
  }
} else := e if {
  id := focusNode["@id"]
  not location(focusNode)
  e := {
    "@type": ["reportSchema:ValidationResultNode", "shacl:ValidationResult"],
    "sourceShapeName": sourceShapeName,
    "focusNode": id,
    "resultMessage": resultMessage,
    "trace": traceLog,
  }
}

# generate the report for violation level
# default value must be added dynamically
report[level] := matches if {
  vs := violation
  level := "violation"
  matches := vs
}

# generate the report for the info level
# default value must be added dynamically
report[level] := matches if {
  vs := info
  level := "info"
  matches := vs
}

# generate the report for the warning level
# default value must be added dynamically
report[level] := matches if {
  vs := warning
  level := "warning"
  matches := vs
}


default warning = []

default info = []
# Path rules

gen_path_set_rule_3 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#supportedOperation"])
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_2 = nodes_tmp2[_]
  nodes = x_2
}

gen_path_set_rule_5 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#supportedOperation"])
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"http://a.ml/vocabularies/apiContract#method",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_2 = nodes_tmp2[_]
  nodes = x_2
}

gen_path_set_rule_6 contains nodes if {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://www.w3.org/ns/shacl#name",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_0 = nodes_tmp2[_]
  nodes = x_0
}

# Constraint rules

violation contains matches if {
  x = target_class("http://a.ml/vocabularies/apiContract#EndPoint")[_]
  #  querying path: apiContract.supportedOperation / apiContract.method
  gen_x_check_2_array = gen_path_set_rule_3 with data.sourceNode as x
  gen_x_check_2_scalar = gen_x_check_2_array[_]
  gen_x_check_2 = as_string(gen_x_check_2_scalar)
  gen_inValues_1 = { "publish","subscribe"}
  not gen_inValues_1[gen_x_check_2]
  _result_0 := traceRuleExecution("in","http://a.ml/vocabularies/apiContract#supportedOperation / http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_x_check_2,"expected": "[\"publish\",\"subscribe\"]"})
  message := "This is the message"
  matches := error("validation1",x, message ,[_result_0])
}

violation contains matches if {
  x = target_class("http://a.ml/vocabularies/apiContract#EndPoint")[_]
  #  querying path: apiContract.supportedOperation / apiContract.method
  gen_propValues_4 = gen_path_set_rule_5 with data.sourceNode as x
  not count(gen_propValues_4) >= 1
  _result_0 := traceRuleExecution("minCount","http://a.ml/vocabularies/apiContract#supportedOperation / http://a.ml/vocabularies/apiContract#method",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_4),"expected": 1})
  message := "This is the message"
  matches := error("validation1",x, message ,[_result_0])
}

violation contains matches if {
  x = target_class("http://a.ml/vocabularies/apiContract#EndPoint")[_]
  #  querying path: shacl.name
  gen_gen_path_set_rule_6_node_7_array = gen_path_set_rule_6 with data.sourceNode as x
  gen_gen_path_set_rule_6_node_7 = gen_gen_path_set_rule_6_node_7_array[_]
  not regex.match(`^put|post$`,as_string(gen_gen_path_set_rule_6_node_7))
  _result_0 := traceRuleExecution("pattern","http://www.w3.org/ns/shacl#name",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"expected": "^put|post$","actual": gen_gen_path_set_rule_6_node_7})
  message := "This is the message"
  matches := error("validation1",x, message ,[_result_0])
}
# Path rules

gen_path_set_rule_9 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#expects"])
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes(x_0["http://a.ml/vocabularies/apiContract#parameter"])
  x_2 = tmp_x_2[_][_]
  tmp_x_3 = nested_nodes(x_2["http://a.ml/vocabularies/shapes#schema"])
  x_3 = tmp_x_3[_][_]
  nodes_tmp = object.get(x_3,"http://www.w3.org/ns/shacl#name",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_4 = nodes_tmp2[_]
  nodes = x_4
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#expects"])
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes(x_0["http://a.ml/vocabularies/apiContract#payload"])
  x_2 = tmp_x_2[_][_]
  tmp_x_3 = nested_nodes(x_2["http://a.ml/vocabularies/shapes#schema"])
  x_3 = tmp_x_3[_][_]
  nodes_tmp = object.get(x_3,"http://www.w3.org/ns/shacl#name",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_4 = nodes_tmp2[_]
  nodes = x_4
}

# Constraint rules

violation contains matches if {
  x = target_class("http://a.ml/vocabularies/apiContract#EndPoint")[_]
  #  querying path: apiContract.expects / (apiContract.parameter / shapes.schema) | (apiContract.payload / shapes.schema) / shacl.name
  gen_propValues_8 = gen_path_set_rule_9 with data.sourceNode as x
  not count(gen_propValues_8) >= 1
  _result_0 := traceRuleExecution("minCount","http://a.ml/vocabularies/apiContract#expects / ((http://a.ml/vocabularies/apiContract#parameter / http://a.ml/vocabularies/shapes#schema) | (http://a.ml/vocabularies/apiContract#payload / http://a.ml/vocabularies/shapes#schema)) / http://www.w3.org/ns/shacl#name",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"condition":">=","actual": count(gen_propValues_8),"expected": 1})
  message := "orPath test"
  matches := error("validation2",x, message ,[_result_0])
}