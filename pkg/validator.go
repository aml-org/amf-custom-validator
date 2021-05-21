package main

import (
	"github.com/aml-org/amfopa/internal/validator"
	"syscall/js"
)

func validateWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()
		res, err := validator.Validate(profileString, dataString, debug)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func main() {
	f := validateWrapper()
	js.Global().Set("__AMF__validateCustomProfile", f)
	<-make(chan bool)
}
