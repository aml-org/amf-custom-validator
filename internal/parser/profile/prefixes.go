package profile

import (
	"errors"
	y "github.com/aml-org/amfopa/internal/parser/yaml"
)

type ProfileContext = map[string]string

// Captures the aliases and JSON-LD URI prefix for the profile
func ParsePrefixes(y *y.Yaml) (ProfileContext, error) {
	ctx := make(ProfileContext)

	if y.IsMap() {
		ks, err := y.GetMapKeys()
		if err != nil {
			return nil, err
		}
		for _, k := range ks {
			v, err := y.Get(k).String()
			if err != nil {
				return nil, err
			}
			ctx[k] = v
		}
		return ctx, nil
	} else {
		return nil, errors.New("context must be a map")
	}
}
