package statements

import "fmt"

type VarGenerator struct {
	vars []string
	counter int
}

func NewVarGenerator() VarGenerator {
	vs :=  []string{"x", "y", "z", "p", "q", "r", "s", "t", "u", "v", "w", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"}
	return VarGenerator{
		vars: vs,
		counter: 0,
	}
}

var globalGenerator = NewVarGenerator()

func Genvar(hint string) string {
	globalGenerator.counter++;
	return fmt.Sprintf("gen_%s_%d", hint, globalGenerator.counter)
}

func GenReset() {
	globalGenerator.counter = 0;
}

func (g VarGenerator) GenExpressionVar(quantification Quantification, cardinality *VariableCardinality) Variable {
	var name string

	if g.counter < len(g.vars) {
		name = g.vars[g.counter]
	} else {
		name = fmt.Sprintf("X%d", g.counter)
	}

	v := Variable{
		Name:           name,
		Quantification: quantification,
		Cardinality:    cardinality,
	}

	g.counter++

	return v
}


