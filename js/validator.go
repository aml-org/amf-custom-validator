// +build js,wasm

package main

import (
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"syscall/js"
)

func validate(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		return "Invalid no of arguments passed"
	}
	profileString := args[0].String()
	dataString := args[1].String()
	debug := args[2].Bool()
	res, err := validator.Validate(profileString, dataString, debug, nil)
	if err != nil {
		return err.Error()
	}
	return res
}

func main() {
	js.Global().Set("_acv_validate", js.FuncOf(validate))
	js.Global().Call("_acv_user_callback")
	fmt.Println("Returned to Go")
	return
}
