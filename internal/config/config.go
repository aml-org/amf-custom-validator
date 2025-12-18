package config

const (
	Override = false
	Debug    = false
)

var regoVersion = "v0" // default Rego version

// SetRegoVersion sets the Rego version to use for code generation
func SetRegoVersion(v string) {
	regoVersion = v
}

// GetRegoVersion returns the current Rego version
func GetRegoVersion() string {
	return regoVersion
}

// IsRegoV1 returns true if the configured Rego version is v1
func IsRegoV1() bool {
	return regoVersion == "v1"
}

// IsRegoV0 returns true if the configured Rego version is v0
func IsRegoV0() bool {
	return regoVersion == "v0"
}
