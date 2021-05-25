package generator

type GeneratedRegoResult interface {
	ConstraintId() string
}

type SimpleRegoResult struct {
	Constraint string           // constraint being evaluated
	Rego       []string         // rego code with the constraint
	Path       string           // path from the parent node to this node
	Variable   string           // variable with the result for the next evaluation
	TraceValue string           // evidence value for tracing
	TraceNode  string           // trace code
	PathRules  []RegoPathResult //path rules to generate the path rule
}

func (r SimpleRegoResult) ConstraintId() string {
	return r.Constraint
}

type BranchRegoResult struct {
	Constraint string
	Branch     []SimpleRegoResult
}

func (r BranchRegoResult) ConstraintId() string {
	return r.Constraint
}
