package commands

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"os"
)

func Compile() {
	helpers.ValidateNArgs(3, "acv compile PROFILE")

	profile := helpers.ReadOrPanic(os.Args[2])

	_, err := validator.ProcessProfile(string(profile), true, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Compile Success!")
	os.Exit(0)
}
