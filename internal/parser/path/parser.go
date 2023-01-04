package path

import (
	"errors"
)

func build(source string, parsed any) PropertyPath {
	switch v := parsed.(type) {
	case IRI:
		return Property{
			BasePath: BasePath{
				source: source,
			},
			Iri:        v.Value,
			Inverse:    v.Inverse,
			Transitive: v.Transitive,
		}
	case AND:
		acc := make([]PropertyPath, 0)
		for _, a := range v.body {
			r := build(source, a)
			acc = append(acc, r)
		}
		return AndPath{
			BasePath: BasePath{
				source: source,
			},
			And: acc,
		}
	case OR:
		acc := make([]PropertyPath, 0)
		for _, o := range v.body {
			r := build(source, o)
			acc = append(acc, r)
		}
		return OrPath{
			BasePath: BasePath{
				source: source,
			},
			Or: acc,
		}
	default:
		panic(errors.New("unknown patch component returned by low-level parser"))
	}
}

func ParsePath(path string) (PropertyPath, error) {
	if path == "" {
		return NullPath{
			BasePath{
				source: path,
			},
		}, nil
	}
	parsed, err := Parse("", []byte(path))
	if err != nil {
		panic(err)
	}

	propertyPath := build(path, parsed)

	return propertyPath, nil
}
