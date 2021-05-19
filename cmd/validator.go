package main

import (
	"github.com/aml-org/amfopa/internal/validator"
	"io/ioutil"
	"os"
)

func main() {
	profilePath := os.Args[1]
	dataPath := os.Args[2]

	profile, err := ioutil.ReadFile(profilePath)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}

	res, err := validator.Validate(string(profile), string(data), true)
	if err != nil {
		panic(err)
	}

	println(res)
}
