package profile

import (
	"errors"
	"fmt"
	y "github.com/aml-org/amf-custom-validator/internal/parser/yaml"
)

func Parse(doc *y.Yaml) (Profile, error) {
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

		customRego, err := doc.Get("rego_extensions").String()
		if err == nil {
			profile.CustomRego = &customRego
		}

		prefixes := doc.Get("prefixes")
		if prefixes.IsFound() {
			context, err := ParsePrefixes(prefixes)
			if err != nil {
				return profile, err
			}
			profile.Prefixes = context
		}

		validations := doc.Get("validations")
		if !validations.IsFound() || !validations.IsMap() {
			return profile, errors.New("validations must be a map of validations")
		}

		violations, err := parseValidationLevel("violation", doc, validations)
		if err != nil {
			return profile, err
		}
		for _, rule := range violations {
			profile.Violation = append(profile.Violation, rule)
		}

		warnings, err := parseValidationLevel("warning", doc, validations)
		if err != nil {
			return profile, err
		}
		for _, rule := range warnings {
			profile.Warning = append(profile.Warning, rule)
		}

		infos, err := parseValidationLevel("info", doc, validations)
		if err != nil {
			return profile, err
		}
		for _, rule := range infos {
			profile.Info = append(profile.Info, rule)
		}
		return profile, nil
	}

	l, c := doc.Pos()
	return profile, errors.New(fmt.Sprintf("expected map at profile YAML document, [%d,%d]", l, c))
}

func parseValidationLevel(level string, profile *y.Yaml, validations *y.Yaml) ([]Rule, error) {
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
				varGenerator := NewVarGenerator()
				r, err := ParseExpression(name, v, level, &varGenerator)
				if err != nil {
					return nil, err
				}
				rules = append(rules, r)
			}
		}
	}
	return rules, nil
}
