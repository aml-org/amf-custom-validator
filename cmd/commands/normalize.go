package commands

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"os"
)

func Normalize() {
	helpers.ValidateNArgs(3, "acv normalize DATA")

	data := helpers.ReadOrPanic(os.Args[2])

	res, err := validator.ProcessInput(string(data), false, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(validator.Encode(res))
	os.Exit(0)
}
