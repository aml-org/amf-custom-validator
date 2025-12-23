package profile_test13

report["profile"] = "Test13"

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

gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#returns"])
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_set_rule_4 contains nodes if {
  init_y_0 = data.sourceNode
  nodes_tmp = object.get(init_y_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  y_0 = nodes_tmp2[_]
  nodes = y_0
}

gen_path_set_rule_5 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://a.ml/vocabularies/apiContract#returns"])
  x_0 = tmp_x_0[_][_]
  nodes = x_0
}

gen_path_set_rule_8 contains nodes if {
  init_z_0 = data.sourceNode
  nodes_tmp = object.get(init_z_0,"http://a.ml/vocabularies/apiContract#statusCode",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  z_0 = nodes_tmp2[_]
  nodes = z_0
}

# Constraint rules

violation contains matches if {
  x = target_class("http://a.ml/vocabularies/apiContract#Operation")[_]
  #  querying path: apiContract.returns
  ys = gen_path_set_rule_1 with data.sourceNode as x
  y_errorAcc0 = []
  ys_br_0 = [ ys_br_0_error|
    y = ys[_]
    #  querying path: apiContract.statusCode
    gen_y_check_3_array = gen_path_set_rule_4 with data.sourceNode as y
    gen_y_check_3_scalar = gen_y_check_3_array[_]
    gen_y_check_3 = as_string(gen_y_check_3_scalar)
    gen_inValues_2 = { "200"}
    not gen_inValues_2[gen_y_check_3]
    _result_0 := traceRuleExecution("in","http://a.ml/vocabularies/apiContract#statusCode",y,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_y_check_3,"expected": "[\"200\"]"})
    message := "error in nested nodes under http://a.ml/vocabularies/apiContract#returns"
    ys_br_0_inner_error := error("nested",y, message ,[_result_0])
    ys_br_0_error = [node_id(y),ys_br_0_inner_error]
  ]
  ys_br_0_errors = { nodeId | n = ys_br_0[_]; nodeId = n[0] }
  ys_br_0_errors_errors = [ node | n = ys_br_0[_]; node = n[1] ]
  y_errorAcc1 = array.concat(y_errorAcc0,ys_br_0_errors_errors)
  y_errorAcc = y_errorAcc1
  # let's accumulate results
  ys_error_node_variables_agg = ys_br_0_errors
  not count(ys) - count(ys_error_node_variables_agg) >= 1
  _result_0 := traceRuleExecution("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(ys_error_node_variables_agg), "successfulNodes":(count(ys)-count(ys_error_node_variables_agg)), "cardinality":1, "subResult": y_errorAcc})
  #  querying path: apiContract.returns
  zs = gen_path_set_rule_5 with data.sourceNode as x
  z_errorAcc0 = []
  zs_br_0 = [ zs_br_0_error|
    z = zs[_]
    #  querying path: apiContract.statusCode
    gen_z_check_7_array = gen_path_set_rule_8 with data.sourceNode as z
    gen_z_check_7_scalar = gen_z_check_7_array[_]
    gen_z_check_7 = as_string(gen_z_check_7_scalar)
    gen_inValues_6 = { "429"}
    not gen_inValues_6[gen_z_check_7]
    _result_0 := traceRuleExecution("in","http://a.ml/vocabularies/apiContract#statusCode",z,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false,"actual": gen_z_check_7,"expected": "[\"429\"]"})
    message := "error in nested nodes under http://a.ml/vocabularies/apiContract#returns"
    zs_br_0_inner_error := error("nested",z, message ,[_result_0])
    zs_br_0_error = [node_id(z),zs_br_0_inner_error]
  ]
  zs_br_0_errors = { nodeId | n = zs_br_0[_]; nodeId = n[0] }
  zs_br_0_errors_errors = [ node | n = zs_br_0[_]; node = n[1] ]
  z_errorAcc1 = array.concat(z_errorAcc0,zs_br_0_errors_errors)
  z_errorAcc = z_errorAcc1
  # let's accumulate results
  zs_error_node_variables_agg = zs_br_0_errors
  not count(zs) - count(zs_error_node_variables_agg) >= 1
  _result_1 := traceRuleExecution("atLeast","http://a.ml/vocabularies/apiContract#returns",x,{"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated":false, "failedNodes":count(zs_error_node_variables_agg), "successfulNodes":(count(zs)-count(zs_error_node_variables_agg)), "cardinality":1, "subResult": z_errorAcc})
  message := "Notify the client when the limit is exceeded by providing the limit number and the time at which the limit will\nbe reset.\n"
  matches := error("lack-of-resources-and-rate-limiting-too-many-requests",x, message ,[_result_0,_result_1])
}