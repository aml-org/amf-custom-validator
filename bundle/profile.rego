package profile_anypoint_best_practices

import data.profile_anypoint_best_practices.report

report["profile"] = "Anypoint Best Practices"

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
	node_predicate_values = nodes_array with data.nodes as object.get(node, predicate, [])
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
	rs = object.get(n, data.property, [])
	rss = nodes_array with data.nodes as rs
	rsss = nested with data.nodes as rss
	r = rsss[_]
}

# collect terminal values
collect_values[r] {
	nodes = data.nodes
	n = nodes[_]
	rs = object.get(n, data.property, [])
	rss = nodes_array with data.nodes as rs
	r = rss[_]
}

# helper to check datatype constraints

check_datatype(x, dt) {
	dt == "http://www.w3.org/2001/XMLSchema#string"
	is_string(x)
}

check_datatype(x, dt) {
	dt == "http://www.w3.org/2001/XMLSchema#integer"
	is_number(x)
}

check_datatype(x, dt) {
	dt == "http://www.w3.org/2001/XMLSchema#float"
	is_number(x)
}

check_datatype(x, dt) {
	dt == "http://www.w3.org/2001/XMLSchema#boolean"
	is_boolean(x)
}

check_datatype(x, dt) {
	is_object(x)
	t = object.get(x, "@type", "")
	t == dt
}

check_datatype(x, dt) = false {
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
		"location": l,
	}
}

# Generates trace when lexical info is not available
trace(constraint, resultPath, focusNode, traceValue) = t {
	not location(focusNode)
	t := {
		"@type": ["reportSchema:TraceMessageNode", "validation:TraceMessage"],
		"component": constraint,
		"resultPath": resultPath,
		"traceValue": traceValue,
	}
}

location(focusNode) = l {
	id := focusNode["@id"]
	location := input["@lexical"][id]
	raw_range := location.range
	uri := location.uri
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
		"trace": traceLog,
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
		"trace": traceLog,
	}
}

# generate the report for violation level
# default value must be added dynamically

# generate the report for the info level
# default value must be added dynamically

# generate the report for the info level
# default value must be added dynamically

default warning = []

default info = []

# Path rules

