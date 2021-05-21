package profile

import "fmt"

type Quantification int

const (
	ForAll Quantification = iota
	Exists
)

type CardinalityOperation int

const (
	LTEQ CardinalityOperation = iota
	LT
	EQ
	NEQ
	GT
	GTEQ
)

type VariableCardinality struct {
	Operator CardinalityOperation
	Value    int
}

func (c VariableCardinality) String() string {
	switch c.Operator {
	case GTEQ:
		return fmt.Sprintf(">= %d", c.Value)
	case GT:
		return fmt.Sprintf("> %d", c.Value)
	case EQ:
		return fmt.Sprintf("= %d", c.Value)
	case NEQ:
		return fmt.Sprintf("<> %d", c.Value)
	case LT:
		return fmt.Sprintf("< %d", c.Value)
	case LTEQ:
		return fmt.Sprintf("<= %d", c.Value)
	default:
		panic("Unknown Cardinality ")
	}
}

type Variable struct {
	Quantification Quantification
	Name           string
	Cardinality    *VariableCardinality
}

func (v Variable) String() string {
	if v.Quantification == ForAll {
		return "∀" + v.Name
	} else if v.Cardinality != nil {
		return "∃" + v.Name + ";" + v.Cardinality.String()
	} else {
		return "∃" + v.Name
	}
}
