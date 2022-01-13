package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
	"testing"
)

func TestGeneratePathPropertyArray(t *testing.T) {
	testGeneratePropertyArray("core.name",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"core:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}`, t)
}

func TestGenerateAndPropertyArray(t *testing.T) {
	testGeneratePropertyArray("x.a / x.b / x.c",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["x:b"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"x:c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`, t)
}

func TestGenerateAndOrPropertyArray(t *testing.T) {
	testGeneratePropertyArray("x.a / x.b | x.c / x.d",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["x:b"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["x:c"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`, t)
}

func TestGenerateAndOrParenthesisPropertyArray(t *testing.T) {
	testGeneratePropertyArray("( x.a / x.b ) | ( x.c / x.d )",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:a"]
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"x:b",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2 = nodes_tmp2[_]
  nodes = x_2
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:c"]
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2 = nodes_tmp2[_]
  nodes = x_2
}`, t)
}

func TestGenerateInversePathPropertyArray(t *testing.T) {
	testGeneratePropertyArray("core.name^",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "core:name" with data.object as init_x_0
  nodes = x_0
}`, t)
}

func TestGenerateAndInversePathPropertyArray(t *testing.T) {
	testGeneratePropertyArray("x.a / x.b^ / x.c",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["x:a"]
  x_0 = tmp_x_0[_][_]
  search_subjects[x_2] with data.predicate as "x:b" with data.object as x_0
  nodes_tmp = object.get(x_2,"x:c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`, t)
}

func TestGenerateOrInversePathPropertyArray(t *testing.T) {
	testGeneratePropertyArray("x.a^ | x.b^",
`gen_path_rule_1[nodes] {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "x:a" with data.object as init_x_0
  nodes = x_0
} {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "x:b" with data.object as init_x_0
  nodes = x_0
}`, t)
}

func testGeneratePropertyArray(pathString string, expected string, t *testing.T) {
	profile.GenReset()
	p, err := path.ParsePath(pathString)
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}
	result := GeneratePropertyArray(p, "x")
	actual := strings.Join(result.rego, "\n")
	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}
