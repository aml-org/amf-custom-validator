package generator

type GeneratedRegoResult interface {
	ConstraintId() string
}

type SimpleRegoResult struct {
	Constraint string
	Rego []string
	Path string
	Value string
	Variable string
	Trace string
	PathRules []RegoPathResult
}

func (r SimpleRegoResult) ConstraintId() string {
	return r.Constraint
}

type BranchRegoResult struct {
	Constraint string
	Branch []SimpleRegoResult
}

func (r BranchRegoResult) ConstraintId() string {
	return r.Constraint
}
