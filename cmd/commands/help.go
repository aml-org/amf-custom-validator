package commands

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Printf(`
Available commands:
  help					Prints help
  validate PROFILE DATA			Validates data in JSON-LD syntax with a specific validation profile
  normalize DATA			Normalizes data to be validated with the generated Rego code
  generate PROFILE			Generates Rego code from a validation profile
`)
	os.Exit(0)
}
