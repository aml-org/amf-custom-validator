package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"strings"
	"testing"
)

func TestGeneratePathPropertyArray(t *testing.T) {
	profile.GenReset()
	p, err := path.ParsePath("core.name")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  init_x_0__test = data.sourceNode
  nodes_tmp = object.get(init_x_0__test,"core:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0__test = nodes_tmp2[_]
  nodes = x_0__test
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}

func TestGenerateAndPropertyArray(t *testing.T) {
	profile.GenReset()
	p, err := path.ParsePath("x.a / x.b / x.c")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  init_x_0__test = data.sourceNode
  tmp_x_0__test = nested_nodes with data.nodes as init_x_0__test["x:a"]
  x_0__test = tmp_x_0__test[_][_]
  tmp_x_2__test = nested_nodes with data.nodes as x_0__test["x:b"]
  x_2__test = tmp_x_2__test[_][_]
  nodes_tmp = object.get(x_2__test,"x:c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3__test = nodes_tmp2[_]
  nodes = x_3__test
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}

func TestGenerateOrPropertyArray(t *testing.T) {
	profile.GenReset()
	p, err := path.ParsePath("x.a / x.b | x.c / x.d")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  init_x_0__test = data.sourceNode
  tmp_x_0__test = nested_nodes with data.nodes as init_x_0__test["x:a"]
  x_0__test = tmp_x_0__test[_][_]
  tmp_x_2__test = nested_nodes with data.nodes as x_0__test["x:b"]
  x_2__test = tmp_x_2__test[_][_]
  nodes_tmp = object.get(x_2__test,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3__test = nodes_tmp2[_]
  nodes = x_3__test
} {
  init_x_0__test = data.sourceNode
  tmp_x_0__test = nested_nodes with data.nodes as init_x_0__test["x:a"]
  x_0__test = tmp_x_0__test[_][_]
  tmp_x_2__test = nested_nodes with data.nodes as x_0__test["x:c"]
  x_2__test = tmp_x_2__test[_][_]
  nodes_tmp = object.get(x_2__test,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3__test = nodes_tmp2[_]
  nodes = x_3__test
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}
