package generator

import (
	"github.com/aml-org/amfopa/internal/parser/path"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"strings"
	"testing"
)

func TestGeneratePathPropertyArray(t *testing.T) {
	statements.GenReset()
	p, err := path.ParsePath("core.name")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  x = data.sourceNode
  nodes_tmp = object.get(x,"core:name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}

func TestGenerateAndPropertyArray(t *testing.T) {
	statements.GenReset()
	p, err := path.ParsePath("x.a / x.b / x.c")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  x = data.sourceNode
  tmp_x_0_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x["x:a"]
  x_0_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_0_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test[][]
  tmp_x_1_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x_0_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test["x:b"]
  x_1_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_1_782e61202f20782e62202f20782e63d41d8cd98f00b204e9800998ecf8427e_test[][]
  nodes_tmp = object.get(x,"x:c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}

func TestGenerateOrPropertyArray(t *testing.T) {
	statements.GenReset()
	p, err := path.ParsePath("x.a / x.b | x.c / x.d")
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}

	result := GeneratePropertyArray(p, "x", "test")

	expected := `gen_path_rule_1[nodes] {
  x = data.sourceNode
  tmp_x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x["x:a"]
  x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test[][]
  tmp_x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test["x:b"]
  x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test[][]
  nodes_tmp = object.get(x,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
} {
  x = data.sourceNode
  tmp_x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x["x:a"]
  x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test[][]
  tmp_x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = nested_nodes with data.nodes as x_0_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test["x:c"]
  x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test = tmp_x_1_782e61202f20782e62207c20782e63202f20782e64d41d8cd98f00b204e9800998ecf8427e_test[][]
  nodes_tmp = object.get(x,"x:d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  nodes = nodes_tmp2[_]
}`
	actual := strings.Join(result.rego, "\n")

	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}
