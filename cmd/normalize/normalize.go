package main

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	dataPath := os.Args[1]

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}

	debug := false
	if len(os.Args) > 2 {
		parsedDebug, ok := strconv.ParseBool(os.Args[3])
		if ok != nil {
			panic("usage normalize FILE_PATH [DEBUG=true|false]")
		}
		debug = parsedDebug
	}

	res, err := validator.ProcessInput(string(data), debug, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(validator.Encode(res))
}
