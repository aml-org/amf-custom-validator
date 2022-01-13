package profile

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/parser/path"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
	"strings"
)

type DatatypeRule struct {
	AtomicStatement
	Argument string
}

func (r DatatypeRule) Negate() Rule {
	negated := r
	negated.Negated = !r.Negated
	return negated
}

func (r DatatypeRule) String() string {
	var negation = ""
	if r.Negated {
		negation = "¬"
	}

	return fmt.Sprintf("%sDatatype(%s,'%s','%s')", negation, r.Variable.Name, r.Path.Source(), r.Argument)
}

func parseDatatype(negated bool, variable Variable, path path.PropertyPath, argument *y.Yaml) (DatatypeRule, error) {
	dt, err := argument.String()
	if err != nil {
		return DatatypeRule{}, err
	}
	r := DatatypeRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Name:    "datatype",
				Negated: negated,
			},
			Variable: variable,
			Path:     path,
		},
		Argument: strings.ReplaceAll(dt, ".", ":"),
	}
	return r, nil
}
