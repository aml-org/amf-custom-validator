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
	dataPath := os.Args[2]

	profile, err := ioutil.ReadFile(profilePath)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}

	debug := false
	if len(os.Args) > 3 {
		parsedDebug, ok := strconv.ParseBool(os.Args[3])
		if ok != nil {
			panic("usage validate PROFILE_PATH FILE_PATH [DEBUG=true|false]")
		}
		debug = parsedDebug
	}

	res, err := validator.Validate(string(profile), string(data), debug, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
