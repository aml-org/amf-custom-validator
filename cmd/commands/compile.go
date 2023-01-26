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

	compiledProfile, err := validator.ProcessProfile(string(profile), true, nil)

	if err != nil {
		_ = fmt.Errorf(err.Error())
		os.Exit(1)
	}
	fmt.Println(compiledProfile)
	os.Exit(0)
}
