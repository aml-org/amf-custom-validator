package path

import "errors"

type PropertyPath interface {
	Source() string
}

type BasePath struct {
	source string
}

func (p BasePath) Source() string {
	return p.source
}

type Property struct {
	BasePath
	Iri        string
	Inverse    bool
	Transitive bool
}

func (p Property) Source() string {
	return p.source
}

type AndPath struct {
	BasePath
	And []PropertyPath
}

func (p AndPath) Source() string {
	return p.source
}

type OrPath struct {
	BasePath
	Or []PropertyPath
}

func (p OrPath) Source() string {
	return p.source
}

type NullPath struct {
	BasePath
}

func (p NullPath) Source() string {
	return p.source
}

func build(source string, parsed interface{}) PropertyPath {
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
