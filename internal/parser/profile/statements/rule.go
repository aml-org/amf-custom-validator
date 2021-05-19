package statements

import (
	"github.com/aml-org/amfopa/internal/parser/path"
	"strings"
)

type Rule interface {
	Negate() Rule
	String() string
	Clone() Rule
}

type BaseStatement struct {
	Negated bool
	Name    string
}

type ComplexStatement struct {
	BaseStatement
	Body []Rule
}

type AtomicStatement struct {
	BaseStatement
	Variable Variable
	Path     path.PropertyPath
}

type Hashable interface {
	ValueHash() string
}

type Profile struct {
	Name        string
	Description *string
	Violation   []Rule
	Warning     []Rule
	Info        []Rule
}

func NewProfile() Profile {
	return Profile{
		Name:        "",
		Description: nil,
		Violation:   make([]Rule, 0),
		Warning:     make([]Rule, 0),
		Info:        make([]Rule, 0),
	}
}

func (p Profile) String() string {
	lines := make([]string, len(p.Warning)+len(p.Info)+len(p.Violation))
	for _, v := range p.Violation {
		lines = append(lines, v.String())
	}

	for _, v := range p.Warning {
		lines = append(lines, v.String())
	}

	for _, v := range p.Info {
		lines = append(lines, v.String())
	}

	return strings.Join(lines, "\n\n")
}
