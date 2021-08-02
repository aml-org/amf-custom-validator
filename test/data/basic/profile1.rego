package test1


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

as_string(x) = x["@id"] {
  is_object(x)
  x["@id"]
}

as_string(x) = json.marshal(x) {
  is_object(x)
  not x["@id"]
}

# Traces one evaluation of a constraint
trace(constraint, path, node, value) = t {
  id := node["@id"]
  raw_lexical := input["@lexical"][id]
  lexical_parts := regex.find_n("\\d+", raw_lexical, 4)
  t := {
    "component": constraint,
    "path": path,
    "focusNode": id,
    "value": value,
	"lexical": {
      "start": {
        "line": lexical_parts[0],
        "column": lexical_parts[1]
      },
      "end": {
        "line": lexical_parts[2],
        "column": lexical_parts[3]
      }
    }
  }
}

trace(constraint, path, node, value) = t {
  id := node["@id"]
  not input["@lexical"][id]
  t := {
    "component": constraint,
    "path": path,
    "focusNode": id,
    "value": value
  }
}

# Builds an error message that can be returned to the calling client software
error(shapeId, target, message, traceLog) = e {
  id := target["@id"]
  e := {
    "shapeId": shapeId,
    "target": id,
    "message": message,
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
  init_x_0__in_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__in_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_5[nodes] {
  init_x_0__maxCount_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__maxCount_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_7[nodes] {
  init_x_0__minCount_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__minCount_,"apiContract:method",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

gen_path_rule_8[nodes] {
  init_x_0__pattern_ = data.sourceNode
  nodes_tmp = object.get(init_x_0__pattern_,"shacl:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}

# Constraint rules

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  gen_x_check_2_array = gen_path_rule_3 with data.sourceNode as x
  gen_x_check_2_scalar = gen_x_check_2_array[_]
  gen_x_check_2 = as_string(gen_x_check_2_scalar)
  gen_inValues_1 = { "publish","subscribe","1","2"}
  not gen_inValues_1[gen_x_check_2]
  _result_0 := trace("in","apiContract.method",x,{"negated":false,"actual": gen_x_check_2,"expected": "[\"publish\",\"subscribe\",\"1\",\"2\"]"})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: shacl.name
  gen_propValues_4 = gen_path_rule_5 with data.sourceNode as x
  not count(gen_propValues_4) <= 1
  _result_0 := trace("maxCount","shacl.name",x,{"negated":false,"condition":"<=","actual": count(gen_propValues_4),"expected": 1})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: apiContract.method
  gen_propValues_6 = gen_path_rule_7 with data.sourceNode as x
  not count(gen_propValues_6) >= 1
  _result_0 := trace("minCount","apiContract.method",x,{"negated":false,"condition":">=","actual": count(gen_propValues_6),"expected": 1})
  matches := error("validation1",x,"This is the message",[_result_0])
}

violation[matches] {
  target_class[x] with data.class as "apiContract:Operation"
  #  querying path: shacl.name
  gen_gen_path_rule_8_node_9_array = gen_path_rule_8 with data.sourceNode as x
  gen_gen_path_rule_8_node_9 = gen_gen_path_rule_8_node_9_array[_]
  not regex.match("^put|post$",gen_gen_path_rule_8_node_9)
  _result_0 := trace("pattern","shacl.name",x,{"negated":false,"argument": gen_gen_path_rule_8_node_9})
  matches := error("validation1",x,"This is the message",[_result_0])
}