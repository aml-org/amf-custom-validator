package generator

import (
	"github.com/aml-org/amf-custom-validator/internal/misc"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	"github.com/aml-org/amf-custom-validator/internal/parser/profile"
	"github.com/aml-org/amf-custom-validator/internal/types"
	"github.com/aml-org/amf-custom-validator/internal/validator/contexts"
	"github.com/open-policy-agent/opa/ast"

	"strings"
	"testing"
)

func TestGeneratePathPropertyArray(t *testing.T) {
	expander := misc.IriExpander{Context: contexts.DefaultAMFContext}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/core#name",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_0 = nodes_tmp2[_]
  nodes = x_0
}`
	testGeneratePropertyArray("core.name", expected, &expander, t, ast.RegoV0)

	expected = `gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  nodes_tmp = object.get(init_x_0,"http://a.ml/vocabularies/core#name",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_0 = nodes_tmp2[_]
  nodes = x_0
}`
	testGeneratePropertyArray("core.name", expected, &expander, t, ast.RegoV1)
}

func TestGenerateAndPropertyArray(t *testing.T) {
	expander := misc.IriExpander{
		Context: types.ObjectMap{
			"x": "http://x.org#",
		},
	}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["http://x.org#b"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"http://x.org#c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`
	testGeneratePropertyArray("x.a / x.b / x.c", expected, &expander, t, ast.RegoV0)

	expected = `gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://x.org#a"])
  x_0 = tmp_x_0[_]
  tmp_x_2 = nested_nodes(x_0["http://x.org#b"])
  x_2 = tmp_x_2[_]
  nodes_tmp = object.get(x_2,"http://x.org#c",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`

	testGeneratePropertyArray("x.a / x.b / x.c", expected, &expander, t, ast.RegoV1)
}

func TestGenerateAndOrPropertyArray(t *testing.T) {
	expander := misc.IriExpander{
		Context: types.ObjectMap{
			"x": "http://x.org#",
		},
	}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["http://x.org#b"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"http://x.org#d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#a"]
  x_0 = tmp_x_0[_][_]
  tmp_x_2 = nested_nodes with data.nodes as x_0["http://x.org#c"]
  x_2 = tmp_x_2[_][_]
  nodes_tmp = object.get(x_2,"http://x.org#d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`
	testGeneratePropertyArray("x.a / x.b | x.c / x.d", expected, &expander, t, ast.RegoV0)
}

func TestGenerateAndOrParenthesisPropertyArray(t *testing.T) {
	expander := misc.IriExpander{
		Context: types.ObjectMap{
			"x": "http://x.org#",
		},
	}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#a"]
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"http://x.org#b",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2 = nodes_tmp2[_]
  nodes = x_2
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#c"]
  x_0 = tmp_x_0[_][_]
  nodes_tmp = object.get(x_0,"http://x.org#d",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_2 = nodes_tmp2[_]
  nodes = x_2
}`
	testGeneratePropertyArray("( x.a / x.b ) | ( x.c / x.d )", expected, &expander, t, ast.RegoV0)

	expected = `gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://x.org#a"])
  x_0 = tmp_x_0[_]
  nodes_tmp = object.get(x_0,"http://x.org#b",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_2 = nodes_tmp2[_]
  nodes = x_2
} {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes(init_x_0["http://x.org#c"])
  x_0 = tmp_x_0[_]
  nodes_tmp = object.get(x_0,"http://x.org#d",[])
  nodes_tmp2 = nodes_array(nodes_tmp)
  x_2 = nodes_tmp2[_]
  nodes = x_2
}`

	testGeneratePropertyArray("( x.a / x.b ) | ( x.c / x.d )", expected, &expander, t, ast.RegoV1)
}

func TestGenerateInversePathPropertyArray(t *testing.T) {
	expander := misc.IriExpander{Context: contexts.DefaultAMFContext}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "http://a.ml/vocabularies/core#name" with data.object as init_x_0
  nodes = x_0
}`
	testGeneratePropertyArray("core.name^", expected, &expander, t, ast.RegoV0)

	expected = `gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  x_0 = search_subjects("http://a.ml/vocabularies/core#name",init_x_0)[_]
  nodes = x_0
}`
	testGeneratePropertyArray("core.name^", expected, &expander, t, ast.RegoV1)
}

func TestGenerateAndInversePathPropertyArray(t *testing.T) {
	expander := misc.IriExpander{
		Context: types.ObjectMap{
			"x": "http://x.org#",
		},
	}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  tmp_x_0 = nested_nodes with data.nodes as init_x_0["http://x.org#a"]
  x_0 = tmp_x_0[_][_]
  search_subjects[x_2] with data.predicate as "http://x.org#b" with data.object as x_0
  nodes_tmp = object.get(x_2,"http://x.org#c",[])
  nodes_tmp2 = nodes_array with data.nodes as nodes_tmp
  x_3 = nodes_tmp2[_]
  nodes = x_3
}`
	testGeneratePropertyArray("x.a / x.b^ / x.c", expected, &expander, t, ast.RegoV0)
}

func TestGenerateOrInversePathPropertyArray(t *testing.T) {
	expander := misc.IriExpander{
		Context: types.ObjectMap{
			"x": "http://x.org#",
		},
	}
	expected := `gen_path_set_rule_1[nodes] {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "http://x.org#a" with data.object as init_x_0
  nodes = x_0
} {
  init_x_0 = data.sourceNode
  search_subjects[x_0] with data.predicate as "http://x.org#b" with data.object as init_x_0
  nodes = x_0
}`
	testGeneratePropertyArray("x.a^ | x.b^", expected, &expander, t, ast.RegoV0)

	expected = `gen_path_set_rule_1 contains nodes if {
  init_x_0 = data.sourceNode
  x_0 = search_subjects("http://x.org#a",init_x_0)[_]
  nodes = x_0
} {
  init_x_0 = data.sourceNode
  x_0 = search_subjects("http://x.org#b",init_x_0)[_]
  nodes = x_0
}`
	testGeneratePropertyArray("x.a^ | x.b^", expected, &expander, t, ast.RegoV1)
}

func testGeneratePropertyArray(pathString string, expected string, iriExpander *misc.IriExpander, t *testing.T, regoVersion ast.RegoVersion) {
	profile.GenReset()
	p, err := path.ParsePath(pathString)
	if err != nil {
		t.Errorf("error parsing path %v", err)
	}
	result := GeneratePropertySet(p, "x", iriExpander, regoVersion)
	actual := strings.Join(result.rego, "\n")
	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Errorf("unexpected rego code, expected:\n%s----actual:\n%s", expected, actual)
	}
}
