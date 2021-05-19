package main

import (
	"github.com/aml-org/amfopa/internal/validator"
	"syscall/js"
)

func validateWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("Inside the function")
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		res, err := validator.Validate(profileString, dataString, true)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func main() {
	f := validateWrapper()
	js.Global().Set("validateCustomProfile", f)
	<-make(chan bool)
}
