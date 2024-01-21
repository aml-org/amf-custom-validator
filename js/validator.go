//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aml-org/amf-custom-validator/internal/validator"
	"github.com/aml-org/amf-custom-validator/internal/validator/config"
	"log"
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

func validateWithConfigurationWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return "Invalid no of arguments passed"
		}
		profileString := args[0].String()
		dataString := args[1].String()
		debug := args[2].Bool()
		if !args[3].IsUndefined() {
			return "ValidationConfiguration not yet supported, please set this parameter as 'undefined'. ValidationConfiguration needs to import functions from JS which is not supported by Go 1.19"
		}
		reportConfig := buildReportConfiguration(args[4])
		res, err := validator.ValidateWithConfiguration(profileString, dataString, debug, nil, config.DefaultValidationConfiguration{}, reportConfig)
		if err != nil {
			return err.Error()
		}
		return res
	})
	return jsonFunc
}

func buildReportConfiguration(value js.Value) config.ReportConfiguration {
	includeReportCreationTime := value.Get("IncludeReportCreationTime").Bool()
	reportSchemaIri := value.Get("ReportSchemaIri").String()
	lexicalSchemaIri := value.Get("LexicalSchemaIri").String()
	return config.ReportConfiguration{
		IncludeReportCreationTime: includeReportCreationTime,
		ReportSchemaIri:           reportSchemaIri,
		LexicalSchemaIri:          lexicalSchemaIri,
	}
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
	fmt.Println("Something really cool is going on here...")
	f := validateWrapper()
	js.Global().Set("__AMF__validateCustomProfile", f)
	log.Println("assigned validateCustomProfile")
	log.Println(js.Global())
	// validate with configuration
	f = validateWithConfigurationWrapper()
	js.Global().Set("__AMF__validateCustomProfileWithConfiguration", f)
	// gen rego
	f = genRegoWrapper()
	js.Global().Set("__AMF__generateRego", f)
	// normalizeInput
	f = normalizeInputWrapper()
	js.Global().Set("__AMF__normalizeInput", f)
	// exit
	f = exitWrapper(c)
	js.Global().Set("__AMF__terminateValidator", f)
	<-c
}
