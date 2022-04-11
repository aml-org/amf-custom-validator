package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
)

type UniqueValuesRule struct {
	AtomicStatement
	Argument bool
}

func (r UniqueValuesRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r UniqueValuesRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%s%s(%s,'%s','%t')", negation, r.Name, r.Variable.Name, r.Path.Source(), r.Argument)
}

func newUniqueValues(negated bool, variable Variable, path path.PropertyPath, argument bool) UniqueValuesRule {
	return UniqueValuesRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: negated,
				Name:    "uniqueValues",
			},
			Variable: variable,
			Path:     path,
		},
		Argument: argument,
	}
}