gen_path_set_rule_1[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#path", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#EndPoint"

	#  querying path: apiContract.path
	gen_gen_path_set_rule_1_node_2_array = gen_path_set_rule_1 with data.sourceNode as x
	gen_gen_path_set_rule_1_node_2 = gen_gen_path_set_rule_1_node_2_array[_]
	not regex.match(`^[a-z\\/\\{\\}-]+$`, gen_gen_path_set_rule_1_node_2)
	_result_0 := trace("pattern", "http://a.ml/vocabularies/apiContract#path", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "expected": "^[a-z\\\\/\\\\{\\\\}-]+$", "actual": gen_gen_path_set_rule_1_node_2})
	message := "Naming conventions for resources:\n- Use lower case (example: /accounts)\n- For resources with more than 2 words\n    - use lowercase for both words (example: /lineitems) or\n    - use kebab-case (aka spinal-case) (example: /line-items)\n"
	matches := error("resource-use-lowercase", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_3[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#header"]
	x_0 = tmp_x_0[_][_]
	nodes = x_0
}

gen_path_set_rule_6[nodes] {
	init_y_0 = data.sourceNode
	nodes_tmp = object.get(init_y_0, "http://a.ml/vocabularies/core#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	y_0 = nodes_tmp2[_]
	nodes = y_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Response"

	#  querying path: apiContract.header
	ys = gen_path_set_rule_3 with data.sourceNode as x
	y_errorAcc0 = []
	ys_br_0 = [ys_br_0_error |
		y = ys[_]

		#  querying path: core.name
		gen_y_check_5_array = gen_path_set_rule_6 with data.sourceNode as y
		gen_y_check_5_scalar = gen_y_check_5_array[_]
		gen_y_check_5 = as_string(gen_y_check_5_scalar)
		gen_inValues_4 = {"content-type"}
		not gen_inValues_4[gen_y_check_5]
		_result_0 := trace("in", "http://a.ml/vocabularies/core#name", y, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_y_check_5, "expected": "[\"content-type\"]"})
		message := "error in nested nodes under http://a.ml/vocabularies/apiContract#header"
		ys_br_0_inner_error := error("nested", y, message, [_result_0])
		ys_br_0_error = [y["@id"], ys_br_0_inner_error]
	]
	ys_br_0_errors = {nodeId | n = ys_br_0[_]; nodeId = n[0]}
	ys_br_0_errors_errors = [node | n = ys_br_0[_]; node = n[1]]
	y_errorAcc1 = array.concat(y_errorAcc0, ys_br_0_errors_errors)
	y_errorAcc = y_errorAcc1

	# let's accumulate results
	ys_error_node_variables_agg = ys_br_0_errors
	not count(ys) - count(ys_error_node_variables_agg) >= 1
	_result_0 := trace("atLeast", "http://a.ml/vocabularies/apiContract#header", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "failedNodes": count(ys_error_node_variables_agg), "successfulNodes": count(ys) - count(ys_error_node_variables_agg), "cardinality": 1, "subResult": y_errorAcc})
	message := "- For the response: use ‘content-type’ header\n"
	matches := error("media-type-headers-response", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_7[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#urlTemplate", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_9[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#urlTemplate", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Server"

	#  querying path: core.urlTemplate
	gen_gen_path_set_rule_7_node_8_array = gen_path_set_rule_7 with data.sourceNode as x
	gen_gen_path_set_rule_7_node_8 = gen_gen_path_set_rule_7_node_8_array[_]
	not regex.match(`/api/[0-9].[0.9]`, gen_gen_path_set_rule_7_node_8)
	_result_0 := trace("pattern", "http://a.ml/vocabularies/core#urlTemplate", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "expected": "/api/[0-9].[0.9]", "actual": gen_gen_path_set_rule_7_node_8})

	#  querying path: core.urlTemplate
	gen_gen_path_set_rule_9_node_10_array = gen_path_set_rule_9 with data.sourceNode as x
	gen_gen_path_set_rule_9_node_10 = gen_gen_path_set_rule_9_node_10_array[_]
	not regex.match(`/api/v[0-9]+`, gen_gen_path_set_rule_9_node_10)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/core#urlTemplate", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "expected": "/api/v[0-9]+", "actual": gen_gen_path_set_rule_9_node_10})
	message := "Include the “api” word and the version of the API in the base Url (e.g. domain/api/v1)"
	matches := error("base-url-pattern-server", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_11[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://www.w3.org/ns/shacl#PropertyShape"

	#  querying path: shacl.name
	gen_gen_path_set_rule_11_node_12_array = gen_path_set_rule_11 with data.sourceNode as x
	gen_gen_path_set_rule_11_node_12 = gen_gen_path_set_rule_11_node_12_array[_]
	not regex.match(`^[a-z]+([A-Z][a-z]+)*$`, gen_gen_path_set_rule_11_node_12)
	_result_0 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "expected": "^[a-z]+([A-Z][a-z]+)*$", "actual": gen_gen_path_set_rule_11_node_12})
	message := "Use camelCase for all the names (fields), preferably don’t use underscores."
	matches := error("camel-case-fields", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_15[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_18[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_17[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_18[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_15[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_20[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_17[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_20[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_x_check_14_array = gen_path_set_rule_15 with data.sourceNode as x
	gen_x_check_14_scalar = gen_x_check_14_array[_]
	gen_x_check_14 = as_string(gen_x_check_14_scalar)
	gen_inValues_13 = {"http://www.w3.org/2001/XMLSchema#dateTime"}
	not gen_inValues_13[gen_x_check_14]
	_result_0 := trace("in", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_14, "expected": "[\"http://www.w3.org/2001/XMLSchema#dateTime\"]"})

	#  querying path: shacl.name
	gen_gen_path_set_rule_18_node_19_array = gen_path_set_rule_18 with data.sourceNode as x
	gen_gen_path_set_rule_18_node_19 = gen_gen_path_set_rule_18_node_19_array[_]
	regex.match(`createdAt`, gen_gen_path_set_rule_18_node_19)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "createdAt", "actual": gen_gen_path_set_rule_18_node_19})
	message := "Use standard date formats: ISO8601\nUse UTC\n  2016-10-27T13:42:21+00:00    (+00:00 is the time zones hour offset)\n  2016-10-27T13:42:21Z     (Z is place holder for local time zone)\n"
	matches := error("date-time-representation", x, message, [_result_0, _result_1])
}

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_propValues_16 = gen_path_set_rule_17 with data.sourceNode as x
	not count(gen_propValues_16) >= 1
	_result_0 := trace("minCount", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_16), "expected": 1})

	#  querying path: shacl.name
	gen_gen_path_set_rule_18_node_19_array = gen_path_set_rule_18 with data.sourceNode as x
	gen_gen_path_set_rule_18_node_19 = gen_gen_path_set_rule_18_node_19_array[_]
	regex.match(`createdAt`, gen_gen_path_set_rule_18_node_19)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "createdAt", "actual": gen_gen_path_set_rule_18_node_19})
	message := "Use standard date formats: ISO8601\nUse UTC\n  2016-10-27T13:42:21+00:00    (+00:00 is the time zones hour offset)\n  2016-10-27T13:42:21Z     (Z is place holder for local time zone)\n"
	matches := error("date-time-representation", x, message, [_result_0, _result_1])
}

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_x_check_14_array = gen_path_set_rule_15 with data.sourceNode as x
	gen_x_check_14_scalar = gen_x_check_14_array[_]
	gen_x_check_14 = as_string(gen_x_check_14_scalar)
	gen_inValues_13 = {"http://www.w3.org/2001/XMLSchema#dateTime"}
	not gen_inValues_13[gen_x_check_14]
	_result_0 := trace("in", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_14, "expected": "[\"http://www.w3.org/2001/XMLSchema#dateTime\"]"})

	#  querying path: shacl.name
	gen_gen_path_set_rule_20_node_21_array = gen_path_set_rule_20 with data.sourceNode as x
	gen_gen_path_set_rule_20_node_21 = gen_gen_path_set_rule_20_node_21_array[_]
	regex.match(`updatedAt`, gen_gen_path_set_rule_20_node_21)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "updatedAt", "actual": gen_gen_path_set_rule_20_node_21})
	message := "Use standard date formats: ISO8601\nUse UTC\n  2016-10-27T13:42:21+00:00    (+00:00 is the time zones hour offset)\n  2016-10-27T13:42:21Z     (Z is place holder for local time zone)\n"
	matches := error("date-time-representation", x, message, [_result_0, _result_1])
}

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_propValues_16 = gen_path_set_rule_17 with data.sourceNode as x
	not count(gen_propValues_16) >= 1
	_result_0 := trace("minCount", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_16), "expected": 1})

	#  querying path: shacl.name
	gen_gen_path_set_rule_20_node_21_array = gen_path_set_rule_20 with data.sourceNode as x
	gen_gen_path_set_rule_20_node_21 = gen_gen_path_set_rule_20_node_21_array[_]
	regex.match(`updatedAt`, gen_gen_path_set_rule_20_node_21)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "updatedAt", "actual": gen_gen_path_set_rule_20_node_21})
	message := "Use standard date formats: ISO8601\nUse UTC\n  2016-10-27T13:42:21+00:00    (+00:00 is the time zones hour offset)\n  2016-10-27T13:42:21Z     (Z is place holder for local time zone)\n"
	matches := error("date-time-representation", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_24[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_27[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_26[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#datatype", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_27[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://www.w3.org/ns/shacl#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_x_check_23_array = gen_path_set_rule_24 with data.sourceNode as x
	gen_x_check_23_scalar = gen_x_check_23_array[_]
	gen_x_check_23 = as_string(gen_x_check_23_scalar)
	gen_inValues_22 = {"http://www.w3.org/2001/XMLSchema#date"}
	not gen_inValues_22[gen_x_check_23]
	_result_0 := trace("in", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_23, "expected": "[\"http://www.w3.org/2001/XMLSchema#date\"]"})

	#  querying path: shacl.name
	gen_gen_path_set_rule_27_node_28_array = gen_path_set_rule_27 with data.sourceNode as x
	gen_gen_path_set_rule_27_node_28 = gen_gen_path_set_rule_27_node_28_array[_]
	regex.match(`^.*[dD]ate.*$`, gen_gen_path_set_rule_27_node_28)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "^.*[dD]ate.*$", "actual": gen_gen_path_set_rule_27_node_28})
	message := "Use standard date formats: ISO8601\nUse calendar date:\n  2016-10-27\n"
	matches := error("date-only-representation", x, message, [_result_0, _result_1])
}

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#ScalarShape"

	#  querying path: shacl.datatype
	gen_propValues_25 = gen_path_set_rule_26 with data.sourceNode as x
	not count(gen_propValues_25) >= 1
	_result_0 := trace("minCount", "http://www.w3.org/ns/shacl#datatype", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_25), "expected": 1})

	#  querying path: shacl.name
	gen_gen_path_set_rule_27_node_28_array = gen_path_set_rule_27 with data.sourceNode as x
	gen_gen_path_set_rule_27_node_28 = gen_gen_path_set_rule_27_node_28_array[_]
	regex.match(`^.*[dD]ate.*$`, gen_gen_path_set_rule_27_node_28)
	_result_1 := trace("pattern", "http://www.w3.org/ns/shacl#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "^.*[dD]ate.*$", "actual": gen_gen_path_set_rule_27_node_28})
	message := "Use standard date formats: ISO8601\nUse calendar date:\n  2016-10-27\n"
	matches := error("date-only-representation", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_30[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#payload"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/shapes#schema", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Response"

	#  querying path: apiContract.payload / shapes.schema
	gen_propValues_29 = gen_path_set_rule_30 with data.sourceNode as x
	not count(gen_propValues_29) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/apiContract#payload / http://a.ml/vocabularies/shapes#schema", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_29), "expected": 1})
	message := "Use schemas or data types in the specification to determine the format of the responses.\n"
	matches := error("use-schemas-responses", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_31[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#payload"]
	x_0 = tmp_x_0[_][_]
	nodes = x_0
}

gen_path_set_rule_33[nodes] {
	init_y_0 = data.sourceNode
	nodes_tmp = object.get(init_y_0, "http://a.ml/vocabularies/shapes#schema", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	y_0 = nodes_tmp2[_]
	nodes = y_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Request"

	#  querying path: apiContract.payload
	ys = gen_path_set_rule_31 with data.sourceNode as x
	y_errorAcc0 = []
	ys_br_0 = [ys_br_0_error |
		y = ys[_]

		#  querying path: shapes.schema
		gen_propValues_32 = gen_path_set_rule_33 with data.sourceNode as y
		not count(gen_propValues_32) >= 1
		_result_0 := trace("minCount", "http://a.ml/vocabularies/shapes#schema", y, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_32), "expected": 1})
		message := "error in nested nodes under http://a.ml/vocabularies/apiContract#payload"
		ys_br_0_inner_error := error("nested", y, message, [_result_0])
		ys_br_0_error = [y["@id"], ys_br_0_inner_error]
	]
	ys_br_0_errors = {nodeId | n = ys_br_0[_]; nodeId = n[0]}
	ys_br_0_errors_errors = [node | n = ys_br_0[_]; node = n[1]]
	y_errorAcc1 = array.concat(y_errorAcc0, ys_br_0_errors_errors)
	y_errorAcc = y_errorAcc1

	# let's accumulate results
	ys_error_node_variables_agg = ys_br_0_errors
	count(ys_error_node_variables_agg) > 0
	_result_0 := trace("nested", "http://a.ml/vocabularies/apiContract#payload", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "failedNodes": count(ys_error_node_variables_agg), "successfulNodes": count(ys) - count(ys_error_node_variables_agg), "subResult": y_errorAcc})
	message := "Use schemas or data types in the specification to determine the format of the requests.\n"
	matches := error("use-schemas-requests", x, message, [_result_0])
}

# Path rules

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Payload"
	schema = find with data.link as x["http://a.ml/vocabularies/shapes#schema"]

	nested_nodes[examples] with data.nodes as object.get(schema, "http://a.ml/vocabularies/apiContract#examples", [])

	examples_from_this_payload = {element |
		example = examples[_]
		sourcemap = find with data.link as object.get(example, "http://a.ml/vocabularies/document-source-maps#sources", [])
		tracked_element = find with data.link as object.get(sourcemap, "http://a.ml/vocabularies/document-source-maps#tracked-element", [])
		values_contains(tracked_element["http://a.ml/vocabularies/document-source-maps#value"], x["@id"])
		element := example
	}

	gen_rego_result_35 := count(examples_from_this_payload) > 0

	gen_rego_result_35 != true
	_result_0 := trace("rego", "", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false})
	message := "Always include examples in requests and responses"
	matches := error("provide-examples", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_37[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#WebAPI"

	#  querying path: core.name
	gen_propValues_36 = gen_path_set_rule_37 with data.sourceNode as x
	not count(gen_propValues_36) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_36), "expected": 1})
	message := "Provide the title for the API"
	matches := error("api-must-have-title", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_39[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#name", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: core.name
	gen_propValues_38 = gen_path_set_rule_39 with data.sourceNode as x
	not count(gen_propValues_38) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#name", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_38), "expected": 1})
	message := "Provide identifiers for the operations"
	matches := error("operations-must-have-identifiers", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_41[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#WebAPI"

	#  querying path: core.description
	gen_propValues_40 = gen_path_set_rule_41 with data.sourceNode as x
	not count(gen_propValues_40) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_40), "expected": 1})
	message := "Provide the description for the API"
	matches := error("api-must-have-description", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_43[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#documentation", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#WebAPI"

	#  querying path: core.documentation
	gen_propValues_42 = gen_path_set_rule_43 with data.sourceNode as x
	not count(gen_propValues_42) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#documentation", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_42), "expected": 1})
	message := "Provide the documentation for the API"
	matches := error("api-must-have-documentation", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_45[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: core.description
	gen_propValues_44 = gen_path_set_rule_45 with data.sourceNode as x
	not count(gen_propValues_44) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_44), "expected": 1})
	message := "Provide descriptions for the operations"
	matches := error("operations-must-have-descriptions", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_47[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Response"

	#  querying path: core.description
	gen_propValues_46 = gen_path_set_rule_47 with data.sourceNode as x
	not count(gen_propValues_46) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_46), "expected": 1})
	message := "Provide descriptions for the responses"
	matches := error("responses-must-have-descriptions", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_49[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_50[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#binding", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Parameter"

	#  querying path: core.description
	gen_propValues_48 = gen_path_set_rule_49 with data.sourceNode as x
	not count(gen_propValues_48) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_48), "expected": 1})

	#  querying path: apiContract.binding
	gen_gen_path_set_rule_50_node_51_array = gen_path_set_rule_50 with data.sourceNode as x
	gen_gen_path_set_rule_50_node_51 = gen_gen_path_set_rule_50_node_51_array[_]
	regex.match(`header`, gen_gen_path_set_rule_50_node_51)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#binding", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "header", "actual": gen_gen_path_set_rule_50_node_51})
	message := "Provide descriptions for request headers"
	matches := error("headers-must-have-descriptions", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_53[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

gen_path_set_rule_54[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#binding", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Parameter"

	#  querying path: core.description
	gen_propValues_52 = gen_path_set_rule_53 with data.sourceNode as x
	not count(gen_propValues_52) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_52), "expected": 1})

	#  querying path: apiContract.binding
	gen_gen_path_set_rule_54_node_55_array = gen_path_set_rule_54 with data.sourceNode as x
	gen_gen_path_set_rule_54_node_55 = gen_gen_path_set_rule_54_node_55_array[_]
	regex.match(`query`, gen_gen_path_set_rule_54_node_55)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#binding", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "query", "actual": gen_gen_path_set_rule_54_node_55})
	message := "Provide descriptions for query params"
	matches := error("query-params-must-have-descriptions", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_57[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/shapes#range"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/shapes#PropertyShape"

	#  querying path: shapes.range / core.description
	gen_propValues_56 = gen_path_set_rule_57 with data.sourceNode as x
	not count(gen_propValues_56) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/shapes#range / http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_56), "expected": 1})
	message := "Provide descriptions for data shapes"
	matches := error("property-shape-ranges-must-have-descriptions", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_59[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#description", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://www.w3.org/ns/shacl#NodeShape"

	#  querying path: core.description
	gen_propValues_58 = gen_path_set_rule_59 with data.sourceNode as x
	not count(gen_propValues_58) >= 1
	_result_0 := trace("minCount", "http://a.ml/vocabularies/core#description", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "condition": ">=", "actual": count(gen_propValues_58), "expected": 1})
	message := "Provide description for data shapes"
	matches := error("payload-shapes-must-have-descriptions", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_62[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/core#mediaType", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Payload"

	#  querying path: core.mediaType
	gen_x_check_61_array = gen_path_set_rule_62 with data.sourceNode as x
	gen_x_check_61_scalar = gen_x_check_61_array[_]
	gen_x_check_61 = as_string(gen_x_check_61_scalar)
	gen_inValues_60 = {"application/json", "application/xml"}
	not gen_inValues_60[gen_x_check_61]
	_result_0 := trace("in", "http://a.ml/vocabularies/core#mediaType", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_61, "expected": "[\"application/json\",\"application/xml\"]"})
	message := "If there is no standard media type and format, try to use as much as possible extensible formats such as JSON\n(application/json) and XML (application/xml), preferably JSON.\n"
	matches := error("preferred-media-type-representations", x, message, [_result_0])
}

# Path rules

gen_path_set_rule_65[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/apiContract#statusCode", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

gen_path_set_rule_66[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#method", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: apiContract.returns / apiContract.statusCode
	gen_x_check_64_array = gen_path_set_rule_65 with data.sourceNode as x
	gen_x_check_64_scalar = gen_x_check_64_array[_]
	gen_x_check_64 = as_string(gen_x_check_64_scalar)
	gen_inValues_63 = {"200", "204", "304", "400", "401", "403", "404", "405", "406", "408", "410", "412", "415", "429", "500", "502", "503", "504", "509", "510", "511", "550", "598", "599"}
	not gen_inValues_63[gen_x_check_64]
	_result_0 := trace("in", "http://a.ml/vocabularies/apiContract#returns / http://a.ml/vocabularies/apiContract#statusCode", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_64, "expected": "[\"200\",\"204\",\"304\",\"400\",\"401\",\"403\",\"404\",\"405\",\"406\",\"408\",\"410\",\"412\",\"415\",\"429\",\"500\",\"502\",\"503\",\"504\",\"509\",\"510\",\"511\",\"550\",\"598\",\"599\"]"})

	#  querying path: apiContract.method
	gen_gen_path_set_rule_66_node_67_array = gen_path_set_rule_66 with data.sourceNode as x
	gen_gen_path_set_rule_66_node_67 = gen_gen_path_set_rule_66_node_67_array[_]
	regex.match(`get`, gen_gen_path_set_rule_66_node_67)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#method", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "get", "actual": gen_gen_path_set_rule_66_node_67})
	message := "These response codes [200,204,304,400,401,403,404,405,406,408,410,412,415,429,500,502,503,504,509,510,511,550,598,\n599]should be used as standard for GET operations, the use of not defined return codes is discouraged and should\nonly be done in exceptional circumstances.\n"
	matches := error("standard-get-status-codes", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_70[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/apiContract#statusCode", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

gen_path_set_rule_71[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#method", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: apiContract.returns / apiContract.statusCode
	gen_x_check_69_array = gen_path_set_rule_70 with data.sourceNode as x
	gen_x_check_69_scalar = gen_x_check_69_array[_]
	gen_x_check_69 = as_string(gen_x_check_69_scalar)
	gen_inValues_68 = {"201", "202", "400", "401", "403", "404", "405", "406", "408", "409", "410", "415", "429", "500", "502", "503", "504", "509", "510", "511", "550", "598", "599"}
	not gen_inValues_68[gen_x_check_69]
	_result_0 := trace("in", "http://a.ml/vocabularies/apiContract#returns / http://a.ml/vocabularies/apiContract#statusCode", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_69, "expected": "[\"201\",\"202\",\"400\",\"401\",\"403\",\"404\",\"405\",\"406\",\"408\",\"409\",\"410\",\"415\",\"429\",\"500\",\"502\",\"503\",\"504\",\"509\",\"510\",\"511\",\"550\",\"598\",\"599\"]"})

	#  querying path: apiContract.method
	gen_gen_path_set_rule_71_node_72_array = gen_path_set_rule_71 with data.sourceNode as x
	gen_gen_path_set_rule_71_node_72 = gen_gen_path_set_rule_71_node_72_array[_]
	regex.match(`post`, gen_gen_path_set_rule_71_node_72)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#method", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "post", "actual": gen_gen_path_set_rule_71_node_72})
	message := "These response codes [201,202,400,401,403,404,405,406,408,409,410,415,429,500,502,503,504,509,510,511,550,598,599]\nshould be used as standard for POST operations, the use of not defined return codes is discouraged and should\nonly be done in exceptional circumstances.\n"
	matches := error("standard-post-status-codes", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_75[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/apiContract#statusCode", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

gen_path_set_rule_76[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#method", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: apiContract.returns / apiContract.statusCode
	gen_x_check_74_array = gen_path_set_rule_75 with data.sourceNode as x
	gen_x_check_74_scalar = gen_x_check_74_array[_]
	gen_x_check_74 = as_string(gen_x_check_74_scalar)
	gen_inValues_73 = {"200", "202", "204", "400", "401", "403", "404", "405", "406", "408", "409", "410", "412", "415", "429", "500", "502", "503", "504", "509", "510", "511", "550", "598", "599"}
	not gen_inValues_73[gen_x_check_74]
	_result_0 := trace("in", "http://a.ml/vocabularies/apiContract#returns / http://a.ml/vocabularies/apiContract#statusCode", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_74, "expected": "[\"200\",\"202\",\"204\",\"400\",\"401\",\"403\",\"404\",\"405\",\"406\",\"408\",\"409\",\"410\",\"412\",\"415\",\"429\",\"500\",\"502\",\"503\",\"504\",\"509\",\"510\",\"511\",\"550\",\"598\",\"599\"]"})

	#  querying path: apiContract.method
	gen_gen_path_set_rule_76_node_77_array = gen_path_set_rule_76 with data.sourceNode as x
	gen_gen_path_set_rule_76_node_77 = gen_gen_path_set_rule_76_node_77_array[_]
	regex.match(`put`, gen_gen_path_set_rule_76_node_77)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#method", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "put", "actual": gen_gen_path_set_rule_76_node_77})
	message := "These response codes [200,202,204,400,401,403,404,405,406,408,409,410,412,415,429,500,502,503,504,509,510,511,550,\n598,599] should be used as standard for PUT operations, the use of not defined return codes is discouraged and\nshould only be done in exceptional circumstances.\n"
	matches := error("standard-put-status-codes", x, message, [_result_0, _result_1])
}

# Path rules

gen_path_set_rule_80[nodes] {
	init_x_0 = data.sourceNode
	tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://a.ml/vocabularies/apiContract#returns"]
	x_0 = tmp_x_0[_][_]
	nodes_tmp = object.get(x_0, "http://a.ml/vocabularies/apiContract#statusCode", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_2 = nodes_tmp2[_]
	nodes = x_2
}

gen_path_set_rule_81[nodes] {
	init_x_0 = data.sourceNode
	nodes_tmp = object.get(init_x_0, "http://a.ml/vocabularies/apiContract#method", [])
	nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
	x_0 = nodes_tmp2[_]
	nodes = x_0
}

# Constraint rules

violation[matches] {
	target_class[x] with data.class as "http://a.ml/vocabularies/apiContract#Operation"

	#  querying path: apiContract.returns / apiContract.statusCode
	gen_x_check_79_array = gen_path_set_rule_80 with data.sourceNode as x
	gen_x_check_79_scalar = gen_x_check_79_array[_]
	gen_x_check_79 = as_string(gen_x_check_79_scalar)
	gen_inValues_78 = {"200", "202", "204", "400", "401", "403", "404", "405", "406", "408", "409", "410", "429", "500", "502", "503", "504", "509", "510", "511", "550", "598", "599"}
	not gen_inValues_78[gen_x_check_79]
	_result_0 := trace("in", "http://a.ml/vocabularies/apiContract#returns / http://a.ml/vocabularies/apiContract#statusCode", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": false, "actual": gen_x_check_79, "expected": "[\"200\",\"202\",\"204\",\"400\",\"401\",\"403\",\"404\",\"405\",\"406\",\"408\",\"409\",\"410\",\"429\",\"500\",\"502\",\"503\",\"504\",\"509\",\"510\",\"511\",\"550\",\"598\",\"599\"]"})

	#  querying path: apiContract.method
	gen_gen_path_set_rule_81_node_82_array = gen_path_set_rule_81 with data.sourceNode as x
	gen_gen_path_set_rule_81_node_82 = gen_gen_path_set_rule_81_node_82_array[_]
	regex.match(`delete`, gen_gen_path_set_rule_81_node_82)
	_result_1 := trace("pattern", "http://a.ml/vocabularies/apiContract#method", x, {"@type": ["reportSchema:TraceValueNode", "validation:TraceValue"], "negated": true, "expected": "delete", "actual": gen_gen_path_set_rule_81_node_82})
	message := "These response codes [200,202,204,400,401,403,404,405,406,408,409,410,429,500,502,503,504,509,510,511,550,598,599]\nshould be used as standard for DELETE operations, the use of not defined return codes is discouraged and should\nonly be done in exceptional circumstances.\n"
	matches := error("standard-delete-status-codes", x, message, [_result_0, _result_1])
}
