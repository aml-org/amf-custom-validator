package profile

import (
	"errors"
	"github.com/aml-org/amfopa/internal/parser/profile/expression"
	"github.com/aml-org/amfopa/internal/parser/profile/statements"
	"github.com/aml-org/amfopa/internal/parser/yaml"
	y "github.com/kylelemons/go-gypsy/yaml"
)

func Parse(doc y.Node) (statements.Profile, error) {
	profile := statements.NewProfile()
	switch m := doc.(type) {
	case y.Map:

		name, err := yaml.GetString(m,"profile")
		if err != nil {
			return profile, err
		}
		profile.Name = name

		description, err := yaml.GetString(m,"description")
		if err == nil {
			profile.Description = &description
		}

		validations, err := parseValidations(m)
		if err != nil {
			return profile, err
		}

		violations, err := parseValidationLevel("violation", m, validations)
		if err != nil {
			return profile, err
		}
		for _,rule := range violations {
			profile.Violation = append(profile.Violation, rule)
		}

		warnings, err := parseValidationLevel("warning", m, validations)
		if err != nil {
			return profile, err
		}
		for _,rule := range warnings {
			profile.Warning = append(profile.Violation, rule)
		}

		infos, err := parseValidationLevel("info", m, validations)
		if err != nil {
			return profile, err
		}
		for _,rule := range infos {
			profile.Info = append(profile.Violation, rule)
		}
		return profile, nil
	default:
		return profile, errors.New("expected map at profile YAML document")
	}
}

func parseValidationLevel(level string, profile y.Map, validations y.Map) ([]statements.Rule, error) {
	validationsLevel := profile.Key(level)
	if validationsLevel != nil {
		switch names := validationsLevel.(type) {
		case y.List:
			rules := make([]statements.Rule,names.Len())
			for i,name := range names {
				switch n := name.(type) {
				case y.Scalar:
					k := n.String()
					switch v := validations.Key(k).(type) {
					case y.Map:
						r, err := expression.Parse(k, v, level)
						if err != nil {
							return nil, err
						}
						rules[i] = r
					default:
						// ignore
					}
				}
			}
			return rules,nil

		default:
			return nil,errors.New("Expected a list of validation names for violation level " + level)
		}
	}

	return make([]statements.Rule,0),nil
}

func parseValidations(node y.Map) (y.Map,error) {
	value := node.Key("validations")
	if value == nil {
		return nil,errors.New("'validations' keyword without values")
	} else {
		switch validations := value.(type) {
		case y.Map:
			return validations,nil
		default:
			return nil, errors.New("expected map of validations not found")
		}
	}
}