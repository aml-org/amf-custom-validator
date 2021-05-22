package profile

import (
	"fmt"
	"github.com/aml-org/amfopa/internal"
	"github.com/aml-org/amfopa/internal/parser/path"
	y "github.com/smallfish/simpleyaml"
	"strings"
)

type DatatypeRule struct {
	AtomicStatement
	Argument string
}

func (r DatatypeRule) Clone() Rule {
	return DatatypeRule{
		AtomicStatement: AtomicStatement{
			BaseStatement: BaseStatement{
				Negated: r.Negated,
				Name:    r.Name,
			},
			Variable: r.Variable,
			Path:     r.Path,
		},
		Argument: r.Argument,
	}
}

func (r DatatypeRule) Negate() Rule {
	cloned := r.Clone()
	switch c := cloned.(type) {
	case DatatypeRule:
		c.Negated = !r.Negated
		return c
	}
	return cloned
}

func (r DatatypeRule) ValueHash() string {
	v := fmt.Sprintf("%s%s", r.Name, r.Argument)
	return internal.HashString(v)
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