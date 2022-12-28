package commands

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"os"
)

func ValidateOptimized() {
	helpers.ValidateNArgs(4, "acv validateOptimized PROFILE DATA")

	profile := helpers.ReadOrPanic(os.Args[2])
	data := helpers.ReadOrPanic(os.Args[3])

	res, err := validator.ValidateWithOptimization(string(profile), string(data), false, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	os.Exit(0)
}
