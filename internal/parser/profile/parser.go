package profile

import (
	"errors"
	y "github.com/smallfish/simpleyaml"
)

func Parse(doc *y.Yaml) (Profile, error) {
	varGenerator := NewVarGenerator()
	profile := NewProfile()
	if doc.IsMap() {
		name, err := doc.Get("profile").String()
		if err != nil {
			return profile, err
		}
		profile.Name = name

		description, err := doc.Get("description").String()
		if err == nil {
			profile.Description = &description
		}

		validations := doc.Get("validations")
		if validations.IsFound() && !validations.IsMap() {
			return profile, errors.New("validations must be a list of validations")
		}

		violations, err := parseValidationLevel("violation", doc, validations, &varGenerator)
		if err != nil {
			return profile, err
		}
		for _, rule := range violations {
			profile.Violation = append(profile.Violation, rule)
		}

		warnings, err := parseValidationLevel("warning", doc, validations, &varGenerator)
		if err != nil {
			return profile, err
		}
		for _, rule := range warnings {
			profile.Warning = append(profile.Violation, rule)
		}

		infos, err := parseValidationLevel("info", doc, validations, &varGenerator)
		if err != nil {
			return profile, err
		}
		for _, rule := range infos {
			profile.Info = append(profile.Violation, rule)
		}
		return profile, nil
	}

	return profile, errors.New("expected map at profile YAML document")
}

func parseValidationLevel(level string, profile *y.Yaml, validations *y.Yaml, varGenerator *VarGenerator) ([]Rule, error) {
	var rules []Rule

	names := profile.Get(level)
	if !names.IsFound() {
		return rules, nil
	}
	size, err := names.GetArraySize()
	if err != nil {
		return rules, nil
	}
	for i := 0; i < size; i++ {
		name, err := names.GetIndex(i).String()
		if err == nil {
			v := validations.Get(name)
			if v.IsFound() {
				r, err := ParseExpression(name, v, level, varGenerator)
				if err != nil {
					return nil, err
				}
				rules = append(rules, r)
			}
		}
	}
	return rules, nil
}
