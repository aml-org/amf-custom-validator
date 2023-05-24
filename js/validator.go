package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"syscall/js"
)

func validateWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
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

func genRegoWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		res, err := validator.GenerateRego(profileString, false, nil)
		if err != nil {
			return err.Error()
		}
		return res.Code
	})
	return jsonFunc
}

func genRegoWASMWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		wasm, err := validator.ProcessProfileWASM(profileString, false, nil)
		if err != nil {
			return err.Error()
		}
		wasm_string := hex.EncodeToString(wasm)
		return wasm_string
	})
	return jsonFunc
}

func normalizeInputWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		dataString := args[0].String()
		normalizedInput, err := validator.ProcessInput(dataString, false, nil)
		if err != nil {
			return err.Error()
		}
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetIndent("", "  ")
		err = enc.Encode(normalizedInput)
		if err != nil {
			return err.Error()
		}

		return b.String()
	})
	return jsonFunc
}

func exitWrapper(c chan bool) js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		c <- true
		return nil
	})
	return jsonFunc
}

func main() {
	c := make(chan bool)
	// validate
	f := validateWrapper()
	js.Global().Set("__AMF__validateCustomProfile", f)
	// gen rego
	f = genRegoWrapper()
	js.Global().Set("__AMF__generateRego", f)
	// gen rego WASM
	f = genRegoWASMWrapper()
	js.Global().Set("__AMF__generateRegoWASM", f)
	// normalizeInput
	f = normalizeInputWrapper()
	js.Global().Set("__AMF__normalizeInput", f)
	// exit
	f = exitWrapper(c)
	js.Global().Set("__AMF__terminateValidator", f)
	<-c
}
