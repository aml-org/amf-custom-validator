package config

import "github.com/open-policy-agent/opa/ast"

const (
	Override = false
	Debug    = false
)

var RegoVersions = [2]ast.RegoVersion{ast.RegoV0, ast.RegoV1}

const RegoVersion1 = ast.RegoV1
