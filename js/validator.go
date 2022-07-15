package main

import (
	"github.com/aml-org/amf-custom-validator/internal/validator"
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
		res, err := validator.Validate(profileString, dataString, debug, nil)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func generateWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		module, err := validator.GenerateRego(profileString, false, nil)
		if err != nil {
			return err.Error()
		}
		return module.Code
	})
	return jsonFunc
}

func normalizeInputWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		inputJsonLd := args[0].String()
		res, err := validator.ProcessInput(inputJsonLd, false, nil)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func exitWrapper(c chan bool) js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c <- true
		return nil
	})
	return jsonFunc
}

func main() {
	c := make(chan bool)
	f := validateWrapper()
	js.Global().Set("__AMF__validateCustomProfile", f)
	f = generateWrapper()
	js.Global().Set("__AMF__generate", f)
	f = normalizeInputWrapper()
	js.Global().Set("__AMF__normalize", f)
	f = exitWrapper(c)
	js.Global().Set("__AMF__terminateValidator", f)
	<-c
}
