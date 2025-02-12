package commands

import (
	"fmt"
	h "github.com/aml-org/amf-custom-validator/cmd/commands/helpers"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"os"
)

func Validate() {

	h.ValidateNsArgs([]int{4, 5}, "acv validate <PROFILE> <DATA> [OUTPUT_FILE]")

	profile := h.ReadOrPanic(os.Args[2])
	data := h.ReadOrPanic(os.Args[3])

	res, err := validator.Validate(string(profile), string(data), false, nil)
	h.CheckError(err)

	if len(os.Args) == 4 {
		// 4 args means output to standard output
		fmt.Println(res)
	}
	if len(os.Args) == 5 {
		// 5 args means output to file path
		path := os.Args[4]
		file := h.OpenOrCreateFile(path)
		h.WriteString(file, res)
	}

	os.Exit(0)
}
