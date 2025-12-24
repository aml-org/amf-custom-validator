package commands

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Printf(`
Available commands:
  help										Prints help
  validate <PROFILE> <DATA> [OUTPUT_FILE]	Validates data in JSON-LD syntax with a specific validation profile
  normalize <DATA>							Normalizes data to be validated with the generated Rego code
  generate <PROFILE>						Generates Rego code from a validation profile
Environment variables:
  REGO_VERSION <v0|v1>                      Rego syntax expected and generated, by default v1
`)
	os.Exit(0)
}
