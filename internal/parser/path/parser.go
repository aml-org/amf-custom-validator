package path

import "strings"

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
	Iri string
	Inverse bool
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


func ParsePath(path string) (PropertyPath, error) {
	parsed := Property{
		BasePath: BasePath{
			source: path,
		},
		Iri: strings.ReplaceAll(path, ".", ":"),
		Inverse: false,
		Transitive: false,
	}

	return parsed, nil
}