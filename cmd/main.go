package main

import (
	"github.com/aml-org/amf-custom-validator/cmd/commands"
	"os"
)

func main() {
	command := os.Args[1]
	switch command {
	case "validate":
		commands.Validate()
	case "generate":
		commands.Generate()
	case "normalize":
		commands.Normalize()
	case "help":
		commands.Help()
	default:
		commands.Fallback(command)
	}
}
