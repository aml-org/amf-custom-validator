package main

import (
	"github.com/aml-org/amfopa/internal/validator"
	"os"
)

func main() {
	profile := os.Args[1]
	data := os.Args[2]

	res, err := validator.Validate(profile, data, true)
	if err != nil {
		panic(err)
	}

	println(res)
}
