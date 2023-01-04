package profile

import "testing"

// message parsing
func TestSingleVariableMessageParsing(t *testing.T) {
	rawExpression := "Endpoint is {{ apiContract.path }}"

	expected := Message{
		Expression: "Endpoint is %v",
		Variables:  []string{"apiContract.path"},
	}

	actual := ParseMessageExpression(rawExpression)
	if expected.Compare(actual) != 0 {
		t.Errorf("Actual did not match expected")
	}
}

func TestDoubleVariableMessageParsing(t *testing.T) {
	rawExpression := "API {{ core.name }} {{ core.version }} is incorrect"

	expected := Message{
		Expression: "API %v %v is incorrect",
		Variables:  []string{"core.name", "core.version"},
	}

	actual := ParseMessageExpression(rawExpression)
	if expected.Compare(actual) != 0 {
		t.Errorf("Actual did not match expected")
	}
}

func TestSingleVariableWithWhitespaceMessageParsing(t *testing.T) {
	rawExpression := "Endpoint is {{              apiContract.path         }}"

	expected := Message{
		Expression: "Endpoint is %v",
		Variables:  []string{"apiContract.path"},
	}

	actual := ParseMessageExpression(rawExpression)
	if expected.Compare(actual) != 0 {
		t.Errorf("Actual did not match expected")
	}
}
