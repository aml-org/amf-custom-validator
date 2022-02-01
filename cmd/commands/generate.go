package commands

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"os"
)

func Generate() {
	helpers.ValidateNArgs(3, "acv generate PROFILE")

	profile := helpers.ReadOrPanic(os.Args[2])

	module, err := validator.GenerateRego(string(profile), false, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(module.Code)
	os.Exit(0)
}
