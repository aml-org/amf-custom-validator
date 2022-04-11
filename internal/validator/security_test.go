package validator

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"testing"
)

const securityPath relativePath = "../../test/data/security"

func tryCompileProfileIn(directory relativePath) (ast.Errors, bool) {
	resolvedDirectory := fmt.Sprintf("%s", directory)
	profile := fmt.Sprintf("%s/profile.yaml", resolvedDirectory)

	_, err := ProcessProfile(read(profile), false, nil)

	errors, ok := err.(ast.Errors)
	return errors, ok
}

func assertUnsafeBuiltInDoesNotCompile(directory relativePath, unsafeBuiltInName string, t *testing.T) {
	errors, ok := tryCompileProfileIn(directory)

	if !ok {
		t.Errorf("Expected ast.Errors to be returned")
		return
	}

	if len(errors) == 0 {
		t.Errorf("OPA compilation errors is empty")
		return
	}

	success := false

	expectedMsg := fmt.Sprintf("unsafe built-in function calls in expression: %s", unsafeBuiltInName)

	for _, e := range errors {
		isTypeError := e.Code == ast.TypeErr
		messageMatchesExpected := e.Message == expectedMsg
		if isTypeError && messageMatchesExpected {
			success = true
			break
		}
	}

	if !success {
		t.Errorf("Failed %s. OPA compilation did not produce 'unsafe built-in' errors for %s", directory, unsafeBuiltInName)
		return
	}
}

func testSecurityCase(subDirectory relativePath, unsafeBuiltInName string, t *testing.T) {
	directory := fmt.Sprintf("%s/%s", securityPath, subDirectory)
	assertUnsafeBuiltInDoesNotCompile(directory, unsafeBuiltInName, t)
}

func TestGraphWalk(t *testing.T) {
	testSecurityCase("graph.walk", "walk", t)
}

func TestHttpSend(t *testing.T) {
	testSecurityCase("http.send", "http.send", t)
}

// Uncomment when bumping OPA version
//func TestNetLookup_ip_addr(t *testing.T) {
//	testSecurityCase("net.lookup_ip_addr", "net.lookup_ip_addr", t)
//}

func TestOpaRuntime(t *testing.T) {
	testSecurityCase("opa.runtime", "opa.runtime", t)
}

func TestRegoParse_module(t *testing.T) {
	testSecurityCase("rego.parse_module", "rego.parse_module", t)
}
