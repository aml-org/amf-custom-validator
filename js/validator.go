// +build js,wasm

package main

import (
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"github.com/aml-org/amf-custom-validator/pkg"
	"github.com/open-policy-agent/opa/rego"
	"syscall/js"
	"unsafe"
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

func compileProfileWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		debug := args[1].Bool()
		compiledProfile, err := pkg.CompileProfile(profileString, debug, nil)
		if err != nil {
			return err.Error()
		}
		return js.ValueOf(unsafe.Pointer(compiledProfile))
	})
	return jsonFunc
}

func validateCompiledWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "Invalid no of arguments passed"
		}

		// Get the memory address
		pointerAddress := uintptr(args[0].Float())

		// Cast the raw memory address to a *rego.PreparedEvalQuery
		compiledProfile := (*rego.PreparedEvalQuery)(unsafe.Pointer(pointerAddress))

		dataString := args[1].String()
		debug := args[2].Bool()
		res, err := pkg.ValidateCompiled(compiledProfile, dataString, debug, nil)
		if err != nil {
			return err.Error()
		}
		return js.ValueOf(res)
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
	validate := validateWrapper()
	validateCompiled := validateCompiledWrapper()
	compileProfile := compileProfileWrapper()
	exit := exitWrapper(c)
	js.Global().Set("__AMF__validate", validate)
	js.Global().Set("__AMF__validateCompiled", validateCompiled)
	js.Global().Set("__AMF__compileProfile", compileProfile)
	js.Global().Set("__AMF__exit", exit)
	<-c
}
