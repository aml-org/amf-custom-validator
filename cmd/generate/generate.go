package main

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	profilePath := os.Args[1]

	profile, err := ioutil.ReadFile(profilePath)
	if err != nil {
		panic(err)
	}

	debug := false
	if len(os.Args) > 2 {
		parsedDebug, ok := strconv.ParseBool(os.Args[2])
		if ok != nil {
			panic("usage generate PROFILE_PATH [DEBUG=true|false]")
		}
		debug = parsedDebug
	}

	err, module := validator.Generate(string(profile), debug)
	if err != nil {
		panic(err)
	}

	fmt.Println(module.Code)
}
